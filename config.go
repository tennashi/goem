package goem

import (
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml"
	"github.com/tennashi/goem/shellpath"
)

type Config struct {
	RootDir string       `toml:"root_dir"`
	Server  ServerConfig `toml:"server"`
}

type ServerConfig struct {
	Port string `toml:"port"`
}

func NewConfig(path string) *Config {
	path = configPath(path)
	file, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer file.Close()
	config := &Config{}
	if err := toml.NewDecoder(file).Decode(config); err != nil {
		return nil
	}

	config.RootDir = shellpath.Resolve(config.RootDir)

	return config
}

func configPath(path string) string {
	if path != "" {
		return path
	}
	path, err := os.UserConfigDir()
	if err != nil {
		return filepath.Join(".", "config.toml")
	}
	return filepath.Join(path, "goemd", "config.toml")
}
