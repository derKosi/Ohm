# Ohm — AI Software Signature Database

This document defines the complete catalog of AI software that Ohm detects.
Each entry specifies detection rules, config locations, and uninstall commands per platform.

---

## 1. Agents & Harnesses

### pi (Coding Agent)
- **Install:** `npm i -g @mariozechner/pi-coding-agent`
- **Binary:** `~/.nvm/versions/node/*/bin/pi` (via nvm), or global npm bin
- **Config:** `~/.pi/` (agent config, skills, sessions, providers)
- **Plugins:** `~/.pi/agent/skills/*/SKILL.md`
- **PATH:** `~/.pi/agent/bin`
- **Uninstall:**
  - Linux/macOS: `npm uninstall -g @mariozechner/pi-coding-agent && rm -rf ~/.pi`
  - Windows: `npm uninstall -g @mariozechner/pi-coding-agent; Remove-Item "$env:USERPROFILE\.pi" -Recurse -Force`
- **Stragglers:** pi-session-search, pi-subagents, pi-zai-provider, pi-memory, pi-zai-usage (related npm packages)

### Claude Code
- **Install:** `npm i -g @anthropic-ai/claude-code` or winget `Anthropic.Claude`
- **Binary:** `claude` in PATH
- **Config:** `~/.claude/`, `~/.claude.json`
- **Project configs:** `CLAUDE.md` files in projects
- **Uninstall:**
  - Linux/macOS: `npm uninstall -g @anthropic-ai/claude-code && rm -rf ~/.claude ~/.claude.json`
  - Windows: `winget uninstall Anthropic.Claude; Remove-Item "$env:USERPROFILE\.claude","$env:USERPROFILE\.claude.json","$env:APPDATA\Claude" -Recurse -Force`

### Aider
- **Install:** `pip install aider-chat`
- **Binary:** `~/.local/bin/aider`
- **Config:** `~/.aider.conf.yml`, `~/.aider*`
- **Uninstall:**
  - Linux/macOS: `pip uninstall aider-chat && rm -f ~/.aider.conf.yml`
  - Windows: `pip uninstall aider-chat; Remove-Item "$env:USERPROFILE\.aider.conf.yml" -Force`

### Codex CLI (OpenAI)
- **Install:** `npm i -g @openai/codex`
- **Binary:** `codex` in PATH
- **Config:** `~/.codex/`
- **Uninstall:**
  - All: `npm uninstall -g @openai/codex && rm -rf ~/.codex`

### Gemini CLI (Google)
- **Install:** `npm i -g @google/gemini-cli`
- **Binary:** `gemini` in PATH
- **Config:** `~/.gemini/`
- **Project configs:** `GEMINI.md` files in projects
- **Uninstall:**
  - All: `npm uninstall -g @google/gemini-cli && rm -rf ~/.gemini`

### Mistral Vibe
- **Binary:** `vibe` in PATH
- **Config:** `~/.vibe/` (config.toml, logs, vibehistory, trusted_folders.toml, update_cache.json)
- **Uninstall:** Remove binary + `rm -rf ~/.vibe`

### Continue
- **Config:** `~/.continue/`
- **VS Code extension:** `continue.continue`
- **JetBrains plugin:** `continue`
- **Uninstall:** Remove extension + `rm -rf ~/.continue`

### Cline (VS Code)
- **Extension:** `saoudrizwan.claude-dev`
- **Config:** VS Code extension storage
- **Uninstall:** Remove VS Code extension

### Amazon Q Developer
- **Install:** `npm i -g @anthropic-ai/amazon-q` or via installer
- **Config:** `~/.amazon-q/`
- **Uninstall:** Remove via CLI or installer

### Cursor
- **Install:** Desktop installer / AppImage
- **Binary:** `/opt/Cursor/`, `~/Applications/Cursor.app/`, `%LOCALAPPDATA%/Cursor/`
- **Config:** `~/.cursor/`, `~/.cursorrules`, `%APPDATA%/Cursor/`
- **Uninstall:**
  - Linux: `rm -rf /opt/Cursor ~/.cursor ~/.cursorrules`
  - macOS: `rm -rf /Applications/Cursor.app ~/.cursor ~/.cursorrules`
  - Windows: via Add/Remove Programs; `Remove-Item "$env:APPDATA\Cursor","$env:USERPROFILE\.cursor" -Recurse -Force`

