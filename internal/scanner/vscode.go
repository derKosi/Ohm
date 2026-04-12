// SPDX-FileCopyrightText: 2026 Mathias Kosinski
// SPDX-License-Identifier: AGPL-3.0-or-later

package scanner

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/derKosi/Ohm/internal/model"
	"github.com/derKosi/Ohm/internal/platform"
)

// Known VS Code AI extension publisher.name pairs.
var vscodeAIExtensions = map[string]string{
	"GitHub.copilot":                    "GitHub Copilot",
	"GitHub.copilot-chat":               "GitHub Copilot Chat",
	"saoudrizwan.claude-dev":            "Cline",
	"Continue.continue":                 "Continue",
	"Codeium.codeium":                   "Codeium",
	"Codeium.windsurf":                  "Windsurf (Codeium)",
	"sourcegraph.cody-ai":               "Cody (Sourcegraph)",
	"Tabnine.tabnine-vscode":            "Tabnine",
	"amazonwebservices.amazon-q-vscode": "Amazon Q",
	"seanvscode.codegpt":                "CodeGPT",
	"FittenTech.fitten-code":            "Fitten Code",
	"visualstudioexptteam.vscodeintellicode": "IntelliCode",
	"JunoLab.juno-vscode":               "Juno AI",
	"tongyi.aliyun-copilot":             "Tongyi Copilot",
}

// scanVSCodeExtensions scans for AI-related VS Code extensions.
func (s *Scanner) scanVSCodeExtensions() {
	// Check common VS Code extensions directories
	extensionsDirs := []string{
		filepath.Join(s.home, ".vscode", "extensions"),
	}

	// Windows additional paths
	if s.os == platform.OSWindows {
		extensionsDirs = append(extensionsDirs,
			filepath.Join(os.Getenv("USERPROFILE"), ".vscode", "extensions"),
		)
	}

	var found []struct {
		id    string
		name  string
		path  string
		size  int64
	}

	for _, extDir := range extensionsDirs {
		entries, err := os.ReadDir(extDir)
		if err != nil {
			continue
		}
		for _, e := range entries {
			if !e.IsDir() {
				continue
			}
			extPath := filepath.Join(extDir, e.Name())
			// VS Code extensions are named like "publisher.name-VERSION"
			dirName := e.Name()
			// Extract publisher.name part (before the dash-version)
			extKey := dirName
			if idx := strings.Index(dirName, "-"); idx > 0 {
				extKey = dirName[:idx]
			}

			displayName, ok := vscodeAIExtensions[extKey]
			if !ok {
				continue
			}

			size := s.dirSize(extPath)
			found = append(found, struct {
				id    string
				name  string
				path  string
				size  int64
			}{
				id:   extKey,
				name: displayName,
				path: extPath,
				size: size,
			})
		}
	}

	if len(found) == 0 {
		return
	}

	var subItems []string
	var totalSize int64
	for _, f := range found {
		subItems = append(subItems, f.name)
		totalSize += f.size
	}

	s.addFinding(model.Finding{
		ID:            "vscode-ai-extensions",
		Category:      model.CatPlugins,
		Name:          "VS Code AI Extensions",
		Path:          filepath.Join(s.home, ".vscode", "extensions"),
		SizeBytes:     totalSize,
		SubItems:      subItems,
		RiskLevel:     model.RiskSafe,
		UninstallCmds: map[string]string{
			"linux":   "code --uninstall-extension <publisher.name>",
			"macos":   "code --uninstall-extension <publisher.name>",
			"windows": "code --uninstall-extension <publisher.name>",
		},
	})
}
