// Package watcher monitora a pasta de entrada e dispara a conversão
// quando novos arquivos terminam de ser copiados.
package watcher

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"

	"holyrics-slide-converter/internal/convert"
)

// Ignorar decide se um arquivo deve ser descartado pelo watcher:
// temporários de editores/navegadores, ocultos e extensões não suportadas.
func Ignorar(name string) bool {
	base := filepath.Base(name)
	if strings.HasPrefix(base, "~$") || strings.HasPrefix(base, ".") {
		return true
	}
	switch strings.ToLower(filepath.Ext(base)) {
	case ".tmp", ".crdownload", ".part", ".download":
		return true
	}
	return !convert.SupportedExt(filepath.Ext(base))
}

// EsperarEstabilizar aguarda o tamanho do arquivo parar de mudar (cópia
// concluída), consultando a cada intervalo, até o limite de timeout.
func EsperarEstabilizar(path string, intervalo, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	var last int64 = -1
	for {
		info, err := os.Stat(path)
		if err != nil {
			return err
		}
		if info.Size() == last && info.Size() > 0 {
			return nil
		}
		last = info.Size()
		if time.Now().After(deadline) {
			return fmt.Errorf("arquivo não estabilizou em %s: %s", timeout, path)
		}
		time.Sleep(intervalo)
	}
}

// Processar é a função chamada para cada arquivo pronto para conversão.
type Processar func(path string) error

// Run monitora dir até o ctx ser cancelado. Arquivos já presentes na pasta
// ao iniciar também são processados. Após converter, o original vai para
// dir/processados; em falha, para dir/erros. maxConcurrent limita quantas
// conversões (proc) rodam em paralelo; valores <= 0 são tratados como 1.
func Run(ctx context.Context, dir string, maxConcurrent int, proc Processar) error {
	for _, sub := range []string{dir, filepath.Join(dir, "processados"), filepath.Join(dir, "erros")} {
		if err := os.MkdirAll(sub, 0o755); err != nil {
			return err
		}
	}

	w, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer w.Close()
	if err := w.Add(dir); err != nil {
		return err
	}

	if maxConcurrent <= 0 {
		maxConcurrent = 1
	}
	sem := make(chan struct{}, maxConcurrent)

	var mu sync.Mutex
	emAndamento := map[string]bool{}

	handle := func(path string) {
		mu.Lock()
		if emAndamento[path] {
			mu.Unlock()
			return
		}
		emAndamento[path] = true
		mu.Unlock()
		defer func() {
			mu.Lock()
			delete(emAndamento, path)
			mu.Unlock()
		}()

		if err := EsperarEstabilizar(path, 500*time.Millisecond, 2*time.Minute); err != nil {
			if !errors.Is(err, os.ErrNotExist) {
				slog.Error("arquivo não estabilizou", "arquivo", filepath.Base(path), "erro", err)
			}
			return
		}

		sem <- struct{}{}
		defer func() { <-sem }()

		if err := proc(path); err != nil {
			slog.Error("falha ao converter", "arquivo", filepath.Base(path), "erro", err)
			mover(path, filepath.Join(dir, "erros"))
			return
		}
		mover(path, filepath.Join(dir, "processados"))
	}

	// Processa arquivos que já estavam na pasta antes do watcher iniciar.
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, e := range entries {
		if e.IsDir() || Ignorar(e.Name()) {
			continue
		}
		go handle(filepath.Join(dir, e.Name()))
	}

	slog.Info("Monitorando (Ctrl+C para sair)", "pasta", dir, "conversoes_simultaneas", maxConcurrent)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case ev, ok := <-w.Events:
			if !ok {
				return nil
			}
			if !ev.Has(fsnotify.Create) && !ev.Has(fsnotify.Write) {
				continue
			}
			if Ignorar(ev.Name) {
				continue
			}
			if info, err := os.Stat(ev.Name); err != nil || info.IsDir() {
				continue
			}
			go handle(ev.Name)
		case err, ok := <-w.Errors:
			if !ok {
				return nil
			}
			slog.Error("watcher", "erro", err)
		}
	}
}

func mover(path, destDir string) {
	dest := filepath.Join(destDir, filepath.Base(path))
	// Evita sobrescrever arquivo homônimo já movido anteriormente.
	if _, err := os.Stat(dest); err == nil {
		ext := filepath.Ext(dest)
		dest = strings.TrimSuffix(dest, ext) + "-" + time.Now().Format("20060102-150405") + ext
	}
	if err := os.Rename(path, dest); err != nil {
		slog.Error("não foi possível mover", "arquivo", filepath.Base(path), "erro", err)
	}
}
