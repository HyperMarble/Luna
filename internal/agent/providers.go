package agent

import "os"

// ProviderName identifies a supported LLM provider.
type ProviderName string

const (
	ProviderAnthropic  ProviderName = "anthropic"
	ProviderOpenAI     ProviderName = "openai"
	ProviderGemini     ProviderName = "gemini"
	ProviderGroq       ProviderName = "groq"
	ProviderCerebras   ProviderName = "cerebras"
	ProviderOpenRouter ProviderName = "openrouter"
	ProviderOllama     ProviderName = "ollama"
)

// ModelEntry is a single selectable model inside a provider.
type ModelEntry struct {
	Provider    ProviderName
	DisplayName string // Provider display label.
	ModelID     string
	ModelLabel  string // Model display label.
	Free        bool   // True if no payment required.
	KeyURL      string // Where to obtain an API key.
	EnvKey      string // Environment variable name for the API key.
}

// ProviderInfo is a provider node used by the tree picker.
type ProviderInfo struct {
	Name        ProviderName
	DisplayName string
	Free        bool   // True if no payment required.
	KeyURL      string // Where to obtain an API key.
	EnvKey      string // Environment variable name for the API key.
	Models      []ModelEntry
}

// Providers returns all providers with their full model lists.
func Providers() []ProviderInfo {
	out := make([]ProviderInfo, len(providerDefs))
	for i, def := range providerDefs {
		models := make([]ModelEntry, len(def.models))
		for j, m := range def.models {
			models[j] = ModelEntry{
				Provider:    def.name,
				DisplayName: def.displayName,
				ModelID:     m.id,
				ModelLabel:  m.label,
				Free:        def.free,
				KeyURL:      def.keyURL,
				EnvKey:      def.envKey,
			}
		}
		out[i] = ProviderInfo{
			Name:        def.name,
			DisplayName: def.displayName,
			Free:        def.free,
			KeyURL:      def.keyURL,
			EnvKey:      def.envKey,
			Models:      models,
		}
	}
	return out
}

// ModelTree returns the full flat list of models across all providers.
func ModelTree() []ModelEntry {
	var out []ModelEntry
	for _, p := range Providers() {
		out = append(out, p.Models...)
	}
	return out
}

// providerDef describes how to build a Provider from environment variables.
type providerDef struct {
	name         ProviderName
	displayName  string
	envKey       string // API key env var.
	baseURL      string // OpenAI-compatible endpoint (empty = Anthropic native).
	defaultModel string
	free         bool   // True if provider has a free tier requiring no payment.
	keyURL       string // URL where the user can obtain an API key.
	models       []modelDef
}

type modelDef struct {
	id    string
	label string
}

