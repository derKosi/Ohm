# ⚡ Ohm

**Resistance against AGI bloat.**

Ohm measures what's there — and helps you remove what you don't need.

A cross-platform (Windows, macOS, Linux) TUI tool that scans your system for AI-related software, shows you what's installed, how much space it eats, and generates an uninstall script you can review before running. **Ohm never deletes anything itself.**

> **Privacy first.** Ohm is 100% offline. No telemetry, no phone-home, no cloud. Your AI config files often contain API keys, project paths, and personal instructions. That data stays on your machine. Period.

---

## Why?

You test AI tools. You install agents, harnesses, runtimes, download models, set up SDKs. Months later your disk is full of 14 GB model files from that one experiment, three competing AI editors you forgot about, and a `.cache/huggingface` directory that quietly grew to 80 GB.

Ohm finds all of it. You pick what goes. Ohm writes the script. You run it when you're ready.

## Privacy-First Architecture

AI tools store sensitive data on disk:

- **API keys** in config files (`ANTHROPIC_API_KEY`, `OPENAI_API_KEY`, custom provider tokens)
- **Project paths** revealing client names, employer info, proprietary structures
- **Conversation history** with PII, code snippets, business logic
- **MCP server configs** with credentials, database connection strings
- **Agent soul files** containing proprietary prompts and business rules
- **Memory/session files** with accumulated context about your work
- **Plugin configs** that may contain third-party API keys

Ohm treats all of this as **confidential**. The tool:
- Never transmits scan results anywhere
- Has no network capability — no HTTP client in the codebase
- Stores all state locally (`~/.ohm/state.json`)
- Generates cleanup scripts that run entirely on your machine
- No analytics, no crash reporting, no update checks that phone home

This is a deliberate design choice and a core differentiator. If you find network code in Ohm, that's a bug.

## Features

- 🔍 **Scan** — Detects AI agents, harnesses, runtimes, editors, SDKs, model caches, config dirs, Docker images, shell modifications, and stragglers
- 📋 **List** — Bubble Tea TUI with categories, sizes, paths, and checkboxes
- 📝 **Script output** — Generates a `.sh` / `.ps1` removal script with all selected items
- 🚫 **Never executes** — Ohm is read-only. The generated script is yours to review and run
- 🔒 **Offline by design** — zero network calls, no telemetry, no data leaves your machine
- 🧹 **Straggler detection** — Finds orphaned model files, stale PATH entries, leftover configs, dead services
- 🧠 **Memory file detection** — Finds agent memory, session history, conversation logs, accumulated context
- 🔑 **Credential-aware** — Flags locations that likely contain API keys or tokens (shows warning, never exposes contents)
- 📁 **Agent instruction files** — Finds AGENTS.md, CLAUDE.md, .cursorrules, .windsurfrules, GEMINI.md, soul files, MCP configs
- 🔌 **Plugin detection** — Finds ComfyUI custom nodes, pi skills, MCP servers, editor extensions
- 🛤️ **PATH & ENV scanning** *(opt-in)* — Detects AI-related entries in PATH, environment variables, shell profiles
- 💾 **State between runs** — Remembers what you uninstalled last time and flags remaining leftovers
- 🖥️ **Cross-platform** — Single Go binary, zero runtime dependencies, works on Windows / macOS / Linux

## Scan Categories

| # | Category | What It Finds |
|---|----------|---------------|
| 1 | **Agents & Harnesses** | pi, Claude Code, Aider, Continue, Cline, Codex CLI, Amazon Q, Cursor Agent, OpenCode, PaperclipAI, ZenFlow, Antigravity, Mistral Vibe, custom harnesses |
| 2 | **AI Editors & IDEs** | Cursor, Windsurf, Augment, Zed (AI features), GitHub Copilot |
| 3 | **Model Runtimes** | Ollama, LM Studio, LocalAI, text-generation-webui, llama.cpp, GPT4All, KoboldCpp |
| 4 | **ComfyUI & Image Models** | ComfyUI install + checkpoints, LoRA adapters, ControlNet, VAEs, CLIP, UNet, embeddings, upscale models, custom nodes |
| 5 | **SDKs & Frameworks** | PyTorch, TensorFlow, HuggingFace, Anthropic SDK, OpenAI SDK, LangChain, LlamaIndex, Playwright, Selenium |
| 6 | **Model Caches** | HuggingFace cache, .gguf files, Ollama model store, safetensors, PyTorch hub cache |
| 7 | **Agent Config & Instructions** | AGENTS.md, CLAUDE.md, .cursorrules, .windsurfrules, GEMINI.md, copilot-instructions.md, soul files, system prompts, CONVENTIONS.md |
| 8 | **Agent Memory & Sessions** | pi sessions, Claude conversation history, Aider chat logs, Continue session data, Vibe history, PaperclipAI context |
| 9 | **MCP Configurations** | .mcp.json files, MCP server configs (often contain API keys and connection strings) |
| 10 | **Plugins & Extensions** | ComfyUI custom nodes, pi skills/plugins, VS Code AI extensions, JetBrains AI plugins |
| 11 | **Config & Data Dirs** | ~/.claude/, ~/.pi/, ~/.aider/, ~/.cursor/, ~/.codex/, ~/.gemini/, ~/.vibe/, ~/.paperclip/, AppData entries, .config dirs |
| 12 | **Docker** | AI-related images and volumes (ollama, vllm, comfyui, localai, etc.) |
| 13 | **Stragglers** | Orphaned files from already-uninstalled tools, leftover model weights, dead services |

