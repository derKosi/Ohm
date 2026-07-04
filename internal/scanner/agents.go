// SPDX-FileCopyrightText: 2026 Mathias Kosinski
// SPDX-License-Identifier: AGPL-3.0-or-later

package scanner

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/derKosi/Ohm/internal/model"
)

// agentDef defines a known AI agent/harness for detection.
type agentDef struct {
	id            string
	name          string
	configDirs    []string
	configFiles   []string
	commands      []string
	npmPackages   []string
	pipPackages   []string
	goPackages    []string
	brewPackages  []string
	risk          model.Risk
	uninstallCmds map[string]string
}

// scanAgents detects AI agents and harnesses.
func (s *Scanner) scanAgents() {
	agents := []agentDef{
		// ── Major Agents ──────────────────────────────────────────────
		{
			id:          "pi",
			name:        "pi (Coding Agent)",
			configDirs:  []string{"~/.pi"},
			commands:    []string{"pi"},
			npmPackages: []string{"@mariozechner/pi-coding-agent", "pi-session-search", "@thesethrose/pi-zai-provider", "@alexanderfortin/pi-zai-usage", "@samfp/pi-memory"},
			risk:        model.RiskCaution,
			uninstallCmds: map[string]string{
				"linux":   "npm uninstall -g @mariozechner/pi-coding-agent pi-session-search @thesethrose/pi-zai-provider @alexanderfortin/pi-zai-usage @samfp/pi-memory",
				"macos":   "npm uninstall -g @mariozechner/pi-coding-agent pi-session-search @thesethrose/pi-zai-provider @alexanderfortin/pi-zai-usage @samfp/pi-memory",
				"windows": "npm uninstall -g @mariozechner/pi-coding-agent pi-session-search @thesethrose/pi-zai-provider @alexanderfortin/pi-zai-usage @samfp/pi-memory",
			},
		},
		{
			id:          "claude-code",
			name:        "Claude Code",
			configDirs:  []string{"~/.claude"},
			configFiles: []string{"~/.claude.json"},
			commands:    []string{"claude"},
			npmPackages: []string{"@anthropic-ai/claude-code"},
			risk:        model.RiskDanger,
			uninstallCmds: map[string]string{
				"linux":   "npm uninstall -g @anthropic-ai/claude-code",
				"macos":   "npm uninstall -g @anthropic-ai/claude-code",
				"windows": "winget uninstall Anthropic.Claude",
			},
		},
		{
			id:          "codex-cli",
			name:        "Codex CLI (OpenAI)",
			configDirs:  []string{"~/.codex"},
			commands:    []string{"codex"},
			npmPackages: []string{"@openai/codex"},
			risk:        model.RiskDanger,
			uninstallCmds: map[string]string{
				"linux":   "npm uninstall -g @openai/codex",
				"macos":   "npm uninstall -g @openai/codex",
				"windows": "npm uninstall -g @openai/codex",
			},
		},
		{
			id:          "gemini-cli",
			name:        "Gemini CLI (Google)",
			configDirs:  []string{"~/.gemini"},
			commands:    []string{"gemini"},
			npmPackages: []string{"@google/gemini-cli"},
			risk:        model.RiskDanger,
			uninstallCmds: map[string]string{
				"linux":   "npm uninstall -g @google/gemini-cli",
				"macos":   "npm uninstall -g @google/gemini-cli",
				"windows": "npm uninstall -g @google/gemini-cli",
			},
		},
		{
			id:         "aider",
			name:       "Aider",
			configFiles: []string{"~/.aider.conf.yml"},
			commands:   []string{"aider"},
			pipPackages: []string{"aider-chat"},
			risk:       model.RiskSafe,
			uninstallCmds: map[string]string{
				"linux":   "pip uninstall aider-chat",
				"macos":   "pip uninstall aider-chat",
				"windows": "pip uninstall aider-chat",
			},
		},
		{
			id:         "mistral-vibe",
			name:       "Mistral Vibe",
			configDirs: []string{"~/.vibe"},
			commands:   []string{"vibe"},
			risk:       model.RiskSafe,
			uninstallCmds: map[string]string{
				"linux":   "# Vibe: remove binary manually + rm -rf ~/.vibe",
				"macos":   "# Vibe: remove binary manually + rm -rf ~/.vibe",
				"windows": "# Vibe: remove binary manually + Remove-Item ~/.vibe -Recurse -Force",
			},
		},
		{
			id:          "paperclipai",
			name:        "PaperclipAI",
			configDirs:  []string{"~/.paperclip"},
			commands:    []string{"paperclip"},
			npmPackages: []string{"paperclipai"},
			risk:        model.RiskCaution,
			uninstallCmds: map[string]string{
				"linux":   "npm uninstall -g paperclipai",
				"macos":   "npm uninstall -g paperclipai",
				"windows": "npm uninstall -g paperclipai",
			},
		},
		{
			id:          "opencode",
			name:        "OpenCode",
			configDirs:  []string{"~/.opencode"},
			commands:    []string{"opencode"},
			npmPackages: []string{"opencode-ai"},
			risk:        model.RiskSafe,
			uninstallCmds: map[string]string{
				"linux":   "npm uninstall -g opencode-ai && rm -f ~/go/bin/opencode",
				"macos":   "npm uninstall -g opencode-ai && rm -f ~/go/bin/opencode",
				"windows": "npm uninstall -g opencode-ai; Remove-Item ~/go/bin/opencode* -Force",
			},
		},
		{
			id:         "continue",
			name:       "Continue",
			configDirs: []string{"~/.continue"},
			risk:       model.RiskDanger,
			uninstallCmds: map[string]string{
				"linux":   "# Continue: remove from editor extensions + rm -rf ~/.continue",
				"macos":   "# Continue: remove from editor extensions + rm -rf ~/.continue",
				"windows": "# Continue: remove from editor extensions + Remove-Item ~/.continue -Recurse -Force",
			},
		},
		{
			id:         "amazon-q",
			name:       "Amazon Q Developer",
			configDirs: []string{"~/.aws/amazonq", "~/.amazon-q"},
			commands:   []string{"q"},
			risk:       model.RiskSafe,
			uninstallCmds: map[string]string{
				"linux":   "# Amazon Q: remove via installer",
				"macos":   "# Amazon Q: remove via installer",
				"windows": "# Amazon Q: remove via Add/Remove Programs",
			},
		},

		// ── NEW: Top missing agents ────────────────────────────────────
		{
			id:         "claw-code",
			name:       "Claw Code (oh-my-codex)",
			configDirs: []string{"~/.claw"},
			configFiles: []string{"~/.claw.json"},
			commands:   []string{"claw"},
			risk:       model.RiskCaution,
			uninstallCmds: map[string]string{
				"linux":   "rm -rf ~/.claw ~/.claw.json",
				"macos":   "rm -rf ~/.claw ~/.claw.json",
				"windows": "Remove-Item ~/.claw,~/.claw.json -Recurse -Force",
			},
		},
		{
			id:         "openhands",
			name:       "OpenHands (OpenDevin)",
			configDirs: []string{"~/.openhands"},
			pipPackages: []string{"openhands-ai"},
			risk:       model.RiskSafe,
			uninstallCmds: map[string]string{
				"linux":   "pip uninstall openhands-ai && rm -rf ~/.openhands",
				"macos":   "pip uninstall openhands-ai && rm -rf ~/.openhands",
				"windows": "pip uninstall openhands-ai; Remove-Item ~/.openhands -Recurse -Force",
			},
		},
		{
			id:         "open-interpreter",
			name:       "Open Interpreter",
			configDirs: []string{"~/.interpreter"},
			commands:   []string{"interpreter"},
			pipPackages: []string{"open-interpreter"},
			risk:       model.RiskCaution,
			uninstallCmds: map[string]string{
				"linux":   "pip uninstall open-interpreter && rm -rf ~/.interpreter",
				"macos":   "pip uninstall open-interpreter && rm -rf ~/.interpreter",
				"windows": "pip uninstall open-interpreter; Remove-Item ~/.interpreter -Recurse -Force",
			},
		},
		{
			id:         "goose",
			name:       "Goose (Block)",
			configDirs: []string{"~/.config/goose"},
			commands:   []string{"goose"},
			risk:       model.RiskSafe,
			uninstallCmds: map[string]string{
				"linux":   "# Goose: remove binary + rm -rf ~/.config/goose",
				"macos":   "# Goose: remove binary + rm -rf ~/.config/goose",
				"windows": "# Goose: remove binary + Remove-Item ~/.config\\goose -Recurse -Force",
			},
		},
		{
			id:         "roo-code",
			name:       "Roo Code CLI",
			configDirs: []string{"~/.roo-code"},
			commands:   []string{"roo"},
			npmPackages: []string{"@roocode/cli"},
			risk:       model.RiskSafe,
			uninstallCmds: map[string]string{
				"linux":   "npm uninstall -g @roocode/cli && rm -rf ~/.roo-code",
				"macos":   "npm uninstall -g @roocode/cli && rm -rf ~/.roo-code",
				"windows": "npm uninstall -g @roocode/cli; Remove-Item ~/.roo-code -Recurse -Force",
			},
		},
		{
			id:         "crush",
			name:       "Crush (Charmbracelet)",
			configDirs: []string{"~/.config/crush"},
			commands:   []string{"crush"},
			brewPackages: []string{"crush"},
			risk:       model.RiskSafe,
			uninstallCmds: map[string]string{
				"linux":   "# Crush: remove binary + rm -rf ~/.config/crush",
				"macos":   "brew uninstall crush && rm -rf ~/.config/crush",
				"windows": "# Crush: remove binary + Remove-Item ~/.config\\crush -Recurse -Force",
			},
		},
		{
			id:         "qwen-code",
			name:       "Qwen Code (Alibaba)",
			configDirs: []string{"~/.qwen"},
			commands:   []string{"qwen-code", "qwen"},
			brewPackages: []string{"qwen-code"},
			risk:       model.RiskSafe,
			uninstallCmds: map[string]string{
				"linux":   "# Qwen Code: remove binary + rm -rf ~/.qwen",
				"macos":   "brew uninstall qwen-code && rm -rf ~/.qwen",
				"windows": "# Qwen Code: remove binary + Remove-Item ~/.qwen -Recurse -Force",
			},
		},
		{
			id:         "kilo-code",
			name:       "Kilo Code CLI",
			configDirs: []string{"~/.kilo-code"},
			commands:   []string{"kilo"},
			npmPackages: []string{"@kilo-org/cli"},
			risk:       model.RiskSafe,
			uninstallCmds: map[string]string{
				"linux":   "npm uninstall -g @kilo-org/cli && rm -rf ~/.kilo-code",
				"macos":   "npm uninstall -g @kilo-org/cli && rm -rf ~/.kilo-code",
				"windows": "npm uninstall -g @kilo-org/cli; Remove-Item ~/.kilo-code -Recurse -Force",
			},
		},
		{
			id:         "plandex",
			name:       "Plandex",
			configDirs: []string{"~/.plandex"},
			commands:   []string{"plandex"},
			brewPackages: []string{"plandex"},
			risk:       model.RiskSafe,
			uninstallCmds: map[string]string{
				"linux":   "# Plandex: remove binary + rm -rf ~/.plandex",
				"macos":   "brew uninstall plandex && rm -rf ~/.plandex",
				"windows": "# Plandex: remove binary + Remove-Item ~/.plandex -Recurse -Force",
			},
		},
		{
			id:         "swe-agent",
			name:       "SWE-agent",
			pipPackages: []string{"swe-agent"},
			risk:       model.RiskSafe,
			uninstallCmds: map[string]string{
				"linux":   "pip uninstall swe-agent",
				"macos":   "pip uninstall swe-agent",
				"windows": "pip uninstall swe-agent",
			},
		},
		{
			id:         "trae-agent",
			name:       "Trae Agent (ByteDance)",
			configFiles: []string{"trae_config.yaml"},
			pipPackages: []string{"trae-agent"},
			risk:       model.RiskSafe,
			uninstallCmds: map[string]string{
				"linux":   "pip uninstall trae-agent",
				"macos":   "pip uninstall trae-agent",
				"windows": "pip uninstall trae-agent",
			},
		},
		{
			id:         "hermes-agent",
			name:       "Hermes Agent (Nous Research)",
			configDirs: []string{"~/.hermes"},
			commands:   []string{"hermes"},
			risk:       model.RiskSafe,
			uninstallCmds: map[string]string{
				"linux":   "# Hermes Agent: remove binary + rm -rf ~/.hermes",
				"macos":   "# Hermes Agent: remove binary + rm -rf ~/.hermes",
				"windows": "# Hermes Agent: remove binary + Remove-Item ~/.hermes -Recurse -Force",
			},
		},
		{
			id:         "kimi-cli",
			name:       "Kimi Code CLI (Moonshot AI)",
			commands:   []string{"kimi"},
			brewPackages: []string{"kimi-cli"},
			risk:       model.RiskSafe,
			uninstallCmds: map[string]string{
				"linux":   "# Kimi CLI: remove binary + rm -rf ~/.kimi",
				"macos":   "brew uninstall kimi-cli",
				"windows": "# Kimi CLI: remove binary",
			},
		},
		{
			id:         "groq-code-cli",
			name:       "Groq Code CLI",
			configDirs: []string{"~/.groq"},
			npmPackages: []string{"@groq/code-cli"},
			risk:       model.RiskSafe,
			uninstallCmds: map[string]string{
				"linux":   "npm uninstall -g @groq/code-cli && rm -rf ~/.groq",
				"macos":   "npm uninstall -g @groq/code-cli && rm -rf ~/.groq",
				"windows": "npm uninstall -g @groq/code-cli; Remove-Item ~/.groq -Recurse -Force",
			},
		},
		{
			id:         "grok-cli",
			name:       "Grok CLI (xAI)",
			commands:   []string{"grok"},
			npmPackages: []string{"@xai/grok-cli"},
			risk:       model.RiskSafe,
			uninstallCmds: map[string]string{
				"linux":   "npm uninstall -g @xai/grok-cli",
				"macos":   "npm uninstall -g @xai/grok-cli",
				"windows": "npm uninstall -g @xai/grok-cli",
			},
		},
		{
			id:         "devon",
			name:       "Devon (Entropy Research)",
			configDirs: []string{"~/.devon"},
			commands:   []string{"devon", "devon-tui"},
			risk:       model.RiskSafe,
			uninstallCmds: map[string]string{
				"linux":   "# Devon: remove binary + rm -rf ~/.devon",
				"macos":   "# Devon: remove binary + rm -rf ~/.devon",
				"windows": "# Devon: remove binary + Remove-Item ~/.devon -Recurse -Force",
			},
		},
		{
			id:         "claurst",
			name:       "Claurst (Claude Code in Rust)",
			configDirs: []string{"~/.claurst"},
			commands:   []string{"claurst"},
			risk:       model.RiskSafe,
			uninstallCmds: map[string]string{
				"linux":   "# Claurst: remove binary + rm -rf ~/.claurst",
				"macos":   "# Claurst: remove binary + rm -rf ~/.claurst",
				"windows": "# Claurst: remove binary + Remove-Item ~/.claurst -Recurse -Force",
			},
		},
		{
			id:         "free-code",
			name:       "Free Code (Claude Code fork)",
			commands:   []string{"free-code"},
			risk:       model.RiskSafe,
			uninstallCmds: map[string]string{
				"linux":   "# Free Code: remove binary",
				"macos":   "# Free Code: remove binary",
				"windows": "# Free Code: remove binary",
			},
		},
		{
			id:         "letta-code",
			name:       "Letta Code (MemGPT)",
			configDirs: []string{"~/.letta"},
			commands:   []string{"letta-code"},
			pipPackages: []string{"letta-code"},
			risk:       model.RiskSafe,
			uninstallCmds: map[string]string{
				"linux":   "pip uninstall letta-code && rm -rf ~/.letta",
				"macos":   "pip uninstall letta-code && rm -rf ~/.letta",
				"windows": "pip uninstall letta-code; Remove-Item ~/.letta -Recurse -Force",
			},
		},
		{
			id:         "forgecode",
			name:       "ForgeCode",
			configDirs: []string{"~/.forgecode"},
			configFiles: []string{"forge.yaml"},
			commands:   []string{"forge"},
			risk:       model.RiskSafe,
			uninstallCmds: map[string]string{
				"linux":   "# ForgeCode: remove binary + rm -rf ~/.forgecode",
				"macos":   "# ForgeCode: remove binary + rm -rf ~/.forgecode",
				"windows": "# ForgeCode: remove binary + Remove-Item ~/.forgecode -Recurse -Force",
			},
		},
		{
			id:         "cline-cli",
			name:       "Cline CLI",
			commands:   []string{"cline"},
			npmPackages: []string{"@cline/cli"},
			risk:       model.RiskSafe,
			uninstallCmds: map[string]string{
				"linux":   "npm uninstall -g @cline/cli",
				"macos":   "npm uninstall -g @cline/cli",
				"windows": "npm uninstall -g @cline/cli",
			},
		},

		// ── NEW: OpenClaw Ecosystem ───────────────────────────────────
		{
			id:          "openclaw",
			name:        "OpenClaw",
			configDirs:  []string{"~/.openclaw"},
			commands:    []string{"openclaw"},
			npmPackages: []string{"openclaw"},
			risk:        model.RiskDanger,
			uninstallCmds: map[string]string{
				"linux":   "npm uninstall -g openclaw && rm -rf ~/.openclaw",
				"macos":   "npm uninstall -g openclaw && rm -rf ~/.openclaw",
				"windows": "npm uninstall -g openclaw; Remove-Item ~/.openclaw -Recurse -Force",
			},
		},
		{
			id:          "nanobot",
			name:        "NanoBot (OpenClaw-lite)",
			configDirs:  []string{"~/.nanobot"},
			commands:    []string{"nanobot"},
			pipPackages: []string{"nanobot"},
			risk:        model.RiskCaution,
			uninstallCmds: map[string]string{
				"linux":   "pip uninstall nanobot && rm -rf ~/.nanobot",
				"macos":   "pip uninstall nanobot && rm -rf ~/.nanobot",
				"windows": "pip uninstall nanobot; Remove-Item ~/.nanobot -Recurse -Force",
			},
		},
		{
			id:         "zeroclaw",
			name:       "ZeroClaw",
			configDirs: []string{"~/.zeroclaw"},
			commands:   []string{"zeroclaw"},
			risk:       model.RiskCaution,
			uninstallCmds: map[string]string{
				"linux":   "# ZeroClaw: remove binary + rm -rf ~/.zeroclaw",
				"macos":   "# ZeroClaw: remove binary + rm -rf ~/.zeroclaw",
				"windows": "# ZeroClaw: remove binary + Remove-Item ~/.zeroclaw -Recurse -Force",
			},
		},
		{
			id:         "picoclaw",
			name:       "PicoClaw",
			configDirs: []string{"~/.picoclaw"},
			commands:   []string{"picoclaw"},
			risk:       model.RiskSafe,
			uninstallCmds: map[string]string{
				"linux":   "# PicoClaw: remove binary + rm -rf ~/.picoclaw",
				"macos":   "# PicoClaw: remove binary + rm -rf ~/.picoclaw",
				"windows": "# PicoClaw: remove binary + Remove-Item ~/.picoclaw -Recurse -Force",
			},
		},
		{
			id:         "nanoclaw",
			name:       "NanoClaw",
			configDirs: []string{"~/.nanoclaw"},
			commands:   []string{"nanoclaw"},
			risk:       model.RiskCaution,
			uninstallCmds: map[string]string{
				"linux":   "# NanoClaw: remove binary + rm -rf ~/.nanoclaw",
				"macos":   "# NanoClaw: remove binary + rm -rf ~/.nanoclaw",
				"windows": "# NanoClaw: remove binary + Remove-Item ~/.nanoclaw -Recurse -Force",
			},
		},
		{
			id:         "ironclaw",
			name:       "IronClaw (NEAR AI)",
			configDirs: []string{"~/.ironclaw"},
			commands:   []string{"ironclaw"},
			risk:       model.RiskSafe,
			uninstallCmds: map[string]string{
				"linux":   "# IronClaw: remove binary + rm -rf ~/.ironclaw",
				"macos":   "# IronClaw: remove binary + rm -rf ~/.ironclaw",
				"windows": "# IronClaw: remove binary + Remove-Item ~/.ironclaw -Recurse -Force",
			},
		},
		{
			id:         "nullclaw",
			name:       "NullClaw",
			configDirs: []string{"~/.nullclaw"},
			commands:   []string{"nullclaw"},
			risk:       model.RiskSafe,
			uninstallCmds: map[string]string{
				"linux":   "# NullClaw: remove binary + rm -rf ~/.nullclaw",
				"macos":   "# NullClaw: remove binary + rm -rf ~/.nullclaw",
				"windows": "# NullClaw: remove binary + Remove-Item ~/.nullclaw -Recurse -Force",
			},
		},
		{
			id:         "moltis",
			name:       "Moltis (OpenClaw-alt)",
			configDirs: []string{"~/.moltis"},
			commands:   []string{"moltis"},
			risk:       model.RiskSafe,
			uninstallCmds: map[string]string{
				"linux":   "# Moltis: remove binary + rm -rf ~/.moltis",
				"macos":   "# Moltis: remove binary + rm -rf ~/.moltis",
				"windows": "# Moltis: remove binary + Remove-Item ~/.moltis -Recurse -Force",
			},
		},
		{
			id:         "clawith",
			name:       "Clawith (OpenClaw for Teams)",
			configDirs: []string{"~/.clawith"},
			commands:   []string{"clawith"},
			risk:       model.RiskSafe,
			uninstallCmds: map[string]string{
				"linux":   "# Clawith: remove binary + rm -rf ~/.clawith",
				"macos":   "# Clawith: remove binary + rm -rf ~/.clawith",
				"windows": "# Clawith: remove binary + Remove-Item ~/.clawith -Recurse -Force",
			},
		},

		// ── Cloud / SaaS AI Agents ────────────────────────────────────
		{
			id:         "replit-agent",
			name:       "Replit Agent",
			configDirs: []string{"~/.replit"},
			commands:   []string{"replit"},
			npmPackages: []string{"replit-client"},
			risk:       model.RiskCaution,
			uninstallCmds: map[string]string{
				"linux":   "# Replit Agent: cloud-based, remove local config with rm -rf ~/.replit",
				"macos":   "# Replit Agent: cloud-based, remove local config with rm -rf ~/.replit",
				"windows": "# Replit Agent: cloud-based, remove local config with Remove-Item ~/.replit -Recurse -Force",
			},
		},
		{
			id:         "devin",
			name:       "Devin (Cognition)",
			configDirs: []string{"~/.devin"},
			commands:   []string{"devin"},
			risk:       model.RiskCaution,
			uninstallCmds: map[string]string{
				"linux":   "# Devin: cloud-based (cognition.ai), remove local config with rm -rf ~/.devin",
				"macos":   "# Devin: cloud-based, remove local config with rm -rf ~/.devin",
				"windows": "# Devin: cloud-based, remove local config with Remove-Item ~/.devin -Recurse -Force",
			},
		},
		{
			id:         "deepseek-cli",
			name:       "DeepSeek CLI",
			configDirs: []string{"~/.deepseek"},
			commands:   []string{"deepseek"},
			npmPackages: []string{"deepseek-cli"},
			risk:       model.RiskCaution,
			uninstallCmds: map[string]string{
				"linux":   "npm uninstall -g deepseek-cli && rm -rf ~/.deepseek",
				"macos":   "npm uninstall -g deepseek-cli && rm -rf ~/.deepseek",
				"windows": "npm uninstall -g deepseek-cli; Remove-Item ~/.deepseek -Recurse -Force",
			},
		},
		{
			id:         "sweep-ai",
			name:       "Sweep AI",
			configDirs: []string{"~/.sweep"},
			pipPackages: []string{"sweepai"},
			risk:       model.RiskSafe,
			uninstallCmds: map[string]string{
				"linux":   "pip uninstall sweepai && rm -rf ~/.sweep",
				"macos":   "pip uninstall sweepai && rm -rf ~/.sweep",
				"windows": "pip uninstall sweepai; Remove-Item ~/.sweep -Recurse -Force",
			},
		},
		{
			id:         "pieces",
			name:       "Pieces for Developers",
			configDirs: []string{"~/.pieces"},
			commands:   []string{"pieces"},
			risk:       model.RiskCaution,
			uninstallCmds: map[string]string{
				"linux":   "# Pieces: remove via Add/Remove Programs or rm -rf ~/.pieces",
				"macos":   "rm -rf /Applications/Pieces.app ~/.pieces",
				"windows": "# Pieces: uninstall via Add/Remove Programs",
			},
		},

		// ── Orchestration & Harnesses ─────────────────────────────────
		{
			id:          "claude-squad",
			name:        "Claude Squad",
			configDirs:  []string{"~/.claude-squad"},
			npmPackages: []string{"claude-squad"},
			risk:        model.RiskSafe,
			uninstallCmds: map[string]string{
				"linux":   "npm uninstall -g claude-squad && rm -rf ~/.claude-squad",
				"macos":   "npm uninstall -g claude-squad && rm -rf ~/.claude-squad",
				"windows": "npm uninstall -g claude-squad; Remove-Item ~/.claude-squad -Recurse -Force",
			},
		},
		{
			id:          "claude-flow",
			name:        "Claude Flow",
			configDirs:  []string{"~/.claude-flow"},
			npmPackages: []string{"claude-flow"},
			risk:        model.RiskSafe,
			uninstallCmds: map[string]string{
				"linux":   "npm uninstall -g claude-flow && rm -rf ~/.claude-flow",
				"macos":   "npm uninstall -g claude-flow && rm -rf ~/.claude-flow",
				"windows": "npm uninstall -g claude-flow; Remove-Item ~/.claude-flow -Recurse -Force",
			},
		},

		// ── 2026 Harness Additions (web-verified signatures) ──────────
		{
			id:          "mimo-code",
			name:        "MiMo Code (Xiaomi)",
			configDirs:  []string{"~/.config/mimocode"},
			commands:    []string{"mimo"},
			npmPackages: []string{"@mimo-ai/cli"},
			risk:        model.RiskDanger,
			uninstallCmds: map[string]string{
				"linux":   "npm uninstall -g @mimo-ai/cli && rm -rf ~/.config/mimocode ~/.local/share/mimocode ~/.cache/mimocode ~/.local/state/mimocode",
				"macos":   "npm uninstall -g @mimo-ai/cli && rm -rf ~/.config/mimocode ~/.local/share/mimocode ~/.cache/mimocode ~/.local/state/mimocode",
				"windows": "npm uninstall -g @mimo-ai/cli; Remove-Item ~/.config/mimocode -Recurse -Force",
			},
		},
		{
			id:          "amp",
			name:        "Amp (Sourcegraph)",
			configDirs:  []string{"~/.config/amp"},
			commands:    []string{"amp"},
			npmPackages: []string{"@ampcode/cli", "@sourcegraph/amp"},
			risk:        model.RiskCaution,
			uninstallCmds: map[string]string{
				"linux":   "npm uninstall -g @ampcode/cli @sourcegraph/amp && rm -rf ~/.config/amp",
				"macos":   "npm uninstall -g @ampcode/cli @sourcegraph/amp && rm -rf ~/.config/amp",
				"windows": "npm uninstall -g @ampcode/cli @sourcegraph/amp; Remove-Item ~/.config/amp -Recurse -Force",
			},
		},
		{
			id:         "droid",
			name:       "Droid (Factory)",
			configDirs: []string{"~/.factory"},
			commands:   []string{"droid"},
			risk:       model.RiskCaution,
			uninstallCmds: map[string]string{
				"linux":   "# Droid: remove binary + rm -rf ~/.factory",
				"macos":   "# Droid: remove binary + rm -rf ~/.factory",
				"windows": "# Droid: remove binary + Remove-Item ~/.factory -Recurse -Force",
			},
		},
		{
			id:          "gptme",
			name:        "gptme",
			configDirs:  []string{"~/.config/gptme"},
			commands:    []string{"gptme"},
			pipPackages: []string{"gptme"},
			risk:        model.RiskCaution,
			uninstallCmds: map[string]string{
				"linux":   "pip uninstall gptme && rm -rf ~/.config/gptme ~/.local/share/gptme",
				"macos":   "pip uninstall gptme && rm -rf ~/.config/gptme ~/.local/share/gptme",
				"windows": "pip uninstall gptme; Remove-Item ~/.config/gptme -Recurse -Force",
			},
		},
		{
			id:          "codebuff",
			name:        "Codebuff",
			configDirs:  []string{"~/.codebuff"},
			commands:    []string{"codebuff"},
			npmPackages: []string{"codebuff"},
			risk:        model.RiskCaution,
			uninstallCmds: map[string]string{
				"linux":   "npm uninstall -g codebuff && rm -rf ~/.codebuff",
				"macos":   "npm uninstall -g codebuff && rm -rf ~/.codebuff",
				"windows": "npm uninstall -g codebuff; Remove-Item ~/.codebuff -Recurse -Force",
			},
		},
		{
			id:       "cursor-agent",
			name:     "Cursor Agent CLI (Anysphere)",
			commands: []string{"cursor-agent"},
			risk:     model.RiskCaution,
			uninstallCmds: map[string]string{
				"linux":   "# Cursor Agent CLI: remove the 'cursor-agent' binary",
				"macos":   "# Cursor Agent CLI: remove the 'cursor-agent' binary",
				"windows": "# Cursor Agent CLI: remove the 'cursor-agent' binary",
			},
		},
	}

	for _, agent := range agents {
		found := false
		var totalSize int64
		var foundPaths []string

		for _, cmd := range agent.commands {
			if s.hasCommand(cmd) {
				found = true
			}
		}

		for _, dir := range agent.configDirs {
			expanded := s.expandPath(dir)
			if s.dirExists(expanded) {
				found = true
				totalSize += s.dirSize(expanded)
				foundPaths = append(foundPaths, expanded)
			}
		}

		for _, file := range agent.configFiles {
			expanded := s.expandPath(file)
			if s.fileExists(expanded) {
				found = true
				totalSize += s.fileSize(expanded)
				foundPaths = append(foundPaths, expanded)
			}
		}

		for _, pkg := range agent.npmPackages {
			if s.hasNpmPackage(pkg) {
				found = true
			}
		}

		for _, pkg := range agent.pipPackages {
			if s.hasPipPackage(pkg) {
				found = true
			}
		}

		for _, pkg := range agent.brewPackages {
			if s.hasBrewPackage(pkg) {
				found = true
			}
		}

		for _, pkg := range agent.goPackages {
			if s.hasGoPackage(pkg) {
				found = true
			}
		}

		if found {
			path := strings.Join(foundPaths, ", ")
			if path == "" {
				path = "(binary found in PATH)"
			}
			s.addFinding(model.Finding{
				ID:            agent.id,
				Category:      model.CatAgents,
				Name:          agent.name,
				Path:          path,
				SizeBytes:     totalSize,
				ConfigPaths:   foundPaths,
				RiskLevel:     agent.risk,
				UninstallCmds: agent.uninstallCmds,
				SubItems:      agent.npmPackages,
			})
		}
	}
}

