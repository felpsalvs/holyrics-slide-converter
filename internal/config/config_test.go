package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestLoadCriaConfigPadrao(t *testing.T) {
	path := filepath.Join(t.TempDir(), "sub", "config.json")

	cfg, created, err := Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if !created {
		t.Fatal("created = false, esperado true na primeira execução")
	}
	if cfg.LarguraPx != 1920 {
		t.Errorf("LarguraPx = %d, esperado 1920", cfg.LarguraPx)
	}
	if cfg.ConversoesSimultaneas != 1 {
		t.Errorf("ConversoesSimultaneas = %d, esperado 1", cfg.ConversoesSimultaneas)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("config.json não foi criado em disco: %v", err)
	}

	cfg2, created2, err := Load(path)
	if err != nil {
		t.Fatalf("segundo Load: %v", err)
	}
	if created2 {
		t.Fatal("created = true na segunda leitura, esperado false")
	}
	if cfg2 != cfg {
		t.Fatalf("config lido difere do criado: %+v != %+v", cfg2, cfg)
	}
}

func TestLoadCamposObrigatorios(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")
	data, _ := json.Marshal(map[string]any{"pasta_entrada": "", "pasta_saida": "/tmp/saida"})
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatal(err)
	}
	if _, _, err := Load(path); err == nil {
		t.Fatal("esperado erro por pasta_entrada explicitamente vazia")
	}
}

func TestLoadLarguraPxZeroUsaDefault(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")
	data, _ := json.Marshal(map[string]any{
		"pasta_entrada": "/tmp/entrada",
		"pasta_saida":   "/tmp/saida",
		"largura_px":    0,
	})
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatal(err)
	}
	cfg, _, err := Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.LarguraPx != 1920 {
		t.Errorf("LarguraPx = %d, esperado 1920", cfg.LarguraPx)
	}
}

func TestValidatePastasIguais(t *testing.T) {
	dir := t.TempDir()
	cfg := Config{PastaEntrada: dir, PastaSaida: dir, LarguraPx: 1920, ConversoesSimultaneas: 1}
	if err := Validate(&cfg); err == nil {
		t.Fatal("esperado erro quando pasta_entrada == pasta_saida")
	}
}

func TestValidateCriaEValidaPastas(t *testing.T) {
	base := t.TempDir()
	cfg := Config{
		PastaEntrada: filepath.Join(base, "entrada"),
		PastaSaida:   filepath.Join(base, "saida"),
		LarguraPx:    1920,
	}
	if err := Validate(&cfg); err != nil {
		t.Fatalf("Validate: %v", err)
	}
	for _, dir := range []string{cfg.PastaEntrada, cfg.PastaSaida} {
		if _, err := os.Stat(dir); err != nil {
			t.Errorf("pasta %s não foi criada: %v", dir, err)
		}
	}
	if cfg.ConversoesSimultaneas != 1 {
		t.Errorf("ConversoesSimultaneas = %d, esperado default 1", cfg.ConversoesSimultaneas)
	}
}

func TestValidateAjustaLarguraPxForaDaFaixa(t *testing.T) {
	base := t.TempDir()
	cfg := Config{
		PastaEntrada: filepath.Join(base, "entrada"),
		PastaSaida:   filepath.Join(base, "saida"),
		LarguraPx:    50, // abaixo de MinLarguraPx
	}
	if err := Validate(&cfg); err != nil {
		t.Fatalf("Validate: %v", err)
	}
	if cfg.LarguraPx != MinLarguraPx {
		t.Errorf("LarguraPx = %d, esperado ajuste para %d", cfg.LarguraPx, MinLarguraPx)
	}

	cfg.LarguraPx = 999999 // acima de MaxLarguraPx
	if err := Validate(&cfg); err != nil {
		t.Fatalf("Validate: %v", err)
	}
	if cfg.LarguraPx != MaxLarguraPx {
		t.Errorf("LarguraPx = %d, esperado ajuste para %d", cfg.LarguraPx, MaxLarguraPx)
	}
}

func TestValidateConversoesSimultaneasInvalidoVaiParaDefault(t *testing.T) {
	base := t.TempDir()
	cfg := Config{
		PastaEntrada:          filepath.Join(base, "entrada"),
		PastaSaida:            filepath.Join(base, "saida"),
		LarguraPx:             1920,
		ConversoesSimultaneas: -3,
	}
	if err := Validate(&cfg); err != nil {
		t.Fatalf("Validate: %v", err)
	}
	if cfg.ConversoesSimultaneas != 1 {
		t.Errorf("ConversoesSimultaneas = %d, esperado default 1", cfg.ConversoesSimultaneas)
	}
}

func TestValidatePastaNaoGravavel(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("rodando como root: restrição de permissão não se aplica")
	}
	base := t.TempDir()
	entrada := filepath.Join(base, "entrada")
	saida := filepath.Join(base, "saida-sem-permissao")
	if err := os.MkdirAll(saida, 0o555); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Chmod(saida, 0o755) })

	cfg := Config{PastaEntrada: entrada, PastaSaida: saida, LarguraPx: 1920}
	if err := Validate(&cfg); err == nil {
		t.Fatal("esperado erro para pasta sem permissão de escrita")
	}
}
