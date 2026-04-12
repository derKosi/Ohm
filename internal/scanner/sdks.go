// SPDX-FileCopyrightText: 2026 Mathias Kosinski
// SPDX-License-Identifier: AGPL-3.0-or-later

package scanner

import (
	"strings"

	"github.com/derKosi/Ohm/internal/model"
)

// scanSDKs detects AI SDKs and frameworks.
func (s *Scanner) scanSDKs() {
	type sdkDef struct {
		id           string
		name         string
		pipPackages  []string
		npmPackages  []string
		risk         model.Risk
		uninstallCmd map[string]string
	}

	sdks := []sdkDef{
		{
			id:   "pytorch",
			name: "PyTorch",
			pipPackages: []string{"torch", "torchvision", "torchaudio"},
			risk: model.RiskSafe,
			uninstallCmd: map[string]string{
				"all": "pip uninstall torch torchvision torchaudio",
			},
		},
		{
			id:   "tensorflow",
			name: "TensorFlow",
			pipPackages: []string{"tensorflow"},
			risk: model.RiskSafe,
			uninstallCmd: map[string]string{
				"all": "pip uninstall tensorflow",
			},
		},
		{
			id:   "huggingface",
			name: "HuggingFace Transformers",
			pipPackages: []string{"transformers", "datasets", "tokenizers"},
			risk: model.RiskSafe,
			uninstallCmd: map[string]string{
				"all": "pip uninstall transformers datasets tokenizers",
			},
		},
		{
			id:   "langchain",
			name: "LangChain",
			pipPackages: []string{"langchain"},
			npmPackages: []string{"langchain"},
			risk: model.RiskSafe,
			uninstallCmd: map[string]string{
				"all": "pip uninstall langchain",
			},
		},
		{
			id:   "llamaindex",
			name: "LlamaIndex",
			pipPackages: []string{"llama-index"},
			risk: model.RiskSafe,
			uninstallCmd: map[string]string{
				"all": "pip uninstall llama-index",
			},
		},
		{
			id:   "openai-sdk",
			name: "OpenAI SDK",
			pipPackages: []string{"openai"},
			npmPackages: []string{"openai"},
			risk: model.RiskSafe,
			uninstallCmd: map[string]string{
				"all": "pip uninstall openai",
			},
		},
		{
			id:   "anthropic-sdk",
			name: "Anthropic SDK",
			pipPackages: []string{"anthropic"},
			npmPackages: []string{"@anthropic-ai/sdk"},
			risk: model.RiskSafe,
			uninstallCmd: map[string]string{
				"all": "pip uninstall anthropic",
			},
		},
		{
			id:   "playwright",
			name: "Playwright (AI-adjacent)",
			pipPackages: []string{"playwright"},
			npmPackages: []string{"playwright"},
			risk: model.RiskCaution,
			uninstallCmd: map[string]string{
				"all": "pip uninstall playwright",
			},
		},
		{
			id:   "selenium",
			name: "Selenium (AI-adjacent)",
			pipPackages: []string{"selenium"},
			risk: model.RiskCaution,
			uninstallCmd: map[string]string{
				"all": "pip uninstall selenium",
			},
		},
	}

	for _, sdk := range sdks {
		found := false

		for _, pkg := range sdk.pipPackages {
			if s.hasPipPackage(pkg) {
				found = true
			}
		}

		for _, pkg := range sdk.npmPackages {
			if s.hasNpmPackage(pkg) {
				found = true
			}
		}

		if found {
			uninstallCmds := make(map[string]string)
			for _, v := range sdk.uninstallCmd {
				uninstallCmds["linux"] = v
				uninstallCmds["macos"] = v
				uninstallCmds["windows"] = v
			}

			s.addFinding(model.Finding{
				ID:            sdk.id,
				Category:      model.CatSDKs,
				Name:          sdk.name,
				InstallMethod: "package",
				SizeBytes:     0, // SDKs don't have easy size calculation
				RiskLevel:     sdk.risk,
				UninstallCmds: uninstallCmds,
				SubItems:      sdk.pipPackages,
			})
		}
	}
}

// joinPaths joins paths with ", ".
func joinPaths(paths []string) string {
	return strings.Join(paths, ", ")
}