### Opt-In Scans (disabled by default, enabled with flags)

| Flag | What It Scans |
|------|---------------|
| `--path` | PATH entries pointing to AI tools, stale entries from removed software |
| `--env` | Environment variables containing API keys, model paths, AI-related config |
| `--shell-profile` | Shell profile modifications (.bashrc, .zshrc, PowerShell $PROFILE, .profile) |
| `--deep` | Full home directory crawl for any AI-related file signatures (slower, more thorough) |

### Detection Methods

Ohm uses two detection layers:

1. **Known software database** — Curated signatures with install paths, config locations, and uninstall commands. See [`docs/SIGNATURES.md`](docs/SIGNATURES.md) for the full catalog. Extensible via `~/.ohm/signatures/` (drop-in YAML files for custom/private tools)
2. **Heuristic detection** — Filesystem fingerprints: known filenames (AGENTS.md, .gguf, safetensors), directory patterns (models/, checkpoints/, loras/), package names with AI-related keywords, running processes with AI-related names

## How It Works

### Real output from a Windows 11 machine:

```
C:\AI> ohm scan --no-tui

  ───┤   ⚡  O H M     ├───
🤖 Agents & Harnesses (9 found, 2.6 GB)
   ⚠️  pi (Coding Agent)              15.3 MB    C:\Users\user\.pi
   🔑 Claude Code                    407.2 MB   C:\Users\user\.claude
      Aider                          0 B        (binary found in PATH)
   🔑 Codex CLI (OpenAI)             138.6 MB   C:\Users\user\.codex
   🔑 Gemini CLI (Google)            725.3 MB   C:\Users\user\.gemini
      Mistral Vibe                   1.9 MB     C:\Users\user\.vibe
   ⚠️  PaperclipAI                    0 B        C:\Users\user\.paperclip
      OpenCode                       0 B        (binary found in PATH)
   🔑 Continue                       1.4 GB     C:\Users\user\.continue

🖥️ AI Editors & IDEs (1 found, 82.4 KB)
   ⚠️  Cursor IDE                     82.4 KB    C:\Users\user\.cursor

📦 SDKs & Frameworks (3 found, 0 B)
      HuggingFace Transformers       0 B
      OpenAI SDK                     0 B
   ⚠️  Playwright (AI-adjacent)       0 B

📄 Agent Config & Instructions (1 found, 10.7 KB)
   ⚠️  Agent Instruction Files        10.7 KB    (scattered across projects)

🧠 Agent Memory & Sessions (4 found, 1.3 GB)
   🔑 Claude Code Memory             407.1 MB   C:\Users\user\.claude
      Mistral Vibe History           6.6 KB     C:\Users\user\.vibe\vibehistory
   🔑 PaperclipAI Context            213.3 MB   C:\Users\user\.paperclip
   ⚠️  Gemini CLI History             725.3 MB   C:\Users\user\.gemini

🔌 MCP Configurations (1 found, 22 B)
   🔑 MCP Configuration Files        22 B

🧩 Plugins & Extensions (1 found)
     - pi skill: paperclip
     - pi skill: para-memory-files

📁 Config & Data Dirs (9 found, 2.6 GB)
   🔑 Claude Config                  407.1 MB
   🔑 Gemini CLI Config              725.3 MB
   🔑 Continue Config                1.4 GB
   ...

Total: 31 items (6.6 GB)
```

