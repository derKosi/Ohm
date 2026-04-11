# Ohm — Roadmap

**Current version**: v0.1.0  
**Status**: Released. 84+ signatures, cross-platform TUI, script generation, zero PII.  
**Repo**: https://github.com/derKosi/Ohm

---

## What's Done ✅

| Feature | Status |
|---------|--------|
| Core scanner (agents, editors, runtimes, SDKs, models, ComfyUI) | ✅ |
| 84+ detection signatures | ✅ |
| Interactive Bubble Tea TUI with cursor, selection, scroll | ✅ |
| Script generation (.sh / .ps1) with warnings | ✅ |
| JSON output (`--json`) | ✅ |
| Text output (`--no-tui`) | ✅ |
| Opt-in scans (`--path`, `--env`, `--shell`, `--deep`) | ✅ |
| MCP config detection | ✅ |
| Agent memory & session detection | ✅ |
| Agent instruction file detection (AGENTS.md, CLAUDE.md, etc.) | ✅ |
| Plugin/extension detection (pi skills, VS Code) | ✅ |
| Docker AI image detection | ✅ |
| Straggler detection | ✅ |
| State persistence (`~/.ohm/state.json`) | ✅ |
| History command | ✅ |
| Scanning animation | ✅ |
| Cross-compile (linux/darwin/windows amd64+arm64) | ✅ |
| GitHub repo + v0.1.0 release with binaries | ✅ |
| `go install github.com/derKosi/Ohm/cmd/ohm@latest` | ✅ |
| PII audit — zero personal data in repo/history | ✅ |

---

## What's Next

### Phase 1: Safety Net 🔴 Critical

> This tool generates `rm -rf` scripts. It needs tests before anyone trusts it.

#### 1.1 Unit Tests
**File**: `internal/scanner/scanner_test.go`, etc.

```
- Test that each agentDef detection works (mock config dir + command)
- Test dirSize on a temp dir with known contents
- Test that finding IDs are unique across all categories
- Test RiskLevel assignment (credential dirs → RiskDanger)
- Test model finding (Ollama `list` output parsing)
- Test straggler detection (compare two ScanResults)
```

**Effort**: ~2-3 hours  
**Blockers**: None — scanner functions are pure-ish, just need temp dirs

#### 1.2 Script Generation Tests
**File**: `internal/generator/script_test.go`

```
- Test bash script format (header, comments, commands, warnings)
- Test PowerShell script format
- Test that credential-containing findings get ⚠️ warnings
- Test empty selection → error
- Test line endings (LF for .sh, CRLF for .ps1)
```

**Effort**: ~1 hour  
**Blockers**: None

#### 1.3 Model Tests
**File**: `internal/model/finding_test.go`

```
- Test FormatBytes (KB, MB, GB, TB edge cases)
- Test ByCategory grouping
- Test SelectedCount / SelectedSize
- Test ScanResult.Count()
```

**Effort**: ~30 min  
**Blockers**: None

---

### Phase 2: Custom Signatures 🟡 High Value

> Let users add their own tools without forking Ohm.

#### 2.1 YAML Signature Format
**New dep**: `gopkg.in/yaml.v3`

```yaml
# ~/.ohm/signatures/my-tool.yaml
name: My Custom Tool
category: agents  # agents|editors|runtimes|sdks|other
detect:
  paths:
    - ~/.mytool/
  files:
    - mytool.conf
  commands:
    - mytool --version
  npm_packages:
    - mytool-cli
  pip_packages:
    - mytool
uninstall:
  linux: "npm uninstall -g mytool-cli && rm -rf ~/.mytool"
  macos: "npm uninstall -g mytool-cli && rm -rf ~/.mytool"
  windows: "npm uninstall -g mytool-cli; Remove-Item ~/.mytool -Recurse -Force"
cleanup:
  - "~/.mytool/"
  - "~/.config/mytool/"
risk: caution  # safe|caution|danger
```

**Files to create/modify**:
- `internal/signature/signature.go` — types
- `internal/signature/loader.go` — YAML loader, dir watcher
- `internal/scanner/scanner.go` — merge custom sigs into scan pipeline

**Effort**: ~3 hours  
**Blockers**: Need to add `gopkg.in/yaml.v3` dependency

#### 2.2 `ohm signatures` Command
```
ohm signatures              # List all known signatures (built-in + custom)
ohm signatures --builtin    # List only built-in signatures
ohm signatures --custom     # List only custom signatures
```

**Effort**: ~1 hour  
**Blockers**: 2.1

#### 2.3 `ohm signatures --add` Interactive
Walk the user through creating a signature:

```
$ ohm signatures --add
  Tool name: My Custom Tool
  Category [agents/editors/runtimes/sdks/other]: agents
  Config directory path [~/.mytool/]: ~/.mytool/
  Detection command [mytool --version]: 
  Risk level [safe/caution/danger]: caution
  Uninstall command (linux): npm uninstall -g mytool-cli && rm -rf ~/.mytool
  
  ✅ Written to ~/.ohm/signatures/my-custom-tool.yaml
```

