package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	APIURL    string `json:"api_url"`
	IconsPath string `json:"icons_path"`
}

func LoadConfig() (*Config, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}

	appDir := filepath.Join(configDir, "OpenLinkHub-system-tray")
	configPath := filepath.Join(appDir, "config.json")

	if err := os.MkdirAll(appDir, 0755); err != nil {
		return nil, err
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		defaultConfig := Config{
			APIURL:    "http://127.0.0.1:27003/api",
			IconsPath: fmt.Sprintf("%s/icons", appDir),
		}

		data, _ := json.MarshalIndent(defaultConfig, "", "  ")
		if err := os.WriteFile(configPath, data, 0644); err != nil {
			return nil, err
		}
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
