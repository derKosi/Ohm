package scanner

import (
	"strings"

	"github.com/derKosi/Ohm/internal/model"
)

// scanDocker detects AI-related Docker images.
func (s *Scanner) scanDocker() {
	if !s.hasCommand("docker") {
		return
	}

	// Check if docker daemon is running
	_, err := s.runCommand("docker", "info")
	if err != nil {
		return
	}

	// Get image list
	out, err := s.runCommand("docker", "images", "--format", "{{.Repository}}:{{.Tag}}\t{{.Size}}")
	if err != nil {
		return
	}

	aiKeywords := []string{
		"ollama", "llm", "gpt", "ai", "comfy", "stable-diffusion",
		"langchain", "vllm", "localai", "text-gen", "llama",
		"mistral", "hugging", "transformers", "litellm", "openai",
		"anthropic", "deepseek", "qwen",
	}

	var aiImages []string
	for _, line := range strings.Split(out, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		imageName := strings.Split(line, "\t")[0]
		lower := strings.ToLower(imageName)
		for _, kw := range aiKeywords {
			if strings.Contains(lower, kw) {
				aiImages = append(aiImages, line)
				break
			}
		}
	}

	if len(aiImages) > 0 {
		s.addFinding(model.Finding{
			ID:        "docker-ai-images",
			Category:  model.CatDocker,
			Name:      "AI Docker Images",
			Path:      "(Docker daemon)",
			SizeBytes: 0, // Docker doesn't report exact bytes in this format
			SubItems:  aiImages,
			RiskLevel: model.RiskSafe,
			UninstallCmds: map[string]string{
				"linux":   "# Remove with: docker rmi <image>",
				"macos":   "# Remove with: docker rmi <image>",
				"windows": "# Remove with: docker rmi <image>",
			},
		})
	}
}
