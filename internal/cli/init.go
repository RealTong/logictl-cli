package cli

import (
	"os"

	"github.com/realtong/logi-cli/internal/app"
	"github.com/spf13/cobra"
)

const starterConfig = `# logi-cli starter configuration
[daemon]
reload_on_change = true

[[devices]]
id = "mx-master-4"
match_vendor_id = 1133
match_product_id = 1234

  [devices.capabilities]
  thumb_button = "button_5"
  wheel_left = "hscroll_left"
  wheel_right = "hscroll_right"

[[actions]]
id = "close_tab"
type = "shortcut"
keys = ["cmd", "w"]

[[profiles]]
id = "chrome"
app_bundle_id = "com.google.Chrome"

  [[profiles.bindings]]
  device = "mx-master-4"
  trigger = "hold(thumb_button)+move(down)"
  action = "close_tab"
`

func newInitCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Create a starter config in the user config directory",
		RunE: func(cmd *cobra.Command, args []string) error {
			paths := app.DefaultPaths()
			if err := os.MkdirAll(paths.ConfigDir, 0o755); err != nil {
				return err
			}

			if err := os.WriteFile(paths.ConfigFile, []byte(starterConfig), 0o644); err != nil {
				return err
			}

			cmd.Printf("wrote starter config to %s\n", paths.ConfigFile)
			return nil
		},
	}
}
