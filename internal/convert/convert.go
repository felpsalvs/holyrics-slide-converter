// Package convert transforma apresentações (PPTX/PPT/ODP) e PDFs em imagens PNG.
package convert

import (
	"fmt"
	"image/png"
	"os"
	"os/exec"
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

// File converte um arquivo suportado em PNGs dentro de outDir/<nome-do-arquivo>/.
// Apresentações passam antes pelo LibreOffice (sofficePath) para virar PDF.
// Retorna a pasta criada e a quantidade de páginas geradas.
func File(path, outDir string, widthPx int, sofficePath string) (string, int, error) {
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
		pdfPath, err = toPDF(path, tmp, sofficePath)
		if err != nil {
			return "", 0, err
		}
	}

	name := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	dest := filepath.Join(outDir, name)
	n, err := pdfToPNG(pdfPath, dest, widthPx)
	if err != nil {
		return "", 0, err
	}
	return dest, n, nil
}

// toPDF converte uma apresentação em PDF via LibreOffice headless e retorna
// o caminho do PDF gerado em outDir.
func toPDF(path, outDir, sofficePath string) (string, error) {
	if sofficePath == "" {
		return "", fmt.Errorf("LibreOffice é necessário para converter %s", filepath.Base(path))
	}
	cmd := exec.Command(sofficePath, "--headless", "--norestore", "--convert-to", "pdf", "--outdir", outDir, path)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("libreoffice falhou: %w: %s", err, strings.TrimSpace(string(out)))
	}
	name := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path)) + ".pdf"
	pdfPath := filepath.Join(outDir, name)
	if _, err := os.Stat(pdfPath); err != nil {
		return "", fmt.Errorf("libreoffice não gerou o PDF esperado (%s): %s", name, strings.TrimSpace(string(out)))
	}
	return pdfPath, nil
}

// pdfToPNG renderiza cada página do PDF como slide-NNN.png em dest,
// com largura alvo widthPx.
func pdfToPNG(pdfPath, dest string, widthPx int) (int, error) {
	doc, err := fitz.New(pdfPath)
	if err != nil {
		return 0, fmt.Errorf("falha ao abrir PDF: %w", err)
	}
	defer doc.Close()

	if err := os.MkdirAll(dest, 0o755); err != nil {
		return 0, err
	}

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
		out, err := os.Create(filepath.Join(dest, fmt.Sprintf("slide-%03d.png", i+1)))
		if err != nil {
			return 0, err
		}
		if err := png.Encode(out, img); err != nil {
			out.Close()
			return 0, fmt.Errorf("página %d: %w", i+1, err)
		}
		if err := out.Close(); err != nil {
			return 0, err
		}
	}
	return total, nil
}
