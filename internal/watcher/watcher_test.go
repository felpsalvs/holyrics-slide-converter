package watcher

import (
	"os"
	"path/filepath"
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
