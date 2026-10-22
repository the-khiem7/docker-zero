# Agent rules for Docker Zero

- Public repository: no usernames, company names, machine names or absolute user paths in code, docs or commits.
- Standard user only: never add a step that needs elevation (UAC, `Program Files`, `HKLM`, services, Windows features). Use `%LOCALAPPDATA%`, `%APPDATA%` and `HKCU`.
- Wails `v3.0.0-beta.30` is pinned for the CLI, the Go module and `@wailsio/runtime`; upgrade all three together.
- Build with `wails3 task build`; build assets live in `packaging/`, not `build/`.
- Before finishing Go work run `wails3 task test` and `wails3 task lint`.
- Scripts sent into WSL go on stdin or as files, never through nested command-line quoting: PowerShell and `wsl.exe` quoting truncates scripts silently.
- Never log or print the Portainer setup token or admin password.
- Container, image and volume management belongs to Portainer; the app only handles system-level actions.
- `docs/baseline/` holds internal planning notes; never reference those files or their entry IDs from code.
