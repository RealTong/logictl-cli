package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/realtong/logi-cli/internal/app"
	"github.com/spf13/cobra"
)

func newInitCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Create a starter config from examples/config.toml",
		RunE: func(cmd *cobra.Command, args []string) error {
			paths := app.DefaultPaths()
			if err := os.MkdirAll(paths.ConfigDir, 0o755); err != nil {
				return err
			}

			examplePath, err := findExampleConfig()
			if err != nil {
				return err
			}

			data, err := os.ReadFile(examplePath)
			if err != nil {
				return err
			}

			if err := os.WriteFile(paths.ConfigFile, data, 0o644); err != nil {
				return err
			}

			cmd.Printf("wrote starter config to %s\n", paths.ConfigFile)
			return nil
		},
	}
}

func findExampleConfig() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		candidate := filepath.Join(dir, "examples", "config.toml")
		if _, err := os.Stat(candidate); err == nil {
			return candidate, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("examples/config.toml not found")
		}
		dir = parent
	}
}