### Windsurf
- **Install:** Desktop installer
- **Config:** `~/.windsurf/`, `~/.windsurfrules`
- **Uninstall:** Remove application + config dirs

### Augment
- **Config:** VS Code / JetBrains extension
- **Uninstall:** Remove extension

### OpenCode
- **Install:** `go install` + `npm i -g opencode-ai`
- **Binary:** `~/go/bin/opencode`, npm bin
- **Config:** `~/.opencode/`
- **Uninstall:** Remove binary + `rm -rf ~/.opencode`

### PaperclipAI
- **Install:** `npm i -g paperclipai`
- **Binary:** `paperclip` in PATH
- **Config:** `~/.paperclip/` (instances, docs, context)
- **Uninstall:** `npm uninstall -g paperclipai && rm -rf ~/.paperclip`

### Antigravity
- **Install:** Unknown (research needed)
- **Config:** Research needed
- **Detection:** Look for binary in PATH, config dir patterns

### Claw Code (oh-my-codex)
- **Binary:** `claw` in PATH
- **Config:** `~/.claw`, `~/.claw.json`
- **Note:** Open-source clean-room rewrite of Claude Code architecture, 72K+ GitHub stars

### OpenHands (OpenDevin)
- **Install:** `pip install openhands-ai`
- **Config:** `~/.openhands`
- **Uninstall:** `pip uninstall openhands-ai && rm -rf ~/.openhands`

### Open Interpreter
- **Install:** `pip install open-interpreter`
- **Binary:** `interpreter` in PATH
- **Config:** `~/.interpreter`

### Goose (Block)
- **Binary:** `goose` in PATH
- **Config:** `~/.config/goose`

### Roo Code CLI
- **Install:** `npm i -g @roocode/cli`
- **Binary:** `roo` in PATH
- **Config:** `~/.roo-code`

### Crush (Charmbracelet)
- **Install:** `brew install crush` or binary
- **Binary:** `crush` in PATH
- **Config:** `~/.config/crush`

### Qwen Code (Alibaba)
- **Install:** `brew install qwen-code` or binary
- **Binary:** `qwen-code`, `qwen` in PATH
- **Config:** `~/.qwen`

### Kilo Code CLI
- **Install:** `npm i -g @kilo-org/cli`
- **Binary:** `kilo` in PATH
- **Config:** `~/.kilo-code`

### Plandex
- **Install:** `brew install plandex` or binary
- **Binary:** `plandex` in PATH
- **Config:** `~/.plandex`

### SWE-agent
- **Install:** `pip install swe-agent`
- **Uninstall:** `pip uninstall swe-agent`

### Trae Agent (ByteDance)
- **Install:** `pip install trae-agent`
- **Config:** `trae_config.yaml`

### Hermes Agent (Nous Research)
- **Binary:** `hermes` in PATH
- **Config:** `~/.hermes`

### Kimi Code CLI (Moonshot AI)
- **Install:** `brew install kimi-cli` or binary
- **Binary:** `kimi` in PATH

### Groq Code CLI
- **Install:** `npm i -g @groq/code-cli`
- **Config:** `~/.groq`

### Grok CLI (xAI)
- **Install:** `npm i -g @xai/grok-cli`
- **Binary:** `grok` in PATH

### Devon (Entropy Research)
- **Binary:** `devon`, `devon-tui` in PATH
- **Config:** `~/.devon`

### Claurst (Claude Code in Rust)
- **Binary:** `claurst` in PATH
- **Config:** `~/.claurst`

### Free Code (Claude Code fork)
- **Binary:** `free-code` in PATH

### Letta Code (MemGPT)
- **Install:** `pip install letta-code`
- **Binary:** `letta-code` in PATH
- **Config:** `~/.letta`

### ForgeCode
- **Binary:** `forge` in PATH
- **Config:** `~/.forgecode`, `forge.yaml`

### Cline CLI
- **Install:** `npm i -g @cline/cli`
- **Binary:** `cline` in PATH

### Claude Squad
- **Install:** `npm i -g claude-squad`
- **Config:** `~/.claude-squad`

### Claude Flow
- **Install:** `npm i -g claude-flow`
- **Config:** `~/.claude-flow`

### Replit Agent
- **Config:** `~/.replit`
- **Note:** Cloud-based AI coding agent by Replit

### Devin (Cognition)
- **Config:** `~/.devin`
- **Note:** Cloud-based AI software engineer by Cognition AI

