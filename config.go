package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/adrg/xdg"
	"github.com/motemen/go-gitconfig"
)

const (
	patEnvName    = "GITHUB_PAT"
	editorEnvName = "EDITOR"
)

type config struct {
	Owner          string `json:"owner" gitconfig:"user.name"`
	Repo           string `json:"repo"`
	Branch         string `json:"branch"`
	CommitterName  string `json:"committer_name" gitconfig:"user.name"`
	CommitterEmail string `json:"committer_email" gitconfig:"user.email"`
	PAT            string `json:"pat"`
	Editor         string `json:"editor"`
}

func newConfig(gitCfg gitconfig.Config) (*config, error) {
	cfg := &config{
		Repo:   "diary",
		Branch: "main",
		PAT:    os.Getenv(patEnvName),
		Editor: os.Getenv(editorEnvName),
	}

	if err := gitCfg.Load(cfg); err != nil {
		return nil, fmt.Errorf("failed to load gitconfig: %w", err)
	}

	return cfg, nil
}

func (cfg *config) getConfigPath(configPath string) string {
	if configPath != "" {
		return configPath
	}
	return path.Join(xdg.ConfigHome, "hubdiary", "config.json")
}

func (cfg *config) loadFile(configPath string) error {
	configPath = cfg.getConfigPath(configPath)

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return fmt.Errorf("config file %s doesn't exist: %w", configPath, err)
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read config file `%s`: %w", configPath, err)
	}

	if err := json.Unmarshal(data, cfg); err != nil {
		return fmt.Errorf("failed to unmarshal config file: `%s`: %w", configPath, err)
	}

	return nil
}
