# logi-cli Design

Date: 2026-04-02
Status: Approved for planning
Primary target: macOS
Initial device target: Logitech MX Master 4

## Summary

`logi-cli` is a lightweight Logitech device behavior engine for macOS built around `go-hid`. It is intended to replace the core input customization use cases of Logi Options+ without account requirements, cloud features, or other non-essential product surface area.

V1 focuses on MX Master 4 and provides:

- A low-overhead background daemon that listens to device events
- A CLI for configuration, diagnostics, and daemon control
- Per-application behavior switching
- Gesture and combination actions
- Action execution through shortcuts, selected system actions, and scripts

V1 does not attempt to write complex mappings directly into device firmware. The system is explicitly host-side and relies on macOS input and accessibility APIs for app-aware behavior and action injection.

## Goals

- Replace the high-value customization flows from Logi Options+ for MX Master 4 on macOS
- Keep runtime behavior stable, predictable, and low-overhead
- Make configuration easy to version, inspect, validate, and recover
- Build abstractions that can later support newer or older Logitech devices without redesigning the rule engine

## Non-Goals

- GUI configuration
- Logitech account integration
- Logi Flow
- Cloud sync
- Broad multi-device support in V1
- Direct device-side persistence for complex behavior rules
- Simultaneous first-class support for Windows and Linux in V1

## Constraints and Assumptions

- `go-hid` is used as the HID communication layer, not as a complete device-behavior platform
- Device-specific report parsing and capability detection will still need project-specific implementation
- Per-app logic, gesture interpretation, and action injection are host-side concerns and are not handled by HID alone
- On macOS, several behaviors will depend on system permissions such as Accessibility and possibly Input Monitoring or Automation
- Product IDs and raw report behavior for MX Master 4 must be discovered from real hardware during implementation; they are not hardcoded in this design

## Recommended Approach

V1 uses a host-side event engine with a `CLI + daemon + declarative config` architecture.

Alternative approaches were considered:

1. Host-side event engine
2. Device-protocol-first configuration writer
3. Hybrid device-side plus host-side system

The recommended approach is the host-side event engine because it best supports per-app rules, gestures, extensibility, and controlled runtime behavior. The device-protocol-first and hybrid approaches both add protocol risk without solving the core V1 use cases better.

## System Architecture

The system is split into five layers.

### 1. HID Layer

Responsibilities:

- Discover attached HID devices
- Match supported Logitech devices
- Open and maintain HID connections
- Read and write raw reports
- Expose device metadata and capability descriptors

This layer only answers "what did the device send" and "how do we talk to this device." It does not interpret user intent.

### 2. Event Pipeline

Responsibilities:

- Convert raw HID reports into normalized internal events
- Handle debouncing and state transitions
- Recognize button presses, holds, releases, wheel actions, and gesture candidates
- Maintain small per-device state machines

Example normalized events:

- `gesture_button_down`
- `gesture_button_hold`
- `horizontal_wheel_left`
- `gesture(gesture_button, down)`

This layer owns timing thresholds and motion interpretation rules.

### 3. Rule Engine

Responsibilities:

- Evaluate normalized events against configuration
- Resolve the active profile using device identity and active application
- Select the best matching binding
- Produce an action execution request

This layer is the core product logic. Rules are written in terms of semantic device capabilities rather than raw bytes.

### 4. Action Executor

Responsibilities:

- Execute keyboard shortcuts
- Execute a bounded set of system actions
- Run user-defined scripts with timeout and isolation controls

The executor consumes a stable internal action model and does not know where the rule came from.

### 5. Control Plane

Responsibilities:

- Manage configuration lifecycle
- Start, stop, inspect, and reload the daemon
- Provide diagnostics and event inspection
- Surface health and error state through CLI commands

The control plane is the operator interface, not the runtime data path.

## Runtime Shape

The runtime model is:

`HID event -> normalize -> match rule -> execute action`

The daemon should remain single-process and single-purpose. It should avoid background scanning loops, GUI dependencies, persistent databases, or remote control surfaces beyond a local control channel for the CLI.

## Configuration Model

Primary config file:

- `~/.config/logi-cli/config.toml`

Supporting files may live under:

- `~/.config/logi-cli/scripts/`
- `~/.config/logi-cli/logs/`
- `~/.config/logi-cli/state/`

TOML is the primary configuration format because it is both machine-friendly and realistic for direct user editing. The CLI is allowed to modify the same config model rather than inventing a second source of truth.

### Core Concepts

#### Device

Defines how a concrete hardware target is matched and how its raw inputs are named in semantic terms.

Suggested fields:

- logical ID
- vendor ID
- product ID
- transport hints if needed
- capability aliases such as `gesture_button`, `gesture_button`, `wheel_left`, `wheel_right`

The rule engine should depend on these aliases, not on report bytes.

#### Action

Defines a reusable action target.

Initial action types:

- `shortcut`
- `system`
- `script`

Examples:

- Shortcut: `cmd+w`
- System action: `mission_control`
- Script: `~/.config/logi-cli/scripts/close-tab.sh`

#### Profile

Defines a scope in which bindings apply.

Initial scopes:

- global
- app-specific by macOS bundle ID, such as `com.google.Chrome`

Application matching should use bundle IDs instead of window titles or process names whenever possible.

#### Binding

Maps a trigger to an action.

Initial trigger categories:

- press
- release
- hold
- horizontal wheel left/right
- combinations
- gesture forms such as `hold(button) + move(direction)`

### Rule Precedence

Bindings are resolved in this order:

1. Exact device + exact app
2. Exact device + global
3. Any device + exact app
4. Any device + global

If multiple bindings still match, an explicit numeric priority may break ties. If no priority is set, the configuration should remain deterministic and reject ambiguous collisions during validation instead of silently choosing one.

### Example Shape

```toml
[daemon]
reload_on_change = true

[[devices]]
id = "mx-master-4"
match_vendor_id = 1133
match_product_id = 0

[devices.capabilities]
gesture_button = "button_5"
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
trigger = "hold(gesture_button)+move(down)"
action = "close_tab"
```

Notes:

- `match_product_id = 0` above is a placeholder shape only, not a real MX Master 4 identifier
- The configuration grammar should be validated strictly before activation

## CLI Surface

V1 should include at least the following commands:

- `logi init`
- `logi validate`
- `logi reload`
- `logi doctor`
- `logi daemon run`
- `logi daemon start`
- `logi daemon stop`
- `logi daemon restart`
- `logi daemon status`
- `logi devices list`
- `logi devices inspect`
- `logi test event`

### Command Intent

- `init`: create the base directory and starter config
- `validate`: parse and verify configuration without applying it
- `reload`: request live config reload from the running daemon
- `doctor`: verify permissions, daemon state, config health, and device visibility
- `devices list`: enumerate visible HID devices
- `devices inspect`: print device metadata and optionally raw report information
- `test event`: print raw or normalized input events to support rule authoring

## Daemon Behavior

The daemon is required because per-app switching and gestures are stateful runtime concerns. A pure one-shot CLI cannot provide this.

### Lifecycle

- Starts in foreground for development and diagnosis
- Can also run under a background service manager
- Loads config on start
- Applies new config only after successful validation
- Exposes a local control interface for status and reload operations

### Performance Rules

- Use blocking or event-driven HID reads, not busy loops
- Do not poll the active application constantly; resolve it on the path where a rule needs it
- Avoid high-frequency global system inspection
- Keep memory bounded and independent of rule count where practical
- Avoid optional subsystems that stay resident without contributing to the runtime path

### Fault Containment

- Script execution must have timeouts
- Action execution failures must not crash the event loop
- Device disconnects should trigger controlled reconnect logic
- Invalid configs should never partially apply

## macOS Integration

The following host capabilities are required or likely required:

- HID device access
- Accessibility permission for input event injection
- Possibly Input Monitoring or Automation depending on exact action implementations

The product must treat permission state as a first-class runtime concern. Missing permissions should be surfaced clearly through `logi doctor` and daemon status output.

## Stability and Diagnostics

This project exists partly because the reference software is considered too heavy and unreliable. Stability is therefore a product requirement, not only an implementation detail.

### Diagnostic Requirements

- Structured logs with conservative default verbosity
- Explicit device connection state
- Clear action execution errors
- Visibility into matched rules during test or debug mode
- Permission status inspection
- Config validation before activation

### Reliability Principles

- Keep the runtime pipeline narrow and observable
- Separate protocol parsing from rule matching
- Separate rule matching from action execution
- Reject ambiguous or malformed configuration early
- Prefer deterministic failure over hidden fallback behavior

## Compatibility Strategy

V1 targets MX Master 4, but the architecture is designed for future device support.

### Compatibility Rules

- Rules are written against semantic capabilities, not raw device bytes
- Each supported device gets its own capability map and report decoder
- Unknown Logitech devices may be detected in a limited observation mode before full support exists
- The rule engine must remain independent from model-specific parsing logic

This creates a path to support:

- Future devices such as a hypothetical MX Master 5
- Older Logitech devices where event patterns differ
- Additional keyboards or mice such as MX Keys in later phases

## V1 Scope

V1 explicitly includes:

- macOS support
- MX Master 4 support
- CLI plus daemon runtime model
- Global and app-specific bindings
- Basic gesture and combination support
- Shortcut, selected system, and script actions
- Validation and diagnostics commands

V1 explicitly excludes:

- GUI
- account or cloud features
- Flow
- broad device catalog support
- direct device-side persistence for complex mappings
- Windows and Linux parity

## Risks

- MX Master 4 report semantics may require device-specific reverse engineering or capture
- macOS permission behavior can vary by action type and OS version
- Application detection and action injection are platform-specific and must be isolated from the cross-device core
- Gesture recognition thresholds may need tuning from real-world use rather than synthetic assumptions

## Implementation Direction

The implementation plan should start with:

1. Repository bootstrap and package layout
2. Minimal daemon skeleton and CLI shell
3. Device enumeration and inspection commands
4. Config parser and validator
5. Normalized event model
6. MX Master 4 adapter
7. Rule engine
8. Action execution layer
9. Diagnostics and permission checks

This order prioritizes observability and hardware understanding before complex user-facing behavior.