### DeepSeek CLI
- **Install:** `npm i -g deepseek-cli`
- **Binary:** `deepseek` in PATH
- **Config:** `~/.deepseek`

### Sweep AI
- **Install:** `pip install sweepai`
- **Config:** `~/.sweep`
- **Note:** AI coding assistant for JetBrains IDEs

### Pieces for Developers
- **Binary:** `pieces` in PATH
- **Config:** `~/.pieces`
- **Note:** On-device AI development assistant

### MiMo Code (Xiaomi)
- **Install:** `npm i -g @mimo-ai/cli` (or `curl -fsSL https://mimo.xiaomi.com/install | bash`)
- **Binary:** `mimo` in PATH
- **Config:** `~/.config/mimocode/` (global), data in `~/.local/share/mimocode/` (incl. `auth.json`), cache `~/.cache/mimocode/`
- **Note:** Terminal-native agent built on an OpenCode fork; persistent memory + sessions. Credential-bearing (`auth.json`)

### Amp (Sourcegraph)
- **Install:** `npm i -g @ampcode/cli` (legacy alias `@sourcegraph/amp`) or brew/choco
- **Binary:** `amp` in PATH
- **Config:** `~/.config/amp/` (`settings.json`)
- **Note:** Agentic coding tool by Sourcegraph, runs in VS Code and the CLI

### Droid (Factory)
- **Install:** `curl -fsSL https://app.factory.ai/cli | sh`
- **Binary:** `droid` in PATH
- **Config:** `~/.factory/` (`settings.json`)
- **Note:** Factory's multi-model CLI coding agent; requires login

### gptme
- **Install:** `pip install gptme`
- **Binary:** `gptme` in PATH
- **Config:** `~/.config/gptme/`, data `~/.local/share/gptme/`
- **Note:** Provider-agnostic terminal agent; executes shell/Python (code-running)

### Codebuff
- **Install:** `npm i -g codebuff`
- **Binary:** `codebuff` in PATH
- **Config:** `~/.codebuff/`
- **Note:** Multi-agent terminal coding assistant

### Cursor Agent CLI (Anysphere)
- **Install:** `curl https://cursor.com/install -fsSL | bash` (or npm `@nothumanwork/cursor-agents-sdk`)
- **Binary:** `cursor-agent` in PATH
- **Note:** Headless/CI agent mode of Cursor; shares `~/.cursor/` with the Cursor IDE

## 2. AI Editors & IDEs

### Cursor IDE
*(covered above in Agents section — dual category)*

### Windsurf (Codeium)
*(covered above)*

### Zed (with AI features)
- **Install:** AppImage, brew, installer
- **Config:** `~/.config/zed/`
- **AI features:** Built-in assistant, context server configs
- **Uninstall:** Remove application + `rm -rf ~/.config/zed`

### Warp Terminal (AI)
- **Install:** Desktop installer, package manager
- **Config:** `~/.warp/`
- **AI features:** AI command search, natural language terminal commands
- **Uninstall:** Remove application + `rm -rf ~/.warp`

### JetBrains AI Assistant
- **Install:** IDE plugin (built into JetBrains IDEs 2024.2+)
- **Config:** `~/.config/JetBrains/`
- **AI features:** Inline completion, chat, commit message generation
- **Uninstall:** Disable AI Assistant plugin in IDE settings

### Tabnine
- **Install:** Editor extension (VS Code, JetBrains, etc.)
- **Config:** `~/.tabnine/`
- **Uninstall:** Remove extension + `rm -rf ~/.tabnine`

### Cody (Sourcegraph)
- **Install:** Editor extension, `cody` CLI
- **Config:** `~/.cody/`, `~/.config/cody/`
- **Uninstall:** Remove extension + `rm -rf ~/.cody ~/.config/cody`

### GitHub Copilot
- **Install:** VS Code extension, JetBrains plugin, CLI (`gh copilot`)
- **Config:** Extension settings
- **Uninstall:** Remove extension/plugin, `gh extension remove github/gh-copilot`

### ZCode (Z.ai / Zhipu)
- **Install:** Desktop app (macOS `.dmg`, Windows installer; Linux beta)
- **App:** `/Applications/ZCode.app` (macOS)
- **Config:** `~/.zcode/` (verified on a real install; contains `coding-plan-cache.json`, `bots-model-cache.v2.json`)
- **Note:** Official GLM-5.2 coding harness; steerable from WeChat/Feishu/Telegram. Subscription SaaS
- **Uninstall:** Remove app + `rm -rf ~/.zcode`

