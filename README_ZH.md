# logictl

[English](README.md) | 简体中文

`logictl` 是一个以 macOS 为优先平台的 Logitech 设备定制工具，采用 `CLI + daemon` 形式，当前重点是替代 `MX Master 4` 在 Logi Options+ 里的输入映射能力。

## 当前范围

- 已接入 `MX Master 4` 的语义事件解码，包括：
  - 左键 / 右键 / 中键
  - 后退 / 前进
  - 手势按钮
  - 主滚轮
  - 拇指轮
  - 滚轮模式切换按钮
  - 触觉面板按压
- TOML 配置文件路径：`~/.config/logictl/config.toml`
- 支持按应用匹配规则，并带有优先级
- 支持快捷键、系统动作和脚本执行
- 提供 `doctor`、`devices`、`test event`、`daemon`、`reload` 等命令

## 首次运行

1. `mkdir -p ./bin`
2. `go build -o ./bin/logictl ./cmd/logictl`
3. `./bin/logictl init`
4. `./bin/logictl validate`
5. `./bin/logictl doctor`
6. `./bin/logictl devices list`
7. `./bin/logictl test event`

## 示例工作流

1. 构建当前 CLI：`go build -o ./bin/logictl ./cmd/logictl`
2. 安装稳定的后台二进制：`./bin/logictl daemon install`
3. 以前台方式运行 daemon：`./bin/logictl daemon run`，或者启动持久化 LaunchAgent：`./bin/logictl daemon start`
4. 在另一个终端中执行：`./bin/logictl daemon status`
5. 修改 `~/.config/logictl/config.toml`
6. 手动执行：`./bin/logictl reload`

## 后台 Daemon 权限

- 开发期间请使用稳定构建产物，例如 `./bin/logictl`，不要用 `go run` 管理 daemon 生命周期。
- `daemon install` 会把当前构建好的二进制复制到稳定的 LaunchAgent 路径：`~/.config/logictl/state/logictl-daemon`
- `daemon start` 和 `daemon restart` 只会使用这个已安装的后台二进制，不会在每次启动时覆盖它。
- 在 macOS 上，`~/.config/logictl/state/logictl-daemon` 需要单独授予 `Input Monitoring` 权限。如果 `daemon start` 报权限错误，请到 `System Settings -> Privacy & Security -> Input Monitoring` 中添加这个精确路径，然后重试。
- 仅重新构建 `./bin/logictl` 不会改变后台 daemon 二进制。通常只有在再次执行 `./bin/logictl daemon install` 之后，才可能需要重新确认权限。
- macOS 没有常规的脚本化方式去直接授予任意本地 CLI 程序 `Input Monitoring` 或 `Accessibility` 权限。`tccutil` 只能重置提示，不能直接授权。
- 配置修改采用手动生效策略。编辑 `~/.config/logictl/config.toml` 后，请执行 `./bin/logictl reload`。

## MX Master 4 BLE 说明

- `MX Master 4` 在 Bluetooth Low Energy 模式下会暴露共享 HID path；如果通过通用 `hidapi` path 直接打开，有可能抢占主鼠标接口。
- 现在 daemon 已经改为在这种布局下使用原生 macOS `IOHIDManager` 采集，而不是直接打开共享 path，因此可以避免鼠标失控。

## 示例绑定

示例配置当前包含：

- `gesture_button_down` 在未形成方向手势时，于松开后触发 `mission_control`
- `hold(gesture_button)+move(up)` 映射到 `app_expose`
- `hold(gesture_button)+move(left)` 映射到 `previous_desktop`
- `hold(gesture_button)+move(right)` 映射到 `next_desktop`
- `back_button_down` 映射到 `cmd+[`
- `forward_button_down` 映射到 `cmd+]`
- `haptic_panel_press` 映射到 `launchpad`
- 在 `com.google.Chrome` 中，`hold(gesture_button)+move(down)` 映射到 `cmd+w`
- `devices.scroll.direction = "standard"`
- `devices.scroll.smooth_scroll = true`

`gesture_button_down` 绑定现在会延迟到 `gesture_button_up` 时再决定是否执行；如果按住期间识别到了 `hold(gesture_button)+move(...)`，那么释放时只会执行方向手势对应的动作，不会再执行默认的 tap 动作。

`wheel_up`、`wheel_down`、`thumb_wheel_left`、`thumb_wheel_right`、`mode_shift_button_press`、`wheel_mode_ratchet` 和 `wheel_mode_free_spin` 也都已经暴露为可绑定的语义触发器。

对于 `MX Master 4`，`direction` 和 `smooth_scroll` 已经是运行中的 daemon 配置项：

- `direction = "standard"`：将主滚轮和拇指轮改成接近 Windows 风格的方向
- `direction = "natural"`：保持原生 macOS 的方向
- `smooth_scroll = true`：把匹配到的滚轮 tick 重写为短的 pixel-scroll burst
- `smooth_scroll = false`：把匹配到的滚轮 tick 重写为单个 line-scroll event

滚轮重写只会在 daemon 运行时启用，并且只作用于 daemon 能够和最近 `MX Master 4` HID 滚轮手势关联上的 scroll event，因此不会影响其他鼠标或触控板的原生滚动路径。

当前示例配置见 [examples/config.toml](examples/config.toml)。

## 验证

- `go test ./...`
- `./bin/logictl doctor`
- `./bin/logictl devices list`

`./bin/logictl doctor` 会额外显示当前正在检查的配置文件路径。

手动验证清单见 [docs/manual-smoke-test.md](docs/manual-smoke-test.md)。
