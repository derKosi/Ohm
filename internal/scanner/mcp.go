package scanner

import (
	"github.com/derKosi/Ohm/internal/model"
)

// scanMCP detects MCP (Model Context Protocol) configuration files.
func (s *Scanner) scanMCP() {
	mcpFiles := findInstructionFiles(s.home, 4, []string{".mcp.json", "mcp.json"})

	if len(mcpFiles) == 0 {
		return
	}

	var totalSize int64
	for _, f := range mcpFiles {
		totalSize += s.fileSize(f)
	}

	s.addFinding(model.Finding{
		ID:        "mcp-configs",
		Category:  model.CatMCP,
		Name:      "MCP Configuration Files",
		Path:      "(scattered across projects)",
		SizeBytes: totalSize,
		SubItems:  mcpFiles,
		RiskLevel: model.RiskDanger, // MCP configs often contain API keys and connection strings
		UninstallCmds: map[string]string{
			"linux":   "# ⚠️ REVIEW FOR API KEYS before removing MCP configs",
			"macos":   "# ⚠️ REVIEW FOR API KEYS before removing MCP configs",
			"windows": "# ⚠️ REVIEW FOR API KEYS before removing MCP configs",
		},
	})
}
