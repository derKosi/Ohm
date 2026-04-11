package scanner

import (
	"os"
	"path/filepath"

	"github.com/derKosi/Ohm/internal/model"
)

// scanPlugins detects AI tool plugins and extensions.
func (s *Scanner) scanPlugins() {
	var subItems []string
	var totalSize int64

	// Pi skills
	piSkillsDir := filepath.Join(s.home, ".pi", "agent", "skills")
	if s.dirExists(piSkillsDir) {
		entries, err := os.ReadDir(piSkillsDir)
		if err == nil {
			for _, e := range entries {
				name := e.Name()
				skillPath := filepath.Join(piSkillsDir, name)
				// Follow symlinks
				info, err := os.Stat(skillPath)
				if err != nil {
					continue
				}
				if info.IsDir() {
					skillFile := filepath.Join(skillPath, "SKILL.md")
					if s.fileExists(skillFile) {
						subItems = append(subItems, "pi skill: "+name)
						totalSize += s.dirSize(skillPath)
					}
				}
			}
		}
	}

	// ComfyUI custom nodes (if ComfyUI found)
	comfyLocations := []string{
		filepath.Join(s.home, "ComfyUI", "custom_nodes"),
		filepath.Join(s.home, "comfyui", "custom_nodes"),
	}
	for _, cn := range comfyLocations {
		if s.dirExists(cn) {
			entries, err := os.ReadDir(cn)
			if err == nil {
				for _, e := range entries {
					if e.IsDir() {
						subItems = append(subItems, "ComfyUI node: "+e.Name())
						totalSize += s.dirSize(filepath.Join(cn, e.Name()))
					}
				}
			}
		}
	}

	// oh-my-opencode / oh-my-openagent
	omoLocations := []string{
		filepath.Join(s.home, ".oh-my-opencode"),
		filepath.Join(s.home, ".config", "opencode", "oh-my-opencode.json"),
		filepath.Join(s.home, ".opencode", "oh-my-opencode.json"),
		filepath.Join(s.home, ".opencode", "oh-my-openagent.json"),
	}
	for _, loc := range omoLocations {
		info, err := os.Stat(loc)
		if err != nil {
			continue
		}
		subItems = append(subItems, "oh-my-opencode: "+loc)
		if info.IsDir() {
			totalSize += s.dirSize(loc)
		} else {
			totalSize += info.Size()
		}
	}

	// oh-my-pi
	ohMyPiDir := filepath.Join(s.home, ".oh-my-pi")
	if s.dirExists(ohMyPiDir) {
		subItems = append(subItems, "oh-my-pi")
		totalSize += s.dirSize(ohMyPiDir)
	}

	// oh-my-codex
	ohMyCodexDir := filepath.Join(s.home, ".oh-my-codex")
	if s.dirExists(ohMyCodexDir) {
		subItems = append(subItems, "oh-my-codex")
		totalSize += s.dirSize(ohMyCodexDir)
	}

	if len(subItems) > 0 {
		// Deduplicate
		subItems = dedupe(subItems)
		s.addFinding(model.Finding{
			ID:            "ai-plugins",
			Category:      model.CatPlugins,
			Name:          "AI Plugins & Extensions",
			Path:          "(see sub-items)",
			SizeBytes:     totalSize,
			SubItems:      subItems,
			RiskLevel:     model.RiskCaution,
			UninstallCmds: map[string]string{
				"linux":   "# Remove plugin directories listed above",
				"macos":   "# Remove plugin directories listed above",
				"windows": "# Remove plugin directories listed above",
			},
		})
	}
}

func dedupe(items []string) []string {
	seen := make(map[string]bool)
	var result []string
	for _, item := range items {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}
	return result
}
