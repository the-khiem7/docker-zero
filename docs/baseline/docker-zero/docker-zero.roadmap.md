---
baseline_schema: "2.0"
pack: "docker-zero"
document: "roadmap"
status: "draft"
updated: "2026-10-22"
code_ref: "175019e"
---

# docker-zero roadmap

| Phase | Depends on | Status | Acceptance |
|---|---|---|---|
| Phase 1 | - | complete | Minimal Wails3 tray exe runs from `%LOCALAPPDATA%\Programs` on this machine |
| Phase 2 | Phase 1 | complete | Repo scaffold builds; lint and tests run |
| Phase 3 | Phase 2 | complete | Go WSL runner and status detection tested against the real distro |
| Phase 4 | Phase 3 | complete | Tray popup and corner box drive Start, Stop, Shut down, Open |
| Phase 5 | Phase 4 | complete | Settings apply and revert each item |
| Phase 6 | Phase 3 | complete | Wizard installs a working setup on a throwaway distro |
| Phase 7 | Phase 5 | complete | Existing manual setup adopted without data loss |
| Phase 8 | Phase 6 | complete | Update Portainer, Reset admin, Uninstall all work |
| Phase 9 | Phase 8, Phase 7 | planned | Release exe and README published |

## Phase 1: Exe policy spike

Status: complete. Built a hello-world Wails3 tray app (`wails3` `v3.0.0-beta.30`, `svelte` template, frameless popup attached to the tray icon) in a throwaway directory, copied it to `%LOCALAPPDATA%\Programs\docker-zero` and ran it (DZ-D10). Evidence: the process started and kept responding (about 32 MB working set), no SmartScreen, AppLocker or antivirus prompt appeared, and a screenshot showed the popup open at the bottom-right corner. The spike exe was removed; no spike code is kept. Template name lesson: DZ-D12.

## Phase 2: Repository scaffold

Status: complete. Wails3 project from the `svelte` template (DZ-D11, DZ-D12) trimmed to Windows, assets moved to `packaging/` (DZ-D13). Files: `main.go` (single-instance tray app with a hidden frameless popup), `internal/paths` (install and config dirs, with tests), `Taskfile.yml` tasks `build`, `dev`, `test`, `lint`, `packaging/`, `frontend/`, `.golangci.yml`, `README.md`, `AGENTS.md`, `.gitignore`. Evidence: `wails3 task build` produced `bin/docker-zero.exe` (10.3 MB); `go test ./...` ok; `wails3 task lint` ran `go vet`, `golangci-lint` `v2.15.0` with 0 issues and `svelte-check` with 0 errors. The gate named `wails3 build`; it is replaced by `wails3 task build` (DZ-D13).

## Phase 3: WSL runner and status

Status: complete. Packages: `internal/wsl` (`Client` with `Available`, `Distros`, `RunningDistros`, `DistroState`, `Start`, `Terminate`, `Shutdown`, `Script` on stdin; UTF-16LE and UTF-8 decoding; hidden console), `internal/status` (`Prober.Probe` returning a `Snapshot` of WSL, distro, Docker service, Portainer container and HTTPS version), `internal/target` (fixed distro, user, container and URLs). Design: DZ-D14. Evidence: `go test ./...` ok for `wsl`, `status`, `paths`; `golangci-lint` 0 issues with and without the `integration` tag; `go test -tags integration ./internal/status` against the real `Ubuntu-24.04` read `WSLAvailable:true Distro:running Docker:active Portainer:running PortainerUp:true PortainerVersion:2.45.2` in 0.58 s. `Start`, `Terminate` and `Shutdown` are unit-tested for their arguments only; their effect on a real distro is Phase 4's gate.

## Phase 4: Day 2 widget

Status: complete. Built: `internal/control` (one action at a time: Start = keepalive hold + `systemctl start docker` + `docker start portainer` + wait up to 90 s for the status endpoint + open browser; Stop distro; Shut down WSL; Open Portainer), `wsl.Hold` (hidden `sleep infinity` session), `internal/trayicon` (status-coloured icon drawn at run time), `internal/ui` (bound `Service`, poll loop, `Shell` with tray menu, popup and corner box), Svelte `Popup`, `Corner` and a shared status store. Design: DZ-D15.

Evidence:

| Check | Result |
|---|---|
| Popup via a second launch, screenshot | Status rows WSL, Ubuntu-24.04, Docker Engine, Portainer `Up, 2.45.2`; four tiles |
| Corner box toggled through UI Automation, screenshot | Box at the bottom-right with Open Portainer and Stop |
| Start clicked through UI Automation on the real distro | `pgrep` showed `/bin/sleep infinity` owned by the app's `wsl.exe`; a browser process started |
| App force-killed with `taskkill /F` | The hold ended (`pgrep` exit 1); the build before the job object left an orphan hold |
| `go test -tags integration ./internal/wsl` on throwaway distro `dz-test` | Terminate stops, Start runs, a hold keeps it running and ends with Terminate, Release ends the sleep inside |
| `wails3 task test`, `wails3 task lint` | All packages ok; `golangci-lint` 0 issues; `svelte-check` 0 errors |
| Stop distro clicked on `Ubuntu-24.04`, user approved | After 8 s only `dz-test` running, no hold `wsl.exe`, Portainer HTTP `000`; popup showed Ubuntu-24.04 `Stopped` and Open Portainer disabled |
| Start after Stop | Portainer HTTP `200` after 7 s, hold back in the distro |
| Shut down WSL clicked, in-popup confirmation answered | Confirmation text and Shut down / Cancel buttons shown; after 10 s no distro running, no hold, Portainer HTTP `000` |
| Start after Shut down | Portainer HTTP `200` after 12 s; `docker ps` showed only `portainer` (the `--rm` MCP containers ended with the distro, as expected) |

Not verified: Open Portainer clicked on its own (Start opened the browser through the same call).

## Phase 5: Settings

Status: complete. Built: `internal/settings` (settings store, `HKCU` Run value, shortcuts, `.wslconfig` editor, Chrome lookup and browser opener, certificate read, trust, untrust and regenerate, self-copy to the install folder), `internal/ui` `Prefs` and `SettingsService`, a framed settings window opened from the tray menu and the popup gear, corner box choice remembered, `--start-distro` launch option, hold taken at launch in keepalive mode. Design: DZ-D16.

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

| Remove, Install (trust), Regenerate on the real store, clicked by the user | Windows showed "Root Certificate Store" delete and "Security Warning" install dialogs naming thumbprint `270FF7DC...99FE11FE`; the user answered them (screenshots). Afterwards the store holds one `CN=portainer`, thumbprint `2AEA72AAA4FB15117E3F8AF0C392AA51DEF54EBF`, 2026-10-22 to 2031-10-21, equal to the SHA-1 fingerprint of `/etc/portainer-certs/portainer.crt`; key `600`, certificate `644`; the old certificate is gone |
| Portainer after Regenerate | Container restarted; `curl.exe` without `-k` to `https://localhost:9443/api/system/status` returned `200` |

Not verified: opening Portainer in Chrome; Apply now (`wsl --shutdown`) from the settings window (the same `ShutdownWSL` call was clicked in the popup in Phase 4). Test state was cleaned up afterwards: install folder and `settings.json` deleted, the hold ended with the app.

## Phase 6: Day 1 wizard

Status: complete. Built: `internal/setup` (`Wizard` with eight resumable steps, `Config`, `Credentials`, `AddToUserPath`), `wsl.Client` `Install` and `Unregister`, `settings.Generate` (create the certificate only when missing), `ui.SetupService` and a framed setup window opened from the tray menu and from the popup (a Set up tile while the distro is missing, a "Set up..." link always), Svelte `Setup`. Design: admin creation and reset DZ-D18, mechanics DZ-D19.

Gate: a full run on a throwaway distro name ends with Portainer login working with the typed credentials; the real distro is not touched. Evidence:

| Check | Result |
|---|---|
| `go test -tags integration ./internal/setup` with `DZ_SETUP_DISTRO=dz-setup-test` | New distro registered, Docker 29.9.2 and Portainer 2.45.2 installed, admin created; first run 2 min 35 s, test 169.8 s |
| Independent sign-in to `https://127.0.0.1:19443/api/auth` with the typed credentials | `200`; a wrong password was not `200` |
| `docker inspect` of the container after the run | No setup token in the environment or the arguments |
| Second run on the same distro | Distro, Docker, user, Portainer and admin steps `skipped`; sign-in still `200` |
| Real setup | `Ubuntu-24.04` untouched, its Portainer answered `200` after the run; `dz-setup-test` unregistered |
| `wails3 task test`, `wails3 task lint` | All packages ok; `golangci-lint` 0 issues; `svelte-check` 0 errors |
| Setup window through UI Automation on the built exe | The popup link opened it; form fields and the eight steps listed; Start with empty fields showed the password-length message and started nothing |

Not verified: the window during a live run and its finish card (a screenshot captured the Windows lock screen); `AddToUserPath` on the real registry (the integration run kept it off); a machine without WSL, or behind a proxy that blocks Docker Hub or the apt hosts; a distro whose `wsl.conf` holds other settings.