## 3. Model Runtimes

### Ollama
- **Install:** apt/yum/brew, also `curl -fsSL https://ollama.com/install.sh | sh`
- **Binary:** `/usr/local/bin/ollama` or `/usr/bin/ollama`
- **Service:** `ollama.service` (systemd), `com.ollama.ollama` (launchd)
- **Models:** `/usr/share/ollama/.ollama/models/`, `~/.ollama/models/`
- **Config:** `/etc/systemd/system/ollama.service`, env vars `OLLAMA_HOST`, `OLLAMA_MODELS`
- **Uninstall:**
  - Linux: `sudo systemctl stop ollama && sudo systemctl disable ollama && sudo rm /usr/local/bin/ollama && sudo rm -rf /usr/share/ollama && sudo rm /etc/systemd/system/ollama.service && systemctl daemon-reload`
  - macOS: Remove from Applications + `rm -rf ~/.ollama`
  - Windows: via Add/Remove Programs

### LM Studio
- **Install:** Desktop installer
- **Binary:** `/opt/LM Studio/`, `~/Applications/LM Studio.app/`, `%LOCALAPPDATA%/LM Studio/`
- **Config:** `~/.cache/lm-studio/`
- **Models:** `~/.cache/lm-studio/models/`
- **Uninstall:** Remove application + `rm -rf ~/.cache/lm-studio`

### GPT4All
- **Install:** Desktop installer, `pip install gpt4all`
- **Config:** `~/.gpt4all/` or `%LOCALAPPDATA%/nomic.ai/GPT4All/`
- **Models:** Within config dir
- **Uninstall:** Remove application + config dir

### text-generation-webui (oobabooga)
- **Install:** Git clone + conda
- **Location:** `~/text-generation-webui/` (or wherever cloned)
- **Models:** `~/text-generation-webui/models/`
- **LoRAs:** `~/text-generation-webui/loras/`
- **Uninstall:** `rm -rf ~/text-generation-webui` + conda env removal

### llama.cpp
- **Install:** Build from source, brew, apt
- **Binary:** `llama-cli`, `llama-server` in PATH or build dir
- **Models:** wherever the user points `--model`
- **Uninstall:** Remove binary, any model files (need discovery)

### KoboldCpp
- **Install:** Binary download or build
- **Models:** wherever configured
- **Uninstall:** Remove binary

### LocalAI
- **Install:** Docker, binary download
- **Config:** usually `~/localai/` or mounted dir
- **Uninstall:** Remove docker image/container or binary

## 4. ComfyUI & Image Models

### ComfyUI
- **Install:** Git clone + Python venv
- **Location:** `~/ComfyUI/` or similar
- **Models structure:**
  - `models/checkpoints/` — Stable Diffusion, SDXL, Flux model weights
  - `models/loras/` — LoRA adapters
  - `models/controlnet/` — ControlNet models
  - `models/vae/` — VAE files
  - `models/clip/` — CLIP models
  - `models/unet/` — UNet models
  - `models/embeddings/` — Textual inversion embeddings
  - `models/upscale_models/` — ESRGAN, etc.
  - `custom_nodes/` — Community plugins (may contain keys/configs)
- **Config:** `extra_model_paths.yaml`
- **Uninstall:** `rm -rf ~/ComfyUI` (or wherever installed)

### Automatic1111 / Stable Diffusion WebUI
- **Install:** Git clone + Python venv
- **Location:** `~/stable-diffusion-webui/`
- **Models:** `models/Stable-diffusion/`, `models/Lora/`, `models/ControlNet/`
- **Uninstall:** `rm -rf ~/stable-diffusion-webui`

### Fooocus
- **Install:** Git clone + Python
- **Location:** `~/Fooocus/`
- **Uninstall:** `rm -rf ~/Fooocus`

### InvokeAI
- **Install:** pip, installer
- **Config:** `~/invokeai/`
- **Uninstall:** Remove install dir

## 5. SDKs & Frameworks

### PyTorch
- **Install:** `pip install torch torchvision torchaudio`
- **Size:** 2-4 GB (CUDA version even larger)
- **Cache:** `~/.cache/torch/`
- **Uninstall:** `pip uninstall torch torchvision torchaudio && rm -rf ~/.cache/torch`

