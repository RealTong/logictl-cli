package macos

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"

	appcore "github.com/realtong/logictl-cli/internal/app"
)

const launchAgentLabel = "io.realtong.logictl"

var runLaunchctl = func(ctx context.Context, args ...string) error {
	return exec.CommandContext(ctx, "launchctl", args...).Run()
}

var currentLaunchctlUID = func() (string, error) {
	current, err := user.Current()
	if err != nil {
		return "", err
	}
	return current.Uid, nil
}

func InstallLaunchAgent(paths appcore.Paths, binary string) error {
	if err := os.MkdirAll(filepath.Dir(paths.PlistFile), 0o755); err != nil {
		return err
	}
	if err := os.MkdirAll(paths.LogDir, 0o755); err != nil {
		return err
	}

	return os.WriteFile(paths.PlistFile, []byte(renderLaunchAgentPlist(paths, binary)), 0o644)
}

func StartLaunchAgent(ctx context.Context, paths appcore.Paths, binary string) error {
	if err := InstallLaunchAgent(paths, binary); err != nil {
		return err
	}

	domain, err := launchctlDomain()
	if err != nil {
		return err
	}

	_ = runLaunchctl(ctx, "bootout", domain, paths.PlistFile)
	if err := runLaunchctl(ctx, "bootstrap", domain, paths.PlistFile); err != nil {
		return err
	}
	return runLaunchctl(ctx, "kickstart", "-k", domain+"/"+launchAgentLabel)
}

func StopLaunchAgent(ctx context.Context, paths appcore.Paths) error {
	domain, err := launchctlDomain()
	if err != nil {
		return err
	}

	return runLaunchctl(ctx, "bootout", domain, paths.PlistFile)
}

func RestartLaunchAgent(ctx context.Context, paths appcore.Paths, binary string) error {
	if err := InstallLaunchAgent(paths, binary); err != nil {
		return err
	}

	domain, err := launchctlDomain()
	if err != nil {
		return err
	}

	if err := runLaunchctl(ctx, "bootout", domain, paths.PlistFile); err != nil {
		return err
	}
	if err := runLaunchctl(ctx, "bootstrap", domain, paths.PlistFile); err != nil {
		return err
	}
	return runLaunchctl(ctx, "kickstart", "-k", domain+"/"+launchAgentLabel)
}

func launchctlDomain() (string, error) {
	uid, err := currentLaunchctlUID()
	if err != nil {
		return "", err
	}
	return "gui/" + uid, nil
}

func renderLaunchAgentPlist(paths appcore.Paths, binary string) string {
	return strings.TrimSpace(fmt.Sprintf(`
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>Label</key>
	<string>%s</string>
	<key>ProgramArguments</key>
	<array>
		<string>%s</string>
		<string>daemon</string>
		<string>run</string>
	</array>
	<key>RunAtLoad</key>
	<true/>
	<key>KeepAlive</key>
	<true/>
	<key>WorkingDirectory</key>
	<string>%s</string>
	<key>EnvironmentVariables</key>
	<dict>
		<key>LOGICTL_CONFIG_FILE</key>
		<string>%s</string>
		<key>LOGICTL_SOCKET_FILE</key>
		<string>%s</string>
	</dict>
	<key>StandardOutPath</key>
	<string>%s</string>
	<key>StandardErrorPath</key>
	<string>%s</string>
</dict>
</plist>
`, launchAgentLabel, binary, paths.ConfigDir, paths.ConfigFile, paths.SocketFile, filepath.Join(paths.LogDir, "daemon.stdout.log"), filepath.Join(paths.LogDir, "daemon.stderr.log")))
}
