package config

import (
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/pelletier/go-toml/v2"
)

// Providers holds persisted API keys for all supported providers.
type Providers struct {
	AnthropicAPIKey  string `toml:"anthropic_api_key"`
	OpenAIAPIKey     string `toml:"openai_api_key"`
	GeminiAPIKey     string `toml:"gemini_api_key"`
	GroqAPIKey       string `toml:"groq_api_key"`
	CerebrasAPIKey   string `toml:"cerebras_api_key"`
	OpenRouterAPIKey string `toml:"openrouter_api_key"`
	OllamaHost       string `toml:"ollama_host"`
}

// Config is the top-level config file structure.
type Config struct {
	Providers Providers `toml:"providers"`
}

var (
	mu      sync.Mutex
	current Config
)

// configPath returns the absolute path to ~/.luna/config.toml.
func configPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".luna", "config.toml"), nil
}

// Load reads ~/.luna/config.toml into the package-level singleton and also
// loads a .env file from the current working directory if one exists.
// Missing files are treated as empty configs, not errors.
func Load() error {
	mu.Lock()
	defer mu.Unlock()

	loadDotEnv()

	path, err := configPath()
	if err != nil {
		return err
	}

	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}

	var c Config
	if err := toml.Unmarshal(data, &c); err != nil {
		return err
	}
	current = c

	// Inject stored keys into the process environment so AutoDetectProvider
	// (which uses os.Getenv) picks them up. Shell env vars take precedence.
	injectEnv("ANTHROPIC_API_KEY", c.Providers.AnthropicAPIKey)
	injectEnv("OPENAI_API_KEY", c.Providers.OpenAIAPIKey)
	injectEnv("GEMINI_API_KEY", c.Providers.GeminiAPIKey)
	injectEnv("GROQ_API_KEY", c.Providers.GroqAPIKey)
	injectEnv("CEREBRAS_API_KEY", c.Providers.CerebrasAPIKey)
	injectEnv("OPENROUTER_API_KEY", c.Providers.OpenRouterAPIKey)
	injectEnv("OLLAMA_HOST", c.Providers.OllamaHost)
	return nil
}

// injectEnv sets an env var only if it is not already set in the shell.
func injectEnv(key, value string) {
	if value != "" && os.Getenv(key) == "" {
		os.Setenv(key, value)
	}
}

// Get returns the current config merged with environment variables.
// Environment variables always take precedence over stored values.
func Get() Config {
	mu.Lock()
	c := current
	mu.Unlock()

	if v := os.Getenv("ANTHROPIC_API_KEY"); v != "" {
		c.Providers.AnthropicAPIKey = v
	}
	if v := os.Getenv("OPENAI_API_KEY"); v != "" {
		c.Providers.OpenAIAPIKey = v
	}
	if v := os.Getenv("GEMINI_API_KEY"); v != "" {
		c.Providers.GeminiAPIKey = v
	}
	if v := os.Getenv("OLLAMA_HOST"); v != "" {
		c.Providers.OllamaHost = v
	}
	return c
}

// SetKey saves a provider API key to ~/.luna/config.toml and sets the
// corresponding environment variable so the running process picks it up
// immediately without a restart.
func SetKey(envKey, value string) error {
	mu.Lock()
	defer mu.Unlock()

	switch envKey {
	case "ANTHROPIC_API_KEY":
		current.Providers.AnthropicAPIKey = value
	case "OPENAI_API_KEY":
		current.Providers.OpenAIAPIKey = value
	case "GEMINI_API_KEY":
		current.Providers.GeminiAPIKey = value
	case "GROQ_API_KEY":
		current.Providers.GroqAPIKey = value
	case "CEREBRAS_API_KEY":
		current.Providers.CerebrasAPIKey = value
	case "OPENROUTER_API_KEY":
		current.Providers.OpenRouterAPIKey = value
	case "OLLAMA_HOST":
		current.Providers.OllamaHost = value
	}

	// Inject into process env so ProviderForModel reads it immediately.
	os.Setenv(envKey, value)

	return save()
}

// KeyForProvider returns the stored value for the given env key.
// Does not merge with environment variables — used for badge/unlock checks.
func KeyForProvider(envKey string) string {
	mu.Lock()
	defer mu.Unlock()

	switch envKey {
	case "ANTHROPIC_API_KEY":
		return current.Providers.AnthropicAPIKey
	case "OPENAI_API_KEY":
		return current.Providers.OpenAIAPIKey
	case "GEMINI_API_KEY":
		return current.Providers.GeminiAPIKey
	case "GROQ_API_KEY":
		return current.Providers.GroqAPIKey
	case "CEREBRAS_API_KEY":
		return current.Providers.CerebrasAPIKey
	case "OPENROUTER_API_KEY":
		return current.Providers.OpenRouterAPIKey
	case "OLLAMA_HOST":
		return current.Providers.OllamaHost
	}
	return ""
}

// loadDotEnv reads KEY=VALUE pairs from .env files and injects missing vars.
// Checks: (1) cwd/.env  (2) ~/.luna/.env
func loadDotEnv() {
	candidates := []string{".env"}
	if home, err := os.UserHomeDir(); err == nil {
		candidates = append(candidates, filepath.Join(home, ".luna", ".env"))
	}
	for _, path := range candidates {
		parseDotEnv(path)
	}
}

func parseDotEnv(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		return
	}
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		k, v, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		injectEnv(strings.TrimSpace(k), strings.Trim(strings.TrimSpace(v), `"'`))
	}
}

// save writes current to disk. Caller must hold mu.
func save() error {
	path, err := configPath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}
	data, err := toml.Marshal(current)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o600)
}
