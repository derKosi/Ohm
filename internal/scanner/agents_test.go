// SPDX-FileCopyrightText: 2026 Mathias Kosinski
// SPDX-License-Identifier: AGPL-3.0-or-later

package scanner

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/derKosi/Ohm/internal/model"
)

func TestScanAgentsDetectsGitHubCopilotCLI(t *testing.T) {
	tempDir := t.TempDir()
	name := "copilot"
	if runtime.GOOS == "windows" {
		name += ".exe"
	}
	copilot := filepath.Join(tempDir, name)
	if err := os.WriteFile(copilot, []byte("#!/bin/sh\n"), 0755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", tempDir)

	s := &Scanner{home: tempDir, sizeCache: make(map[string]int64)}
	s.scanAgents()

	for _, finding := range s.findings {
		if finding.ID != "github-copilot-cli" {
			continue
		}
		if finding.Name != "GitHub Copilot CLI" {
			t.Errorf("name = %q, want GitHub Copilot CLI", finding.Name)
		}
		if finding.RiskLevel != model.RiskDanger {
			t.Errorf("risk = %v, want RiskDanger", finding.RiskLevel)
		}
		if got := finding.UninstallCmds["linux"]; got != "npm uninstall -g @github/copilot" {
			t.Errorf("linux uninstall command = %q", got)
		}
		return
	}

	t.Fatal("GitHub Copilot CLI was not detected from PATH")
}
