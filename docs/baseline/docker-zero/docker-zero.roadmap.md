---
baseline_schema: "2.0"
pack: "docker-zero"
document: "roadmap"
status: "draft"
updated: "2026-10-28"
code_ref: "uncommitted"
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
| Phase 9 | Phase 8, Phase 7, Phase 10 | active | Release workflow builds and publishes the exe on tag push (DZ-D32); README published |
| Phase 10 | Phase 8 | active | Brand logo used as the app icon, UI upgraded, illustrated README, landing page |

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
| Auto start and command | `HKCU Run` value kept `--start-distro`; `docker-zero\bin\docker.cmd` written and run (`docker ps` listed `portainer`); user PATH 22 entries before and after, old entry replaced by the new one |
| Idle line | `.wslconfig` unchanged; stored idle mode became `wslconfig` |
| Portainer data | Container id and start time identical before and after, `portainer_data` listed, status endpoint `200` |
| `wails3 task test`, `wails3 task lint`, `wails3 task build` | All ok; 0 lint issues; 0 svelte-check errors |

Not verified: the card on a machine with a stopped distro or with only some legacy items; a failing item in the window; a locked Startup folder. The real account now runs the adopted setup; the legacy items were backed up by hand outside the repository.

Revision summary (release run, 2026-10-22), legacy items restored from the manual backup:

| Check | Result |
|---|---|
| Distro stopped, only `Portainer.lnk` and the `docker-wsl` folder restored | Card listed exactly those two items; the distro stayed `Stopped` through detection |
| `docker-wsl\bin\docker.cmd` held open by another process, Take over | Shortcut item replaced; the docker item showed "remove ...docker.cmd: The process cannot access the file because it is being used by another process."; the other legacy files were removed |
| Take over again after the lock ended | Legacy folder gone; first build lost the result card on a stopped distro, fixed in `21ee8a5` (DZ-D21), then "Replaced 1 item" shown |

Not verified: a locked Startup folder.

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

Revision summary (release run, 2026-10-22), the settings window driven by UI Automation on the real `Ubuntu-24.04`:

| Check | Result |
|---|---|
| Update Portainer, a newer `lts` existed | Container recreated, Portainer `2.45.2` to `2.45.3`, admin still present (`/api/users/admin/check` `204`) |
| Update again | "Portainer already runs the latest image."; container id unchanged |
| Update with Docker Hub names pointed at `127.0.0.9` in the distro's `/etc/hosts` | Window showed the pull error (`connection refused`); container id, start time and state unchanged; hosts file restored |
| Reset admin with username `nobody-dz` | Password replaced, then "signing in as "nobody-dz" failed; check the username"; sign-in as `admin` with the new password `200` |
| Reset admin as `admin` | "Password replaced."; new password `200`, another password `422` |
| Uninstall all without backup | Five steps done or skipped; Windows "Root Certificate Store" delete dialog shown and answered Yes (not by the test script); certificate, distro, Run value, both shortcuts, PATH entry (22 to 21 entries), `.wslconfig` and settings gone; after Quit the process ended and the install folder was deleted |
| Left behind by that uninstall | `%APPDATA%\docker-zero.exe`, the WebView2 data folder; fixed in `8ca3a6c` (DZ-D25) |
| Second uninstall with the fixed build | Install folder and settings folder with `webview` removed after Quit; trust step skipped (certificate was not trusted) |

Not verified: the backup option from the window; a declined certificate dialog.

## Phase 9: Release

Status: active. Gate met at `da09172`; publishing deferred until Phase 10 is done (DZ-D27), so the exe below will be rebuilt and its hash will not name the published file. Built: README rewritten for users with five screenshots in `docs/images/` (account paths cut out or restated), MIT `LICENSE` (DZ-D22), app icon in the tray tile's style replacing the Wails placeholder (`494b924`), and the fixes the release run found: DZ-D23 (`7446db6`), DZ-D21 card (`21ee8a5`), DZ-D25 (`8ca3a6c`), DZ-D24 (`f171d3a`). Release exe: `bin/docker-zero.exe` built by `wails3 task build` at `da09172`, 12 240 384 bytes, SHA-256 `c0f34fc6027dfb499d29dcc263736e9cc6fdb61065188775556aed3d27a3fe21`, version info `0.1.0`.