### TensorFlow
- **Install:** `pip install tensorflow`
- **Cache:** `~/.cache/tensorflow/`
- **Uninstall:** `pip uninstall tensorflow`

### HuggingFace Transformers
- **Install:** `pip install transformers`
- **Cache:** `~/.cache/huggingface/` (hub models, token)
- **Config:** `~/.huggingface/` (token file)
- **Uninstall:** `pip uninstall transformers datasets tokenizers && rm -rf ~/.cache/huggingface ~/.huggingface`

### LangChain
- **Install:** `pip install langchain` or `npm i langchain`
- **Uninstall:** `pip uninstall langchain` or `npm uninstall langchain`

### LlamaIndex
- **Install:** `pip install llama-index`
- **Uninstall:** `pip uninstall llama-index`

### OpenAI SDK
- **Install:** `pip install openai` or `npm i openai`
- **Uninstall:** `pip uninstall openai` or `npm uninstall openai`

### Anthropic SDK
- **Install:** `pip install anthropic` or `npm i @anthropic-ai/sdk`
- **Uninstall:** `pip uninstall anthropic` or `npm uninstall @anthropic-ai/sdk`

### Playwright (AI-adjacent)
- **Install:** `npm i playwright` or `pip install playwright`
- **Browsers:** `~/.cache/ms-playwright/` (1+ GB)
- **Note:** Often installed only for AI agent browser automation
- **Uninstall:** `npm uninstall playwright` or `pip uninstall playwright && rm -rf ~/.cache/ms-playwright`

### Selenium (AI-adjacent)
- **Install:** `pip install selenium`
- **Note:** When present alongside AI tools, likely AI-driven testing
- **Uninstall:** `pip uninstall selenium`

## 6. Model Caches

### HuggingFace Hub Cache
- **Path:** `~/.cache/huggingface/hub/`
- **Size:** Can be 10s to 100s of GB
- **Detection:** Directory exists and contains `models--*` subdirs

### PyTorch Hub Cache
- **Path:** `~/.cache/torch/hub/`
- **Detection:** Directory exists

### GPT4All Models
- **Path:** `~/.gpt4all/` or `%LOCALAPPDATA%/nomic.ai/GPT4All/`

### GGUF Files (general)
- **Detection:** `find / -name "*.gguf"` (user home + common locations)

### Safetensors Files
- **Detection:** `find / -name "*.safetensors"` (user home + common locations)

## 7. Agent Config & Instructions

### AGENTS.md (pi)
- **Detection:** Recursive scan from home/projects
- **Contains:** Agent behavioral instructions, tool configs, project rules

### CLAUDE.md (Claude Code)
- **Detection:** Recursive scan from home/projects

### .cursorrules (Cursor)
- **Detection:** Recursive scan from home/projects

### .windsurfrules (Windsurf)
- **Detection:** Recursive scan from home/projects

### GEMINI.md (Gemini CLI)
- **Detection:** Recursive scan from home/projects

### copilot-instructions.md (GitHub Copilot)
- **Detection:** Recursive scan from home/projects

### CONVENTIONS.md (generic)
- **Detection:** Recursive scan from home/projects

### Soul/Character Files
- **Patterns:** `soul.md`, `system-prompt.md`, `character.json`, `personality.md`
- **Detection:** Heuristic scan in known AI tool dirs

## 8. Agent Memory & Sessions

### pi Sessions
- **Path:** `~/.pi/sessions/`
- **Contains:** Full conversation history, tool call results
- **Size:** Can accumulate significantly

### Claude Projects/Memory
- **Path:** `~/.claude/projects/`, `~/.claude/todos/`
- **Contains:** Project-scoped memory, conversation history

### Aider Chat History
- **Path:** `.aider.chat.history.md` (per-project)

### Continue Session Data
- **Path:** `~/.continue/sessions/`

### Vibe History
- **Path:** `~/.vibe/vibehistory/`

### PaperclipAI Context
- **Path:** `~/.paperclip/context.json`, `~/.paperclip/instances/`

## 9. MCP Configurations

### .mcp.json
- **Detection:** Recursive scan from home/projects
- **Contains:** MCP server configs — often has API keys, connection strings, DB URLs
- **Security:** HIGH — treat as credential-containing

