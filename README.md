# Docker Zero

A small Windows tray app that sets up and runs Docker Engine inside a WSL2 distro, with Portainer CE as the GUI, for people on business machines who cannot use Docker Desktop and have no admin rights.

Status: early development. Nothing here is ready for use yet.

## What it does

- Day 1: a wizard installs an `Ubuntu-24.04` WSL2 distro, Docker Engine and Portainer CE on `https://localhost:9443`, all as a standard user.
- Day 2: a tray icon with an Action Center style popup starts and stops the distro, opens Portainer and shows live status.
- Settings: auto start at sign-in, shortcuts, how the distro is kept running, browser choice and a locally trusted certificate.

The WSL Windows feature must already be enabled; enabling it needs admin rights, so the app only detects it.

## Build

Requirements: Go 1.25 or newer, Node.js with npm, and the Wails CLI `v3.0.0-beta.30`:

```sh
go install github.com/wailsapp/wails/v3/cmd/wails3@v3.0.0-beta.30
```

Commands:

| Task | Command |
|---|---|
| Build `bin/docker-zero.exe` | `wails3 task build` |
| Run with live reload | `wails3 task dev` |
| Unit tests | `wails3 task test` |
| `go vet`, `golangci-lint`, Svelte type check | `wails3 task lint` |

Wails build assets (icon, manifest, version info) live in `packaging/` instead of the Wails default `build/`, so use `wails3 task build` rather than `wails3 build`.

## Layout

| Path | Contents |
|---|---|
| `main.go` | App wiring: tray, popup window, services |
| `internal/` | Go packages: paths, WSL runner, status, settings |
| `frontend/` | Svelte 5 + TypeScript UI |
| `packaging/` | Wails build assets and Taskfiles |
| `docs/guides/docker-engine-wsl/` | The manual procedure this app automates |