Release workflow (DZ-D32): written in `.github/workflows/release.yml` (tag `v*`, `windows-latest`, version stamp from the tag (DZ-D33), `wails3 task build`, `wails3 task test`, `gh release create` with the exe and its SHA-256, a tag with a suffix is a pre-release). Observed on GitHub: tag `v0.1.0-alpha` failed at the test step, which ran before the build (DZ-D32); tag `v0.1.1-alpha` passed every step in 3 min 5 s and created a pre-release with `docker-zero.exe` (11 691 008 bytes) and `docker-zero.exe.sha256`; the downloaded exe hashed to the published SHA-256 (`72acb45a9bd60ed8409edb31c040b467b8a7bc4958ece87ae0eeb3f6f0d67251`) and contains `0.1.1-alpha`, not `0.0.0-dev`. Not verified: the exe running, and on a managed machine (DZ-D12). The pre-release and both test tags were then deleted from GitHub with `gh`, so no release exists. Next: the final gate runs on the exe downloaded from the real release (DZ-D29).

Gate: the release exe runs on a clean profile. This account was cleaned with Uninstall all, so no distro, settings, install folder, Run value, shortcuts, PATH entry or `.wslconfig`; `dz-test` kept. Evidence:

| Check | Result |
|---|---|
| Launch from a scratch folder | Setup window opened by itself in 1.2 s |
| Setup through the window (username `admin`, random password not stored) | All eight steps done, "All set" after 454 s (an earlier run took about 12 min, mostly the distro download); Docker `29.9.2`, Portainer `2.45.3`, `/api/users/admin/check` `204` |
| `docker` from a fresh environment | `docker version` `29.9.2`, `docker ps` listed `portainer` |
| Auto start, start distro, both shortcuts and the `.wslconfig` mode turned on | Exe copied to the install folder with the same hash; Run value with `--start-distro`; `.wslconfig` `[general]` plus `instanceIdleTimeout=-1` |
| Relaunch from the install folder | No setup window (DZ-D24) |
| `wails3 task test`, `wails3 task lint`, `golangci-lint` with the integration tag | All ok; 0 issues; 0 svelte-check errors |

Side effects left for the user: `wsl --unregister` made `dz-test` the default distro, set back to `Ubuntu-24.04` by hand; the Portainer admin password is a random one nobody knows, so the admin resets it with Reset admin password; the new certificate (SHA-1 `1BFBD28089E4A85A0E024D4C382C2084C3B8F240`, to 2031-10-21) is not trusted yet.

Not verified: SmartScreen or policy behaviour on another machine; the release download from GitHub; a machine without WSL.

## Phase 10: Design

Status: active. Brand, app UI and front door before the first release (DZ-D27). Items run in order; each starts only after the user accepts the review point of the one before. Exception: item 2 was built before item 1 at the user's request while the logo brief was pending. Item 1 done: logo chosen after seven rounds (DZ-D31). Items 5 and 6 were built ahead of items 3 and 4 at the user's request; item 6 is built and awaits review (DZ-D35).

| Item | Output | Review point |
|---|---|---|
| 1. Brand logo | Signature note and two or three rough logo directions on the local board (`frontdoor-illustrate`); the chosen mark as a vector source, rendered to `packaging/appicon.png`, `packaging/windows/icon.ico` (`wails3 generate icons`) and the status tray icons of `internal/trayicon` | User picks or mixes a direction; the mark reads at 16, 20, 24 and 32 px; tray states stay distinguishable on light and dark taskbars |
| 2. UI inventory | One HTML page listing every window and state (popup, corner box, settings, maintenance, setup with the take-over card, confirmations, errors), each with a screenshot: mocked frontend for every state, real exe for the main windows, committed (DZ-D28) | User confirms the list is complete |
| 3. UI evaluation and design | Findings per screen; redesigned mockups in the logo's signature for the screens that need it | User approves each changed screen |
| 4. Apply the design | Svelte and CSS changes under `frontend/`, the window sizes in `internal/ui` if a design needs them; no behaviour change | `wails3 task test`, `wails3 task lint`, `wails3 task build` pass; every inventory state re-shot from the built exe in light and dark and matched to its approved mockup |
| 5. README banner | Hero banner and figures as light and dark WebP pairs (`frontdoor-readme` with `frontdoor-illustrate`), README rewritten story-first with no UI screenshots (DZ-D34), no-affiliation statement (DZ-D30) | User approves the README plan before building; the skill's verify step passes; no account path in any image |
| 6. Landing page | One `index.html` for GitHub Pages (`frontdoor-landing`), social preview image, favicon from the mark, no-affiliation statement (DZ-D30) | User picks a direction and the placement; the skill's axe and Lighthouse checks; the user switches on Pages in the repository settings |

Tech debt to solve in items 3 and 4 (research first, then design):

