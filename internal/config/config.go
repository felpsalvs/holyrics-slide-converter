// Package config carrega e cria o arquivo de configuração do conversor.
package config

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
)

// MinLarguraPx e MaxLarguraPx delimitam a faixa sã para LarguraPx, evitando
// PNGs gigantes/estouro de memória por erro de digitação no config.
const (
	MinLarguraPx = 100
	MaxLarguraPx = 8000
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
	// ConversoesSimultaneas limita quantas conversões rodam em paralelo.
	// LibreOffice headless não é seguro com muitas instâncias simultâneas.
	ConversoesSimultaneas int `json:"conversoes_simultaneas"`
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
		PastaEntrada:          filepath.Join(base, "entrada"),
		PastaSaida:            filepath.Join(base, "convertidos"),
		LarguraPx:             1920,
		ConversoesSimultaneas: 1,
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

// Validate confere se cfg é seguro para operar: pastas distintas e
// graváveis, e ajusta em silêncio (com aviso no log) valores fora de uma
// faixa sã. Deve ser chamado antes do watcher subir, para que uma pasta
// sem permissão seja detectada na inicialização e não no primeiro evento
// ao vivo durante um culto.
func Validate(cfg *Config) error {
	if cfg.PastaEntrada == cfg.PastaSaida {
		return fmt.Errorf("pasta_entrada e pasta_saida não podem ser a mesma pasta: %s", cfg.PastaEntrada)
	}
	for _, dir := range []string{cfg.PastaEntrada, cfg.PastaSaida} {
		if err := ensureWritable(dir); err != nil {
			return err
		}
	}
	if cfg.LarguraPx < MinLarguraPx || cfg.LarguraPx > MaxLarguraPx {
		clamped := clamp(cfg.LarguraPx, MinLarguraPx, MaxLarguraPx)
		slog.Warn("largura_px fora da faixa permitida, ajustando", "valor", cfg.LarguraPx, "ajustado_para", clamped)
		cfg.LarguraPx = clamped
	}
	if cfg.ConversoesSimultaneas <= 0 {
		cfg.ConversoesSimultaneas = 1
	}
	return nil
}

func ensureWritable(dir string) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("pasta %s: %w", dir, err)
	}
	probe := filepath.Join(dir, ".holyrics-write-test")
	if err := os.WriteFile(probe, []byte("x"), 0o644); err != nil {
		return fmt.Errorf("pasta %s não é gravável: %w", dir, err)
	}
	return os.Remove(probe)
}

func clamp(v, min, max int) int {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
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