## Phase 7: Adopt existing setup

Status: complete. Built: `internal/adopt` (`Adopter.Detect` and `Adopt`, `Windows` machine with four replaceable items), `settings.ReadShortcutTarget`, `setup.WriteLauncher`, `ui.AdoptService`, a once-per-profile offer at launch (stored flag `adoptPrompted`) and an Existing setup found card in the setup window. Design: DZ-D8 and DZ-D21.

Gate: run on this workstation; Portainer data intact and the legacy shortcut replaced. Evidence, from the built exe launched on the real account after a manual backup of the legacy items:

| Check | Result |
|---|---|
| Read-only detection (`TestDetectOnThisMachine`) | Distro running, container running, certificate present, four items found |
| First launch | The setup window opened on its own with the card and four checked items |
| Take over clicked through UI Automation | `Replaced 4 items`; the card then read `Nothing left to replace` |
| Legacy launchers | Desktop and Start Menu `Portainer.lnk`, the Startup keepalive shortcut and the `docker-wsl` folder gone; Docker Zero shortcuts present |
| Auto start and command | `HKCU Run` value kept `--start-distro`; `docker-zeroin\docker.cmd` written and run (`docker ps` listed `portainer`); user PATH 22 entries before and after, old entry replaced by the new one |
| Idle line | `.wslconfig` unchanged; stored idle mode became `wslconfig` |
| Portainer data | Container id and start time identical before and after, `portainer_data` listed, status endpoint `200` |
| `wails3 task test`, `wails3 task lint`, `wails3 task build` | All ok; 0 lint issues; 0 svelte-check errors |

Not verified: the card on a machine with a stopped distro or with only some legacy items; a failing item in the window; a locked Startup folder. The real account now runs the adopted setup; the legacy items were backed up by hand outside the repository.

## Phase 8: Maintenance and uninstall

Status: complete. Built: `setup.Wizard` methods `UpdatePortainer`, `ResetAdmin` and `Uninstall` (five steps, `Host` hooks, `FinishUninstall`), `wsl.Client.Export`, `RemoveFromUserPath`, `NewHost`, `ui.MaintainService` and `ui.UninstallService` behind one `opGuard`, Svelte `Maintenance` (Portainer card, Danger zone with two confirmations and a backup choice) inside the settings window. Design: DZ-D20, scope DZ-D7.

Gate: each action run on a throwaway distro. Evidence:

| Check | Result |
|---|---|
| `go test -tags integration -run TestMaintenance` with `DZ_SETUP_DISTRO=dz-maint-test` | Passed in 577 s |
| Update with the image current | `Portainer already runs the latest image`; container id unchanged |
| Update from an older local image | Container recreated on `portainer/portainer-ce:lts`; sign-in with the old admin `200`, so the data volume was kept |
| Reset admin | New password `200`, old password not `200`, container `running` afterwards |
| Uninstall with a backup | `wsl --export` tar over 1 MiB; distro `not-installed`; `docker.cmd` gone; all seven Host hooks ran |
| Real setup | `Ubuntu-24.04` untouched; `dz-maint-test` gone |
| `TestUserPathRoundTrip` on the real HKCU PATH | Add then remove restored value and type exactly; this also verifies `AddToUserPath` |
| `wails3 task test`, `wails3 task lint`, `wails3 task build` | All ok; 0 lint issues; 0 svelte-check errors |

Not verified: the Maintenance section in a live window; Windows trust removal and the HKCU Run and shortcut removals on the real account (hooks were recorders in the integration run); `FinishUninstall` deleting the install folder after exit; a failed pull or an unreachable Docker Hub during update.

## Phase 9: Release

Status: planned. Release build, README with screenshots, licence file. Gate: the release exe runs on a clean profile.

## Risks

| Risk | State |
|---|---|
| Unsigned exe blocked on managed machines | Passed on this machine in Phase 1 (DZ-D12); other machines unverified |
| Wails3 is pre-release | Pinned version (DZ-D9) |
| `.wslconfig` and `wsl --shutdown` affect every distro of the account | Confirmation dialogs; keepalive mode is the default (DZ-D6) |
| Trusting a root certificate on a managed device may raise a security alert | Optional step, user-initiated (DZ-D6) |
| Docker Hub anonymous pull terms for organisations | Unresolved in the source pack (PT-Q6) |

## Next action

1. Phase 9: release build, README with screenshots, licence file; gate is the release exe running on a clean profile.