| Debt | Evidence | Research to do |
|---|---|---|
| The setup window has no visible way back once setup is done | Reachable only from the tray menu entry "Set up Docker Zero..." and the small "Set up..." link at the foot of the popup; the Settings window has no entry; the window opens by itself only while the distro is not installed (DZ-D24) | Which name and place suit a returning user (a "Repair setup" action in Settings, a banner in the popup when a part is missing); what the window should show on a finished setup (the "Nothing left to replace" observation of item 2) |
| Setup does not notice a missing or stale part | On this workstation `docker-zero\bin` was on the user PATH while the folder and `docker.cmd` were gone; every shell reported `docker` not found and the popup showed Running (DZ-D36) | Cause (uninstall, a manual delete, or a skipped step 8); whether the popup should check the wrapper files and PATH and offer Repair |

Item 1 evidence (`3aa5256` icons, `c824f21` gallery):

| Check | Result |
|---|---|
| Logo | Round 7 direction 62 chosen after seven rounds (DZ-D31); `design/brand/logo.svg`, `design/brand/logo-tile.svg` |
| App icon | `packaging/appicon.png` 1024 px rendered from the tile; `packaging/windows/icon.ico` 16 to 256 px by `wails3 generate icons` (its macOS step fails on the missing `build/darwin`, after the `.ico` is written) |
| Tray icon | `internal/trayicon` embeds `stopped`, `partial`, `running`, `failed` PNGs at 32 px; `PNG(level)` unchanged |
| Checks | `wails3 task test` ok, `wails3 task lint` 0 issues, `svelte-check` 0 errors; icons viewed at 32 and 16 px on light and dark |
| Gallery | `design/brand/gallery/` (72 sheets, `directions.json`, `directions.html`); the `.frontdoor/` copy deleted; `AGENTS.md` rule added, scoped to the logo rounds by DZ-D34 |

Not verified: the icons in the real Windows tray and in Explorer from a built exe.

Item 2 evidence, awaiting the user's review:

| Check | Result |
|---|---|
| Mock harness | `frontend/vite.inventory.config.ts` (`npm run inventory`, port 9246) swaps the six bound service modules for `frontend/inventory/mock/*`; scenarios in `frontend/inventory/fixtures.ts`; each mock `satisfies` the generated binding's signatures, so a new or changed Go method fails `svelte-check` |
| Render | `design/ui-inventory/shoot.py` (Playwright on the installed Edge): 50 states over popup, corner box, settings, maintenance and setup, light and dark, WebP at 2x, 4.4 MB |
| Page | `design/ui-inventory/index.html`: filters by window, screenshot theme and review state, Keep / Revise / Discuss with a note per state kept in the browser, Copy review as Markdown; no horizontal scroll at 390 px; no page errors |
| Real references | `design/ui-inventory/real.ps1` captured the popup, corner box, settings and setup windows of the installed exe; account paths painted over |
| Observations on the page | 10; among them a failed Start in the corner box shows nothing (the error line is squeezed to zero height), and Setup opened after setup shows Existing setup found / Nothing left to replace above a form that would run setup again |
| `wails3 task lint`, `wails3 task test`, `wails3 task build` | 0 issues, 0 svelte-check errors; tests ok; build ok, the production bundle does not import the mocks |

Not verified: the mocks against the Go services at run time beyond the type check; the page in browsers other than Edge.

Item 5 evidence, built and awaiting the user's review:

| Check | Result |
|---|---|
| Figures | `design/readme/source/figures.html` and `icons.html` render `hero`, `story-licence`, `story-commands`, `story-gui`, `setup`, `tray`, `keep`, `safe`, seven `nav-*` buttons and six `ico-*` icons, each light and dark, to `design/readme/asset/` (44 WebP including the `cta` download button, 856 KB on disk, none over 60 KB); the banners and a sample of the buttons and icons viewed at full size; contrast of the smallest text pairs computed, all 4.5 or above after the button teal was darkened to `#1e7468` |
| Download button | `cta` item in `figures.html` (button teal `#1e7468`, cream text 5.15, badge on a cream coin, new `down` glyph), `cta.webp` and `cta-dark.webp` 1345x248, 16 KB each, viewed at full size in both themes; README places it above the navigation, links `releases/latest/download/docker-zero.exe` and `.sha256`, and the Install step points to it; `releases/latest` answers 404 today (no release exists since the test pre-release and tags were deleted), so the link is dead until the first release without a suffix |
| README | Rewritten story-first; every figure directly under its heading, one banner per story paragraph; five screenshots deleted; image paths and the seven navigation anchors checked to exist; no en or em dash; no account path |
| Labels | Setup steps from `internal/setup/setup.go`; popup rows and buttons from `frontend/src/Popup.svelte`; form text from `frontend/src/Setup.svelte`; tray colours from `internal/trayicon/*.png` |
| Licence claim | Quoted from Docker's licence page as fetched: paid subscription for professional use in larger organizations; linked, no thresholds stated |
| `wails3 task test` | All packages ok |

