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
| Phase 4 | Phase 3 | active | Tray popup and corner box drive Start, Stop, Shut down, Open |
| Phase 5 | Phase 4 | active | Settings apply and revert each item |
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

Status: active; gate not met: Stop distro and Shut down WSL not yet clicked against the real distro (see below), waiting on the user. Built: `internal/control` (one action at a time: Start = keepalive hold + `systemctl start docker` + `docker start portainer` + wait up to 90 s for the status endpoint + open browser; Stop distro; Shut down WSL; Open Portainer), `wsl.Hold` (hidden `sleep infinity` session), `internal/trayicon` (status-coloured icon drawn at run time), `internal/ui` (bound `Service`, poll loop, `Shell` with tray menu, popup and corner box), Svelte `Popup`, `Corner` and a shared status store. Design: DZ-D15.

Evidence:

| Check | Result |
|---|---|
| Popup via a second launch, screenshot | Status rows WSL, Ubuntu-24.04, Docker Engine, Portainer `Up, 2.45.2`; four tiles |
| Corner box toggled through UI Automation, screenshot | Box at the bottom-right with Open Portainer and Stop |
| Start clicked through UI Automation on the real distro | `pgrep` showed `/bin/sleep infinity` owned by the app's `wsl.exe`; a browser process started |
| App force-killed with `taskkill /F` | The hold ended (`pgrep` exit 1); the build before the job object left an orphan hold |
| `go test -tags integration ./internal/wsl` on throwaway distro `dz-test` | Terminate stops, Start runs, a hold keeps it running and ends with Terminate, Release ends the sleep inside |
| `wails3 task test`, `wails3 task lint` | All packages ok; `golangci-lint` 0 issues; `svelte-check` 0 errors |

Not verified: Stop distro and Shut down WSL clicked in the app against `Ubuntu-24.04`, because the user's live containers (several `hashicorp/terraform-mcp-server` sessions) run there; `wsl --shutdown` was not run anywhere because it stops every distro of the account. Their code paths call the `wsl.exe` commands checked above.

## Phase 5: Settings

Status: active; gate not met: certificate Install, Remove and Regenerate not yet run against the real store (see below), waiting on the user. Built: `internal/settings` (settings store, `HKCU` Run value, shortcuts, `.wslconfig` editor, Chrome lookup and browser opener, certificate read, trust, untrust and regenerate, self-copy to the install folder), `internal/ui` `Prefs` and `SettingsService`, a framed settings window opened from the tray menu and the popup gear, corner box choice remembered, `--start-distro` launch option, hold taken at launch in keepalive mode. Design: DZ-D16.

Evidence, each item applied and reverted through the settings window by UI Automation on this machine:

| Item | Observed |
|---|---|
| Auto start on, with start distro | Run value `Docker Zero` = quoted install path plus `--start-distro`; exe copied to the install folder (12 032 000 bytes) |
| Auto start off | Run value gone; the seven other Run values untouched |
| Desktop and Start Menu shortcuts on | Both `Docker Zero.lnk` target the installed exe, icon `exe,0` |
| Shortcuts off | Both files gone |
| Browser Chrome, then System default | `settings.json` `browser` `chrome`, then `default` |
| Idle WSL default | `.wslconfig` removed (it held only the line), `settings.json` `default`, hold gone |
| Idle `.wslconfig` | File recreated as `[general]` plus `instanceIdleTimeout=-1` |
| Idle keepalive | Line removed, hold back in the distro; the original `.wslconfig` was then restored byte for byte (hash equal) |
| Corner box on, then off | `settings.json` `cornerBox` `true`, then absent |
| Certificate read | Thumbprint `270FF7DCB695691011F84F72E477E4D799FE11FE`, expires 2031-10-14, trusted `true` (matches PT-D19) |
| Regenerate script on throwaway `dz-test` | New key `600`, certificate `644`, parsed; reported "not configured" because no Portainer container there |

Not verified: Install (trust), Remove and Regenerate against the real store and the real Portainer, because each store change needs the user to answer Windows' confirmation dialog and Regenerate restarts the live Portainer; opening Portainer in Chrome; Apply now (`wsl --shutdown`). Test state was cleaned up afterwards: install folder and `settings.json` deleted, the hold ended with the app.

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

1. Phase 7: adopt the existing manual setup (Phase 6 stays blocked on the Portainer admin question).
2. Ask the user before clicking Stop distro and Shut down WSL in the app against `Ubuntu-24.04`; both stop the live containers there.
3. Ask the user to run Install (trust), Remove and Regenerate in the settings window and answer Windows' dialogs; Regenerate restarts the live Portainer.
4. Ask the user whether `packaging/windows/wails.exe.manifest` may be committed; the commit scan flags its public Common-Controls token.
