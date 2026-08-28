package api

import (
	"fmt"
	"os"
	"path/filepath"
)

var appIconData []byte

func init() {
	exePath, err := os.Executable()
	if err == nil {
		exeDir := filepath.Dir(exePath)
		candidates := []string{
			filepath.Join(exeDir, "icon.png"),
			filepath.Join(exeDir, "..", "icon.png"),
			filepath.Join(exeDir, "Resources", "icon.png"),
		}
		for _, p := range candidates {
			if data, err := os.ReadFile(p); err == nil && len(data) > 0 {
				appIconData = data
				return
			}
		}
	}

	wd, _ := os.Getwd()
	iconPath := filepath.Join(wd, "icon.png")
	if data, err := os.ReadFile(iconPath); err == nil && len(data) > 0 {
		appIconData = data
		return
	}

	fmt.Fprintln(os.Stderr, "WARNING: icon.png not found")
}
