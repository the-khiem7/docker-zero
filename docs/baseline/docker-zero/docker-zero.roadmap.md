---
baseline_schema: "2.0"
pack: "docker-zero"
document: "roadmap"
status: "draft"
updated: "2026-10-22"
code_ref: "uncommitted"
---

# docker-zero roadmap

| Phase | Depends on | Status | Acceptance |
|---|---|---|---|
| Phase 1 | - | complete | Minimal Wails3 tray exe runs from `%LOCALAPPDATA%\Programs` on this machine |
| Phase 2 | Phase 1 | complete | Repo scaffold builds; lint and tests run |
| Phase 3 | Phase 2 | complete | Go WSL runner and status detection tested against the real distro |
| Phase 4 | Phase 3 | planned | Tray popup and corner box drive Start, Stop, Shut down, Open |
| Phase 5 | Phase 4 | planned | Settings apply and revert each item |
| Phase 6 | Phase 3, DZ-Q2 | blocked | Wizard installs a working setup on a throwaway distro |
| Phase 7 | Phase 5 | planned | Existing manual setup adopted without data loss |
| Phase 8 | Phase 6 | planned | Update Portainer, Reset admin, Uninstall all work |
| Phase 9 | Phase 8, Phase 7 | planned | Release exe and README published |

## Phase 1: Exe policy spike

Status: complete. Built a hello-world Wails3 tray app (`wails3` `v3.0.0-beta.30`, `svelte` template, frameless popup attached to the tray icon) in a throwaway directory, copied it to `%LOCALAPPDATA%\Programs\docker-zero` and ran it (DZ-D10). Evidence: the process started and kept responding (about 32 MB working set), no SmartScreen, AppLocker or antivirus prompt appeared, and a screenshot showed the popup open at the bottom-right corner. The spike exe was removed; no spike code is kept. Template name lesson: DZ-D12.

## Phase 2: Repository scaffold

Status: complete. Wails3 project from the `svelte` template (DZ-D11, DZ-D12) trimmed to Windows, assets moved to `packaging/` (DZ-D13). Files: `main.go` (single-instance tray app with a hidden frameless popup), `internal/paths` (install and config dirs, with tests), `Taskfile.yml` tasks `build`, `dev`, `test`, `lint`, `packaging/`, `frontend/`, `.golangci.yml`, `README.md`, `AGENTS.md`, `.gitignore`. Evidence: `wails3 task build` produced `bin/docker-zero.exe` (10.3 MB); `go test ./...` ok; `wails3 task lint` ran `go vet`, `golangci-lint` `v2.15.0` with 0 issues and `svelte-check` with 0 errors. The gate named `wails3 build`; it is replaced by `wails3 task build` (DZ-D13).

## Phase 3: WSL runner and status

Status: complete. Packages: `internal/wsl` (`Client` with `Available`, `Distros`, `RunningDistros`, `DistroState`, `Start`, `Terminate`, `Shutdown`, `Script` on stdin; UTF-16LE and UTF-8 decoding; hidden console), `internal/status` (`Prober.Probe` returning a `Snapshot` of WSL, distro, Docker service, Portainer container and HTTPS version), `internal/target` (fixed distro, user, container and URLs). Design: DZ-D14. Evidence: `go test ./...` ok for `wsl`, `status`, `paths`; `golangci-lint` 0 issues with and without the `integration` tag; `go test -tags integration ./internal/status` against the real `Ubuntu-24.04` read `WSLAvailable:true Distro:running Docker:active Portainer:running PortainerUp:true PortainerVersion:2.45.2` in 0.58 s. `Start`, `Terminate` and `Shutdown` are unit-tested for their arguments only; their effect on a real distro is Phase 4's gate.

## Phase 4: Day 2 widget

Status: planned. Tray icon, Action Center style popup, optional corner box, actions and polling per DZ-D3 and DZ-D4. Gate: each action observed changing the real distro state.

## Phase 5: Settings

Status: planned. Items per DZ-D6. Gate: each setting applied and reverted, with the resulting registry value, file or shortcut observed.

## Phase 6: Day 1 wizard

Status: blocked on DZ-Q2 (admin init and reset mechanism). Steps per DZ-D5. Gate: a full run on a throwaway distro name ends with Portainer login working with the typed credentials; the real distro is not touched.

## Phase 7: Adopt existing setup

Status: planned. Detection and migration per DZ-D8. Gate: run on this workstation; Portainer data intact and the legacy shortcut replaced.

## Phase 8: Maintenance and uninstall

Status: planned. Update Portainer, Reset admin, Uninstall all per DZ-D7. Gate: each run on a throwaway distro.

## Phase 9: Release

Status: planned. Release build, README with screenshots, licence file. Gate: the release exe runs on a clean profile.

## Risks

| Risk | State |
|---|---|
| Unsigned exe blocked on managed machines | Passed on this machine in Phase 1 (DZ-D12); other machines unverified |
| Wails3 is pre-release | Pinned version (DZ-D9) |
| `packaging/windows/wails.exe.manifest` uncommitted: the commit secret scan flags the public Common-Controls assembly token on line 6 | Open; a clean clone cannot build until the user allows committing it |
| `.wslconfig` and `wsl --shutdown` affect every distro of the account | Confirmation dialogs; keepalive mode is the default (DZ-D6) |
| Trusting a root certificate on a managed device may raise a security alert | Optional step, user-initiated (DZ-D6) |
| Docker Hub anonymous pull terms for organisations | Unresolved in the source pack (PT-Q6) |

## Next action

1. Phase 4: build the tray popup and corner box with Start, Stop, Shut down, Open and live status.
