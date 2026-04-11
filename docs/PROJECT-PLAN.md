# Ohm — Full Project Plan

## Overview

Ohm is a cross-platform (Linux, macOS, Windows) TUI tool written in Go that scans a system for AI-related software, presents findings in an interactive interface, and generates cleanup scripts without executing them.

**Core principle: Read-only, offline, privacy-first.**

## Architecture

```
ohm/
├── cmd/
│   └── ohm/
│       └── main.go              # CLI entry point (built-in)
├── internal/
│   ├── scanner/
│   │   ├── scanner.go           # Scanner orchestrator
│   │   ├── agents.go            # Agent/harness detection
│   │   ├── editors.go           # AI editor detection
│   │   ├── runtimes.go          # Model runtime detection
│   │   ├── comfyui.go           # ComfyUI + image model detection
│   │   ├── sdks.go              # SDK/framework detection
│   │   ├── models.go            # Model cache detection
│   │   ├── instructions.go      # Agent instruction file detection
│   │   ├── memory.go            # Session/memory file detection
│   │   ├── mcp.go               # MCP config detection
│   │   ├── configs.go           # Config dir detection
│   │   ├── docker.go            # Docker image detection
│   │   ├── stragglers.go        # Straggler/orphan detection
│   │   ├── plugins.go           # Plugin/extension detection
│   │   ├── vscode.go            # VS Code extension detection
│   │   └── optin.go              # PATH/ENV/Shell (opt-in)
│   ├── model/
│   │   ├── finding.go           # Scan result type
│   │   └── state.go             # Persistent state between runs
│   ├── generator/
│   │   └── script.go            # Script generation (sh/ps1)
│   ├── platform/
│   │   └── detect.go            # OS detection + path helpers
│   └── tui/                     # (embedded in cmd/ohm/main.go)
├── signatures/
│   └── *.yaml                   # Built-in tool signatures
├── docs/
│   ├── SIGNATURES.md            # Signature catalog (this file)
│   └── PROJECT-PLAN.md          # This file
├── test-fixtures/
│   ├── windows11-ai-inventory.md
│   └── linux-vm-ai-inventory.md
├── go.mod
├── go.sum
├── README.md
├── .gitignore
└── Makefile
```

## Dependencies

| Package | Purpose | Version |
|---------|---------|---------|
| `github.com/charmbracelet/bubbletea` | TUI framework | v1.x |
| `github.com/charmbracelet/lipgloss` | Terminal styling | v1.x |

**No network/HTTP libraries. No telemetry. No oauth. No cobra (built-in CLI). No yaml.v3 yet (will be needed for custom signatures).** If a PR adds `net/http`, it gets rejected.

## Scan Pipeline

```
1. OS Detection (Linux/macOS/Windows)
2. Home Dir Resolution
3. Package Manager Discovery (npm, pip, brew, apt, winget, choco, go, cargo, uv)
4. Parallel Category Scanners:
   ├── Agents & Harnesses (binary in PATH + config dirs + package lists)
   ├── AI Editors (installed apps + AppImage + brew cask + Windows programs)
   ├── Model Runtimes (binaries + services + model dirs)
   ├── ComfyUI (install dir + model subdirs by type)
   ├── SDKs & Frameworks (pip list + npm list + cargo + go list)
   ├── Model Caches (known cache paths + .gguf/.safetensors discovery)
   ├── Agent Instructions (recursive find for known filenames)
   ├── Agent Memory & Sessions (known session dirs)
   ├── MCP Configs (recursive find .mcp.json)
   ├── Config Dirs (check all known ~/.tool dirs)
   ├── Docker (docker images + filter)
   └── Stragglers (compare current findings with previous state)
5. Opt-in Scanners (if flags set):
   ├── PATH entries
   ├── Environment variables
   └── Shell profile modifications
6. Size Calculation (du for dirs, stat for files)
7. Categorization & Deduplication
8. Presentation (Bubble Tea TUI)
```

## Data Model

```go
type Finding struct {
    ID          string   // unique identifier
    Category    Category // agents, editors, runtimes, etc.
    Name        string   // human-readable name
    Version     string   // detected version (if available)
    InstallMethod string // npm, pip, binary, brew, etc.
    Path        string   // primary location
    SizeBytes   int64    // disk usage
    ConfigPaths []string // associated config/data paths
    Contains    []string // sub-findings (e.g., models in a runtime)
    RiskLevel   Risk     // safe, caution, danger (credential-containing)
    UninstallCmds map[string]string // platform -> uninstall command
    Selected    bool     // user selected for removal
}

type Risk int
const (
    RiskSafe    Risk = iota // standard removal
    RiskCaution             // may affect other software
    RiskDanger              // contains credentials/keys
)

type Category int
const (
    CatAgents Category = iota
    CatEditors
    CatRuntimes
    CatComfyUI
    CatSDKs
    CatModelCaches
    CatInstructions
    CatMemory
    CatMCP
    CatConfigs
    CatDocker
    CatStragglers
)
```

## Script Generation

The generator produces platform-appropriate scripts:

- **Linux/macOS:** `.sh` (bash)
- **Windows:** `.ps1` (PowerShell)

