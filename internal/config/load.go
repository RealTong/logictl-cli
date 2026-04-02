package config

import (
	"fmt"
	"strings"

	"github.com/BurntSushi/toml"
)

func LoadFile(path string) (*Config, error) {
	var cfg Config
	md, err := toml.DecodeFile(path, &cfg)
	if err != nil {
		return nil, err
	}

	if undecoded := md.Undecoded(); len(undecoded) > 0 {
		keys := make([]string, 0, len(undecoded))
		for _, key := range undecoded {
			keys = append(keys, key.String())
		}
		return nil, fmt.Errorf("config contains unsupported keys: %s", strings.Join(keys, ", "))
	}

	return &cfg, nil
}
