// SPDX-FileCopyrightText: 2026 Mathias Kosinski
// SPDX-License-Identifier: AGPL-3.0-or-later

package scanner

import (
	"os"
	"path/filepath"

	"github.com/derKosi/Ohm/internal/model"
)

// scanMemory detects agent memory, session, and conversation files.
func (s *Scanner) scanMemory() {
	type memDef struct {
		id       string
		name     string
		paths    []string
		risk     model.Risk
	}

	memories := []memDef{
		{
			id:   "pi-sessions",
			name: "pi Sessions",
			paths: []string{
				filepath.Join(s.home, ".pi", "sessions"),
			},
			risk: model.RiskCaution,
		},
		{
			id:   "claude-memory",
			name: "Claude Code Memory",
			paths: []string{
				filepath.Join(s.home, ".claude"),
			},
			risk: model.RiskDanger,
		},
		{
			id:   "vibe-history",
			name: "Mistral Vibe History",
			paths: []string{
				filepath.Join(s.home, ".vibe", "vibehistory"),
			},
			risk: model.RiskSafe,
		},
		{
			id:   "paperclip-context",
			name: "PaperclipAI Context",
			paths: []string{
				filepath.Join(s.home, ".paperclip", "context.json"),
				filepath.Join(s.home, ".paperclip", "instances"),
			},
			risk: model.RiskDanger,
		},
		{
			id:   "gemini-history",
			name: "Gemini CLI History",
			paths: []string{
				filepath.Join(s.home, ".gemini"),
			},
			risk: model.RiskCaution,
		},
		{
			id:   "aider-history",
			name: "Aider Chat History",
			paths: []string{
				filepath.Join(s.home, ".aider.chat.history.md"),
			},
			risk: model.RiskCaution,
		},
	}

	for _, mem := range memories {
		var foundPaths []string
		var totalSize int64

		for _, p := range mem.paths {
			info, err := os.Stat(p)
			if err != nil {
				continue
			}
			foundPaths = append(foundPaths, p)
			if info.IsDir() {
				totalSize += s.dirSize(p)
			} else {
				totalSize += info.Size()
			}
		}

		if len(foundPaths) > 0 {
			s.addFinding(model.Finding{
				ID:          mem.id,
				Category:    model.CatMemory,
				Name:        mem.name,
				Path:        joinPaths(foundPaths),
				SizeBytes:   totalSize,
				ConfigPaths: foundPaths,
				RiskLevel:   mem.risk,
				UninstallCmds: map[string]string{
					"linux":   "# Review before removing: " + joinPaths(foundPaths),
					"macos":   "# Review before removing: " + joinPaths(foundPaths),
					"windows": "# Review before removing: " + joinPaths(foundPaths),
				},
			})
		}
	}
}
