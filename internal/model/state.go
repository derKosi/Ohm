// SPDX-FileCopyrightText: 2026 Mathias Kosinski
// SPDX-License-Identifier: AGPL-3.0-or-later

package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"time"
)

// State represents persistent state between Ohm runs.
// State persists scan results between runs.
type State struct {
	Version  int       `json:"version"`
	LastScan time.Time `json:"last_scan"`
	// ScanOpts records the scanner options of the last scan, so `ohm generate`
	// can re-run an equivalent scan instead of trusting persisted finding
	// bytes (names/uninstall commands are emitted into an executable script).
	ScanOpts ScanOpts   `json:"scan_opts,omitempty"`
	Removed  []Removed  `json:"removed,omitempty"`
	Findings []Finding  `json:"findings,omitempty"`
}

// ScanOpts mirrors scanner.Options for state persistence. Kept as a separate
// type in model so the state schema does not depend on the scanner package.
type ScanOpts struct {
	Path  bool `json:"path,omitempty"`
	Env   bool `json:"env,omitempty"`
	Shell bool `json:"shell,omitempty"`
	Deep  bool `json:"deep,omitempty"`
}

// Removed tracks a previously removed item for straggler detection.
type Removed struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	RemovedAt time.Time `json:"removed_at"`
}

// StatePath returns the path to the state file.
func StatePath() (string, error) {
	// The state file feeds `ohm generate`'s emitted script bytes, so which
	// file counts as "the user's own state" is a security decision. $HOME /
	// %USERPROFILE% are per-process environment values (direnv .envrc, repo
	// wrappers, sudo -E): resolving the state location through them lets
	// anyone controlling the environment of a single run point `ohm generate`
	// at an attacker-placed state file. Anchor the location to the account's
	// home directory from the OS user database instead. Note: this fails
	// closed for generate when the current uid has no user-database entry
	// (some minimal containers) — acceptable for a trust anchor; the scan
	// save sites already ignore LoadState errors.
	u, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("resolving current user: %w", err)
	}
	if u.HomeDir == "" {
		return "", errors.New("current user has no home directory")
	}
	return filepath.Join(u.HomeDir, ".ohm", "state.json"), nil
}

// LoadState loads state from disk.
func LoadState() (*State, error) {
	path, err := StatePath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &State{Version: 1}, nil
		}
		return nil, err
	}

	var state State
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, err
	}
	return &state, nil
}

// Save persists state to disk.
func (s *State) Save() error {
	path, err := StatePath()
	if err != nil {
		return err
	}

	dir := filepath.Dir(path)
	// State can carry verbatim credential lines captured by `ohm scan --shell`
	// (findings[].sub_items). Keep both the directory and the file private to
	// the owning user; umask can only tighten these, never loosen them.
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}

	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(path, data, 0600); err != nil {
		return err
	}
	// os.WriteFile applies perm only at creation; tighten files that an older
	// Ohm version left world-readable at 0644.
	return os.Chmod(path, 0600)
}

// MarkRemoved records a finding as removed.
func (s *State) MarkRemoved(finding Finding) {
	s.Removed = append(s.Removed, Removed{
		ID:        finding.ID,
		Name:      finding.Name,
		RemovedAt: time.Now(),
	})
}

// IsRemoved checks if a finding was previously removed.
func (s *State) IsRemoved(id string) bool {
	for _, r := range s.Removed {
		if r.ID == id {
			return true
		}
	}
	return false
}
