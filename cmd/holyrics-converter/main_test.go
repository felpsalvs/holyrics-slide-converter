package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"holyrics-slide-converter/internal/config"
)

func writeConfig(t *testing.T, cfg config.Config) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "config.json")
	data, err := json.Marshal(cfg)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestRunOnceSucesso(t *testing.T) {
	out := filepath.Join(t.TempDir(), "saida")
	cfgPath := writeConfig(t, config.Config{
		PastaEntrada: t.TempDir(),
		PastaSaida:   out,
		LarguraPx:    640,
	})

	var stdout, stderr bytes.Buffer
	code := run([]string{"--config", cfgPath, "--once", "../../tests/fixtures/sample.pdf"}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("exit code = %d, stderr = %s", code, stderr.String())
	}
	entries, err := os.ReadDir(filepath.Join(out, "sample"))
	if err != nil {
		t.Fatalf("pasta de saída não criada: %v", err)
	}
	if len(entries) == 0 {
		t.Fatal("nenhum PNG gerado")
	}
}

func TestRunOnceArquivoInexistente(t *testing.T) {
	cfgPath := writeConfig(t, config.Config{
		PastaEntrada: t.TempDir(),
		PastaSaida:   t.TempDir(),
		LarguraPx:    640,
	})

	var stdout, stderr bytes.Buffer
	code := run([]string{"--config", cfgPath, "--once", "nao-existe.pdf"}, &stdout, &stderr)

	if code != exitErro {
		t.Fatalf("exit code = %d, esperado %d", code, exitErro)
	}
}

func TestRunOnceExtensaoNaoSuportada(t *testing.T) {
	naoSuportado := filepath.Join(t.TempDir(), "arquivo.txt")
	if err := os.WriteFile(naoSuportado, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	cfgPath := writeConfig(t, config.Config{
		PastaEntrada: t.TempDir(),
		PastaSaida:   t.TempDir(),
		LarguraPx:    640,
	})

	var stdout, stderr bytes.Buffer
	code := run([]string{"--config", cfgPath, "--once", naoSuportado}, &stdout, &stderr)

	if code != exitErro {
		t.Fatalf("exit code = %d, esperado %d", code, exitErro)
	}
}

func TestRunConfigRecemCriado(t *testing.T) {
	cfgPath := filepath.Join(t.TempDir(), "sub", "config.json")

	var stdout, stderr bytes.Buffer
	code := run([]string{"--config", cfgPath}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("exit code = %d, stderr = %s", code, stderr.String())
	}
	if _, err := os.Stat(cfgPath); err != nil {
		t.Fatalf("config.json não foi criado: %v", err)
	}
}

func TestRunLogLevelInvalido(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := run([]string{"--log-level", "nivel-invalido"}, &stdout, &stderr)

	if code != exitUso {
		t.Fatalf("exit code = %d, esperado %d", code, exitUso)
	}
}
