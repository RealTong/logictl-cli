MX Master 4 raw HID report fixtures captured from a real Bluetooth device.

Capture source files:
- `/tmp/logictl-gesture-button-down.txt`
- `/tmp/logictl-gesture-button-hold-move-down.txt`

Hardware context:
- Path: `DevSrvsID:4295271338`
- VID:PID: `046d:b042`
- Product: `MX Master 4`
- Serial: `5625C585`
- Transport: `Bluetooth`

Known-good evidence in this task:
- Historical fixture files `thumb-button-down.txt` and `thumb-button-hold-move-down.txt`
  were originally captured for the same physical control and should now be read as
  `gesture_button` samples.

Known limits:
- These fixtures are only sufficient to confidently cover gesture-button down, hold, and hold-plus-move-down behavior.
- Other button states or gesture directions should not be inferred without more captures.
