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

1. `mkdir -p ./bin`
2. `go build -o ./bin/logi ./cmd/logi`
3. `./bin/logi init`
4. `./bin/logi validate`
5. `./bin/logi doctor`
6. `./bin/logi devices list`
7. `./bin/logi test event`

## Example Workflow

1. Build the current CLI with `go build -o ./bin/logi ./cmd/logi`
2. Install the stable background binary with `./bin/logi daemon install`
3. Start the foreground daemon with `./bin/logi daemon run`, or start the persistent LaunchAgent with `./bin/logi daemon start`
4. In another terminal, run `./bin/logi daemon status`
5. Reload config changes with `./bin/logi reload`

## Background Daemon Permissions

- Use a built binary such as `./bin/logi`; do not use `go run` for daemon lifecycle commands during development.
- `daemon install` copies the current built binary into the stable LaunchAgent path at `~/.config/logi-cli/state/logi-launchagent`.
- `daemon start` and `daemon restart` only use that installed background binary; they do not overwrite it.
- On macOS, `~/.config/logi-cli/state/logi-launchagent` needs its own `Input Monitoring` permission. If `daemon start` reports a permission error, add that exact path to `System Settings -> Privacy & Security -> Input Monitoring`, then retry.
- Rebuilding `./bin/logi` alone does not change the background daemon binary. Permissions usually only need to be revisited after you run `./bin/logi daemon install` again.
- macOS does not provide a normal scriptable way to auto-grant `Input Monitoring` or `Accessibility` to an arbitrary local CLI binary. `tccutil` can reset prompts, but it cannot grant access.

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
- `./bin/logi doctor`
- `./bin/logi devices list`

The manual validation checklist is in [docs/manual-smoke-test.md](docs/manual-smoke-test.md).
