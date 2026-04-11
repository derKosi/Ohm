package scanner

import (
	"os"
	"path/filepath"

	"github.com/derKosi/Ohm/internal/model"
)

// scanComfyUI detects ComfyUI installations and image models.
func (s *Scanner) scanComfyUI() {
	// Known ComfyUI install locations
	comfyLocations := []string{
		filepath.Join(s.home, "ComfyUI"),
		filepath.Join(s.home, "comfyui"),
		filepath.Join(s.home, "stable-diffusion-webui"),
		filepath.Join(s.home, "Fooocus"),
		filepath.Join(s.home, "invokeai"),
	}

	modelSubDirs := []struct {
		name string
		path string
	}{
		{"Checkpoints", "models/checkpoints"},
		{"LoRA Adapters", "models/loras"},
		{"ControlNet Models", "models/controlnet"},
		{"VAE", "models/vae"},
		{"CLIP Models", "models/clip"},
		{"UNet Models", "models/unet"},
		{"Embeddings", "models/embeddings"},
		{"Upscale Models", "models/upscale_models"},
	}

	for _, loc := range comfyLocations {
		if !s.dirExists(loc) {
			continue
		}

		totalSize := s.dirSize(loc)
		var subItems []string
		var configPaths []string

		for _, sub := range modelSubDirs {
			subPath := filepath.Join(loc, sub.path)
			if s.dirExists(subPath) {
				subSize := s.dirSize(subPath)
				if subSize > 0 {
					subItems = append(subItems, sub.name+" ("+model.FormatBytes(subSize)+")")
					configPaths = append(configPaths, subPath)
				}
			}
		}

		// Check custom nodes
		customNodes := filepath.Join(loc, "custom_nodes")
		if s.dirExists(customNodes) {
			entries, err := os.ReadDir(customNodes)
			if err == nil && len(entries) > 0 {
				for _, e := range entries {
					if e.IsDir() {
						subItems = append(subItems, "Custom Node: "+e.Name())
					}
				}
			}
		}

		name := "ComfyUI"
		if filepath.Base(loc) == "stable-diffusion-webui" {
			name = "Stable Diffusion WebUI (A1111)"
		} else if filepath.Base(loc) == "Fooocus" {
			name = "Fooocus"
		} else if filepath.Base(loc) == "invokeai" {
			name = "InvokeAI"
		}

		s.addFinding(model.Finding{
			ID:          "comfyui-" + filepath.Base(loc),
			Category:    model.CatComfyUI,
			Name:        name,
			Path:        loc,
			SizeBytes:   totalSize,
			ConfigPaths: configPaths,
			SubItems:    subItems,
			RiskLevel:   model.RiskSafe,
			UninstallCmds: map[string]string{
				"linux":   "rm -rf " + loc,
				"macos":   "rm -rf " + loc,
				"windows": "Remove-Item '" + loc + "' -Recurse -Force",
			},
		})
	}
}
