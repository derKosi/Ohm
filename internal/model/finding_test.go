// SPDX-FileCopyrightText: 2026 Mathias Kosinski
// SPDX-License-Identifier: AGPL-3.0-or-later

package model

import (
	"testing"
)

func TestFormatBytes(t *testing.T) {
	tests := []struct {
		name string
		bytes int64
		want  string
	}{
		{"zero", 0, "0 B"},
		{"bytes under KB", 512, "512 B"},
		{"exactly 1 KB", 1024, "1.0 KB"},
		{"kilobytes", 1536, "1.5 KB"},
		{"exactly 1 MB", 1024 * 1024, "1.0 MB"},
		{"megabytes", 100 * 1024 * 1024, "100.0 MB"},
		{"exactly 1 GB", 1024 * 1024 * 1024, "1.0 GB"},
		// 13.4 GB ≈ 14_386_912_460 bytes (13.4 * 1024^3, rounded)
		{"13.4 GB (PaperclipAI-style)", 14386912460, "13.4 GB"},
		{"exactly 1 TB", 1024 * 1024 * 1024 * 1024, "1.0 TB"},
		{"several TB", 5 * 1024 * 1024 * 1024 * 1024, "5.0 TB"},
		{"negative (shouldn't happen but shouldn't panic)", -1, "-1 B"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatBytes(tt.bytes)
			if got != tt.want {
				t.Errorf("FormatBytes(%d) = %q, want %q", tt.bytes, got, tt.want)
			}
		})
	}
}

func TestFindingFormatSize(t *testing.T) {
	f := Finding{SizeBytes: 5 * 1024 * 1024}
	if got := f.FormatSize(); got != "5.0 MB" {
		t.Errorf("Finding.FormatSize() = %q, want %q", got, "5.0 MB")
	}
}

func TestCategoryString(t *testing.T) {
	tests := []struct {
		cat  Category
		want string
	}{
		{CatAgents, "Agents & Harnesses"},
		{CatRuntimes, "Model Runtimes"},
		{CatStragglers, "Stragglers"},
	}
	for _, tt := range tests {
		if got := tt.cat.String(); got != tt.want {
			t.Errorf("Category(%d).String() = %q, want %q", tt.cat, got, tt.want)
		}
	}
}

func TestCategoryStringOutOfRange(t *testing.T) {
	cat := Category(999)
	if got := cat.String(); got != "Unknown" {
		t.Errorf("out-of-range Category.String() = %q, want %q", got, "Unknown")
	}
}

func TestRiskIcon(t *testing.T) {
	tests := []struct {
		risk Risk
		want string
	}{
		{RiskSafe, "  "},
		{RiskCaution, "⚠️ "},
		{RiskDanger, "🔑"},
	}
	for _, tt := range tests {
		if got := tt.risk.Icon(); got != tt.want {
			t.Errorf("Risk(%d).Icon() = %q, want %q", tt.risk, got, tt.want)
		}
	}
}

func TestScanResultCount(t *testing.T) {
	sr := &ScanResult{
		Findings: []Finding{
			{ID: "a"},
			{ID: "b"},
			{ID: "c"},
		},
	}
	if got := sr.Count(); got != 3 {
		t.Errorf("Count() = %d, want 3", got)
	}
}

func TestScanResultCountEmpty(t *testing.T) {
	sr := &ScanResult{}
	if got := sr.Count(); got != 0 {
		t.Errorf("empty Count() = %d, want 0", got)
	}
}

func TestScanResultSelectedCount(t *testing.T) {
	sr := &ScanResult{
		Findings: []Finding{
			{ID: "a", Selected: false},
			{ID: "b", Selected: true},
			{ID: "c", Selected: true},
			{ID: "d", Selected: false},
		},
	}
	if got := sr.SelectedCount(); got != 2 {
		t.Errorf("SelectedCount() = %d, want 2", got)
	}
}

func TestScanResultSelectedSize(t *testing.T) {
	sr := &ScanResult{
		Findings: []Finding{
			{ID: "a", SizeBytes: 100, Selected: false},
			{ID: "b", SizeBytes: 200, Selected: true},
			{ID: "c", SizeBytes: 300, Selected: true},
			{ID: "d", SizeBytes: 400, Selected: false},
		},
	}
	// Only b (200) + c (300) are selected
	if got := sr.SelectedSize(); got != 500 {
		t.Errorf("SelectedSize() = %d, want 500", got)
	}
}

func TestScanResultTotalSize(t *testing.T) {
	sr := &ScanResult{
		Findings: []Finding{
			{ID: "a", SizeBytes: 100},
			{ID: "b", SizeBytes: 200},
			{ID: "c", SizeBytes: 300},
		},
	}
	if got := sr.TotalSize(); got != 600 {
		t.Errorf("TotalSize() = %d, want 600", got)
	}
}

func TestByCategory(t *testing.T) {
	sr := &ScanResult{
		Findings: []Finding{
			{ID: "ollama", Category: CatRuntimes, SizeBytes: 1000},
			{ID: "claude", Category: CatAgents, SizeBytes: 500},
			{ID: "cursor", Category: CatAgents, SizeBytes: 200},
			{ID: "lm-studio", Category: CatRuntimes, SizeBytes: 300},
		},
	}

	groups := sr.ByCategory()

	// Should produce exactly 2 groups (Agents and Runtimes)
	if len(groups) != 2 {
		t.Fatalf("ByCategory() returned %d groups, want 2", len(groups))
	}

	// Find the Agents group (order depends on first occurrence in Findings)
	var agentsGroup, runtimesGroup *CategoryGroup
	for i := range groups {
		switch groups[i].Category {
		case CatAgents:
			agentsGroup = &groups[i]
		case CatRuntimes:
			runtimesGroup = &groups[i]
		}
	}

	if agentsGroup == nil || runtimesGroup == nil {
		t.Fatalf("missing expected category groups")
	}

	// Agents group should have 2 findings (claude at idx 1, cursor at idx 2)
	if agentsGroup.Count() != 2 {
		t.Errorf("agents group count = %d, want 2", agentsGroup.Count())
	}
	if agentsGroup.TotalSize(sr.Findings) != 700 {
		t.Errorf("agents group total size = %d, want 700", agentsGroup.TotalSize(sr.Findings))
	}

	// Runtimes group should have 2 findings (ollama at idx 0, lm-studio at idx 3)
	if runtimesGroup.Count() != 2 {
		t.Errorf("runtimes group count = %d, want 2", runtimesGroup.Count())
	}
	if runtimesGroup.TotalSize(sr.Findings) != 1300 {
		t.Errorf("runtimes group total size = %d, want 1300", runtimesGroup.TotalSize(sr.Findings))
	}
}

func TestByCategoryEmpty(t *testing.T) {
	sr := &ScanResult{}
	groups := sr.ByCategory()
	if len(groups) != 0 {
		t.Errorf("empty ByCategory() returned %d groups, want 0", len(groups))
	}
}
