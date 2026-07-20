// Package soffice localiza o executável do LibreOffice na máquina.
package soffice

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// candidates retorna os caminhos conhecidos do soffice para o sistema operacional.
func candidates(goos string) []string {
	switch goos {
	case "windows":
		return []string{
			`C:\Program Files\LibreOffice\program\soffice.exe`,
			`C:\Program Files (x86)\LibreOffice\program\soffice.exe`,
		}
	case "darwin":
		return []string{
			"/Applications/LibreOffice.app/Contents/MacOS/soffice",
		}
	default:
		return []string{
			"/usr/bin/soffice",
			"/usr/bin/libreoffice",
			"/snap/bin/libreoffice",
		}
	}
}

// Find retorna o caminho do soffice. Se configured não for vazio, valida e usa
// esse caminho; caso contrário procura no PATH e nos locais padrão de instalação.
func Find(configured string) (string, error) {
	return find(configured, runtime.GOOS, os.Stat, exec.LookPath)
}

func find(configured, goos string, stat func(string) (os.FileInfo, error), lookPath func(string) (string, error)) (string, error) {
	if configured != "" {
		if _, err := stat(configured); err != nil {
			return "", fmt.Errorf("caminho_soffice configurado não encontrado: %s", configured)
		}
		return configured, nil
	}
	if p, err := lookPath("soffice"); err == nil {
		return p, nil
	}
	for _, c := range candidates(goos) {
		if _, err := stat(c); err == nil {
			return c, nil
		}
	}
	return "", fmt.Errorf("LibreOffice não encontrado. Instale em https://pt-br.libreoffice.org/ ou defina caminho_soffice no config.json")
}

// ToPDF converte uma apresentação (path) em PDF via LibreOffice headless e
// retorna o caminho do PDF gerado dentro de outDir. Cada chamada roda com um
// perfil de usuário isolado e descartável: instâncias headless do
// LibreOffice não são seguras para rodar concorrentemente sob o mesmo
// perfil (disputa de lock).
func ToPDF(sofficePath, path, outDir string) (string, error) {
	if sofficePath == "" {
		return "", fmt.Errorf("LibreOffice é necessário para converter %s", filepath.Base(path))
	}
	profileDir, err := os.MkdirTemp("", "holyrics-converter-profile-*")
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(profileDir)

	userInstallation := "-env:UserInstallation=" + fileURL(profileDir)
	cmd := exec.Command(sofficePath, "--headless", "--norestore", userInstallation, "--convert-to", "pdf", "--outdir", outDir, path)
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

// fileURL converte um caminho de sistema de arquivos absoluto em uma URL
// file:// aceita pelo LibreOffice, inclusive em caminhos Windows (C:\...).
func fileURL(dir string) string {
	abs, err := filepath.Abs(dir)
	if err != nil {
		abs = dir
	}
	abs = filepath.ToSlash(abs)
	if !strings.HasPrefix(abs, "/") {
		abs = "/" + abs
	}
	return "file://" + abs
}