var providerDefs = []providerDef{
	{
		name:         ProviderCerebras,
		displayName:  "Cerebras",
		envKey:       "CEREBRAS_API_KEY",
		baseURL:      "https://api.cerebras.ai/v1",
		defaultModel: "gpt-oss-120b",
		free:         true,
		keyURL:       "cloud.cerebras.ai",
		models: []modelDef{
			{"llama-3.3-70b", "Llama 3.3 70B"},
			{"llama-4-scout-17b-16e-instruct", "Llama 4 Scout 17B"},
			{"llama3.1-8b", "Llama 3.1 8B"},
			{"gpt-oss-120b", "OpenAI GPT OSS 120B"},
			{"qwen-3-235b-a22b-instruct-2507", "Qwen 3 235B Instruct"},
			{"zai-glm-4.7", "Z.ai GLM 4.7"},
		},
	},
	{
		name:         ProviderGroq,
		displayName:  "Groq",
		envKey:       "GROQ_API_KEY",
		baseURL:      "https://api.groq.com/openai/v1",
		defaultModel: "llama-3.3-70b-versatile",
		free:         true,
		keyURL:       "console.groq.com",
		models: []modelDef{
			{"gpt-oss-120b", "GPT OSS 120B"},
			{"gpt-oss-20b", "GPT OSS 20B"},
			{"qwen3-32b", "Qwen 3 32B"},
			{"meta-llama/llama-4-scout-17b-16e-instruct", "Llama 4 Scout"},
			{"moonshotai/kimi-k2-instruct", "Kimi K2"},
			{"llama-3.3-70b-versatile", "Llama 3.3 70B"},
			{"llama-3.1-8b-instant", "Llama 3.1 8B"},
		},
	},
	{
		name:         ProviderOpenRouter,
		displayName:  "OpenRouter",
		envKey:       "OPENROUTER_API_KEY",
		baseURL:      "https://openrouter.ai/api/v1",
		defaultModel: "meta-llama/llama-3.3-70b-instruct:free",
		free:         true,
		keyURL:       "openrouter.ai/keys",
		models: []modelDef{
			{"meta-llama/llama-3.3-70b-instruct:free", "Llama 3.3 70B (free)"},
			{"google/gemini-2.0-flash-001:free", "Gemini 2.0 Flash (free)"},
			{"anthropic/claude-sonnet-4-6", "Claude Sonnet 4.6"},
			{"__custom__", "+ Custom model"},
		},
	},
	{
		name:         ProviderAnthropic,
		displayName:  "Anthropic",
		envKey:       "ANTHROPIC_API_KEY",
		baseURL:      "",
		defaultModel: "claude-sonnet-4-6",
		free:         false,
		keyURL:       "console.anthropic.com",
		models: []modelDef{
			{"claude-opus-4-6", "Claude Opus 4.6"},
			{"claude-sonnet-4-6", "Claude Sonnet 4.6"},
			{"claude-haiku-4-5", "Claude Haiku 4.5"},
		},
	},
	{
		name:         ProviderOpenAI,
		displayName:  "OpenAI",
		envKey:       "OPENAI_API_KEY",
		baseURL:      "https://api.openai.com/v1",
		defaultModel: "gpt-4o",
		free:         false,
		keyURL:       "platform.openai.com",
		models: []modelDef{
			{"gpt-4o", "GPT-4o"},
			{"gpt-4o-mini", "GPT-4o Mini"},
			{"o3", "o3"},
		},
	},
	{
		name:         ProviderGemini,
		displayName:  "Gemini",
		envKey:       "GEMINI_API_KEY",
		baseURL:      "https://generativelanguage.googleapis.com/v1beta/openai",
		defaultModel: "gemini-2.0-flash",
		free:         false,
		keyURL:       "aistudio.google.com",
		models: []modelDef{
			{"gemini-2.5-pro", "Gemini 2.5 Pro"},
			{"gemini-2.0-flash", "Gemini 2.0 Flash"},
		},
	},
	{
		name:         ProviderOllama,
		displayName:  "Ollama",
		envKey:       "OLLAMA_HOST",
		baseURL:      "http://localhost:11434/v1",
		defaultModel: "llama3.2",
		free:         false,
		keyURL:       "ollama.com",
		models: []modelDef{
			{"llama3.2", "Llama 3.2"},
			{"mistral", "Mistral"},
			{"codellama", "Code Llama"},
		},
	},
}

// AutoDetectProvider returns the first provider whose API key is set.
// Priority: Cerebras → Groq → OpenRouter → Anthropic → OpenAI → Gemini → Ollama → Stub.
func AutoDetectProvider() Provider {
	for _, def := range providerDefs {
		key := os.Getenv(def.envKey)

		if def.name == ProviderOllama {
			if os.Getenv("LUNA_PROVIDER") == "ollama" {
				host := key
				if host == "" {
					host = "http://localhost:11434"
				}
				model := modelOrDefault("", def.defaultModel)
				return newOpenAIProvider(host+"/v1", "", model, systemPrompt)
			}
			continue
		}

		if key == "" {
			continue
		}

		model := modelOrDefault(os.Getenv("LUNA_MODEL"), def.defaultModel)

		if def.name == ProviderAnthropic {
			return NewClaudeProvider(model)
		}

		baseURL := def.baseURL
		if def.name == ProviderGemini {
			baseURL = baseURL + "?key=" + key
		}

		return newOpenAIProvider(baseURL, key, model, systemPrompt)
	}

	return NewStubProvider()
}

// ProviderForModel builds a Provider for a specific provider+model combination.
func ProviderForModel(name ProviderName, modelID string) Provider {
	for _, def := range providerDefs {
		if def.name != name {
			continue
		}
		modelID = modelOrDefault(modelID, def.defaultModel)
		if def.name == ProviderAnthropic {
			return NewClaudeProvider(modelID)
		}
		key := os.Getenv(def.envKey)
		baseURL := def.baseURL
		if def.name == ProviderGemini {
			baseURL = baseURL + "?key=" + key
		}
		if def.name == ProviderOllama {
			if host := os.Getenv("OLLAMA_HOST"); host != "" {
				baseURL = host + "/v1"
			}
		}
		return newOpenAIProvider(baseURL, key, modelID, systemPrompt)
	}
	return NewStubProvider()
}

func modelOrDefault(override, fallback string) string {
	if override != "" {
		return override
	}
	return fallback
}
