package cli

import (
	"github.com/realtong/logi-cli/internal/app"
	"github.com/realtong/logi-cli/internal/config"
	"github.com/spf13/cobra"
)

func newValidateCmd() *cobra.Command {
	var configPath string

	cmd := &cobra.Command{
		Use:   "validate",
		Short: "Validate a logi-cli config file",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.LoadFile(configPath)
			if err != nil {
				return err
			}
			return config.Validate(cfg)
		},
	}

	cmd.Flags().StringVar(&configPath, "config", app.DefaultPaths().ConfigFile, "config file to validate")
	return cmd
}
