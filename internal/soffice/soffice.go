// Package soffice localiza o executável do LibreOffice na máquina.
package soffice

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
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
