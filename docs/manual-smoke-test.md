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
2. Press the thumb button once. Confirm `thumb_button_down` and `thumb_button_up` appear.
3. Hold the thumb button, move the mouse down, then release. Confirm `thumb_button_hold` and `hold(thumb_button)+move(down)` appear.

## Daemon Control Plane

1. Run `go run ./cmd/logi daemon run`
2. In another terminal, run `go run ./cmd/logi daemon status`
3. Confirm the output is `running`
4. In the second terminal, run `go run ./cmd/logi reload`
5. Confirm the output is `reload requested`

## LaunchAgent

1. Run `go run ./cmd/logi daemon start`
2. Run `go run ./cmd/logi daemon status`
3. Confirm the daemon reports `running`
4. Run `go run ./cmd/logi daemon restart`
5. Run `go run ./cmd/logi daemon stop`
