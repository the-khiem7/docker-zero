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
| Phase 1 | - | planned | Minimal Wails3 tray exe runs from `%LOCALAPPDATA%\Programs` on this machine |
| Phase 2 | Phase 1 | planned | Repo scaffold builds; lint and tests run |
| Phase 3 | Phase 2 | planned | Go WSL runner and status detection tested against the real distro |
| Phase 4 | Phase 3 | planned | Tray popup and corner box drive Start, Stop, Shut down, Open |
| Phase 5 | Phase 4 | planned | Settings apply and revert each item |
| Phase 6 | Phase 3, DZ-Q2 | blocked | Wizard installs a working setup on a throwaway distro |
| Phase 7 | Phase 5 | planned | Existing manual setup adopted without data loss |
| Phase 8 | Phase 6 | planned | Update Portainer, Reset admin, Uninstall all work |
| Phase 9 | Phase 8, Phase 7 | planned | Release exe and README published |

## Phase 1: Exe policy spike

Status: planned. Build a hello-world Wails3 tray app with `wails3` `v3.0.0-beta.30`, copy it to `%LOCALAPPDATA%\Programs\docker-zero`, run it, and record whether SmartScreen, AppLocker or Defender intervene (DZ-D10). Gate: the tray icon appears and a popup window opens. A block stops the plan until a delivery form is chosen.

## Phase 2: Repository scaffold

Status: planned. Wails3 project from the `svelte-ts` template (DZ-D11), Go module, `go vet` and `golangci-lint`, `go test`, `README.md`, `AGENTS.md`, `.gitignore`. Gate: `wails3 build` produces the exe; tests and lint pass.

## Phase 3: WSL runner and status

Status: planned. Go package that runs `wsl.exe` with UTF-16 output decoded, sends scripts on stdin or as files (PT-D15), reports distro state and Portainer HTTPS health. Gate: unit tests plus a run against the real `Ubuntu-24.04`.

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
| Unsigned exe blocked on managed machines | Open; tested in Phase 1 (DZ-D10) |
| Wails3 is pre-release | Pinned version (DZ-D9) |
| `.wslconfig` and `wsl --shutdown` affect every distro of the account | Confirmation dialogs; keepalive mode is the default (DZ-D6) |
| Trusting a root certificate on a managed device may raise a security alert | Optional step, user-initiated (DZ-D6) |
| Docker Hub anonymous pull terms for organisations | Unresolved in the source pack (PT-Q6) |

## Next action

1. Phase 1: after the user approves it, build and run the Wails3 tray spike from `%LOCALAPPDATA%\Programs\docker-zero` on this machine.
