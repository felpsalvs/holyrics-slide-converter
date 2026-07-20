package convert

import (
	"image/png"
	"os"
	"path/filepath"
	"testing"

	"holyrics-slide-converter/internal/soffice"
)

func TestSupportedExt(t *testing.T) {
	for _, ext := range []string{".pdf", ".pptx", ".PPTX", ".ppt", ".odp"} {
		if !SupportedExt(ext) {
			t.Errorf("SupportedExt(%q) = false, esperado true", ext)
		}
	}
	for _, ext := range []string{".txt", ".png", ".docx", ""} {
		if SupportedExt(ext) {
			t.Errorf("SupportedExt(%q) = true, esperado false", ext)
		}
	}
}

func TestFilePDF(t *testing.T) {
	out := t.TempDir()
	dest, n, err := File(filepath.Join("..", "..", "tests", "fixtures", "sample.pdf"), out, 1920, nil)
	if err != nil {
		t.Fatalf("File: %v", err)
	}
	if n != 3 {
		t.Fatalf("esperado 3 páginas, veio %d", n)
	}
	if dest != filepath.Join(out, "sample") {
		t.Fatalf("pasta de destino inesperada: %s", dest)
	}
	for _, name := range []string{"slide-001.png", "slide-002.png", "slide-003.png"} {
		f, err := os.Open(filepath.Join(dest, name))
		if err != nil {
			t.Fatalf("%s não gerado: %v", name, err)
		}
		img, err := png.Decode(f)
		f.Close()
		if err != nil {
			t.Fatalf("%s: PNG inválido: %v", name, err)
		}
		w := img.Bounds().Dx()
		if w < 1900 || w > 1940 {
			t.Errorf("%s: largura %dpx fora do alvo de 1920px", name, w)
		}
	}
}

func TestFileUnsupported(t *testing.T) {
	if _, _, err := File("nota.txt", t.TempDir(), 1920, nil); err == nil {
		t.Fatal("esperado erro para extensão não suportada")
	}
}

func TestFilePPTXNeedsSoffice(t *testing.T) {
	toPDF := func(path, outDir string) (string, error) {
		return soffice.ToPDF("", path, outDir)
	}
	if _, _, err := File("slides.pptx", t.TempDir(), 1920, toPDF); err == nil {
		t.Fatal("esperado erro quando LibreOffice não está configurado")
	}
}

// TestFileODP roda a conversão completa ODP->PDF->PNG quando o LibreOffice
// está instalado na máquina; caso contrário é pulado. Se o LibreOffice
// estiver instalado mas a fixture faltar, falha em vez de pular: um Skip
// silencioso aqui já deu falso senso de cobertura no passado.
func TestFileODP(t *testing.T) {
	sofficePath, err := soffice.Find("")
	if err != nil {
		t.Skipf("LibreOffice não instalado: %v", err)
	}
	src := filepath.Join("..", "..", "tests", "fixtures", "sample.odp")
	if _, err := os.Stat(src); err != nil {
		t.Fatalf("fixture %s ausente (LibreOffice está instalado, então este caminho deveria ser testado): %v", src, err)
	}
	out := t.TempDir()
	toPDF := func(path, outDir string) (string, error) {
		return soffice.ToPDF(sofficePath, path, outDir)
	}
	dest, n, err := File(src, out, 1280, toPDF)
	if err != nil {
		t.Fatalf("File: %v", err)
	}
	if n < 1 {
		t.Fatalf("nenhuma página gerada")
	}
	if _, err := os.Stat(filepath.Join(dest, "slide-001.png")); err != nil {
		t.Fatalf("slide-001.png não gerado: %v", err)
	}
}
