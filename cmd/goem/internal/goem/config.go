package goem

import (
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml"
)

type config struct {
	Maildir string `toml:"maildir"`
}

func loadConfig(path string) (*config, error) {
	if path == "" {
		cd, err := os.UserConfigDir()
		if err != nil {
			return nil, err
		}
		path = filepath.Join(cd, "goem", "config.toml")
	}
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	cfg := &config{}
	if err := toml.NewDecoder(f).Decode(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
