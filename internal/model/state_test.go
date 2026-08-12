// SPDX-FileCopyrightText: 2026 Mathias Kosinski
// SPDX-License-Identifier: AGPL-3.0-or-later

package model

import (
	"testing"
)

func TestMarkRemovedAndIsRemoved(t *testing.T) {
	s := &State{Version: 1}

	// Nothing removed initially
	if s.IsRemoved("ollama") {
		t.Error("IsRemoved(ollama) = true on empty state, want false")
	}

	// Mark a finding as removed
	s.MarkRemoved(Finding{ID: "ollama", Name: "Ollama"})

	// Now it should be removed
	if !s.IsRemoved("ollama") {
		t.Error("IsRemoved(ollama) = false after MarkRemoved, want true")
	}

	// Other IDs still not removed
	if s.IsRemoved("claude-code") {
		t.Error("IsRemoved(claude-code) = true, want false")
	}

	// Removed list should have exactly 1 entry
	if len(s.Removed) != 1 {
		t.Fatalf("len(Removed) = %d, want 1", len(s.Removed))
	}
	if s.Removed[0].ID != "ollama" {
		t.Errorf("Removed[0].ID = %q, want %q", s.Removed[0].ID, "ollama")
	}
	if s.Removed[0].Name != "Ollama" {
		t.Errorf("Removed[0].Name = %q, want %q", s.Removed[0].Name, "Ollama")
	}
}

func TestIsRemovedMultipleEntries(t *testing.T) {
	s := &State{Version: 1}
	s.MarkRemoved(Finding{ID: "a", Name: "A"})
	s.MarkRemoved(Finding{ID: "b", Name: "B"})
	s.MarkRemoved(Finding{ID: "c", Name: "C"})

	for _, id := range []string{"a", "b", "c"} {
		if !s.IsRemoved(id) {
			t.Errorf("IsRemoved(%q) = false, want true", id)
		}
	}

	if s.IsRemoved("d") {
		t.Error("IsRemoved(\"d\") = true, want false")
	}
}

func TestMarkRemovedPreservesExistingEntries(t *testing.T) {
	s := &State{Version: 1, Removed: []Removed{{ID: "old", Name: "Old"}}}

	s.MarkRemoved(Finding{ID: "new", Name: "New"})

	if len(s.Removed) != 2 {
		t.Fatalf("len(Removed) = %d, want 2", len(s.Removed))
	}
	if s.Removed[0].ID != "old" {
		t.Errorf("Removed[0].ID = %q, want %q (existing entry should be preserved)", s.Removed[0].ID, "old")
	}
	if s.Removed[1].ID != "new" {
		t.Errorf("Removed[1].ID = %q, want %q (new entry should be appended)", s.Removed[1].ID, "new")
	}
}