Each script:
1. Header with date, machine info, total size to free
2. Warning comment at the top
3. Grouped by category with comments
4. Each removal command preceded by `echo "Removing <name>..."`
5. Credential-containing removals flagged with `# ⚠️ WARNING: This may contain API keys`
6. Final summary line

## Persistent State

File: `~/.ohm/state.json`

```json
{
  "version": 1,
  "last_scan": "2026-04-11T15:30:00Z",
  "last_removal": "2026-04-10T10:00:00Z",
  "findings": [
    {
      "id": "ollama",
      "name": "Ollama",
      "removed": false,
      "first_seen": "2026-04-01T00:00:00Z"
    }
  ],
  "removed": [
    {
      "id": "aider",
      "name": "Aider",
      "removed_at": "2026-04-10T10:05:00Z",
      "stragglers_remaining": ["/home/user/.aider.conf.yml"]
    }
  ]
}
```

State enables:
- `ohm stragglers` — compare current scan against previously-removed items
- `ohm history` — show removal log
- Size savings over time

## Build & Release

### Makefile Targets
```makefile
build          # Build for current platform
build-all      # Cross-compile linux/darwin/windows amd64+arm64
test           # Run tests
lint           # golangci-lint
install        # go install
clean          # Remove build artifacts
```

### GoReleaser (`.goreleaser.yml`)
- Cross-compile: linux/amd64, linux/arm64, darwin/amd64, darwin/arm64, windows/amd64, windows/arm64
- Archive formats: .tar.gz (unix), .zip (windows)
- Checksums
- Changelog from git log

### GitHub Actions
- On tag push `v*`: build all platforms, create GitHub Release, upload binaries
- On PR: build + test + lint

## Milestones

### Phase 1: Foundation ✅
- [x] Project planning
- [x] Signature catalog (84+ entries)
- [x] Go module + dependencies (Bubble Tea, Lip Gloss)
- [x] Data model types (Finding, Category, Risk, ScanResult)
- [x] OS detection + path resolution
- [x] Test fixtures (Windows + Linux)

### Phase 2: Core Scanners ✅
- [x] Package manager discovery (npm, pip, brew)
- [x] Agent/harness scanner (47 agents)
- [x] Config dir scanner
- [x] Model runtime scanner (8 runtimes)
- [x] Editor scanner (7 editors)
- [x] SDK/framework scanner (9 SDKs)
- [x] Size calculation
- [x] Basic CLI output (--no-tui)

### Phase 3: Extended Scanners ✅
- [x] ComfyUI scanner
- [x] Model cache scanner (.gguf, safetensors, HuggingFace)
- [x] Instruction file scanner (AGENTS.md, CLAUDE.md, etc.)
- [x] Memory/session scanner
- [x] MCP config scanner
- [x] Docker scanner
- [x] Plugin/extension scanner (pi skills, VS Code)
- [x] Straggler detection
- [x] Opt-in: PATH scanner
- [x] Opt-in: ENV scanner
- [x] Opt-in: Shell profile scanner

### Phase 4: TUI ✅
- [x] Bubble Tea main app with viewport scrolling
- [x] Category grouping
- [x] Checkbox selection (space, a toggle all)
- [x] Size display per item and per category
- [x] Risk indicators (🔑 danger, ⚠️ caution)
- [x] Scanning animation

### Phase 5: Script Generation ✅
- [x] Bash script generator
- [x] PowerShell script generator
- [x] Comments + warnings + credential flags
- [x] Size summary

### Phase 6: State & History ✅
- [x] JSON state persistence
- [x] Straggler detection via state comparison
- [x] History command
- [x] `ohm stragglers` command

### Phase 7: Release ✅
- [x] Makefile (build, build-all, test, lint, clean)
- [x] Cross-compile 6 platforms
- [x] GitHub repo + v0.1.0 release
- [x] `go install` working
- [x] PII audit (clean)
- [x] README with real scan outputs

### Phase 8: Remaining Work
- [ ] Unit tests (scanner, generator, model)
- [ ] GitHub Actions CI
- [ ] Custom YAML signatures (see [docs/ROADMAP.md](ROADMAP.md))
- [ ] `ohm signatures` command
- [ ] GoReleaser
- [ ] Self-update command
- [ ] TUI expand/collapse, detail view, search/filter
- [ ] Known bug fixes (see [docs/ROADMAP.md](ROADMAP.md))

See **[`docs/ROADMAP.md`](ROADMAP.md)** for detailed specs and priority order.

## Testing Strategy

Since we **never execute uninstall commands**, testing is:

1. **Scanner tests:** Mock filesystem with test fixtures, verify detection
2. **Script generation tests:** Verify output format, platform correctness
3. **TUI tests:** Bubble Tea test harness (tea.TestProgram)
4. **Integration:** Run `ohm scan` on this machine, compare output against `linux-vm-ai-inventory.md`
5. **Cross-platform:** CI runs on Linux + macOS + Windows runners

## Safety Guardrails

- No `os.Remove`, `os.RemoveAll`, or `exec.Command("rm")` anywhere in the codebase
- Linter rule: reject any code that modifies the filesystem (only `os.Stat`, `os.ReadDir`, `filepath.Walk`)
- Script generator writes to a file the user must manually execute
- `go vet` + custom linter to catch destructive patterns
