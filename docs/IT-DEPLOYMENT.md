# IT Deployment Guide

Rolling out Ohm across a fleet of developer machines. This guide covers
silent installation, agent-based deployment, compliance reporting, and the
licensing situation for IT departments.

---

## Licensing for IT (read this first)

**You do not need a commercial license to run Ohm on your company's machines.**

Under AGPL-3.0, running `ohm scan` on any number of employee laptops, workstations,
or servers — personal, company, or enterprise — is free. This includes:

- Scanning every developer machine for AI software
- Collecting `--json` output into internal compliance reports
- Building internal tools that shell out to `ohm scan`

A commercial license is only required if you **embed Ohm's Go packages into a
proprietary product** or **offer Ohm as part of a paid SaaS**. For the 99% case
(fleet scanning for hygiene/security), AGPL is sufficient.

See [COMMERCIAL-LICENSE.md](../COMMERCIAL-LICENSE.md) for the full breakdown and
pricing tiers (Solo / Team / Business / Enterprise) if you do need embedding rights.

---

## Silent Installation

Both installers are **non-interactive** — they produce output to stdout but never
prompt for input. This makes them safe to call from deployment agents, scheduled
tasks, and management platforms.

### Windows (PowerShell, no admin required)

```powershell
# Silent install — downloads latest release, installs to %USERPROFILE%\.ohm\bin,
# adds to user PATH. No prompts, no elevation.
irm https://github.com/derKosi/Ohm/releases/latest/download/install.ps1 | iex
```

For **system-wide** installation (all users, requires admin), download the
`Windows_x86_64.zip` from the [latest release](https://github.com/derKosi/Ohm/releases/latest),
extract `ohm.exe` to `C:\Program Files\Ohm\`, and add that path to the system PATH
via Group Policy or Intune.

### Linux / macOS

```bash
# Silent install — installs to /usr/local/bin (may need sudo) or ~/.local/bin
curl -fsSL https://github.com/derKosi/Ohm/releases/latest/download/install.sh | sh
```

For system-wide installation without sudo, the script automatically falls back to
`~/.local/bin` and prints the PATH export the user needs.

---

## Agent-Based Deployment

### Microsoft Intune (Windows)

1. **Download** the `Windows_x86_64.zip` from the latest release.
2. **Package** as an `.intunewin` wrapper (Intune Win32 App Packaging Tool):
   - Install command: `powershell.exe -ExecutionPolicy Bypass -File install.ps1`
   - Or extract `ohm.exe` directly to `C:\Program Files\Ohm\` via your wrapping script.
3. **Assign** to the developer device group.
4. Ohm runs per-user (no system service needed), so install in user context.

### SCCM / Configuration Manager (Windows)

```powershell
# Distribution point package:
#   Source: the Windows zip + install.ps1
#   Program: powershell.exe -ExecutionPolicy Bypass -File install.ps1
#   Run with user rights (not admin) — installs to user profile
```

Deploy to the "All Developer Workstations" collection. Ohm needs no agent, no
service, no background process — it's a single binary the user runs on demand.

### Ansible (Linux / macOS)

```yaml
# playbook.yml — installs Ohm on all dev machines
- hosts: developer_workstations
  tasks:
    - name: Install Ohm
      shell: curl -fsSL https://github.com/derKosi/Ohm/releases/latest/download/install.sh | sh
      args:
        creates: /usr/local/bin/ohm

    - name: Verify installation
      command: ohm version
      register: ohm_check
      changed_when: false

    - name: Show installed version
      debug:
        var: ohm_check.stdout
```

### Jamf Pro (macOS)

1. Upload the `Darwin_x86_64.tar.gz` (or `arm64` for Apple Silicon) as a package.
2. **Script** (post-install): extract `ohm` to `/usr/local/bin/`, `chmod +x`.
3. Scope to the developer Mac fleet.
4. No daemon needed — Ohm is on-demand only.

---

## Compliance Reporting

Ohm's `--json` output is designed for ingestion into compliance/security tooling.
Run it across the fleet, collect the JSON, aggregate into a dashboard or SIEM.

### Per-machine scan

```bash
ohm scan --json --all > ohm-report.json
```

The JSON contains version, timestamp, platform, hostname, finding count, and the
full findings array — but **no file contents, no credential values, no PII**.
API-key locations are flagged with `risk_level: "danger"` but the key values
themselves are never included.

### Fleet-wide collection (example: Ansible)

```yaml
- hosts: developer_workstations
  tasks:
    - name: Scan for AI software
      command: ohm scan --json --all
      register: scan_result
      changed_when: false

    - name: Collect reports centrally
      copy:
        content: "{{ scan_result.stdout }}"
        dest: "/srv/compliance/ohm-{{ inventory_hostname }}-{{ ansible_date_time.date }}.json"
      delegate_to: compliance-server
```

### What the JSON tells IT

| Field | What it reveals |
|-------|-----------------|
| `count` | Total AI software items on the machine |
| `findings[].category` | Agents, runtimes, editors, SDKs, models, configs, etc. |
| `findings[].size_bytes` | Disk consumption per item (for cleanup prioritization) |
| `findings[].risk_level` | `"danger"` = likely contains API keys/credentials |
| `findings[].path` | Where the item lives on disk |

This lets IT answer questions like:
- "Which machines have AI agents with exposed API keys?"
- "How much disk space is model-cache bloat consuming fleet-wide?"
- "Which devs are running unapproved local LLM runtimes?"

---

## Privacy Guarantee for IT

Ohm is **100% offline**. There is:

- No telemetry, no phone-home, no crash reporting
- No HTTP client in the codebase (verifiable: `grep -r "net/http" cmd/ internal/`)
- No data leaving the machine, ever

IT can verify this independently. The binary makes zero network connections —
scanning is entirely local filesystem inspection. This makes Ohm safe to deploy
in air-gapped environments, regulated industries, and zero-trust networks.

---

## Uninstall

### Windows
```powershell
Remove-Item -Recurse -Force "$env:USERPROFILE\.ohm"
# Remove from PATH via System > Environment Variables, or:
[Environment]::SetEnvironmentVariable('PATH', ($env:PATH -replace [regex]::Escape("$env:USERPROFILE\.ohm\bin;"), ''), 'User')
```

### Linux / macOS
```bash
sudo rm -f /usr/local/bin/ohm
# or: rm -f ~/.local/bin/ohm
rm -rf ~/.ohm/state.json  # optional: remove scan history
```

---

## Questions?

- **Licensing**: [COMMERCIAL-LICENSE.md](../COMMERCIAL-LICENSE.md) or derKosi@kosiit.de
- **Security/Privacy**: [README.md](../README.md) § "Privacy-First Architecture"
- **Issues**: https://github.com/derKosi/Ohm/issues
