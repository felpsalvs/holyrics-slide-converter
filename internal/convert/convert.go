// Package convert transforma apresentações (PPTX/PPT/ODP) e PDFs em imagens PNG.
package convert

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	"github.com/gen2brain/go-fitz"
)

// SupportedExt informa se a extensão (com ponto, minúscula ou não) é conversível.
func SupportedExt(ext string) bool {
	switch strings.ToLower(ext) {
	case ".pdf", ".pptx", ".ppt", ".odp":
		return true
	}
	return false
}

// ToPDFFunc converte uma apresentação em PDF, retornando o caminho do PDF
// gerado dentro de outDir. Normalmente é soffice.ToPDF associado a um
// caminho de LibreOffice.
type ToPDFFunc func(path, outDir string) (string, error)

// PageWriter recebe as páginas renderizadas de um PDF, uma a uma, na ordem.
// Existe para permitir destinos além do filesystem local no futuro (ex.: a
// API do Holyrics), sem alterar pdfToPNG.
type PageWriter interface {
	WritePage(index int, img image.Image) error
	Close() error
}

// File converte um arquivo suportado em PNGs dentro de outDir/<nome-do-arquivo>/.
// Apresentações passam antes por toPDF para virar PDF.
// Retorna a pasta criada e a quantidade de páginas geradas.
func File(path, outDir string, widthPx int, toPDF ToPDFFunc) (string, int, error) {
	ext := strings.ToLower(filepath.Ext(path))
	if !SupportedExt(ext) {
		return "", 0, fmt.Errorf("extensão não suportada: %s", ext)
	}

	pdfPath := path
	if ext != ".pdf" {
		tmp, err := os.MkdirTemp("", "holyrics-converter-*")
		if err != nil {
			return "", 0, err
		}
		defer os.RemoveAll(tmp)
		pdfPath, err = toPDF(path, tmp)
		if err != nil {
			return "", 0, err
		}
	}

	name := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	dest := filepath.Join(outDir, name)
	w, err := newDirWriter(dest)
	if err != nil {
		return "", 0, err
	}
	n, convErr := pdfToPNG(pdfPath, widthPx, w)
	closeErr := w.Close()
	if convErr != nil {
		return "", 0, convErr
	}
	if closeErr != nil {
		return "", 0, closeErr
	}
	return dest, n, nil
}

// pdfToPNG renderiza cada página do PDF e entrega a w, na ordem, com largura
// alvo widthPx.
func pdfToPNG(pdfPath string, widthPx int, w PageWriter) (int, error) {
	doc, err := fitz.New(pdfPath)
	if err != nil {
		return 0, fmt.Errorf("falha ao abrir PDF: %w", err)
	}
	defer doc.Close()

	total := doc.NumPage()
	for i := 0; i < total; i++ {
		bounds, err := doc.Bound(i)
		if err != nil {
			return 0, fmt.Errorf("página %d: %w", i+1, err)
		}
		dpi := 144.0
		if bounds.Dx() > 0 {
			dpi = 72.0 * float64(widthPx) / float64(bounds.Dx())
		}
		img, err := doc.ImageDPI(i, dpi)
		if err != nil {
			return 0, fmt.Errorf("página %d: %w", i+1, err)
		}
		if err := w.WritePage(i, img); err != nil {
			return 0, fmt.Errorf("página %d: %w", i+1, err)
		}
	}
	return total, nil
}

// dirWriter é a implementação padrão de PageWriter: grava slide-NNN.png em
// outDir, o comportamento histórico de File.
type dirWriter struct {
	dir string
}

func newDirWriter(dir string) (*dirWriter, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}
	return &dirWriter{dir: dir}, nil
}

func (w *dirWriter) WritePage(index int, img image.Image) error {
	out, err := os.Create(filepath.Join(w.dir, fmt.Sprintf("slide-%03d.png", index+1)))
	if err != nil {
		return err
	}
	if err := png.Encode(out, img); err != nil {
		out.Close()
		return err
	}
	return out.Close()
}

func (w *dirWriter) Close() error { return nil }
