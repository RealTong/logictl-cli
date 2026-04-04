# Manual Smoke Test

## First Run

1. `mkdir -p ./bin`
2. `go build -o ./bin/logictl ./cmd/logictl`
3. `./bin/logictl init`
4. `./bin/logictl validate`
5. Grant `Accessibility` and `Input Monitoring` if `./bin/logictl doctor` reports them missing.
6. `./bin/logictl doctor`
7. `./bin/logictl devices list`
8. `./bin/logictl test event --path '<mx-master-4-path>'`

## MX Master 4 Semantic Capture

1. Run `./bin/logictl test event --path '<mx-master-4-path>'`
2. Press the gesture button once. Confirm `gesture_button_down` and `gesture_button_up` appear.
3. Hold the gesture button, move the mouse down, then release. Confirm `gesture_button_hold` and `hold(gesture_button)+move(down)` appear.
4. Press the back and forward buttons. Confirm `back_button_down/up` and `forward_button_down/up` appear.
5. Spin the main wheel and thumb wheel. Confirm `wheel_up/down` and `thumb_wheel_left/right` appear.
6. Press the wheel mode-shift button. Confirm `mode_shift_button_press` plus either `wheel_mode_ratchet` or `wheel_mode_free_spin` appear.
7. Press the haptic panel. Confirm `haptic_panel_press` appears.

## Daemon Control Plane

1. Run `./bin/logictl daemon run`
2. In another terminal, run `./bin/logictl daemon status`
3. Confirm the output is `running`
4. In the second terminal, run `./bin/logictl reload`
5. Confirm the output is `reload requested`
6. With `devices.scroll.direction = "standard"`, confirm the main wheel and thumb wheel now move in the Windows-style direction while the daemon is running.
7. Flip `devices.scroll.direction` back to `"natural"`, run `./bin/logictl reload`, and confirm the native macOS direction returns.
8. Toggle `devices.scroll.smooth_scroll` between `true` and `false`, reload, and confirm scrolling switches between short smooth bursts and line-like ticks.

## LaunchAgent

1. Run `./bin/logictl daemon install`
2. Grant `Input Monitoring` to `~/.config/logictl/state/logictl-daemon` if macOS prompts for it or if `start` later reports a permission error.
3. Run `./bin/logictl daemon start`
4. Run `./bin/logictl daemon status`
5. Confirm the daemon reports `running`
6. Run `./bin/logictl daemon restart`
7. Run `./bin/logictl daemon stop`