### Claude MCP Settings
- **Path:** `~/.claude/mcp.json` or similar

### Cursor MCP
- **Path:** `~/.cursor/mcp.json` or similar

## 10. Docker AI Images

### Common AI Images
- `ollama/ollama`
- `vllm/vllm`
- `ghcr.io/berriai/litellm`
- `localai/localai`
- `comfyui/comfyui`
- `goauthentik/authentik` (if AI-adjacent)
- Any image with tags: `llm`, `ai`, `gpt`, `stable-diffusion`, `comfyui`

### Detection
- `docker images` + filter for AI-related names/tags
- Also check stopped containers and volumes

## 11. Config & Data Directories (Consolidated)

### All Known Config Dirs
```
~/.pi/                    pi agent
~/.claude/                Claude Code
~/.claude.json            Claude main config
~/.aider.conf.yml         Aider config
~/.codex/                 Codex CLI
~/.cursor/                Cursor
~/.cursorrules            Cursor rules
~/.gemini/                Gemini CLI
~/.opencode/              OpenCode
~/.ollama/                Ollama user data
~/.continue/              Continue
~/.cline/                 Cline (VS Code)
~/.amazon-q/              Amazon Q
~/.augment/               Augment
~/.windsurf/              Windsurf
~/.windsurfrules          Windsurf rules
~/.paperclip/             PaperclipAI
~/.vibe/                  Mistral Vibe
~/.gpt4all/               GPT4All
~/.huggingface/           HuggingFace token
~/.lm-studio/             LM Studio
%APPDATA%/Claude/         Claude desktop (Windows)
%APPDATA%/Cursor/         Cursor IDE (Windows)
%LOCALAPPDATA%/LM Studio/ LM Studio (Windows)
%LOCALAPPDATA%/Programs/  Various AI tools (Windows)
```

## 11b. OpenClaw Ecosystem

### OpenClaw
- **Install:** `npm i -g openclaw`
- **Binary:** `openclaw` in PATH
- **Config:** `~/.openclaw/`
- **Note:** Open-source personal AI assistant, 72K+ GitHub stars, NVIDIA NemoClaw integration

### NanoBot
- **Install:** `pip install nanobot`
- **Binary:** `nanobot` in PATH
- **Config:** `~/.nanobot/`
- **Note:** Lightweight OpenClaw alternative, MCP agent framework

### ZeroClaw
- **Binary:** `zeroclaw` in PATH
- **Config:** `~/.zeroclaw/`
- **Note:** Ultra-lightweight AI agent runtime, Rust-based

### PicoClaw
- **Binary:** `picoclaw` in PATH
- **Config:** `~/.picoclaw/`
- **Note:** Tiny AI assistant in Go, designed for Raspberry Pi and embedded

### NanoClaw
- **Binary:** `nanoclaw` in PATH
- **Config:** `~/.nanoclaw/`
- **Note:** Lightweight OpenClaw with WhatsApp/Telegram/Slack integration

### IronClaw (NEAR AI)
- **Binary:** `ironclaw` in PATH
- **Config:** `~/.ironclaw/`
- **Note:** OpenClaw-inspired in Rust, focused on privacy and security by NEAR AI

### NullClaw
- **Binary:** `nullclaw` in PATH
- **Config:** `~/.nullclaw/`
- **Note:** Smallest fully autonomous AI assistant, written in Zig

### Moltis
- **Binary:** `moltis` in PATH
- **Config:** `~/.moltis/`
- **Note:** Secure persistent personal agent server in Rust, multi-provider LLMs

### Clawith
- **Binary:** `clawith` in PATH
- **Config:** `~/.clawith/`
- **Note:** OpenClaw for Teams — multi-agent collaboration platform

## 12. Stragglers

### After Agent Uninstall
- Config dirs that weren't removed
- Session/memory files
- Project-level instruction files (AGENTS.md, CLAUDE.md, etc.)

### After Model Runtime Uninstall
- Downloaded model files (often in separate location from runtime)
- Cache directories
- Service files (systemd, launchd)

### After Editor Uninstall
- Extension storage
- Workspace storage
- Rule files in projects

### General
- Stale PATH entries pointing to removed tools
- Environment variables for removed tools (OLLAMA_HOST, etc.)
- Shell profile modifications (eval commands, PATH additions)
- npm/pip global packages that were dependencies of removed tools
- Orphaned Docker volumes from removed containers
