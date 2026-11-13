# cannon-ext-ssh

Cannon extension that provides a browser-based terminal in admin (xterm.js + local PTY).

## Requirements

- Cannon with admin extension WebSocket proxy and relaxed extension-admin CSP (see [cannon](https://github.com/rob121/cannon) `internal/admin/crud_extension_apps.go`).
- Linux or macOS (PTY).

## Permissions

Registers `terminal.use` (`cannon-ext-ssh.terminal.use`). Marked **dangerous** in the permission catalog.

## Build

```bash
go build -o cannon-ext-ssh .
```

Install the binary into your Cannon extensions directory and enable it in **Admin → Extensions**.

## Smoke test

1. Sign in as an admin with `terminal.use`.
2. Open **Extensions → SSH Terminal**.
3. Run `whoami` and confirm it matches the OS user running extensions.
4. Resize the window and run `echo ok`.
