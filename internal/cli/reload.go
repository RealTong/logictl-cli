package cli

import (
	"github.com/realtong/logictl-cli/internal/daemon"
	"github.com/spf13/cobra"
)

func newReloadCmd(app *daemon.App) *cobra.Command {
	return &cobra.Command{
		Use:   "reload",
		Short: "Request a config reload from the running daemon",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			status, err := app.Reload()
			if err != nil {
				return err
			}

			if status.Message == "" {
				cmd.Println("reload requested")
				return nil
			}

			cmd.Println(status.Message)
			return nil
		},
	}
}
