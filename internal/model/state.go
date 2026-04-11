package model

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

// State represents persistent state between Ohm runs.
type State struct {
	Version   int        `json:"version"`
	LastScan  time.Time  `json:"last_scan"`
	Removed   []Removed  `json:"removed,omitempty"`
	Findings  []Finding  `json:"findings,omitempty"`
}

// Removed tracks a previously removed item for straggler detection.
type Removed struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	RemovedAt time.Time `json:"removed_at"`
}

// StatePath returns the path to the state file.
func StatePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".ohm", "state.json"), nil
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
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
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