Not verified: the README on github.com in light, dark and phone width (including the picture elements); the text fragment of the Docker link in a browser; the figures against the redesigned UI.

Item 6 evidence, built and committed but not pushed, awaiting the user's review (DZ-D35):

| Check | Result |
|---|---|
| Files | `index.html` (50 363 bytes, inline CSS and JS, no analytics, no third-party script) and `.nojekyll` in `b394f9b`; `design/landing/source/og.html` and `design/landing/asset/og.png` (1200x630, 79 869 bytes, viewed) in `4176b1f`; favicon is `design/brand/logo-tile.svg`; the working tree was clean at `12e1644` |
| Structure, links, metadata | `shoot.py check` at 390, 820 and 1280 px in light and dark, with axe in both schemes: 0 FAIL, 0 WARN; title, description, Open Graph, Twitter, icon and both theme colours present |
| Controls | Start, Stop distro, the compare buttons, Start setup, Replay and two ring tabs each changed the screen with and without reduced motion |
| Ambient leaps | Over 20 s after the first pod: with motion, one orca appeared by itself; with reduced motion, none; a click on the water produced one orca in both modes |
| Lighthouse (localhost) | Mobile: Performance 94 in two runs (First Contentful Paint and Largest Contentful Paint lost points), Accessibility 100, Best Practices 100, SEO 100. Desktop: 99, 100, 100, 100 |
| Looked at | Hero stopped and running at 1280 and 390 px, the sections at 1280 px in light and dark, the compare slider on a phone |
| Text | Claims from the README and the app's labels; no en or em dash; the repository address appears only in URLs |

Not verified: the page on its real Pages URL and base path; the social preview in a share card; the Docker licence link's text fragment in a browser; the feel of the motion (only still frames were seen); Edge only, no other browser; a screen reader; a phone device. Pages is not switched on and no release exists, so the download button and the SHA-256 link answer 404 until `v0.1.0` (DZ-D34).

## Risks

| Risk | State |
|---|---|
| Unsigned exe blocked on managed machines | Passed on this machine in Phase 1 (DZ-D12); other machines unverified |
| Wails3 is pre-release | Pinned version (DZ-D9) |
| `.wslconfig` and `wsl --shutdown` affect every distro of the account | Confirmation dialogs; keepalive mode is the default (DZ-D6) |
| Trusting a root certificate on a managed device may raise a security alert | Optional step, user-initiated (DZ-D6) |
| Docker Hub anonymous pull terms for organisations | Unresolved in the source pack (PT-Q6) |
| README figures name the current UI's labels; the redesign (items 3 and 4) can leave them stale | Re-check the setup and tray figures after item 4 (DZ-D34) |
| README download button is a dead link until a release without a suffix exists | Accepted by the user; publish `v0.1.0` and open the button from the README on github.com (DZ-D34) |
| Landing page download button and SHA-256 link are dead until a release without a suffix exists; the page demos the current UI's labels and the redesign (items 3 and 4) can leave them stale | Switch Pages on only after `v0.1.0`; re-check the page's demos after item 4 (DZ-D35, DZ-D34) |
| Landing page `og:image` and canonical URL name the Pages address, which resolves only after Pages is on; sharing debuggers cache the image | Re-scrape after the first deploy (DZ-D35) |
| Logo or name read as Docker Inc. or Portainer branding | Phase 10 logo uses an orca, not Docker's whale graphic (DZ-D31); README and landing page state no affiliation (DZ-D30); Docker's trademark guidelines not checked, the user checks them before the release |

## Next action

1. Phase 10 item 2: the user reviews `design/ui-inventory/index.html` and confirms the list of states is complete.
2. Phase 10 item 5 (README) was built ahead of items 3 and 4 at the user's request (DZ-D34): the user reviewed `README.md` and `design/readme/asset/`, and the README is committed; the user pushes it. Items 3, 4 and 6 follow in order, one commit per item after its checks pass; after item 4, re-check the setup and tray figures against the redesigned UI.
3. Phase 10 item 6 (landing page) was built ahead of items 3 and 4 (DZ-D35): the user reviews `index.html` in a browser; `index.html`, `.nojekyll` and `design/landing/` are committed, and the user pushes them. The user answers DZ-Q6 (link the page from the README and the repository Website field). The user switches Pages on (Settings, Pages, Deploy from a branch, `main`, `/ (root)`) only after a release without a suffix exists. Re-check the page's demos after item 4.
4. Phase 9: after Phase 10, the user pushes tag `v0.1.0` on the commit that ends it; the workflow builds the exe and creates the release. Download that exe and re-run the full clean-profile gate on it (DZ-D29). The agent's earlier attempt to publish was refused by its permission mode, so the user pushes the tag or allows it.
5. Then mark Phase 9 complete with the release URL.
