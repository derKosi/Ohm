// SPDX-FileCopyrightText: 2026 Mathias Kosinski
// SPDX-License-Identifier: AGPL-3.0-or-later

package generator

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/derKosi/Ohm/internal/model"
	"github.com/derKosi/Ohm/internal/platform"
)

// runInTempDir runs a function with the working directory set to a temp dir,
// restoring the original cwd afterwards. Generator writes to cwd, so we isolate.
func runInTempDir(t *testing.T, fn func()) {
	t.Helper()
	orig, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd: %v", err)
	}
	dir := t.TempDir()
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("Chdir: %v", err)
	}
	defer os.Chdir(orig)
	fn()
}

func TestGenerateBasicStructure(t *testing.T) {
	runInTempDir(t, func() {
		result := &model.ScanResult{
			Findings: []model.Finding{
				{
					ID:        "ollama",
					Name:      "Ollama",
					Category:  model.CatRuntimes,
					SizeBytes: 1000,
					Selected:  true,
					UninstallCmds: map[string]string{
						"linux":   "rm -rf ~/.ollama",
						"macos":   "rm -rf ~/.ollama",
						"windows": "Remove-Item ~/.ollama -Recurse -Force",
					},
				},
			},
			Platform: "linux/amd64",
			Hostname: "test-host",
		}

		path, err := Generate(result)
		if err != nil {
			t.Fatalf("Generate() error: %v", err)
		}

		// File should exist
		info, err := os.Stat(path)
		if err != nil {
			t.Fatalf("generated file does not exist: %v", err)
		}
		if info.Size() == 0 {
			t.Error("generated file is empty")
		}

		content, _ := os.ReadFile(path)

		// Should contain the header warning about review
		if !strings.Contains(string(content), "REVIEW") {
			t.Error("script missing REVIEW warning")
		}

		// Should mention the finding name
		if !strings.Contains(string(content), "Ollama") {
			t.Error("script missing finding name 'Ollama'")
		}

		// Should mention the platform/hostname from the result
		if !strings.Contains(string(content), "test-host") {
			t.Error("script missing hostname")
		}

		// Should contain the item count
		if !strings.Contains(string(content), "Items to remove: 1") || !strings.Contains(string(content), "Total items: 1") {
			t.Error("script missing item count")
		}
	})
}

func TestGenerateRiskWarnings(t *testing.T) {
	runInTempDir(t, func() {
		result := &model.ScanResult{
			Findings: []model.Finding{
				{
					ID:        "claude-code",
					Name:      "Claude Code",
					Category:  model.CatAgents,
					SizeBytes: 500,
					Selected:  true,
					RiskLevel: model.RiskDanger,
					UninstallCmds: map[string]string{
						"linux":   "rm -rf ~/.claude",
						"macos":   "rm -rf ~/.claude",
						"windows": "Remove-Item ~/.claude -Recurse -Force",
					},
				},
				{
					ID:        "cursor",
					Name:      "Cursor",
					Category:  model.CatEditors,
					SizeBytes: 200,
					Selected:  true,
					RiskLevel: model.RiskCaution,
					UninstallCmds: map[string]string{
						"linux":   "# Cursor: remove manually",
						"macos":   "# Cursor: remove manually",
						"windows": "# Cursor: remove manually",
					},
				},
			},
		}

		path, err := Generate(result)
		if err != nil {
			t.Fatalf("Generate() error: %v", err)
		}

		content, _ := os.ReadFile(path)
		str := string(content)

		// RiskDanger should produce a credential warning
		if !strings.Contains(str, "API keys or credentials") {
			t.Error("script missing credential warning for RiskDanger finding")
		}

		// RiskCaution should produce a caution line
		if !strings.Contains(str, "Caution") {
			t.Error("script missing Caution note for RiskCaution finding")
		}
	})
}

