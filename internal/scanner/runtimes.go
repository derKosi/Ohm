// SPDX-FileCopyrightText: 2026 Mathias Kosinski
// SPDX-License-Identifier: AGPL-3.0-or-later

package scanner

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/derKosi/Ohm/internal/model"
	"github.com/derKosi/Ohm/internal/platform"
)

// scanRuntimes detects AI model runtimes.
func (s *Scanner) scanRuntimes() {
	type runtimeDef struct {
		id            string
		name          string
		configDirs    []string
		commands      []string
		pipPackages   []string
		risk          model.Risk
		uninstallCmds map[string]string
		checkModels   func(s *Scanner) []string
	}

	runtimes := []runtimeDef{
		{
			id:   "ollama",
			name: "Ollama",
			configDirs: []string{
				"~/.ollama",
				"/usr/share/ollama",
			},
			commands: []string{"ollama"},
			risk:     model.RiskSafe,
			uninstallCmds: map[string]string{
				"linux":   "sudo systemctl stop ollama && sudo systemctl disable ollama && sudo rm -f /usr/local/bin/ollama && sudo rm -rf /usr/share/ollama && sudo rm -f /etc/systemd/system/ollama.service && systemctl daemon-reload",
				"macos":   "rm -rf /Applications/Ollama.app ~/.ollama",
				"windows": "# Ollama: uninstall via Add/Remove Programs",
			},
			checkModels: func(s *Scanner) []string {
				out, err := s.runCommand("ollama", "list")
				if err != nil {
					return nil
				}
				var models []string
				for _, line := range strings.Split(out, "\n") {
					line = strings.TrimSpace(line)
					if line != "" && !strings.HasPrefix(line, "NAME") {
						parts := strings.Fields(line)
						if len(parts) > 0 {
							models = append(models, parts[0])
						}
					}
				}
				return models
			},
		},
		{
			id:   "lm-studio",
			name: "LM Studio",
			configDirs: func() []string {
				dirs := []string{"~/.cache/lm-studio"}
				if ld := platform.LocalAppData(); ld != "" {
					dirs = append(dirs, filepath.Join(ld, "LM Studio"))
				}
				return dirs
			}(),
			commands: []string{"lm-studio", "lms"},
			risk:     model.RiskSafe,
			uninstallCmds: map[string]string{
				"linux":   "# LM Studio: remove AppImage + rm -rf ~/.cache/lm-studio",
				"macos":   "rm -rf '/Applications/LM Studio.app' ~/.cache/lm-studio",
				"windows": "# LM Studio: uninstall via Add/Remove Programs",
			},
		},
		{
			id:   "gpt4all",
			name: "GPT4All",
			configDirs: []string{
				"~/.gpt4all",
			},
			commands: []string{"gpt4all"},
			risk:     model.RiskSafe,
			uninstallCmds: map[string]string{
				"linux":   "# GPT4All: remove binary + rm -rf ~/.gpt4all",
				"macos":   "rm -rf /Applications/GPT4All.app ~/.gpt4all",
				"windows": "# GPT4All: uninstall via Add/Remove Programs",
			},
		},
		// ── NEW: Additional runtimes ────────────────────────────────
		{
			id:          "localai",
			name:        "LocalAI",
			configDirs:  []string{"~/.localai"},
			commands:    []string{"localai", "local-ai"},
			risk:        model.RiskSafe,
			uninstallCmds: map[string]string{
				"linux":   "# LocalAI: remove binary + rm -rf ~/.localai",
				"macos":   "# LocalAI: remove binary + rm -rf ~/.localai",
				"windows": "# LocalAI: remove binary + Remove-Item ~/.localai -Recurse -Force",
			},
		},
		{
			id:          "vllm",
			name:        "vLLM",
			pipPackages: []string{"vllm"},
			risk:        model.RiskSafe,
			uninstallCmds: map[string]string{
				"linux":   "pip uninstall vllm",
				"macos":   "pip uninstall vllm",
				"windows": "pip uninstall vllm",
			},
		},
		{
			id:          "koboldcpp",
			name:        "KoboldCpp",
			configDirs:  []string{"~/.koboldcpp"},
			commands:    []string{"koboldcpp"},
			risk:        model.RiskSafe,
			uninstallCmds: map[string]string{
				"linux":   "# KoboldCpp: remove binary + rm -rf ~/.koboldcpp",
				"macos":   "# KoboldCpp: remove binary + rm -rf ~/.koboldcpp",
				"windows": "# KoboldCpp: uninstall via Add/Remove Programs",
			},
		},
		{
			id:          "text-gen-webui",
			name:        "Text Generation WebUI (oobabooga)",
			pipPackages: []string{"text-generation-webui"},
			risk:        model.RiskSafe,
			uninstallCmds: map[string]string{
				"linux":   "pip uninstall text-generation-webui && rm -rf ~/text-generation-webui",
				"macos":   "pip uninstall text-generation-webui && rm -rf ~/text-generation-webui",
				"windows": "pip uninstall text-generation-webui; Remove-Item ~/text-generation-webui -Recurse -Force",
			},
		},
		{
			id:          "tabby",
			name:        "Tabby (self-hosted)",
			commands:    []string{"tabby"},
			risk:        model.RiskSafe,
			uninstallCmds: map[string]string{
				"linux":   "# Tabby: remove binary + rm -rf ~/.tabby",
				"macos":   "# Tabby: remove binary + rm -rf ~/.tabby",
				"windows": "# Tabby: remove binary + Remove-Item ~/.tabby -Recurse -Force",
			},
		},
	}

	for _, rt := range runtimes {
		found := false
		var totalSize int64
		var foundPaths []string

		for _, cmd := range rt.commands {
			if s.hasCommand(cmd) {
				found = true
			}
		}

		for _, dir := range rt.configDirs {
			expanded := s.expandPath(dir)
			if s.dirExists(expanded) {
				found = true
				size := s.dirSize(expanded)
				totalSize += size
				foundPaths = append(foundPaths, expanded)
			}
		}

		for _, pkg := range rt.pipPackages {
			if s.hasPipPackage(pkg) {
				found = true
			}
		}

		if found {
			var subItems []string
			if rt.checkModels != nil {
				subItems = rt.checkModels(s)
			}

			s.addFinding(model.Finding{
				ID:            rt.id,
				Category:      model.CatRuntimes,
				Name:          rt.name,
				Path:          joinPaths(foundPaths),
				SizeBytes:     totalSize,
				ConfigPaths:   foundPaths,
				SubItems:      subItems,
				RiskLevel:     rt.risk,
				UninstallCmds: rt.uninstallCmds,
			})
		}
	}
}

// findFilesByExt finds files with a given extension under listed roots up to maxDepth.
func findFilesByExt(roots []string, ext string, maxDepth int) []string {
	var results []string
	for _, root := range roots {
		if _, err := os.Stat(root); err != nil {
			continue
		}
		_ = maxDepth // used via closure below
		filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil
			}
			rel, _ := filepath.Rel(root, path)
			d := len(strings.Split(rel, string(filepath.Separator)))
			if d > maxDepth {
				if info.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}
			if !info.IsDir() && strings.HasSuffix(info.Name(), ext) {
				results = append(results, path)
			}
			return nil
		})
	}
	return results
}