### Real output from a Linux machine:

```
$ ohm scan --no-tui

  ───┤   ⚡  O H M     ├───
🤖 Agents & Harnesses (5 found, 13.3 GB)
   ⚠️  pi (Coding Agent)              92.3 MB    /home/user/.pi
   🔑 Claude Code                    278 B      /home/user/.claude, /home/user/.claude.json
   🔑 Gemini CLI (Google)            5.3 MB     /home/user/.gemini
      Mistral Vibe                   523.0 KB   /home/user/.vibe
   ⚠️  PaperclipAI                    13.2 GB    /home/user/.paperclip

⚙️ Model Runtimes (1 found, 261.6 MB)
      Ollama                         261.6 MB   /usr/share/ollama

💾 Model Caches (1 found, 1.2 GB)
      Playwright Browsers            1.2 GB     /home/user/.cache/ms-playwright

📄 Agent Config & Instructions (1 found, 55.8 KB)
   ⚠️  Agent Instruction Files        55.8 KB    (scattered across projects)

🧠 Agent Memory & Sessions (4 found, 13.2 GB)
   🔑 Claude Code Memory             242 B      /home/user/.claude
      Mistral Vibe History           5.0 KB     /home/user/.vibe/vibehistory
   🔑 PaperclipAI Context            13.2 GB    /home/user/.paperclip
   ⚠️  Gemini CLI History             5.3 MB     /home/user/.gemini

🧩 Plugins & Extensions (1 found, 4.3 KB)
   ⚠️  AI Plugins & Extensions        4.3 KB     (see sub-items)

📁 Config & Data Dirs (6 found, 13.3 GB)
   🔑 Claude Config                  242 B      /home/user/.claude
   🔑 Claude Config (JSON)           36 B       /home/user/.claude.json
   ⚠️  pi Config                      92.3 MB    /home/user/.pi
   ⚠️  Gemini CLI Config              5.3 MB     /home/user/.gemini
      Mistral Vibe Config            523.0 KB   /home/user/.vibe
   🔑 PaperclipAI Config             13.2 GB    /home/user/.paperclip

Total: 19 items (41.3 GB)
```

### Interactive TUI mode (press space to select, g to generate script):

```
  ───┤   ⚡  O H M     ├───
🔒 All scanning is local. No data leaves this machine.

🤖 Agents & Harnesses (5 found, 13.3 GB)
  [ ] ⚠️  pi (Coding Agent)              92.3 MB    /home/user/.pi
  [ ] 🔑 Claude Code                    278 B      /home/user/.claude, /home/user/.claude.json
  [ ] 🔑 Gemini CLI (Google)            5.3 MB     /home/user/.gemini
  [ ]    Mistral Vibe                   523.0 KB   /home/user/.vibe
  [x] ⚠️  PaperclipAI                    13.2 GB    /home/user/.paperclip

⚙️ Model Runtimes (1 found, 261.6 MB)
  [ ]    Ollama                         261.6 MB   /usr/share/ollama

💾 Model Caches (1 found, 1.2 GB)
  [ ]    Playwright Browsers            1.2 GB     /home/user/.cache/ms-playwright

📄 Agent Config & Instructions (1 found, 55.8 KB)
  [ ] ⚠️  Agent Instruction Files        55.8 KB    (scattered across projects)

🧠 Agent Memory & Sessions (4 found, 13.2 GB)
  [ ] 🔑 Claude Code Memory             242 B      /home/user/.claude
  [ ]    Mistral Vibe History           5.0 KB     /home/user/.vibe/vibehistory
  [x] 🔑 PaperclipAI Context            13.2 GB    ...text.json, /home/user/.paperclip/instances
  [ ] ⚠️  Gemini CLI History             5.3 MB     /home/user/.gemini

🧩 Plugins & Extensions (1 found, 4.3 KB)
  [ ] ⚠️  AI Plugins & Extensions        4.3 KB     (see sub-items)

📁 Config & Data Dirs (6 found, 13.3 GB)
  [ ] 🔑 Claude Config                  242 B      /home/user/.claude
  [ ] 🔑 Claude Config (JSON)           36 B       /home/user/.claude.json
  [ ] ⚠️  pi Config                      92.3 MB    /home/user/.pi
  [ ] ⚠️  Gemini CLI Config              5.3 MB     /home/user/.gemini
  [ ]    Mistral Vibe Config            523.0 KB   /home/user/.vibe
> [x] 🔑 PaperclipAI Config             13.2 GB    /home/user/.paperclip


Total: 19 items (41.3 GB) | Selected: 3 items (39.6 GB)
↑/k up • ↓/j down • pgup/pgdn • space select • a toggle all • g generate • q quit
MIT License © 2026 Mathias Kosinski · Built with Pi Harness + GLM-5.1
```

