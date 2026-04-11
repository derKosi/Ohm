package scanner

import "github.com/derKosi/Ohm/internal/model"

// Instruction files to search for.
var instructionFiles = []string{
	"AGENTS.md",
	"CLAUDE.md",
	".cursorrules",
	".windsurfrules",
	"GEMINI.md",
	"copilot-instructions.md",
	"CONVENTIONS.md",
	".cursorrules",
	".clinerules",
	".aider.conf.yml",
}

// scanInstructions detects agent instruction and config files.
func (s *Scanner) scanInstructions() {
	found := findInstructionFiles(s.home, 4, instructionFiles)

	if len(found) == 0 {
		return
	}

	var totalSize int64
	for _, f := range found {
		totalSize += s.fileSize(f)
	}

	s.addFinding(model.Finding{
		ID:          "agent-instructions",
		Category:    model.CatInstructions,
		Name:        "Agent Instruction Files",
		Path:        "(scattered across projects)",
		SizeBytes:   totalSize,
		SubItems:    found,
		RiskLevel:   model.RiskCaution,
		UninstallCmds: map[string]string{
			"linux":   "# Review and remove individual files listed above",
			"macos":   "# Review and remove individual files listed above",
			"windows": "# Review and remove individual files listed above",
		},
	})
}
