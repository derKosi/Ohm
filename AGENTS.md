# AGENTS.md

This file guides AI agents (Claude Code, Hermes, OpenClaw, Copilot, etc.)
that want to install, use, or integrate Ohm. It follows the [AGENTS.md convention](https://agents.md/).

---

## What Ohm is

Ohm is a privacy-first CLI scanner that detects AI-related software (agents,
runtimes, editors, SDKs, model caches, configs, credentials, Docker images)
on a system. It produces a reviewable cleanup script. **Ohm never deletes
anything itself** — it is strictly read-only and 100% offline (no network
stack, no telemetry).

- Binary: `ohm` (Go, single static binary, no runtime dependencies)
- Platforms: Linux, macOS, Windows (amd64 + arm64)
- License: AGPL-3.0 + commercial dual-license
- Repo: https://github.com/derKosi/Ohm

---

## Installation (for agents installing Ohm on a machine)

Ohm ships prebuilt binaries. **Go is NOT required** to install or run Ohm.

### Linux / macOS

```sh
curl -fsSL https://github.com/derKosi/Ohm/releases/latest/download/install.sh | sh
```

Installs to `/usr/local/bin/ohm` (falls back to `~/.local/bin/ohm` if
`/usr/local/bin` is not writable). Non-interactive, safe for scripted use.

### Windows (PowerShell)

```powershell
irm https://github.com/derKosi/Ohm/releases/latest/download/install.ps1 | iex
```

Installs to `%USERPROFILE%\.ohm\bin\ohm.exe`, adds to user PATH. No admin
rights required. Non-interactive.

### Direct binary download (if curl/irm are unavailable)

Download the matching archive from
https://github.com/derKosi/Ohm/releases/latest, extract the binary, place it
on the PATH.

### Verify installation

```sh
ohm version
```

---

## Usage (for agents running Ohm)

### Basic scan (interactive TUI)

```sh
ohm scan
```

### Non-interactive scan (for scripts/agents/CI)

```sh
ohm scan --no-tui        # human-readable text output
ohm scan --json          # machine-readable JSON
ohm scan --all           # include PATH, ENV, shell profiles, deep crawl
```

The `--json` output is the recommended interface for agents. Structure:

```json
{
  "version": "0.1.x",
  "scanned_at": "2026-08-12T...",
  "count": 19,
  "findings": [
    {
      "id": "claude-code",
      "name": "Claude Code",
      "category": "Agents & Harnesses",
      "path": "/home/user/.claude",
      "size_bytes": 407200000,
      "risk_level": "danger",
      "subitems": []
    }
  ]
}
```

**Important:** the JSON contains **paths and sizes, never file contents or
credential values.** `risk_level: "danger"` means the location likely contains
API keys, but the values are never exposed. Safe to log, aggregate, and transmit.

### Generate cleanup script

```sh
ohm scan                 # select items in TUI (or via state.json)
ohm generate             # writes ohm-cleanup-<date>.sh or .ps1
```

The generated script is **reviewable** and must be **executed by the user
or agent explicitly** — Ohm does not run it. The script contains `rm -rf` /
`Remove-Item` commands that the user/agent should read before executing.

### Other commands

```sh
ohm stragglers           # scan only for leftover files from removed tools
ohm history              # show what was removed in previous runs
ohm version              # show version
```

---

## License (critical for agents)

Ohm is **dual-licensed**:

### AGPL-3.0 (open-source) — what agents and users can do freely

- **Run** `ohm scan` on any machine (personal, company, enterprise) — **no license needed**
- **Read** the `--json` output and use it in internal tools, reports, dashboards
- **Call** Ohm as a subprocess from another program (shell out, parse JSON)
- **Modify** Ohm for your own use (must share source if you distribute)

### Commercial license — what requires purchasing

You need a commercial license if you:
- **Import** Ohm's Go packages (`github.com/derKosi/Ohm/internal/...`) into your own code
- **Embed** Ohm's scanner in a proprietary product
- **Offer** Ohm (modified or not) as part of a paid SaaS

For details and pricing: [COMMERCIAL-LICENSE.md](COMMERCIAL-LICENSE.md) or derKosi@kosiit.de.

**Rule of thumb:** If you're shelling out to `ohm scan` and reading its output,
you're fine under AGPL. If you're importing the Go code as a library, you need
a commercial license.

---

## Guardrails for agents

### DO

- Run `ohm scan --json` to inventory AI software on a machine
- Use the JSON output to inform cleanup or compliance decisions
- Present the generated cleanup script to the user for review before executing
- Aggregate JSON across machines for fleet reporting

### DON'T

- **Do not** execute the generated cleanup script without human review.
  The script contains destructive `rm -rf` commands. Ohm generates it; the
  user/agent decides what to run.
- **Do not** transmit the full JSON over unencrypted channels if policy
  requires it — although it contains no credential values, it does reveal
  what AI tools and paths exist on a machine.
- **Do not** import Ohm's Go packages into proprietary code without a
  commercial license. Shell out to the binary instead.
- **Do not** expect Ohm to make network connections. If you observe network
  traffic from the `ohm` binary, that is a bug — report it.

---

## For agents contributing to Ohm

If you are an AI agent modifying Ohm's source code:

- **Language:** Go 1.24+. Match existing code style (`gofmt`, no external deps without discussion).
- **Build:** `go build ./cmd/ohm`
- **Test:** `go test ./...` (currently no test files — adding tests is welcome)
- **Vet:** `go vet ./...` must pass
- **Commits:** sign off with `-s` (DCO), conventional commit messages (`feat:`, `fix:`, `docs:`)
- **Detection entries:** when adding AI tool signatures, verify the tool is real (web search),
  add to the appropriate `internal/scanner/<category>.go`, and update `docs/SIGNATURES.md`.
- **Privacy invariant:** never add `net/http`, telemetry, or any network capability. The
  "100% offline" guarantee is the core product property and must hold.
- **Read-only invariant:** never add `os.Remove`, `os.RemoveAll`, or `exec.Command("rm")`
  to the scanner. Ohm generates scripts; it does not execute destructive operations.
- See [CONTRIBUTING.md](CONTRIBUTING.md) for the full DCO and contribution process.

---

## Architecture (quick reference for agents navigating the codebase)

```
cmd/ohm/main.go              CLI entry point, TUI, command dispatch
internal/scanner/            Detection logic (one file per category)
  scanner.go                 Core Scanner struct, orchestration, dedup
  agents.go                  AI agents & harnesses (51 signatures)
  editors.go                 AI editors & IDEs
  runtimes.go                Model runtimes (Ollama, llama.cpp, KTransformers, ...)
  comfyui.go                 ComfyUI installations + model subdirs
  sdks.go                    SDK detection (pip/npm/brew/go queries)
  models.go                  Model caches (.gguf, .safetensors, HF cache)
  configs.go                 Config & data directories
  mcp.go                     MCP configuration files
  memory.go                  Agent memory & session files
  instructions.go            Agent instruction files (AGENTS.md, CLAUDE.md, ...)
  plugins.go                 Plugins & extensions (pi skills, ComfyUI nodes)
  vscode.go                  VS Code AI extensions
  docker.go                  Docker AI images
  stragglers.go              Orphaned files from removed tools
  deep.go                    --deep scan (home crawl for missed items)
  optin.go                   --path, --env, --shell opt-in scans
internal/platform/detect.go  Platform detection, DirSize, path helpers
internal/model/              Finding/State data structures, categories
internal/generator/script.go Cleanup script generation (.sh / .ps1)
docs/                        Human-readable documentation
.goreleaser.yml              Release pipeline (6 targets, archives, checksums)
.github/workflows/           CI (build/vet/test) + auto-release on v* tags
install.sh / install.ps1     Platform installers (uploaded as release assets)
```

---

## Contact

- **Licensing questions:** derKosi@kosiit.de
- **Bugs/features:** https://github.com/derKosi/Ohm/issues
- **Source:** https://github.com/derKosi/Ohm
