package watcher

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

func TestIgnorar(t *testing.T) {
	ignorados := []string{
		"~$culto.pptx", ".DS_Store", "baixando.pdf.crdownload",
		"foto.png", "nota.txt", "arquivo.tmp", "video.part",
	}
	for _, name := range ignorados {
		if !Ignorar(name) {
			t.Errorf("Ignorar(%q) = false, esperado true", name)
		}
	}
	aceitos := []string{"culto.pptx", "Estudo.PDF", "louvor.odp", "antigo.ppt"}
	for _, name := range aceitos {
		if Ignorar(name) {
			t.Errorf("Ignorar(%q) = true, esperado false", name)
		}
	}
}

func TestEsperarEstabilizar(t *testing.T) {
	path := filepath.Join(t.TempDir(), "a.pdf")
	if err := os.WriteFile(path, []byte("conteudo"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := EsperarEstabilizar(path, 10*time.Millisecond, time.Second); err != nil {
		t.Fatalf("arquivo estável retornou erro: %v", err)
	}
}

func TestEsperarEstabilizarCrescendo(t *testing.T) {
	path := filepath.Join(t.TempDir(), "a.pdf")
	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	done := make(chan struct{})
	go func() {
		defer close(done)
		defer f.Close()
		for i := 0; i < 5; i++ {
			f.WriteString("bloco de dados")
			f.Sync()
			time.Sleep(20 * time.Millisecond)
		}
	}()
	if err := EsperarEstabilizar(path, 30*time.Millisecond, 5*time.Second); err != nil {
		t.Fatalf("EsperarEstabilizar: %v", err)
	}
	<-done
	info, _ := os.Stat(path)
	if info.Size() != int64(5*len("bloco de dados")) {
		t.Fatalf("estabilizou antes da cópia terminar (size=%d)", info.Size())
	}
}

func TestEsperarEstabilizarInexistente(t *testing.T) {
	if err := EsperarEstabilizar(filepath.Join(t.TempDir(), "x.pdf"), time.Millisecond, time.Second); err == nil {
		t.Fatal("esperado erro para arquivo inexistente")
	}
}

// TestRunLimitaConcorrencia confirma que Run nunca chama proc mais vezes em
// paralelo do que o limite configurado, mesmo com vários arquivos prontos
// ao mesmo tempo — proteção contra N instâncias soffice headless disputando
// lock de perfil simultaneamente.
func TestRunLimitaConcorrencia(t *testing.T) {
	dir := t.TempDir()
	const arquivos = 4
	for i := 0; i < arquivos; i++ {
		name := filepath.Join(dir, fmt.Sprintf("arquivo-%d.pdf", i))
		if err := os.WriteFile(name, []byte("conteudo"), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	const limite = 2
	var (
		mu    sync.Mutex
		atual int
		pico  int
	)
	release := make(chan struct{})
	proc := func(path string) error {
		mu.Lock()
		atual++
		if atual > pico {
			pico = atual
		}
		mu.Unlock()
		<-release
		mu.Lock()
		atual--
		mu.Unlock()
		return nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() { done <- Run(ctx, dir, limite, proc) }()

	// Dá tempo para todos os arquivos estabilizarem e disputarem o semáforo.
	time.Sleep(1200 * time.Millisecond)
	close(release)
	cancel()
	<-done

	mu.Lock()
	defer mu.Unlock()
	if pico > limite {
		t.Fatalf("pico de concorrência = %d, esperado <= %d", pico, limite)
	}
	if pico < limite {
		t.Fatalf("pico de concorrência = %d, esperado exatamente %d (limite deveria ser atingido com %d arquivos)", pico, limite, arquivos)
	}
}
