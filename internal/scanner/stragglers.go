package scanner

import (
	"github.com/derKosi/Ohm/internal/model"
)

// scanStragglers detects leftover files from previously removed tools.
func (s *Scanner) scanStragglers() {
	// Load previous state to compare
	// For now, check for common orphan patterns

	type stragglerDef struct {
		id       string
		name     string
		paths    []string
		risk     model.Risk
		implies  string // tool that was likely removed
	}

	stragglers := []stragglerDef{
		{
			id:   "orphan-aider-config",
			name: "Aider config (Aider not installed)",
			paths: []string{"~/.aider.conf.yml"},
			implies: "aider",
			risk: model.RiskSafe,
		},
		{
			id:   "orphan-codex-config",
			name: "Codex CLI config (Codex not installed)",
			paths: []string{"~/.codex"},
			implies: "codex-cli",
			risk: model.RiskSafe,
		},
		{
			id:   "orphan-cursor-config",
			name: "Cursor config (Cursor not installed)",
			paths: []string{"~/.cursor", "~/.cursorrules"},
			implies: "cursor",
			risk: model.RiskSafe,
		},
		{
			id:   "orphan-opencode-config",
			name: "OpenCode config (OpenCode not installed)",
			paths: []string{"~/.opencode"},
			implies: "opencode",
			risk: model.RiskSafe,
		},
	}

	for _, sg := range stragglers {
		// Check if the implied tool is installed
		if s.hasCommand(sg.implies) {
			continue // Tool still installed, not a straggler
		}

		// Also check if tool was found in agents scan
		found := false
		for _, f := range s.findings {
			if f.ID == sg.implies {
				found = true
				break
			}
		}
		if found {
			continue
		}

		// Check if straggler paths exist
		var foundPaths []string
		var totalSize int64
		for _, p := range sg.paths {
			expanded := s.expandPath(p)
			if s.dirExists(expanded) {
				foundPaths = append(foundPaths, expanded)
				totalSize += s.dirSize(expanded)
			} else if s.fileExists(expanded) {
				foundPaths = append(foundPaths, expanded)
				totalSize += s.fileSize(expanded)
			}
		}

		if len(foundPaths) > 0 {
			s.addFinding(model.Finding{
				ID:          sg.id,
				Category:    model.CatStragglers,
				Name:        sg.name,
				Path:        joinPaths(foundPaths),
				SizeBytes:   totalSize,
				ConfigPaths: foundPaths,
				RiskLevel:   sg.risk,
				UninstallCmds: map[string]string{
					"linux":   "rm -rf " + joinPaths(foundPaths),
					"macos":   "rm -rf " + joinPaths(foundPaths),
					"windows": "Remove-Item '" + joinPaths(foundPaths) + "' -Recurse -Force",
				},
			})
		}
	}
}
