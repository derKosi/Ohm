package scanner

import (
	"github.com/derKosi/Ohm/internal/model"
)

// scanEditors detects AI editors and IDEs.
func (s *Scanner) scanEditors() {
	editors := []struct {
		id            string
		name          string
		configDirs    []string
		commands      []string
		risk          model.Risk
		uninstallCmds map[string]string
	}{
		{
			id:   "cursor",
			name: "Cursor IDE",
			configDirs: []string{
				"~/.cursor",
				"~/.cursorrules",
			},
			commands: []string{"cursor"},
			risk:     model.RiskCaution,
			uninstallCmds: map[string]string{
				"linux":   "# Cursor: remove AppImage/binary + rm -rf ~/.cursor ~/.cursorrules",
				"macos":   "rm -rf /Applications/Cursor.app ~/.cursor ~/.cursorrules",
				"windows": "# Cursor: uninstall via Add/Remove Programs",
			},
		},
		{
			id:   "windsurf",
			name: "Windsurf (Codeium)",
			configDirs: []string{
				"~/.windsurf",
				"~/.windsurfrules",
			},
			commands: []string{"windsurf"},
			risk:     model.RiskCaution,
			uninstallCmds: map[string]string{
				"linux":   "# Windsurf: remove AppImage/binary + rm -rf ~/.windsurf ~/.windsurfrules",
				"macos":   "rm -rf /Applications/Windsurf.app ~/.windsurf ~/.windsurfrules",
				"windows": "# Windsurf: uninstall via Add/Remove Programs",
			},
		},
		{
			id:   "zed",
			name: "Zed Editor",
			configDirs: []string{
				"~/.config/zed",
			},
			commands: []string{"zed"},
			risk:     model.RiskSafe,
			uninstallCmds: map[string]string{
				"linux":   "# Zed: remove binary + rm -rf ~/.config/zed",
				"macos":   "rm -rf /Applications/Zed.app ~/.config/zed",
				"windows": "# Zed: remove binary",
			},
		},
		{
			id:   "warp",
			name: "Warp Terminal (AI)",
			configDirs: []string{
				"~/.warp",
			},
			commands: []string{"warp"},
			risk:     model.RiskCaution,
			uninstallCmds: map[string]string{
				"linux":   "# Warp: remove AppImage/binary + rm -rf ~/.warp",
				"macos":   "rm -rf /Applications/Warp.app ~/.warp",
				"windows": "# Warp: uninstall via Add/Remove Programs",
			},
		},
		{
			id:   "jetbrains-ai",
			name: "JetBrains AI Assistant",
			configDirs: []string{
				"~/.config/JetBrains",
			},
			risk:     model.RiskCaution,
			uninstallCmds: map[string]string{
				"linux":   "# JetBrains AI: disable AI Assistant plugin in IDE settings",
				"macos":   "# JetBrains AI: disable AI Assistant plugin in IDE settings",
				"windows": "# JetBrains AI: disable AI Assistant plugin in IDE settings",
			},
		},
		{
			id:   "tabnine",
			name: "Tabnine",
			configDirs: []string{
				"~/.tabnine",
			},
			commands: []string{"tabnine"},
			risk:     model.RiskSafe,
			uninstallCmds: map[string]string{
				"linux":   "# Tabnine: remove from editor extensions + rm -rf ~/.tabnine",
				"macos":   "# Tabnine: remove from editor extensions + rm -rf ~/.tabnine",
				"windows": "# Tabnine: remove from editor extensions + Remove-Item ~/.tabnine -Recurse -Force",
			},
		},
		{
			id:   "cody",
			name: "Cody (Sourcegraph)",
			configDirs: []string{
				"~/.config/cody",
				"~/.cody",
			},
			commands: []string{"cody"},
			risk:     model.RiskCaution,
			uninstallCmds: map[string]string{
				"linux":   "# Cody: remove from editor extensions + rm -rf ~/.cody ~/.config/cody",
				"macos":   "# Cody: remove from editor extensions + rm -rf ~/.cody ~/.config/cody",
				"windows": "# Cody: remove from editor extensions + Remove-Item ~/.cody -Recurse -Force",
			},
		},
	}

	for _, editor := range editors {
		found := false
		var totalSize int64
		var foundPaths []string

		for _, cmd := range editor.commands {
			if s.hasCommand(cmd) {
				found = true
			}
		}

		for _, dir := range editor.configDirs {
			expanded := s.expandPath(dir)
			if s.dirExists(expanded) {
				found = true
				size := s.dirSize(expanded)
				totalSize += size
				foundPaths = append(foundPaths, expanded)
			}
		}

		if found {
			s.addFinding(model.Finding{
				ID:            editor.id,
				Category:      model.CatEditors,
				Name:          editor.name,
				Path:          joinPaths(foundPaths),
				SizeBytes:     totalSize,
				ConfigPaths:   foundPaths,
				RiskLevel:     editor.risk,
				UninstallCmds: editor.uninstallCmds,
			})
		}
	}
}
