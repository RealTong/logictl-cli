# logictl

English | [简体中文](README_ZH.md)

`logictl` is a macOS-first CLI and daemon for Logitech device customization, focused on replacing the input remapping parts of Logi Options+ for `MX Master 4`.

## Current Scope

- `MX Master 4` semantic event decoding for:
  - left / right / middle
  - back / forward
  - gesture button
  - main wheel
  - thumb wheel
  - wheel mode-shift button
  - haptic panel press
- TOML config at `~/.config/logictl/config.toml`
- Per-app rule matching with precedence
- Shortcut, system-action, and script execution
- `doctor`, `devices`, `test event`, `daemon`, and `reload` commands

## First Run

1. `mkdir -p ./bin`
2. `go build -o ./bin/logictl ./cmd/logictl`
3. `./bin/logictl init`
4. `./bin/logictl validate`
5. `./bin/logictl doctor`
6. `./bin/logictl devices list`
7. `./bin/logictl test event`

## Example Workflow

1. Build the current CLI with `go build -o ./bin/logictl ./cmd/logictl`
2. Install the stable background binary with `./bin/logictl daemon install`
3. Start the foreground daemon with `./bin/logictl daemon run`, or start the persistent LaunchAgent with `./bin/logictl daemon start`
4. In another terminal, run `./bin/logictl daemon status`
5. Edit `~/.config/logictl/config.toml`
6. Apply config changes explicitly with `./bin/logictl reload`

## Background Daemon Permissions

- Use a built binary such as `./bin/logictl`; do not use `go run` for daemon lifecycle commands during development.
- `daemon install` copies the current built binary into the stable LaunchAgent path at `~/.config/logictl/state/logictl-daemon`.
- `daemon start` and `daemon restart` only use that installed background binary; they do not overwrite it.
- On macOS, `~/.config/logictl/state/logictl-daemon` needs its own `Input Monitoring` permission. If `daemon start` reports a permission error, add that exact path to `System Settings -> Privacy & Security -> Input Monitoring`, then retry.
- Rebuilding `./bin/logictl` alone does not change the background daemon binary. Permissions usually only need to be revisited after you run `./bin/logictl daemon install` again.
- macOS does not provide a normal scriptable way to auto-grant `Input Monitoring` or `Accessibility` to an arbitrary local CLI binary. `tccutil` can reset prompts, but it cannot grant access.
- Config changes are manual by design. After editing `~/.config/logictl/config.toml`, run `./bin/logictl reload`.

## MX Master 4 BLE Notes

- `MX Master 4` over Bluetooth Low Energy exposes a shared HID path that can seize the primary pointer when opened through generic `hidapi` path access.
- The daemon now uses native macOS `IOHIDManager` capture for this layout instead of opening the shared path directly, which avoids freezing mouse movement.

## Example Binding

The sample config includes:

- `gesture_button_down` to `mission_control` when the gesture button is released without a directional gesture
- `hold(gesture_button)+move(up)` to `app_expose`
- `hold(gesture_button)+move(left)` to `previous_desktop`
- `hold(gesture_button)+move(right)` to `next_desktop`
- `back_button_down` to `cmd+[`
- `forward_button_down` to `cmd+]`
- `haptic_panel_press` to `launchpad` as a lightweight default
- `hold(gesture_button)+move(down)` in `com.google.Chrome` to `cmd+w`
- `devices.scroll.direction = "standard"`
- `devices.scroll.smooth_scroll = true`

`gesture_button_down` bindings are deferred until `gesture_button_up`. If a directional `hold(gesture_button)+move(...)` gesture is recognized while the button is held, that directional action replaces the deferred tap action and fires on release instead.

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
- `./bin/logictl doctor`
- `./bin/logictl devices list`

`./bin/logictl doctor` also shows the exact config path currently being inspected.

The manual validation checklist is in [docs/manual-smoke-test.md](docs/manual-smoke-test.md).
