// SPDX-FileCopyrightText: 2026 Mathias Kosinski
// SPDX-License-Identifier: AGPL-3.0-or-later

package scanner

import (
	"os"
	"path/filepath"

	"github.com/derKosi/Ohm/internal/model"
)

// scanModelCaches detects cached model files.
func (s *Scanner) scanModelCaches() {
	type cacheDef struct {
		id       string
		name     string
		paths    []string
		risk     model.Risk
		cleanup  map[string]string
	}

	caches := []cacheDef{
		{
			id:   "huggingface-cache",
			name: "HuggingFace Hub Cache",
			paths: []string{
				filepath.Join(s.home, ".cache", "huggingface"),
			},
			risk: model.RiskSafe,
			cleanup: map[string]string{
				"linux":   "rm -rf ~/.cache/huggingface",
				"macos":   "rm -rf ~/.cache/huggingface",
				"windows": "Remove-Item \"$env:USERPROFILE\\.cache\\huggingface\" -Recurse -Force",
			},
		},
		{
			id:   "pytorch-cache",
			name: "PyTorch Hub Cache",
			paths: []string{
				filepath.Join(s.home, ".cache", "torch"),
			},
			risk: model.RiskSafe,
			cleanup: map[string]string{
				"linux":   "rm -rf ~/.cache/torch",
				"macos":   "rm -rf ~/.cache/torch",
				"windows": "Remove-Item \"$env:USERPROFILE\\.cache\\torch\" -Recurse -Force",
			},
		},
		{
			id:   "playwright-browsers",
			name: "Playwright Browsers",
			paths: []string{
				filepath.Join(s.home, ".cache", "ms-playwright"),
			},
			risk: model.RiskSafe,
			cleanup: map[string]string{
				"linux":   "rm -rf ~/.cache/ms-playwright",
				"macos":   "rm -rf ~/.cache/ms-playwright",
				"windows": "Remove-Item \"$env:USERPROFILE\\.cache\\ms-playwright\" -Recurse -Force",
			},
		},
	}

	for _, cache := range caches {
		var foundPaths []string
		var totalSize int64

		for _, p := range cache.paths {
			if s.dirExists(p) {
				size := s.dirSize(p)
				if size > 0 {
					foundPaths = append(foundPaths, p)
					totalSize += size
				}
			}
		}

		if len(foundPaths) > 0 {
			s.addFinding(model.Finding{
				ID:            cache.id,
				Category:      model.CatModelCaches,
				Name:          cache.name,
				Path:          joinPaths(foundPaths),
				SizeBytes:     totalSize,
				ConfigPaths:   foundPaths,
				RiskLevel:     cache.risk,
				UninstallCmds: cache.cleanup,
			})
		}
	}

	// Scan for .gguf files
	ggufFiles := findFilesByExt([]string{s.home}, ".gguf", 4)
	if len(ggufFiles) > 0 {
		var totalSize int64
		for _, f := range ggufFiles {
			if info, err := os.Stat(f); err == nil {
				totalSize += info.Size()
			}
		}
		s.addFinding(model.Finding{
			ID:       "gguf-files",
			Category: model.CatModelCaches,
			Name:     "GGUF Model Files",
			Path:     filepath.Dir(ggufFiles[0]) + "/*",
			SizeBytes: totalSize,
			SubItems: ggufFiles,
			RiskLevel: model.RiskSafe,
			UninstallCmds: map[string]string{
				"linux":   "# Remove .gguf files manually",
				"macos":   "# Remove .gguf files manually",
				"windows": "# Remove .gguf files manually",
			},
		})
	}

	// Scan for .safetensors files
	safetensorFiles := findFilesByExt([]string{s.home}, ".safetensors", 4)
	if len(safetensorFiles) > 0 {
		var totalSize int64
		for _, f := range safetensorFiles {
			if info, err := os.Stat(f); err == nil {
				totalSize += info.Size()
			}
		}
		s.addFinding(model.Finding{
			ID:       "safetensors-files",
			Category: model.CatModelCaches,
			Name:     "Safetensors Model Files",
			Path:     filepath.Dir(safetensorFiles[0]) + "/*",
			SizeBytes: totalSize,
			SubItems: safetensorFiles,
			RiskLevel: model.RiskSafe,
			UninstallCmds: map[string]string{
				"linux":   "# Remove .safetensors files manually",
				"macos":   "# Remove .safetensors files manually",
				"windows": "# Remove .safetensors files manually",
			},
		})
	}
}
