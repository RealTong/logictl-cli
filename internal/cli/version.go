package cli

import "github.com/spf13/cobra"

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the logictl version",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println("logictl development")
		},
	}
}
