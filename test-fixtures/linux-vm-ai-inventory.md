# Ohm — Linux VM AI Inventory (Test Fixture)

Source: Ohm scan on Ubuntu development VM
Date: 2026-04-11
Purpose: Test case for Ohm scan target detection (Linux platform)

---

## Summary

Total AI software footprint: **~41.3 GB** across 19 items

## Interactive TUI Scan

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

## Details

### Installed AI Tools

| Tool | How Installed | Version | Size |
|------|--------------|---------|------|
| pi | npm global | @mariozechner/pi-coding-agent@0.66.1 | 92.3 MB (config) |
| Gemini CLI | npm global | @google/gemini-cli@0.37.0 | 5.3 MB (config) |
| Claude Code | npm global | v1.1.x | 278 B (config only) |
| PaperclipAI | npm global | 2026.403.0 | 13.2 GB |
| Mistral Vibe | unknown | unknown | 523 KB |
| Ollama | systemd service + apt | latest | 261.6 MB |

### Related Tools (AI-Adjacent)

| Tool | Size | Note |
|------|------|------|
| Playwright browsers | 1.2 GB | Installed for AI agent browser testing |

### Ollama Models

| Model | Size |
|-------|------|
| nomic-embed-text:latest | 274 MB |

### Pi Agent Plugins/Skills

| Plugin | Source |
|--------|--------|
| caveman | Local skill |
| paperclip | Local skill |
| paperclip-create-agent | Local skill |
| paperclip-create-plugin | Local skill |
| para-memory-files | Local skill |

### Agent Instruction Files

| File | Locations |
|------|-----------|
| AGENTS.md | ~/.nvm/, ~/AGENTS.md, ~/projects/backend/, ~/tools/cli/ |
| CLAUDE.md | ~/.nvm/, ~/projects/backend/, ~/projects/frontend/ |
| .cursorrules | ~/projects/webapp/ |

### System Services (AI-related)

| Service | Status |
|---------|--------|
| ollama.service | active (running), enabled on boot |

### Stragglers

| Straggler | Detail |
|-----------|--------|
| AI entry in PATH | ~/.pi/agent/bin |