// hasNpmPackage checks if an npm global package is installed.
func (s *Scanner) hasNpmPackage(pkg string) bool {
	if !s.hasCommand("npm") {
		return false
	}
	out, err := s.runCommand("npm", "list", "-g", pkg, "--depth=0")
	if err != nil {
		return false
	}
	return strings.Contains(out, pkg) && !strings.Contains(out, "(empty)")
}

// hasPipPackage checks if a pip package is installed.
func (s *Scanner) hasPipPackage(pkg string) bool {
	// Check if pip exists first to avoid slow timeouts on machines without pip
	if !s.hasCommand("pip3") && !s.hasCommand("pip") {
		return false
	}
	for _, cmd := range []string{"pip3", "pip"} {
		out, err := s.runCommand(cmd, "show", pkg)
		if err == nil && out != "" {
			return true
		}
	}
	return false
}

// hasBrewPackage checks if a Homebrew package is installed.
func (s *Scanner) hasBrewPackage(pkg string) bool {
	if !s.hasCommand("brew") {
		return false
	}
	out, err := s.runCommand("brew", "list", pkg)
	if err != nil {
		return false
	}
	return strings.TrimSpace(out) != ""
}

// hasGoPackage checks if a Go package is installed.
func (s *Scanner) hasGoPackage(pkg string) bool {
	if !s.hasCommand("go") {
		return false
	}
	out, err := s.runCommand("go", "list", "-m", pkg)
	if err != nil {
		return false
	}
	return strings.TrimSpace(out) != ""
}

// findInstructionFiles recursively finds instruction config files from a root directory.
func findInstructionFiles(root string, maxDepth int, filenames []string) []string {
	var results []string

	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		rel, relErr := filepath.Rel(root, path)
		if relErr != nil {
			return nil
		}
		depth := len(strings.Split(rel, string(filepath.Separator)))
		if depth > maxDepth {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if !info.IsDir() {
			for _, target := range filenames {
				if info.Name() == target {
					results = append(results, path)
					break
				}
			}
		}
		return nil
	})

	return results
}
