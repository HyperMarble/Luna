package config

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/pelletier/go-toml/v2"
)

// Providers holds persisted API keys for all paid providers.
type Providers struct {
	AnthropicAPIKey string `toml:"anthropic_api_key"`
	OpenAIAPIKey    string `toml:"openai_api_key"`
	GeminiAPIKey    string `toml:"gemini_api_key"`
	OllamaHost      string `toml:"ollama_host"`
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

// Load reads ~/.luna/config.toml into the package-level singleton.
// A missing file is treated as an empty config, not an error.
func Load() error {
	mu.Lock()
	defer mu.Unlock()

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
	return nil
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
	case "OLLAMA_HOST":
		return current.Providers.OllamaHost
	}
	return ""
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
