// SPDX-FileCopyrightText: 2026 Mathias Kosinski
// SPDX-License-Identifier: AGPL-3.0-or-later

package model

import (
	"fmt"
	"time"
)

// Category represents a scan category for AI software.
type Category int

const (
	CatAgents Category = iota
	CatEditors
	CatRuntimes
	CatComfyUI
	CatSDKs
	CatModelCaches
	CatInstructions
	CatMemory
	CatMCP
	CatPlugins
	CatConfigs
	CatDocker
	CatStragglers
)

func (c Category) String() string {
	names := []string{
		"Agents & Harnesses",
		"AI Editors & IDEs",
		"Model Runtimes",
		"ComfyUI & Image Models",
		"SDKs & Frameworks",
		"Model Caches",
		"Agent Config & Instructions",
		"Agent Memory & Sessions",
		"MCP Configurations",
		"Plugins & Extensions",
		"Config & Data Dirs",
		"Docker",
		"Stragglers",
	}
	if int(c) < len(names) {
		return names[c]
	}
	return "Unknown"
}

func (c Category) Icon() string {
	icons := []string{
		"🤖",
		"🖥️",
		"⚙️",
		"🎨",
		"📦",
		"💾",
		"📄",
		"🧠",
		"🔌",
		"🧩",
		"📁",
		"🐳",
		"👻",
	}
	if int(c) < len(icons) {
		return icons[c]
	}
	return "❓"
}

// Risk represents the risk level of removing a finding.
type Risk int

const (
	RiskSafe    Risk = iota // Standard removal, no concerns
	RiskCaution             // May affect other software
	RiskDanger              // Contains credentials, keys, or sensitive data
)

func (r Risk) String() string {
	names := []string{"Safe", "Caution", "Danger"}
	if int(r) < len(names) {
		return names[r]
	}
	return "Unknown"
}

func (r Risk) Icon() string {
	icons := []string{"  ", "⚠️ ", "🔑"}
	if int(r) < len(icons) {
		return icons[r]
	}
	return "❓"
}

// Finding represents a single discovered AI software item.
type Finding struct {
	ID            string            `json:"id"`
	Category      Category          `json:"category"`
	Name          string            `json:"name"`
	Version       string            `json:"version,omitempty"`
	InstallMethod string            `json:"install_method,omitempty"`
	Path          string            `json:"path"`
	SizeBytes     int64             `json:"size_bytes"`
	ConfigPaths   []string          `json:"config_paths,omitempty"`
	SubItems      []string          `json:"sub_items,omitempty"`
	RiskLevel     Risk              `json:"risk_level"`
	UninstallCmds map[string]string `json:"uninstall_cmds,omitempty"`
	Selected      bool              `json:"selected"`
}

// FormatSize returns a human-readable size string.
func (f *Finding) FormatSize() string {
	return FormatBytes(f.SizeBytes)
}

// FormatBytes converts bytes to human-readable string.
func FormatBytes(b int64) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
		TB = GB * 1024
	)
	switch {
	case b >= TB:
		return fmt.Sprintf("%.1f TB", float64(b)/float64(TB))
	case b >= GB:
		return fmt.Sprintf("%.1f GB", float64(b)/float64(GB))
	case b >= MB:
		return fmt.Sprintf("%.1f MB", float64(b)/float64(MB))
	case b >= KB:
		return fmt.Sprintf("%.1f KB", float64(b)/float64(KB))
	default:
		return fmt.Sprintf("%d B", b)
	}
}

// ScanResult holds all findings from a scan.
type ScanResult struct {
	Findings  []Finding  `json:"findings"`
	ScannedAt time.Time  `json:"scanned_at"`
	Platform  string     `json:"platform"`
	Hostname  string     `json:"hostname"`
	Warnings  []string   `json:"warnings,omitempty"`
}

// TotalSize returns the total size of all findings.
func (sr *ScanResult) TotalSize() int64 {
	var total int64
	for _, f := range sr.Findings {
		total += f.SizeBytes
	}
	return total
}

// SelectedSize returns the total size of selected findings.
func (sr *ScanResult) SelectedSize() int64 {
	var total int64
	for _, f := range sr.Findings {
		if f.Selected {
			total += f.SizeBytes
		}
	}
	return total
}

// ByCategory groups findings by category, returning indices into Findings.
func (sr *ScanResult) ByCategory() []CategoryGroup {
	seen := make(map[Category]bool)
	var result []CategoryGroup
	for _, f := range sr.Findings {
		if !seen[f.Category] {
			seen[f.Category] = true
			result = append(result, CategoryGroup{Category: f.Category})
		}
	}
	for i := range sr.Findings {
		for gi := range result {
			if sr.Findings[i].Category == result[gi].Category {
				result[gi].FindingIdxs = append(result[gi].FindingIdxs, i)
				break
			}
		}
	}
	return result
}

// CategoryGroup is a group of findings in the same category.
type CategoryGroup struct {
	Category    Category
	FindingIdxs []int // indices into ScanResult.Findings
}

// TotalSize returns the total size of findings in this group.
func (cg CategoryGroup) TotalSize(all []Finding) int64 {
	var total int64
	for _, idx := range cg.FindingIdxs {
		total += all[idx].SizeBytes
	}
	return total
}

// Count returns count of findings in this group.
func (cg CategoryGroup) Count() int {
	return len(cg.FindingIdxs)
}

// Count returns the total number of findings.
func (sr *ScanResult) Count() int {
	return len(sr.Findings)
}

// SelectedCount returns the number of selected findings.
func (sr *ScanResult) SelectedCount() int {
	count := 0
	for _, f := range sr.Findings {
		if f.Selected {
			count++
		}
	}
	return count
}
