// Command holyrics-converter monitora uma pasta e converte apresentações e
// PDFs em imagens PNG para a biblioteca do Holyrics.
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"holyrics-slide-converter/internal/config"
	"holyrics-slide-converter/internal/convert"
	"holyrics-slide-converter/internal/soffice"
	"holyrics-slide-converter/internal/watcher"
)

const (
	exitOK   = 0
	exitErro = 1
	exitUso  = 2
)

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("holyrics-converter", flag.ContinueOnError)
	fs.SetOutput(stderr)
	once := fs.String("once", "", "converte um único arquivo e sai, sem monitorar a pasta")
	configPath := fs.String("config", "", "caminho alternativo para o config.json")
	logLevel := fs.String("log-level", "info", "nível de log: debug, info, warn, error")
	if err := fs.Parse(args); err != nil {
		return exitUso
	}

	level, err := parseLevel(*logLevel)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return exitUso
	}
	logger := slog.New(slog.NewTextHandler(stdout, &slog.HandlerOptions{Level: level}))
	slog.SetDefault(logger)

	path := *configPath
	if path == "" {
		p, err := config.DefaultPath()
		if err != nil {
			fmt.Fprintf(stderr, "não foi possível determinar o config.json: %v\n", err)
			return exitErro
		}
		path = p
	}

	cfg, created, err := config.Load(path)
	if err != nil {
		fmt.Fprintf(stderr, "config inválido: %v\n", err)
		return exitErro
	}
	if created {
		fmt.Fprintf(stdout, "Config criado em %s. Ajuste pasta_entrada e pasta_saida e rode novamente.\n", path)
		return exitOK
	}

	sofficePath, err := soffice.Find(cfg.CaminhoSoffice)
	if err != nil {
		slog.Warn("LibreOffice não encontrado; apenas PDFs poderão ser convertidos", "erro", err)
	}

	if *once != "" {
		dest, n, err := convert.File(*once, cfg.PastaSaida, cfg.LarguraPx, sofficePath)
		if err != nil {
			fmt.Fprintf(stderr, "falha ao converter %s: %v\n", *once, err)
			return exitErro
		}
		fmt.Fprintf(stdout, "%d página(s) geradas em %s\n", n, dest)
		return exitOK
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	proc := func(path string) error {
		_, _, err := convert.File(path, cfg.PastaSaida, cfg.LarguraPx, sofficePath)
		return err
	}

	if err := watcher.Run(ctx, cfg.PastaEntrada, proc); err != nil && !errors.Is(err, context.Canceled) {
		fmt.Fprintf(stderr, "watcher encerrado com erro: %v\n", err)
		return exitErro
	}
	return exitOK
}

func parseLevel(s string) (slog.Level, error) {
	var level slog.Level
	if err := level.UnmarshalText([]byte(s)); err != nil {
		return 0, fmt.Errorf("nível de log inválido: %s", s)
	}
	return level, nil
}
