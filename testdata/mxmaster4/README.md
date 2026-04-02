MX Master 4 raw HID report fixtures captured from a real Bluetooth device.

Capture source files:
- `/tmp/logi-thumb-button-down.txt`
- `/tmp/logi-thumb-button-hold-move-down.txt`

Hardware context:
- Path: `DevSrvsID:4295271338`
- VID:PID: `046d:b042`
- Product: `MX Master 4`
- Serial: `5625C585`
- Transport: `Bluetooth`

Known-good evidence in this task:
- `thumb-button-down.txt` contains a thumb-button press sequence.
- `thumb-button-hold-move-down.txt` contains a thumb-button hold with downward movement and a later release.

Known limits:
- These fixtures are only sufficient to confidently cover thumb-button down, hold, and hold-plus-move-down behavior.
- Other button states or gesture directions should not be inferred without more captures.