func TestGenerateOnlySelectedFindings(t *testing.T) {
	runInTempDir(t, func() {
		result := &model.ScanResult{
			Findings: []model.Finding{
				{
					ID:        "selected-tool",
					Name:      "Selected Tool",
					Category:  model.CatAgents,
					SizeBytes: 100,
					Selected:  true,
					UninstallCmds: map[string]string{
						"linux": "echo removing selected",
					},
				},
				{
					ID:        "unselected-tool",
					Name:      "Unselected Tool",
					Category:  model.CatAgents,
					SizeBytes: 200,
					Selected:  false,
					UninstallCmds: map[string]string{
						"linux": "echo should not appear",
					},
				},
			},
		}

		path, err := Generate(result)
		if err != nil {
			t.Fatalf("Generate() error: %v", err)
		}

		content, _ := os.ReadFile(path)
		str := string(content)

		// Selected finding should be in the script
		if !strings.Contains(str, "Selected Tool") {
			t.Error("script missing selected finding")
		}

		// Unselected finding should NOT be in the script
		if strings.Contains(str, "Unselected Tool") {
			t.Error("script contains unselected finding (should be excluded)")
		}
	})
}

func TestGenerateEmptySelection(t *testing.T) {
	runInTempDir(t, func() {
		result := &model.ScanResult{
			Findings: []model.Finding{
				{ID: "a", Name: "A", Category: model.CatAgents, Selected: false},
				{ID: "b", Name: "B", Category: model.CatAgents, Selected: false},
			},
		}

		path, err := Generate(result)
		if err != nil {
			t.Fatalf("Generate() error: %v", err)
		}

		content, _ := os.ReadFile(path)
		str := string(content)

		// Script should still be generated but report zero items
		if !strings.Contains(str, "Items to remove: 0") || !strings.Contains(str, "Total items: 0") {
			t.Error("script should report zero items for empty selection")
		}
	})
}

func TestGenerateFileExtension(t *testing.T) {
	runInTempDir(t, func() {
		result := &model.ScanResult{
			Findings: []model.Finding{
				{ID: "x", Name: "X", Category: model.CatAgents, Selected: true},
			},
		}

		path, err := Generate(result)
		if err != nil {
			t.Fatalf("Generate() error: %v", err)
		}

		wantExt := platform.ScriptExtension() // .sh on unix, .ps1 on windows
		if !strings.HasSuffix(filepath.Base(path), wantExt) {
			t.Errorf("generated file extension = %q, want suffix %q", filepath.Base(path), wantExt)
		}
	})
}

func TestGenerateCommentCommandsPreserved(t *testing.T) {
	// Some uninstall commands are comments (starting with #) — those should be
	// written verbatim without an "echo Removing..." prefix.
	runInTempDir(t, func() {
		result := &model.ScanResult{
			Findings: []model.Finding{
				{
					ID:        "manual-remove",
					Name:      "Manual Remove Tool",
					Category:  model.CatAgents,
					Selected:  true,
					UninstallCmds: map[string]string{
						"linux":   "# Manual Remove: remove app + config manually",
						"macos":   "# Manual Remove: remove app + config manually",
						"windows": "# Manual Remove: remove app + config manually",
					},
				},
			},
		}

		path, err := Generate(result)
		if err != nil {
			t.Fatalf("Generate() error: %v", err)
		}

		content, _ := os.ReadFile(path)
		str := string(content)

		if !strings.Contains(str, "# Manual Remove:") {
			t.Error("script should preserve comment-style uninstall commands verbatim")
		}
	})
}

func TestGeneratePlatformSpecificCommand(t *testing.T) {
	// On Linux, only the "linux" key should be used; on macOS "macos"; on Windows "windows".
	runInTempDir(t, func() {
		result := &model.ScanResult{
			Findings: []model.Finding{
				{
					ID:        "platform-tool",
					Name:      "Platform Tool",
					Category:  model.CatAgents,
					Selected:  true,
					UninstallCmds: map[string]string{
						"linux":   "rm -rf ~/.platform-tool-linux",
						"macos":   "rm -rf ~/.platform-tool-mac",
						"windows": "Remove-Item ~/.platform-tool-win -Recurse -Force",
					},
				},
			},
		}

		path, err := Generate(result)
		if err != nil {
			t.Fatalf("Generate() error: %v", err)
		}

		content, _ := os.ReadFile(path)
		str := string(content)

		// Determine which platform-specific command should appear based on the test runner's OS.
		var wantSubstring string
		switch runtime.GOOS {
		case "linux":
			wantSubstring = "platform-tool-linux"
		case "darwin":
			wantSubstring = "platform-tool-mac"
		case "windows":
			wantSubstring = "platform-tool-win"
		default:
			t.Skip("test running on unsupported platform")
		}

		if !strings.Contains(str, wantSubstring) {
			t.Errorf("script should contain platform-specific command for %s, looking for %q", runtime.GOOS, wantSubstring)
		}
	})
}
