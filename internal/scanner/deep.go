// SPDX-FileCopyrightText: 2026 Mathias Kosinski
// SPDX-License-Identifier: AGPL-3.0-or-later

package scanner

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/derKosi/Ohm/internal/model"
)

// deepDirNames are directory-basename patterns that strongly indicate AI tooling
// but may live in non-standard locations the curated signatures don't cover.
var deepDirNames = []string{
	"models", "checkpoints", "loras", "lora",
	"comfyui", "stable-diffusion", "text-generation-webui",
	"ollama", "lm-studio", "koboldcpp", "gpt4all",
	".cursor", ".claude", ".codex", ".gemini", ".aider",
	".continue", ".vibe", ".paperclip", ".opencode",
}

// deepFileExts are file extensions that indicate AI model weights or caches
// beyond the .gguf / .safetensors already covered by scanModelCaches.
var deepFileExts = []string{
	".ckpt", ".pt", ".bin", ".onnx", ".tflite", ".mlmodel",
}

// scanDeep does a thorough crawl of the user's home directory looking for
// AI-related directories and files the curated signatures may have missed.
// It is opt-in (--deep) because it touches many more inodes than the default scan.
func (s *Scanner) scanDeep() {
	// Known paths already covered by other categories — don't re-report them.
	knownPaths := make(map[string]bool)
	s.mu.Lock()
	for _, f := range s.findings {
		for _, p := range f.ConfigPaths {
			knownPaths[filepath.Clean(p)] = true
		}
		if f.Path != "" && f.Path != "(scattered across projects)" && f.Path != "(see sub-items)" {
			knownPaths[filepath.Clean(f.Path)] = true
		}
	}
	s.mu.Unlock()

	// Directories found that look AI-related.
	type deepHit struct {
		path string
		size int64
	}
	var hits []deepHit

	// Skip well-known heavy/irrelevant subtrees to keep the walk bearable.
	skipDirs := map[string]bool{
		"node_modules": true, ".git": true, "Trash": true, "Caches": true,
		"Library": true, // macOS system dir — covered via explicit paths elsewhere
	}

	filepath.Walk(s.home, func(path string, info os.FileInfo, err error) error {
		if err != nil || info == nil {
			return nil
		}
		if info.IsDir() {
			name := info.Name()

			// Prune heavy/irrelevant subtrees.
			if skipDirs[name] {
				return filepath.SkipDir
			}
			// Don't follow symlinked dirs.
			if info.Mode()&os.ModeSymlink != 0 {
				return filepath.SkipDir
			}

			lower := strings.ToLower(name)
			for _, pat := range deepDirNames {
				if lower == pat || strings.HasPrefix(lower, pat) {
					abs, _ := filepath.Abs(path)
					if !knownPaths[filepath.Clean(abs)] {
						hits = append(hits, deepHit{path: abs, size: s.dirSize(abs)})
						return filepath.SkipDir // don't descend into matched dirs
					}
				}
			}
			return nil
		}

		// File-level check for model-weight extensions.
		for _, ext := range deepFileExts {
			if strings.HasSuffix(strings.ToLower(info.Name()), ext) {
				abs, _ := filepath.Abs(path)
				if !knownPaths[filepath.Clean(abs)] {
					hits = append(hits, deepHit{path: abs, size: info.Size()})
				}
				break
			}
		}
		return nil
	})

	if len(hits) == 0 {
		return
	}

	// Aggregate into a single finding with sub-items so the TUI stays readable.
	subs := make([]string, 0, len(hits))
	var totalSize int64
	for _, h := range hits {
		subs = append(subs, formatDeepPath(h.path, h.size))
		totalSize += h.size
	}

	s.addFinding(model.Finding{
		ID:        "deep-scan",
		Name:      "Deep-scan AI paths",
		Category:  model.CatStragglers,
		Path:      "(home directory crawl)",
		SizeBytes: totalSize,
		SubItems:  subs,
		RiskLevel: model.RiskCaution,
		UninstallCmds: map[string]string{
			"linux":   "# Review each path above and remove what you don't need",
			"macos":   "# Review each path above and remove what you don't need",
			"windows": "# Review each path above and remove what you don't need",
		},
	})
}

// formatDeepPath renders a discovered path + human-readable size for the sub-items list.
func formatDeepPath(path string, size int64) string {
	return path + " (" + model.FormatBytes(size) + ")"
}
