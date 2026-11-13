# Web SSH Terminal

This extension opens a **local shell** in Cannon admin. Commands run as the **same OS user** that runs Cannon extensions (not site visitors).

## Security

- Grant **Use web terminal** only to trusted administrators.
- Anyone with permission has the equivalent of interactive SSH access on the host.
- There is no command logging in this version.

## Usage

1. Ensure Cannon includes the admin WebSocket proxy (0.3.19+).
2. Install and start `cannon-ext-ssh`.
3. Open **Extensions → SSH Terminal** in admin.
4. The terminal connects automatically; resize the browser window as needed.

## Troubleshooting

- **Forbidden**: your user needs `cannon-ext-ssh.terminal.use` (or administrator `*`).
- **Session closed immediately**: check that the extension is running and Cannon can reach its socket.
- **failed to start shell**: Cannon/extensions must run as a Linux user with a real shell (`/bin/bash` or `/bin/sh`). Service accounts using `nologin` cannot open a PTY — use a dedicated user with a login shell or set `SHELL=/bin/bash` in the process environment.
- **Docker**: the container needs a TTY-capable environment (`/dev/ptmx`); many minimal images work if bash is installed and the process user can run it.
