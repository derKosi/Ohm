# Ohm — Salt Analysis & Value Proposition

**Date**: 2026-04-11  
**Status**: Build ✅ | 3.7K lines Go | Cross-platform (Win/Mac/Linux)  
**Binary**: 5MB single binary, no dependencies

---

## Salt Analysis

### 🧂 Core Strengths (high salt)
1. **Clear, unique niche** — "Resistance against AGI bloat" — AI software scanner/cleaner. No other tool does this.
2. **84+ signatures** — detects agents (Claude Code, Cursor, Copilot, Aider, OpenClaw, Devin, Replit, etc.), editors (Windsurf, Warp, JetBrains AI, Tabnine, Cody), runtimes (Ollama, LM Studio), VS Code extensions, models, MCP servers, Docker images, SDKs, and more.
3. **Privacy-first architecture** — 100% offline, no telemetry, no HTTP client. Explicitly documented as a core differentiator.
4. **Cross-platform** — Windows, macOS, Linux from a single Go binary.
5. **Rich TUI** — terminal UI with categories, sizes, uninstall script generation. Polished user experience.
6. **Never deletes itself** — generates reviewable cleanup scripts. Safe by design.
7. **OpenClaw ecosystem detection** — covers obscure/rare tools beyond mainstream.

### 🧂 Medium Salt
8. **State management** — `~/.ohm/state.json` tracks scan history.
9. **Categorized findings** — agents, runtimes, editors, models, configs, Docker, MCP, plugins, SDKs.
10. **Test fixtures** — sample outputs for regression testing.

### 🧂 Low Salt (needs work)
11. **No tests** — 0 test files across 6 packages. Critical gap for a tool that generates destructive scripts.
12. **No remote repo** — no GitHub remote configured.
13. **No CI/CD** — no automated builds or releases.
14. **No auto-update** — version management is manual.

---

## Value Proposition

**For**: AI developers, ML engineers, and power users who experiment with AI tools  
**Who**: Have 50-200+ GB of AI-related software, models, and cache files scattered across their system  
**Ohm**: Is a privacy-first, offline terminal tool that scans your system for all AI-related software and generates a reviewable cleanup script  
**That**: Finds hidden AI bloat (14GB model files, 80GB HuggingFace caches, forgotten agent configs) and helps you reclaim disk space safely  
**Unlike**: Manual cleanup (you don't know what's where), `du -sh` (no AI-specific knowledge), or general uninstallers (no AI signature database)

**One-liner**: *Find every AI tool on your system. Reclaim your disk. Stay private.*

---

## Saltiest Quick Wins (No Oversight Needed)

### 🔥 Tier 1: Immediate Value

| # | Task | Impact | Effort |
|---|------|--------|--------|
| 1 | ~~Remove committed binary~~ | Already gitignored ✅ | — |
| 2 | Add unit tests for scanner signatures (at least model finding, path detection) | Safety for destructive tool | 2hr |
| 3 | ~~Add `--json` output flag~~ | Already implemented ✅ | — |
| 4 | Add GitHub Actions CI (build + lint for Win/Mac/Linux) | Quality assurance | 30min |
| 5 | Add GitHub remote and push | Visibility | 5min |

### ⭐ Tier 2: Polish

| # | Task | Impact | Effort |
|---|------|--------|--------|
| 6 | Add `ohm update` self-update command | User experience | 1hr |
| 7 | Add signature auto-update from GitHub repo | Keep up with new AI tools | 2hr |
| 8 | Add space savings summary (total reclaimable GB) | Motivation | 30min |

---

## Metrics Snapshot
- **Language**: Go (single binary)
- **Lines**: ~3.7K
- **Signatures**: 84+ (47 agents, 7 editors, 8 runtimes, 9 SDKs, + ComfyUI, model caches, MCP, Docker, plugins, configs, memory, instructions)
- **Platforms**: Windows, macOS, Linux
- **Tests**: 0 (critical gap)
- **Binary**: ~3.5MB, zero dependencies
