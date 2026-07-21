package soffice

import (
	"errors"
	"os"
	"testing"
)

var errNotFound = errors.New("não encontrado")

func statOnly(existing ...string) func(string) (os.FileInfo, error) {
	set := map[string]bool{}
	for _, p := range existing {
		set[p] = true
	}
	return func(p string) (os.FileInfo, error) {
		if set[p] {
			return nil, nil
		}
		return nil, errNotFound
	}
}

func noLookPath(string) (string, error) { return "", errNotFound }

func TestFindConfigured(t *testing.T) {
	p, err := find("/x/soffice", "darwin", statOnly("/x/soffice"), noLookPath)
	if err != nil || p != "/x/soffice" {
		t.Fatalf("find = %q, %v", p, err)
	}
	if _, err := find("/x/soffice", "darwin", statOnly(), noLookPath); err == nil {
		t.Fatal("esperado erro para caminho configurado inexistente")
	}
}

func TestFindPath(t *testing.T) {
	lp := func(name string) (string, error) {
		if name == "soffice" {
			return "/usr/local/bin/soffice", nil
		}
		return "", errNotFound
	}
	p, err := find("", "linux", statOnly(), lp)
	if err != nil || p != "/usr/local/bin/soffice" {
		t.Fatalf("find = %q, %v", p, err)
	}
}

func TestFindWellKnown(t *testing.T) {
	cases := map[string]string{
		"windows": `C:\Program Files\LibreOffice\program\soffice.exe`,
		"darwin":  "/Applications/LibreOffice.app/Contents/MacOS/soffice",
		"linux":   "/usr/bin/libreoffice",
	}
	for goos, want := range cases {
		p, err := find("", goos, statOnly(want), noLookPath)
		if err != nil || p != want {
			t.Fatalf("%s: find = %q, %v (esperado %q)", goos, p, err, want)
		}
	}
}

func TestFindMissing(t *testing.T) {
	if _, err := find("", "darwin", statOnly(), noLookPath); err == nil {
		t.Fatal("esperado erro quando nada é encontrado")
	}
}

func TestFileURL(t *testing.T) {
	got := fileURL("/tmp/perfil-123")
	if got != "file:///tmp/perfil-123" {
		t.Fatalf("fileURL = %q, esperado file:///tmp/perfil-123", got)
	}
}

func TestToPDFSemSoffice(t *testing.T) {
	if _, err := ToPDF("", "slides.pptx", t.TempDir()); err == nil {
		t.Fatal("esperado erro quando sofficePath está vazio")
	}
}
