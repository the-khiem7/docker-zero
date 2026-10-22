# Docker Zero

A small Windows tray app that sets up and runs Docker Engine inside a WSL2 distro, with Portainer CE as the GUI, for people on business machines who cannot use Docker Desktop and have no admin rights.

<p align="center">
  <img src="docs/images/popup.png" alt="Tray popup with the WSL, distro, Docker Engine and Portainer status and the Start, Open Portainer, Stop distro and Shut down WSL buttons" width="360">
</p>

## What it does

| | |
|---|---|
| Day 1 | A setup window installs an `Ubuntu-24.04` WSL2 distro, Docker Engine and Portainer CE on `https://localhost:9443`, creates the Portainer admin and adds a `docker` command to your PATH. |
| Day 2 | A tray icon with a popup starts and stops the distro, opens Portainer and shows live status. An optional corner box shows the same status on the desktop. |
| Settings | Auto start at sign-in, Desktop and Start Menu shortcuts, how the distro is kept running, which browser opens Portainer, and a locally trusted certificate. |
| Maintenance | Update Portainer to the latest long-term-support image, reset the admin password, or remove everything. |
| Take over | If you already set up Docker Engine and Portainer by hand with the [manual guide](docs/guides/docker-engine-wsl/README.md), the app replaces your launchers with its own and keeps the distro, containers and data. |

Everything runs as a standard user: no UAC prompt, no Windows service, nothing under `Program Files`. Containers, images and volumes are managed in Portainer; the app only handles the system around them.

## Requirements

- Windows 10 or 11 with the WSL Windows feature already enabled. Enabling it needs admin rights, so the app only detects it; if it is missing, ask your IT team to enable WSL.
- Internet access to download the Ubuntu distro, Docker Engine and the Portainer image, and a few GB of free disk space for the distro.

## Install

1. Download `docker-zero.exe` from the [latest release](../../releases/latest).
2. Run it. The exe is not code-signed, so Windows SmartScreen may ask you to confirm the first time.
3. The setup window opens by itself. Choose a Portainer admin username and a password of at least 12 characters, then click **Start setup**. It takes several minutes, mostly for downloads.

<p align="center">
  <img src="docs/images/setup.png" alt="Setup window with the Portainer admin account form and the list of setup steps" width="400">
  <img src="docs/images/setup-done.png" alt="Setup window with every step done and an Open Portainer button" width="400">
</p>

When setup is done, open a new terminal to use the `docker` command, and click **Open Portainer** to sign in.

The exe can stay where you downloaded it. Turning on auto start or a shortcut copies it to `%LOCALAPPDATA%\Programs\docker-zero` first, so the shortcut keeps working after the Downloads folder is cleaned.

## Use

Click the tray icon to open the popup. The dot on the icon shows the overall state: grey for stopped, amber for partly running, green for running, red for a failure.

| Button | Effect |
|---|---|
| Start | Starts the distro, Docker Engine and Portainer. |
| Open Portainer | Opens `https://localhost:9443` in the chosen browser. |
| Stop distro | Stops `Ubuntu-24.04` and every container in it. Other WSL distros keep running. |
| Shut down WSL | Runs `wsl --shutdown`, which stops every WSL distro of your account. The app asks first. |

The gear icon opens the settings.

<p align="center">
  <img src="docs/images/settings.png" alt="Settings window with auto start, shortcuts and the options to keep the distro running" width="400">
  <img src="docs/images/maintenance.png" alt="Settings window with the certificate, Portainer maintenance and danger zone sections" width="400">
</p>

### Keep the distro running

WSL stops an idle distro, and Portainer with it. Pick one of three modes in the settings:

| Mode | How | Scope |
|---|---|---|
| Keep running while Docker Zero runs (default) | A hidden session holds the distro up | Only this app |
| Never stop idle distros | Sets `instanceIdleTimeout=-1` in `%USERPROFILE%\.wslconfig` | Every WSL2 distro of your account, after WSL restarts |
| WSL default | Nothing | WSL stops the distro when nothing uses it |

### Certificate

Portainer uses a certificate created inside the distro for `localhost`. Until you trust it, the browser shows a warning. **Install (trust)** adds it to the trusted root certificates of your Windows account only; Windows asks you to confirm. The private key never leaves the distro. **Regenerate** creates a new certificate and restarts Portainer.

Trusting a certificate is a security decision. This one cannot sign other certificates (`CA:FALSE`) and is valid only for `localhost`, `127.0.0.1` and `::1`, so its key could at most impersonate those local addresses; the key is readable only by root inside your distro. On a managed device, security tooling may still report a change to the trusted roots.

### Maintenance

- **Update Portainer** pulls the latest `lts` image and recreates the container with the same data volume. If the image is already current, nothing restarts.
- **Reset admin password** replaces the password of the existing admin account. Portainer stops for a few seconds. Docker Zero never stores the password.

### Uninstall

**Uninstall all...** in the danger zone removes, in order: an optional backup of the distro to a `.tar` file you choose, the trusted certificate, the `Ubuntu-24.04` distro with every image, container and volume in it, auto start, shortcuts, the `docker` command and its PATH entry, the `.wslconfig` line the app added, and the settings. The app's own folder is deleted after you click **Quit**. Other WSL distros are not touched.

## Build from source

Requirements: Go 1.26 or newer, Node.js with npm, and the Wails CLI `v3.0.0-beta.30`:

```sh
go install github.com/wailsapp/wails/v3/cmd/wails3@v3.0.0-beta.30
```

| Task | Command |
|---|---|
| Build `bin/docker-zero.exe` | `wails3 task build` |
| Run with live reload | `wails3 task dev` |
| Unit tests | `wails3 task test` |
| `go vet`, `golangci-lint`, Svelte type check | `wails3 task lint` |
| Read-only status probe against this machine | `go test -tags integration ./internal/status` |

Wails build assets (icon, manifest, version info) live in `packaging/` instead of the Wails default `build/`, so use `wails3 task build` rather than `wails3 build`.

<details>
<summary>Layout</summary>

| Path | Contents |
|---|---|
| `main.go` | App wiring: tray, windows, services |
| `internal/wsl` | Runs `wsl.exe`, reads distro state, sends scripts on stdin |
| `internal/status` | Probes WSL, the distro, Docker Engine and Portainer |
| `internal/control` | Start, stop and the keepalive session |
| `internal/setup` | Setup wizard, maintenance and uninstall |
| `internal/adopt` | Takes over a setup made with the manual guide |
| `internal/settings` | Auto start, shortcuts, `.wslconfig`, browser, certificate |
| `internal/ui` | Services bound to the windows |
| `internal/trayicon` | Draws the tray icon at runtime |
| `frontend/` | Svelte 5 + TypeScript UI |
| `packaging/` | Wails build assets and Taskfiles |
| `docs/guides/docker-engine-wsl/` | The manual procedure this app automates |

</details>

## License

[MIT](LICENSE)