**Effort**: ~2 hours  
**Blockers**: 2.1

---

### Phase 3: CI/CD & Release Automation 🟡

#### 3.1 GitHub Actions CI
**File**: `.github/workflows/ci.yml`

```
on: [push, pull_request]
jobs:
  build:
    strategy:
      matrix:
        os: [ubuntu-latest, macos-latest, windows-latest]
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with: go-version: '1.24'
      - run: go build ./...
      - run: go vet ./...
      - run: go test ./...
```

**Effort**: ~30 min  
**Blockers**: None

#### 3.2 GoReleaser
**File**: `.goreleaser.yml`

```
- On tag push v*: build all 6 platforms
- Archive: .tar.gz (unix), .zip (windows)
- Checksums file
- Changelog from git log
- Upload to GitHub Release
```

**Effort**: ~1 hour  
**Blockers**: None

---

### Phase 4: Polish 🟢

#### 4.1 Self-Update
```
ohm update    # Check GitHub releases, download new binary, replace self
```

Implementation: HTTP GET to `api.github.com/repos/derKosi/Ohm/releases/latest`, compare versions, download, replace binary. This is the **only** network call Ohm would ever make — must be opt-in and documented.

**Effort**: ~2 hours  
**Blockers**: Adds `net/http` — needs careful security review

#### 4.2 Signature Auto-Update
Download new signatures from GitHub without updating the binary itself.

```
~/.ohm/signatures/remote/    # auto-synced from repo
```

**Effort**: ~2 hours  
**Blockers**: Same as 4.1

#### 4.3 Expand/Collapse in TUI
Currently all categories are always expanded. Add:
- `enter` to expand/collapse a category
- Show collapsed count: `🤖 Agents & Harnesses (5 found) ▸`

**Effort**: ~2 hours  
**Blockers**: None

#### 4.4 Detail View
Press `d` on a finding to show full details:
- All config paths with sizes
- Sub-items (models in a runtime, files in ComfyUI)
- Risk explanation
- Uninstall command preview

**Effort**: ~2 hours  
**Blockers**: None

#### 4.5 Search/Filter
Press `/` to filter findings by name or path.

**Effort**: ~1.5 hours  
**Blockers**: None

---

### Phase 5: Known Bugs & Edge Cases

| Issue | Impact | Fix |
|-------|--------|-----|
| LM Studio size returns 0 on Windows | Missing data | Windows service/AppData path detection |
| `amazon-q` config dir is `~/.aws/amazonq/` not `~/.amazon-q/` | False negative | Update path in agents.go |
| Aider `hasPipPackage("aider-chat")` slow if pip not installed | Slowness | Check pip exists first |
| No dedup between categories (Ollama in runtimes + model caches) | Confusing counts | Cross-category dedup by path |
| `--deep` scan has no max-depth limit | Can scan entire home dir | Default max-depth of 5 |
| Windows `--shell` checks Unix profiles (.bashrc) | False findings | Skip Unix profiles on Windows |

---

## Suggested Priority Order

```
1. Unit tests (1.1-1.3)              — Safety first
2. GitHub Actions CI (3.1)            — Automated quality gate
3. Custom signatures YAML (2.1-2.2)  — Most requested feature
4. Known bugs (Phase 5)              — Fix what's broken
5. GoReleaser (3.2)                  — Automated releases
6. TUI polish (4.3-4.5)             — Nice to have
7. Self-update (4.1-4.2)            — Last, adds network code
```

---

## Adding New Signatures

When a new AI tool launches (it happens weekly), add it in two places:

### 1. Code: `internal/scanner/<category>.go`

Add a new entry to the slice in the relevant scanner function:

```go
{
    id:         "new-tool",
    name:       "New Tool (Vendor)",
    configDirs: []string{"~/.newtool"},
    commands:   []string{"newtool"},
    npmPackages: []string{"newtool-cli"},
    risk:       model.RiskCaution,
    uninstallCmds: map[string]string{
        "linux":   "npm uninstall -g newtool-cli && rm -rf ~/.newtool",
        "macos":   "npm uninstall -g newtool-cli && rm -rf ~/.newtool",
        "windows": "npm uninstall -g newtool-cli; Remove-Item ~/.newtool -Recurse -Force",
    },
},
```

### 2. Docs: `docs/SIGNATURES.md`

Add the tool to the right category section with install, config, and uninstall details.

### Checklist for new signatures

- [ ] Real tool (verify with web search — no made-up entries)
- [ ] Correct `id` (kebab-case, unique across all categories)
- [ ] All 3 platform uninstall commands (linux, macos, windows)
- [ ] Appropriate `risk` level (has API keys → `RiskDanger`, configs → `RiskCaution`, otherwise `RiskSafe`)
- [ ] Added to `SIGNATURES.md`
- [ ] Run `./ohm scan` on a machine with the tool to verify detection
