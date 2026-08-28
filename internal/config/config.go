package config

import (
	"os"
	"path/filepath"

	"bling/internal/model"

	"gopkg.in/yaml.v3"
)

func LoadConfig() *model.AppConfig {
	config := &model.AppConfig{
		GUI: model.GUIConfig{
			Title:  "鼠标精灵",
			Width:  400,
			Height: 280,
		},
		Mover: model.MoverConfig{
			Interval: 60,
		},
	}

	configPath := findConfigFile()
	data, err := os.ReadFile(configPath)
	if err != nil {
		return config
	}

	yaml.Unmarshal(data, config)
	return config
}

func findConfigFile() string {
	candidates := []string{
		"configs/config.yaml",
		"config.yaml",
	}

	if exePath, err := os.Executable(); err == nil {
		exeDir := filepath.Dir(exePath)
		exeParent := filepath.Dir(exeDir)

		candidates = append(candidates,
			filepath.Join(exeDir, "configs", "config.yaml"),
			filepath.Join(exeDir, "config.yaml"),
			filepath.Join(exeParent, "configs", "config.yaml"),
			filepath.Join(exeParent, "config.yaml"),
			filepath.Join(exeParent, "Resources", "config.yaml"),
			filepath.Join(exeParent, "Resources", "configs", "config.yaml"),
		)
	}

	for _, path := range candidates {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	return "configs/config.yaml"
}
