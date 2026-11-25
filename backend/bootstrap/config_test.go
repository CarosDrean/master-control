package bootstrap

import (
	"testing"
	"time"
)

func env(m map[string]string) func(string) string {
	return func(k string) string { return m[k] }
}

func TestLoadConfigDefaults(t *testing.T) {
	cfg, err := loadConfig(env(nil))
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Addr() != "127.0.0.1:8080" {
		t.Errorf("addr = %q, quiero loopback por defecto", cfg.Addr())
	}
	if !cfg.TokenGenerated || len(cfg.Token) < 32 {
		t.Errorf("sin MC_TOKEN debe generarse un token efímero, got %q", cfg.Token)
	}
	if cfg.MaxTimeout != defaultMaxTimeout {
		t.Errorf("max timeout = %v", cfg.MaxTimeout)
	}
}

func TestLoadConfigValues(t *testing.T) {
	cfg, err := loadConfig(env(map[string]string{
		"MC_HOST":            "100.64.0.1",
		"MC_PORT":            "9090",
		"MC_TOKEN":           " secreto ",
		"MC_ALLOWED_ORIGINS": "http://a:5173/, ,http://b:5173",
		"MC_MAX_TIMEOUT":     "30",
	}))
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Addr() != "100.64.0.1:9090" || cfg.Token != "secreto" || cfg.TokenGenerated {
		t.Errorf("config inesperada: %+v", cfg)
	}
	if len(cfg.AllowedOrigins) != 2 || cfg.AllowedOrigins[0] != "http://a:5173" {
		t.Errorf("origins = %v", cfg.AllowedOrigins)
	}
	if cfg.MaxTimeout != 30*time.Second {
		t.Errorf("max timeout = %v", cfg.MaxTimeout)
	}
}

func TestLoadConfigInvalid(t *testing.T) {
	for _, kv := range []map[string]string{
		{"MC_PORT": "abc"},
		{"MC_PORT": "70000"},
		{"MC_MAX_TIMEOUT": "0"},
	} {
		if _, err := loadConfig(env(kv)); err == nil {
			t.Errorf("esperaba error para %v", kv)
		}
	}
}

func TestHealthPath(t *testing.T) {
	if got := healthPath("master-control"); got != "/master-control/health" {
		t.Errorf("got %q", got)
	}
	if got := healthPath(""); got != "/health" {
		t.Errorf("got %q", got)
	}
}
