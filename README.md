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
