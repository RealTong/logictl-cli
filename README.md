# logi-cli

`logi-cli` is a macOS-first CLI and daemon for Logitech device customization, focused on replacing the input remapping parts of Logi Options+ for `MX Master 4`.

## Current Scope

- `MX Master 4` semantic event decoding
- TOML config at `~/.config/logi-cli/config.toml`
- Per-app rule matching with precedence
- Shortcut, system-action, and script execution
- `doctor`, `devices`, `test event`, `daemon`, and `reload` commands

## First Run

1. `go run ./cmd/logi init`
2. `go run ./cmd/logi validate`
3. `go run ./cmd/logi doctor`
4. `go run ./cmd/logi devices list`
5. `go run ./cmd/logi test event`

## Example Workflow

1. Start the foreground daemon with `go run ./cmd/logi daemon run`
2. In another terminal, run `go run ./cmd/logi daemon status`
3. Reload config changes with `go run ./cmd/logi reload`
4. Install the persistent LaunchAgent with `go run ./cmd/logi daemon start`

## Background Daemon Permissions

- `daemon run` runs in the current foreground process, so it uses the permissions already granted to your terminal or host app.
- `daemon start` installs a LaunchAgent and stages a stable background binary at `~/.config/logi-cli/state/logi-launchagent`.
- On macOS, that staged binary needs its own `Input Monitoring` permission. If `daemon start` reports a permission error, add `~/.config/logi-cli/state/logi-launchagent` to `System Settings -> Privacy & Security -> Input Monitoring`, then retry.

## MX Master 4 BLE Notes

- `MX Master 4` over Bluetooth Low Energy exposes a shared HID path that can seize the primary pointer when opened through generic `hidapi` path access.
- The daemon now uses native macOS `IOHIDManager` capture for this layout instead of opening the shared path directly, which avoids freezing mouse movement.

## Example Binding

The sample config maps:

- `thumb_button_down` to `mission_control`
- `hold(thumb_button)+move(down)` in `com.google.Chrome` to `cmd+w`

See [examples/config.toml](examples/config.toml) for the current example file.

## Verification

- `go test ./...`
- `go run ./cmd/logi doctor`
- `go run ./cmd/logi devices list`

The manual validation checklist is in [docs/manual-smoke-test.md](docs/manual-smoke-test.md).
