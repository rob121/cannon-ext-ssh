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
- **failed to start shell / permission denied**: The OS is blocking fork/exec of a shell (SELinux, AppArmor, or strict systemd `SystemCallFilter`). Check `journalctl` / `ausearch -m avc -ts recent`. Run Cannon under a user allowed to execute `/bin/bash`, or set `CANNON_EXT_SSH_SHELL=/path/to/bash` in the **Cannon** service environment (extensions inherit it).
- **nologin user**: use a dedicated Linux user with login shell ` /bin/bash`, not `www-data` with `nologin`.
- **Docker**: the container needs a TTY-capable environment (`/dev/ptmx`); many minimal images work if bash is installed and the process user can run it.
