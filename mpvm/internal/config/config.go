package config

import (
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
)

// XDG-like paths but always using ~/.config, ~/.local/state, ~/.cache
func getConfigHome() string {
	if env := os.Getenv("XDG_CONFIG_HOME"); env != "" {
		return env
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config")
}

func getStateHome() string {
	if env := os.Getenv("XDG_STATE_HOME"); env != "" {
		return env
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".local", "state")
}

func getCacheHome() string {
	if env := os.Getenv("XDG_CACHE_HOME"); env != "" {
		return env
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".cache")
}

type Config struct {
	VM     VMConfig     `toml:"vm"`
	Git    GitConfig    `toml:"git"`
	SSH    SSHConfig    `toml:"ssh"`
	Claude ClaudeConfig `toml:"claude"`
}

type VMConfig struct {
	Name   string `toml:"name"`
	CPUs   int    `toml:"cpus"`
	Memory string `toml:"memory"`
	Disk   string `toml:"disk"`
}

type GitConfig struct {
	UserName  string `toml:"user_name"`
	UserEmail string `toml:"user_email"`
}

type SSHConfig struct {
	SigningKey string `toml:"signing_key"`
}

type ClaudeConfig struct {
	Marketplaces []string `toml:"marketplaces"`
}

func DefaultConfig() *Config {
	home, _ := os.UserHomeDir()
	return &Config{
		VM: VMConfig{
			Name:   "coding-agent",
			CPUs:   2,
			Memory: "4G",
			Disk:   "20G",
		},
		Git: GitConfig{},
		SSH: SSHConfig{
			SigningKey: filepath.Join(home, ".ssh", "id_ed25519.pub"),
		},
		Claude: ClaudeConfig{
			Marketplaces: []string{
				"anthropics/claude-plugins-official",
				"anthropics/skills",
			},
		},
	}
}

func ConfigDir() string {
	return filepath.Join(getConfigHome(), "mpvm")
}

func ConfigPath() string {
	return filepath.Join(ConfigDir(), "config.toml")
}

func StateDir() string {
	return filepath.Join(getStateHome(), "mpvm")
}

func CacheDir() string {
	return filepath.Join(getCacheHome(), "mpvm")
}

func Load() (*Config, error) {
	cfg := DefaultConfig()

	configPath := ConfigPath()
	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return nil, err
	}

	if err := toml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) Save() error {
	configDir := ConfigDir()
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	data, err := toml.Marshal(c)
	if err != nil {
		return err
	}

	return os.WriteFile(ConfigPath(), data, 0644)
}

func ExpandPath(path string) string {
	if path == "" {
		return ""
	}

	if path[0] == '~' {
		home, _ := os.UserHomeDir()
		path = filepath.Join(home, path[1:])
	}

	return filepath.FromSlash(path)
}
