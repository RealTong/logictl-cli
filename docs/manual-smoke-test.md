# Manual Smoke Test

## First Run

1. `go run ./cmd/logi init`
2. `go run ./cmd/logi validate`
3. Grant `Accessibility` and `Input Monitoring` if `go run ./cmd/logi doctor` reports them missing.
4. `go run ./cmd/logi doctor`
5. `go run ./cmd/logi devices list`
6. `go run ./cmd/logi test event --path '<mx-master-4-path>'`

## MX Master 4 Semantic Capture

1. Run `go run ./cmd/logi test event --path '<mx-master-4-path>'`
2. Press the gesture button once. Confirm `gesture_button_down` and `gesture_button_up` appear.
3. Hold the gesture button, move the mouse down, then release. Confirm `gesture_button_hold` and `hold(gesture_button)+move(down)` appear.
4. Press the back and forward buttons. Confirm `back_button_down/up` and `forward_button_down/up` appear.
5. Spin the main wheel and thumb wheel. Confirm `wheel_up/down` and `thumb_wheel_left/right` appear.
6. Press the wheel mode-shift button. Confirm `mode_shift_button_press` plus either `wheel_mode_ratchet` or `wheel_mode_free_spin` appear.
7. Press the haptic panel. Confirm `haptic_panel_press` appears.

## Daemon Control Plane

1. Run `go run ./cmd/logi daemon run`
2. In another terminal, run `go run ./cmd/logi daemon status`
3. Confirm the output is `running`
4. In the second terminal, run `go run ./cmd/logi reload`
5. Confirm the output is `reload requested`
6. With `devices.scroll.direction = "standard"`, confirm the main wheel and thumb wheel now move in the Windows-style direction while the daemon is running.
7. Flip `devices.scroll.direction` back to `"natural"`, run `go run ./cmd/logi reload`, and confirm the native macOS direction returns.
8. Toggle `devices.scroll.smooth_scroll` between `true` and `false`, reload, and confirm scrolling switches between short smooth bursts and line-like ticks.

## LaunchAgent

1. Run `go run ./cmd/logi daemon start`
2. Run `go run ./cmd/logi daemon status`
3. Confirm the daemon reports `running`
4. Run `go run ./cmd/logi daemon restart`
5. Run `go run ./cmd/logi daemon stop`
