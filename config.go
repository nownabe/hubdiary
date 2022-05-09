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
	patEnvName = "GITHUB_PAT"
)

type config struct {
	Repo  string `json:"repo"`
	User  string `json:"user"`
	Email string `json:"email"`
	PAT   string `json:"pat"`
}

type gitConfig struct {
	UserName  string `gitconfig:"user.name"`
	UserEmail string `gitconfig:"user.email"`
}

func newConfig(gitCfg gitconfig.Config, configPath, repo, user, email, pat string) (*config, error) {
	cfg, err := configFromFile(configPath)
	if err != nil {
		return nil, err
	}

	if err := applyDefaultConfig(gitCfg, cfg); err != nil {
		return nil, err
	}

	if repo != "" {
		cfg.Repo = repo
	}

	if user != "" {
		cfg.User = user
	}

	if email != "" {
		cfg.Email = email
	}

	if pat != "" {
		cfg.PAT = pat
	}

	return cfg, nil
}

func applyDefaultConfig(gitCfg gitconfig.Config, cfg *config) error {
	var gc gitConfig

	if err := gitCfg.Load(&gc); err != nil {
		return fmt.Errorf("failed to load gitconfig: %w", err)
	}

	if cfg.User == "" {
		cfg.User = gc.UserName
	}

	if cfg.Email == "" {
		cfg.Email = gc.UserEmail
	}

	if cfg.Repo == "" {
		cfg.Repo = cfg.User + "/diary"
	}

	if cfg.PAT == "" {
		cfg.PAT = os.Getenv(patEnvName)
	}

	return nil
}

func configFromFile(configPath string) (*config, error) {
	if configPath == "" {
		configPath = path.Join(xdg.ConfigHome, "hubdiary", "config.json")
	}

	cfg := &config{}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return cfg, nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file `%s`: %w", configPath, err)
	}

	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config file: `%s`: %w", configPath, err)
	}

	return cfg, nil
}
