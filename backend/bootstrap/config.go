package bootstrap

import (
	"crypto/rand"
	"encoding/hex"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

// Valores por defecto: solo loopback y un timeout máximo acotado.
const (
	defaultHost       = "127.0.0.1"
	defaultPort       = 8080
	defaultMaxTimeout = 120 * time.Second
)

// Config agrupa la configuración del servidor leída del entorno (backend/.env).
type Config struct {
	Host           string
	Port           int
	Token          string
	TokenGenerated bool
	AllowedOrigins []string
	MaxTimeout     time.Duration
}

// Addr es la dirección host:puerto en la que escucha el servidor.
func (c Config) Addr() string {
	return net.JoinHostPort(c.Host, strconv.Itoa(c.Port))
}

// loadConfig lee MC_HOST, MC_PORT, MC_TOKEN, MC_ALLOWED_ORIGINS y MC_MAX_TIMEOUT.
// Sin MC_TOKEN genera uno efímero: el servidor nunca arranca sin autenticación.
func loadConfig(getenv func(string) string) (Config, error) {
	cfg := Config{
		Host:       strings.TrimSpace(getenv("MC_HOST")),
		Port:       defaultPort,
		Token:      strings.TrimSpace(getenv("MC_TOKEN")),
		MaxTimeout: defaultMaxTimeout,
	}

	if cfg.Host == "" {
		cfg.Host = defaultHost
	}

	if raw := strings.TrimSpace(getenv("MC_PORT")); raw != "" {
		port, err := strconv.Atoi(raw)
		if err != nil || port <= 0 || port > 65535 {
			return Config{}, &configError{key: "MC_PORT", value: raw}
		}
		cfg.Port = port
	}

	if raw := strings.TrimSpace(getenv("MC_MAX_TIMEOUT")); raw != "" {
		seconds, err := strconv.Atoi(raw)
		if err != nil || seconds <= 0 {
			return Config{}, &configError{key: "MC_MAX_TIMEOUT", value: raw}
		}
		cfg.MaxTimeout = time.Duration(seconds) * time.Second
	}

	for _, origin := range strings.Split(getenv("MC_ALLOWED_ORIGINS"), ",") {
		origin = strings.TrimRight(strings.TrimSpace(origin), "/")
		if origin != "" {
			cfg.AllowedOrigins = append(cfg.AllowedOrigins, origin)
		}
	}

	if cfg.Token == "" {
		token, err := randomToken()
		if err != nil {
			return Config{}, err
		}
		cfg.Token = token
		cfg.TokenGenerated = true
	}

	return cfg, nil
}

func randomToken() (string, error) {
	buf := make([]byte, 24)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}

type configError struct {
	key   string
	value string
}

func (e *configError) Error() string {
	return "valor inválido para " + e.key + ": " + e.value
}

// getenv permite inyectar el entorno en tests.
var getenv = os.Getenv
