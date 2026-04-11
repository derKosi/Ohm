package scanner

import (
	"os"
	"strings"

	"github.com/derKosi/Ohm/internal/model"
)

// scanPATH scans PATH for AI-related entries.
func (s *Scanner) scanPATH() {
	path := os.Getenv("PATH")
	if path == "" {
		return
	}

	aiKeywords := []string{
		"ollama", "claude", "pi", "aider", "gemini", "codex",
		"cursor", "openai", "paperclip", "vibe", "mistral",
		"comfyui", "huggingface", "lm-studio",
	}

	var aiEntries []string
	pathSep := ":"
	if s.os == 2 { // Windows
		pathSep = ";"
	}
	for _, entry := range strings.Split(path, pathSep) {
		lower := strings.ToLower(entry)
		for _, kw := range aiKeywords {
			if strings.Contains(lower, kw) {
				aiEntries = append(aiEntries, entry)
				break
			}
		}
	}

	if len(aiEntries) > 0 {
		s.addFinding(model.Finding{
			ID:        "path-ai-entries",
			Category:  model.CatStragglers,
			Name:      "AI Entries in PATH",
			Path:      "(PATH environment variable)",
			SizeBytes: 0,
			SubItems:  aiEntries,
			RiskLevel: model.RiskCaution,
			UninstallCmds: map[string]string{
				"linux":   "# Edit ~/.bashrc or ~/.zshrc to remove these PATH entries",
				"macos":   "# Edit ~/.zshrc to remove these PATH entries",
				"windows": "# Edit System Environment Variables to remove these PATH entries",
			},
		})
	}
}

// scanENV scans environment variables for AI-related entries.
func (s *Scanner) scanENV() {
	aiEnvVars := []string{
		"OPENAI_API_KEY",
		"ANTHROPIC_API_KEY",
		"GEMINI_API_KEY",
		"MISTRAL_API_KEY",
		"OLLAMA_HOST",
		"OLLAMA_MODELS",
		"HF_TOKEN",
		"HUGGINGFACE_TOKEN",
		"COHERE_API_KEY",
		"DEEPSEEK_API_KEY",
		"GROQ_API_KEY",
		"TOGETHER_API_KEY",
		"REPLICATE_API_TOKEN",
		"AZURE_OPENAI_API_KEY",
		"AZURE_OPENAI_ENDPOINT",
		"VOYAGE_API_KEY",
		"PINECONE_API_KEY",
		"CHROMA_API_KEY",
	}

	var found []string
	for _, v := range aiEnvVars {
		if os.Getenv(v) != "" {
			found = append(found, v+"=<REDACTED>")
		}
	}

	if len(found) > 0 {
		s.addFinding(model.Finding{
			ID:        "env-ai-vars",
			Category:  model.CatStragglers,
			Name:      "AI Environment Variables",
			Path:      "(environment)",
			SizeBytes: 0,
			SubItems:  found,
			RiskLevel: model.RiskDanger,
			UninstallCmds: map[string]string{
				"linux":   "# Edit ~/.bashrc or ~/.zshrc to remove these variables",
				"macos":   "# Edit ~/.zshrc to remove these variables",
				"windows": "# Edit System Environment Variables to remove these",
			},
		})
	}
}

// scanShellProfiles scans shell profiles for AI-related modifications.
func (s *Scanner) scanShellProfiles() {
	profiles := []string{}
	home := s.home

	switch s.os {
	case 0: // Linux
		profiles = []string{
			home + "/.bashrc",
			home + "/.zshrc",
			home + "/.profile",
			home + "/.bash_profile",
		}
	case 1: // macOS
		profiles = []string{
			home + "/.zshrc",
			home + "/.bashrc",
			home + "/.bash_profile",
			home + "/.profile",
		}
	case 2: // Windows
		profiles = []string{
			home + "/Documents/PowerShell/Microsoft.PowerShell_profile.ps1",
			home + "/Documents/WindowsPowerShell/Microsoft.PowerShell_profile.ps1",
		}
	}

	aiKeywords := []string{
		"ollama", "claude", "aider", "gemini", "codex",
		"cursor", "openai", "paperclip", "vibe", "mistral",
		"pi-coding", "comfyui", "huggingface",
	}

	for _, profile := range profiles {
		data, err := os.ReadFile(profile)
		if err != nil {
			continue
		}

		var aiLines []string
		for _, line := range strings.Split(string(data), "\n") {
			lower := strings.ToLower(line)
			for _, kw := range aiKeywords {
				if strings.Contains(lower, kw) && !strings.HasPrefix(strings.TrimSpace(line), "#") {
					aiLines = append(aiLines, strings.TrimSpace(line))
					break
				}
			}
		}

		if len(aiLines) > 0 {
			s.addFinding(model.Finding{
				ID:        "shell-profile-" + profile,
				Category:  model.CatStragglers,
				Name:      "AI Entries in " + profile,
				Path:      profile,
				SizeBytes: 0,
				SubItems:  aiLines,
				RiskLevel: model.RiskCaution,
				UninstallCmds: map[string]string{
					"linux":   "# Edit " + profile + " to remove AI-related lines",
					"macos":   "# Edit " + profile + " to remove AI-related lines",
					"windows": "# Edit " + profile + " to remove AI-related lines",
				},
			})
		}
	}
}
