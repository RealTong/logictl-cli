package cli

import (
	"os"

	"github.com/realtong/logictl-cli/examples"
	"github.com/realtong/logictl-cli/internal/app"
	"github.com/spf13/cobra"
)

func starterConfigContent() []byte {
	return examples.ConfigTOML
}

func newInitCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Create a starter config in the user config directory",
		RunE: func(cmd *cobra.Command, args []string) error {
			paths := app.DefaultPaths()
			if err := os.MkdirAll(paths.ConfigDir, 0o755); err != nil {
				return err
			}

			if _, err := os.Stat(paths.ConfigFile); err == nil {
				return os.ErrExist
			} else if !os.IsNotExist(err) {
				return err
			}

			if err := os.WriteFile(paths.ConfigFile, starterConfigContent(), 0o644); err != nil {
				return err
			}

			cmd.Printf("wrote starter config to %s\n", paths.ConfigFile)
			return nil
		},
	}
}
