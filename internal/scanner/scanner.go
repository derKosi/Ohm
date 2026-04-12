// SPDX-FileCopyrightText: 2026 Mathias Kosinski
// SPDX-License-Identifier: AGPL-3.0-or-later

package scanner

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"github.com/derKosi/Ohm/internal/model"
	"github.com/derKosi/Ohm/internal/platform"
)

// Options controls what the scanner looks for.
type Options struct {
	ScanPATH     bool
	ScanENV      bool
	ScanShell    bool
	ScanDeep     bool
	CustomSigDir string
}

// Scanner discovers AI software on the system.
type Scanner struct {
	os       platform.OS
	home     string
	opts     Options
	findings []model.Finding
	mu       sync.Mutex
}

// New creates a new Scanner.
func New(opts Options) *Scanner {
	return &Scanner{
		os:   platform.Detect(),
		home: platform.HomeDir(),
		opts: opts,
	}
}

// Scan runs all scanners and returns findings.
func (s *Scanner) Scan() *model.ScanResult {
	s.findings = nil

	hostname, _ := os.Hostname()

	// Run category scanners
	s.scanAgents()
	s.scanEditors()
	s.scanRuntimes()
	s.scanComfyUI()
	s.scanSDKs()
	s.scanModelCaches()
	s.scanInstructions()
	s.scanMemory()
	s.scanMCP()
	s.scanPlugins()
	s.scanVSCodeExtensions()
	s.scanConfigDirs()
	s.scanDocker()
	s.scanStragglers()

	// Opt-in scanners
	if s.opts.ScanPATH {
		s.scanPATH()
	}
	if s.opts.ScanENV {
		s.scanENV()
	}
	if s.opts.ScanShell {
		s.scanShellProfiles()
	}

	// Dedup across categories: if a config dir finding's path is already covered
	// by a more specific finding (e.g. Ollama in runtimes + config dirs), remove the duplicate.
	s.dedupFindings()

	result := &model.ScanResult{
		Findings:  s.findings,
		ScannedAt: time.Now(),
		Platform:  runtime.GOOS + "/" + runtime.GOARCH,
		Hostname:  hostname,
	}

	// Platform-specific warnings
	s.detectPlatformWarnings(result)

	return result
}

// addFinding safely adds a finding.
func (s *Scanner) addFinding(f model.Finding) {
	s.mu.Lock()
	s.findings = append(s.findings, f)
	s.mu.Unlock()
}

// finding creates a basic finding with common defaults.
func (s *Scanner) finding(id, name string, category model.Category) model.Finding {
	return model.Finding{
		ID:       id,
		Name:     name,
		Category: category,
	}
}

// dirExists checks if a directory exists.
func (s *Scanner) dirExists(path string) bool {
	return platform.IsDir(path)
}

// fileExists checks if a file exists.
func (s *Scanner) fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// dirSize gets directory size.
func (s *Scanner) dirSize(path string) int64 {
	return platform.DirSize(path)
}

// fileSize gets file size.
func (s *Scanner) fileSize(path string) int64 {
	info, err := os.Stat(path)
	if err != nil {
		return 0
	}
	return info.Size()
}

// expandPath resolves ~ to home directory.
func (s *Scanner) expandPath(path string) string {
	if len(path) > 0 && path[0] == '~' {
		return filepath.Join(s.home, path[1:])
	}
	return path
}

// hasCommand checks if a command exists.
func (s *Scanner) hasCommand(name string) bool {
	return platform.HasCommand(name)
}

// runCommand runs a command and returns output.
func (s *Scanner) runCommand(name string, args ...string) (string, error) {
	return platform.RunCommand(name, args...)
}

// dedupFindings removes config-dir findings that duplicate a more specific category finding.
// For example, Ollama may appear in both Model Runtimes and Config & Data Dirs.
// We keep the more specific one (runtimes) and drop the config-dir duplicate.
func (s *Scanner) dedupFindings() {
	// Build a set of config paths from non-config-dir findings
	pathOwners := make(map[string]string) // path -> finding ID
	for _, f := range s.findings {
		if f.Category == model.CatConfigs {
			continue
		}
		for _, p := range f.ConfigPaths {
			pathOwners[p] = f.ID
		}
	}

	// Now filter config-dir findings: remove those whose paths are all covered by another category
	var deduped []model.Finding
	for _, f := range s.findings {
		if f.Category != model.CatConfigs {
			deduped = append(deduped, f)
			continue
		}
		// Check if all config paths are owned by another finding
		allDupe := len(f.ConfigPaths) > 0
		for _, p := range f.ConfigPaths {
			if _, owned := pathOwners[p]; !owned {
				allDupe = false
				break
			}
		}
		if !allDupe {
			deduped = append(deduped, f)
		}
	}
	s.findings = deduped
}

// detectPlatformWarnings adds informational warnings about the scan environment.
func (s *Scanner) detectPlatformWarnings(result *model.ScanResult) {
	// Windows: detect WSL and suggest scanning inside it too
	if s.os == platform.OSWindows {
		if _, err := exec.LookPath("wsl.exe"); err == nil {
			result.Warnings = append(result.Warnings,
				"WSL detected. Ohm cannot scan inside WSL from Windows. Run \"ohm scan\" inside WSL for a complete picture.")
		}
	}
}