### Generated cleanup script:

```
$ ohm generate
  📝 Written: ohm-cleanup-2026-04-11.sh
  ⚠️  Review the script before running.
```

```bash
#!/usr/bin/env bash
# Ohm Cleanup Script
# ⚠️ REVIEW BEFORE RUNNING
# Total items: 2 | Estimated space to free: 545.8 MB

set -euo pipefail
echo "Removing Claude Code..."
npm uninstall -g @anthropic-ai/claude-code

echo "Removing Codex CLI (OpenAI)..."
npm uninstall -g @openai/codex
rm -rf ~/.codex
```

## Installation

```bash
# From source (requires Go 1.22+)
go install github.com/derKosi/Ohm@latest

# Or download binary from GitHub releases
# Single file, no dependencies, no installer, no internet required.
```

## Usage

```bash
ohm scan              # Scan system for AI software
ohm scan --path       # Also check PATH for AI tool entries
ohm scan --env        # Also check environment variables
ohm scan --shell      # Also check shell profiles
ohm scan --deep       # Thorough filesystem crawl
ohm scan --no-tui     # Text output (no TUI)
ohm scan --json       # JSON output for scripting/CI
ohm generate          # Generate cleanup script from last selection
ohm stragglers        # Scan only for leftover files from removed tools
ohm history           # Show what was removed in previous runs
ohm version           # Show version
```

## Custom Signatures *(coming soon)*

Ohm ships with a built-in database of 84+ known AI tools. Custom YAML signatures (`~/.ohm/signatures/*.yaml`) for private or niche tools are planned for a future release.

## Safety

- Ohm is **read-only**. It scans, lists, and writes scripts. It never deletes, uninstalls, or modifies anything.
- Generated scripts include comments explaining every action.
- Scripts require manual execution — you stay in control.
- Credential-containing files are flagged with ⚠️ but their contents are never displayed.
- No `os.Remove`, `os.RemoveAll`, or `exec.Command("rm")` anywhere in the codebase.

## Tech Stack

| Component | Choice |
|-----------|--------|
| Language | Go 1.22+ |
| TUI | [Bubble Tea](https://github.com/charmbracelet/bubbletea) |
| Styling | [Lip Gloss](https://github.com/charmbracelet/lipgloss) |
| CLI | Built-in (no framework) |
| Storage | Local JSON state file (`~/.ohm/state.json`) |
| Signatures | Built-in + drop-in YAML (`~/.ohm/signatures/*.yaml`) |
| Binary | Single static binary, zero dependencies, no network stack |

## Third-Party Libraries

Ohm uses the following open-source libraries:

| Library | License | Purpose |
|---------|---------|--------|
| [Bubble Tea](https://github.com/charmbracelet/bubbletea) | MIT | TUI framework |
| [Lip Gloss](https://github.com/charmbracelet/lipgloss) | MIT | Terminal styling |
| [Bubbles](https://github.com/charmbracelet/bubbles) | MIT | TUI components |
| [yaml.v3](https://github.com/go-yaml/yaml) | Apache 2.0 | YAML parsing |

## License

[MIT License](LICENSE) — © 2026 Mathias Kosinski

Free to use, modify, and distribute. Attribution required (keep the copyright notice).

## Documentation

| File | Description |
|------|-------------|
| [`docs/ROADMAP.md`](docs/ROADMAP.md) | Remaining work, priorities, custom signatures spec |
| [`docs/SIGNATURES.md`](docs/SIGNATURES.md) | Complete AI software signature catalog |
| [`docs/SALT-ANALYSIS.md`](docs/SALT-ANALYSIS.md) | Value proposition and salt analysis |
| [`docs/PROJECT-PLAN.md`](docs/PROJECT-PLAN.md) | Architecture, data model, scan pipeline |
| [`test-fixtures/`](test-fixtures/) | Real-world inventory files for testing |

---

*Ohm — Resistance against AGI bloat.*
*Designed with help of [Pi Harness](https://github.com/mariozechner/pi-coding-agent) and GLM-5.1.*
*[MIT License](LICENSE) — © 2026 Mathias Kosinski*
