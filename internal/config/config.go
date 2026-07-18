// Package config carrega e cria o arquivo de configuração do conversor.
package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Config define as opções do conversor, salvas em config.json.
type Config struct {
	// PastaEntrada é a pasta monitorada: arquivos .pptx/.ppt/.odp/.pdf
	// colocados nela são convertidos automaticamente.
	PastaEntrada string `json:"pasta_entrada"`
	// PastaSaida é onde as imagens são geradas, em uma subpasta por arquivo.
	// Aponte para dentro da biblioteca de imagens do Holyrics.
	PastaSaida string `json:"pasta_saida"`
	// LarguraPx é a largura alvo das imagens geradas (px).
	LarguraPx int `json:"largura_px"`
	// CaminhoSoffice permite fixar o caminho do LibreOffice (soffice).
	// Vazio = detectar automaticamente.
	CaminhoSoffice string `json:"caminho_soffice"`
}

// DefaultPath retorna o caminho padrão do config.json, ao lado da pasta
// de trabalho do usuário (~/HolyricsConverter/config.json).
func DefaultPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, "HolyricsConverter", "config.json"), nil
}

func defaults() (Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return Config{}, err
	}
	base := filepath.Join(home, "HolyricsConverter")
	return Config{
		PastaEntrada: filepath.Join(base, "entrada"),
		PastaSaida:   filepath.Join(base, "convertidos"),
		LarguraPx:    1920,
	}, nil
}

// Load lê o config do caminho informado. Se o arquivo não existir, cria um
// com valores padrão e retorna created=true para o chamador orientar o usuário.
func Load(path string) (cfg Config, created bool, err error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		cfg, err = defaults()
		if err != nil {
			return Config{}, false, err
		}
		if err := save(path, cfg); err != nil {
			return Config{}, false, err
		}
		return cfg, true, nil
	}
	if err != nil {
		return Config{}, false, err
	}
	cfg, err = defaults()
	if err != nil {
		return Config{}, false, err
	}
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Config{}, false, fmt.Errorf("config inválido em %s: %w", path, err)
	}
	if cfg.PastaEntrada == "" || cfg.PastaSaida == "" {
		return Config{}, false, fmt.Errorf("config %s: pasta_entrada e pasta_saida são obrigatórias", path)
	}
	if cfg.LarguraPx <= 0 {
		cfg.LarguraPx = 1920
	}
	return cfg, false, nil
}

func save(path string, cfg Config) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(data, '\n'), 0o644)
}
