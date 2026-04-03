# logi-cli

`logi-cli` is a macOS-first CLI and daemon for Logitech device customization, focused on replacing the input remapping parts of Logi Options+ for `MX Master 4`.

## Current Scope

- `MX Master 4` semantic event decoding for:
  - left / right / middle
  - back / forward
  - gesture button
  - main wheel
  - thumb wheel
  - wheel mode-shift button
  - haptic panel press
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
- During development, repeated `go run` rebuilds can cause macOS to treat the staged daemon binary as a new untrusted executable. If permissions appear to "disappear" after a rebuild, re-enable `Input Monitoring` for `~/.config/logi-cli/state/logi-launchagent`, or prefer a stable `go build -o ./bin/logi ./cmd/logi` workflow.

## MX Master 4 BLE Notes

- `MX Master 4` over Bluetooth Low Energy exposes a shared HID path that can seize the primary pointer when opened through generic `hidapi` path access.
- The daemon now uses native macOS `IOHIDManager` capture for this layout instead of opening the shared path directly, which avoids freezing mouse movement.

## Example Binding

The sample config includes:

- `gesture_button_down` to `mission_control`
- `hold(gesture_button)+move(up)` to `app_expose`
- `hold(gesture_button)+move(left)` to `previous_desktop`
- `hold(gesture_button)+move(right)` to `next_desktop`
- `back_button_down` to `cmd+[`
- `forward_button_down` to `cmd+]`
- `haptic_panel_press` to `launchpad` as a lightweight default
- `hold(gesture_button)+move(down)` in `com.google.Chrome` to `cmd+w`
- `devices.scroll.direction = "standard"`
- `devices.scroll.smooth_scroll = true`

`wheel_up`, `wheel_down`, `thumb_wheel_left`, `thumb_wheel_right`, `mode_shift_button_press`, `wheel_mode_ratchet`, and `wheel_mode_free_spin` are also exposed as semantic triggers for custom bindings.

`direction` and `smooth_scroll` are now live daemon settings for `MX Master 4`:

- `direction = "standard"` inverts vertical and horizontal wheel output relative to the native macOS session event
- `direction = "natural"` leaves the native direction unchanged
- `smooth_scroll = true` rewrites matched wheel ticks into short pixel-scroll bursts
- `smooth_scroll = false` rewrites matched wheel ticks as single line-scroll events

The rewrite path only activates while the daemon is running and only for scroll events that the daemon can correlate with recent `MX Master 4` HID wheel gestures, which keeps other mice and the trackpad on the native path.

See [examples/config.toml](examples/config.toml) for the current example file.

## Verification

- `go test ./...`
- `go run ./cmd/logi doctor`
- `go run ./cmd/logi devices list`

The manual validation checklist is in [docs/manual-smoke-test.md](docs/manual-smoke-test.md).
