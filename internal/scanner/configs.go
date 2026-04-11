package scanner

import (
	"os"
	"path/filepath"

	"github.com/derKosi/Ohm/internal/model"
)

// scanConfigDirs detects known AI tool config directories.
func (s *Scanner) scanConfigDirs() {
	type configDef struct {
		id   string
		name string
		path string
		risk model.Risk
	}

	configs := []configDef{
		// ── Original configs ────────────────────────────────────────
		{"claude-config", "Claude Config", filepath.Join(s.home, ".claude"), model.RiskDanger},
		{"claude-json", "Claude Config (JSON)", filepath.Join(s.home, ".claude.json"), model.RiskDanger},
		{"pi-config", "pi Config", filepath.Join(s.home, ".pi"), model.RiskCaution},
		{"aider-config", "Aider Config", filepath.Join(s.home, ".aider.conf.yml"), model.RiskSafe},
		{"codex-config", "Codex CLI Config", filepath.Join(s.home, ".codex"), model.RiskCaution},
		{"cursor-config", "Cursor Config", filepath.Join(s.home, ".cursor"), model.RiskCaution},
		{"gemini-config", "Gemini CLI Config", filepath.Join(s.home, ".gemini"), model.RiskCaution},
		{"opencode-config", "OpenCode Config", filepath.Join(s.home, ".opencode"), model.RiskSafe},
		{"ollama-config", "Ollama User Data", filepath.Join(s.home, ".ollama"), model.RiskSafe},
		{"continue-config", "Continue Config", filepath.Join(s.home, ".continue"), model.RiskDanger},
		{"amazon-q-config", "Amazon Q Config", filepath.Join(s.home, ".amazon-q"), model.RiskCaution},
		{"vibe-config", "Mistral Vibe Config", filepath.Join(s.home, ".vibe"), model.RiskSafe},
		{"paperclip-config", "PaperclipAI Config", filepath.Join(s.home, ".paperclip"), model.RiskDanger},
		{"gpt4all-config", "GPT4All Config", filepath.Join(s.home, ".gpt4all"), model.RiskSafe},
		{"huggingface-token", "HuggingFace Token", filepath.Join(s.home, ".huggingface"), model.RiskDanger},

		// ── NEW: Agent configs ──────────────────────────────────────
		{"claw-config", "Claw Code Config", filepath.Join(s.home, ".claw"), model.RiskCaution},
		{"openhands-config", "OpenHands Config", filepath.Join(s.home, ".openhands"), model.RiskSafe},
		{"interpreter-config", "Open Interpreter Config", filepath.Join(s.home, ".interpreter"), model.RiskCaution},
		{"goose-config", "Goose Config", filepath.Join(s.home, ".config", "goose"), model.RiskSafe},
		{"roo-code-config", "Roo Code Config", filepath.Join(s.home, ".roo-code"), model.RiskSafe},
		{"crush-config", "Crush Config", filepath.Join(s.home, ".config", "crush"), model.RiskSafe},
		{"qwen-config", "Qwen Code Config", filepath.Join(s.home, ".qwen"), model.RiskSafe},
		{"kilo-code-config", "Kilo Code Config", filepath.Join(s.home, ".kilo-code"), model.RiskSafe},
		{"plandex-config", "Plandex Config", filepath.Join(s.home, ".plandex"), model.RiskSafe},
		{"hermes-config", "Hermes Agent Config", filepath.Join(s.home, ".hermes"), model.RiskSafe},
		{"groq-config", "Groq Code CLI Config", filepath.Join(s.home, ".groq"), model.RiskSafe},
		{"devon-config", "Devon Config", filepath.Join(s.home, ".devon"), model.RiskSafe},
		{"claurst-config", "Claurst Config", filepath.Join(s.home, ".claurst"), model.RiskSafe},
		{"letta-config", "Letta Code Config", filepath.Join(s.home, ".letta"), model.RiskSafe},
		{"forgecode-config", "ForgeCode Config", filepath.Join(s.home, ".forgecode"), model.RiskSafe},

		// ── NEW: OpenClaw ecosystem configs ─────────────────────────
		{"openclaw-config", "OpenClaw Config", filepath.Join(s.home, ".openclaw"), model.RiskDanger},
		{"nanobot-config", "NanoBot Config", filepath.Join(s.home, ".nanobot"), model.RiskCaution},
		{"zeroclaw-config", "ZeroClaw Config", filepath.Join(s.home, ".zeroclaw"), model.RiskCaution},
		{"picoclaw-config", "PicoClaw Config", filepath.Join(s.home, ".picoclaw"), model.RiskSafe},
		{"nanoclaw-config", "NanoClaw Config", filepath.Join(s.home, ".nanoclaw"), model.RiskCaution},
		{"ironclaw-config", "IronClaw Config", filepath.Join(s.home, ".ironclaw"), model.RiskSafe},
		{"nullclaw-config", "NullClaw Config", filepath.Join(s.home, ".nullclaw"), model.RiskSafe},
		{"moltis-config", "Moltis Config", filepath.Join(s.home, ".moltis"), model.RiskSafe},
		{"clawith-config", "Clawith Config", filepath.Join(s.home, ".clawith"), model.RiskSafe},

		// ── NEW: Orchestration configs ──────────────────────────────
		{"claude-squad-config", "Claude Squad Config", filepath.Join(s.home, ".claude-squad"), model.RiskSafe},
		{"claude-flow-config", "Claude Flow Config", filepath.Join(s.home, ".claude-flow"), model.RiskSafe},
	}

	// Track which config dirs were already reported as part of other findings
	reported := make(map[string]bool)

	for _, cfg := range configs {
		info, err := os.Stat(cfg.path)
		if err != nil {
			continue
		}

		// Skip if already reported
		if reported[cfg.id] {
			continue
		}
		reported[cfg.id] = true

		var size int64
		if info.IsDir() {
			size = s.dirSize(cfg.path)
		} else {
			size = info.Size()
		}

		s.addFinding(model.Finding{
			ID:        cfg.id,
			Category:  model.CatConfigs,
			Name:      cfg.name,
			Path:      cfg.path,
			SizeBytes: size,
			RiskLevel: cfg.risk,
			UninstallCmds: map[string]string{
				"linux":   "rm -rf " + cfg.path,
				"macos":   "rm -rf " + cfg.path,
				"windows": "Remove-Item '" + cfg.path + "' -Recurse -Force",
			},
		})
	}
}
