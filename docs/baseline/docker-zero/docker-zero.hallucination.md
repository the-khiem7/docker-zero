---
baseline_schema: "2.0"
pack: "docker-zero"
document: "hallucination"
status: "draft"
updated: "2026-10-23"
code_ref: "uncommitted"
---

# docker-zero hallucination journal

## Entry index

| ID | About | Status | Related |
|---|---|---|---|
| DZ-D1 | Brand, audience and repository identity | closed | DZ-D2 |
| DZ-D2 | Imported `portainer` pack and guide are reference only | closed | DZ-D1, DZ-D8 |
| DZ-D3 | Widget form: tray popup plus optional corner box | closed | DZ-D4 |
| DZ-D4 | Day 2 actions: system-level only, two stop modes, live status | closed | DZ-D3 |
| DZ-D5 | Day 1 wizard covers every step; WSL feature not enabled by the app | closed | DZ-D18, DZ-D6 |
| DZ-D6 | Settings: auto start, shortcuts, idle mode, browser, certificate | closed | DZ-D5, DZ-D7 |
| DZ-D7 | Portainer `lts`, Update button, admin reset, uninstall danger zone | closed | DZ-D18, DZ-D6 |
| DZ-D8 | Adopt an existing manual setup | closed | DZ-D2 |
| DZ-D9 | Wails3 and portable self-installing exe | closed | DZ-D11, DZ-D10 |
| DZ-D10 | Exe policy spike runs before app code | closed | DZ-D9 |
| DZ-D11 | Frontend framework: Svelte | closed | DZ-Q1, DZ-D9, DZ-D12 |
| DZ-D12 | Wails3 Svelte template name and tray spike result | closed | DZ-D10, DZ-D11 |
| DZ-D13 | Scaffold shape: `packaging/` folder, pinned tools, Windows only | closed | DZ-D9 |
| DZ-D14 | How the WSL runner reads state and probes status | closed | DZ-D4 |
| DZ-D15 | Day 2 widget mechanics: hold, job object, polling, confirm, second launch | closed, job object part superseded | DZ-D3, DZ-D4, DZ-D6, DZ-D16 |
| DZ-D16 | Settings mechanics and the hold-only job object | closed, certificate read amended by DZ-D23 | DZ-D6, DZ-D9, DZ-D15, DZ-D23 |
| DZ-D17 | Lesson: a phase marked complete with its gate unmet | closed | DZ-D15, DZ-D16 |
| DZ-D18 | Portainer admin: created through the init API with an app-chosen setup token, reset with the helper and a hash | closed | DZ-Q2, DZ-D5, DZ-D7 |
| DZ-D19 | Setup wizard mechanics: resumable steps, token through an env file, pinned certificate, no `docker.ps1` | closed | DZ-D5, DZ-D18 |
| DZ-D20 | Maintenance mechanics: update skips a current image, reset sends a hash, uninstall order and self-removal | closed, settings removal superseded by DZ-D25 | DZ-D7, DZ-D18, DZ-D19, DZ-D25 |
| DZ-D21 | Adoption mechanics: four replaceable items, read-only detection, replace before remove, one offer | closed, launch offer superseded by DZ-D24 | DZ-D8, DZ-D6, DZ-D20, DZ-D24 |
| DZ-D22 | Repository licence: MIT | closed | DZ-Q3, DZ-D1 |
| DZ-D23 | Nothing at launch runs a command inside a stopped distro | closed | DZ-D16, DZ-D14 |
| DZ-D24 | Launch opens the setup window when nothing is installed; the take-over offer needs launchers to replace | closed | DZ-D21, DZ-D5 |
| DZ-D25 | WebView2 data beside the settings, both folders removed after exit | closed | DZ-D20, DZ-D9 |
| DZ-D26 | Lesson: Windows line endings and a shared WebView2 folder broke checks during the release run | closed | DZ-D13, DZ-D25 |
| DZ-D27 | The first release waits for the design phase | closed | DZ-D22, DZ-D29 |
| DZ-D28 | UI inventory: mocked frontend for every state, real exe for the main windows, committed | closed | DZ-Q4, DZ-D27 |
| DZ-D29 | The rebuilt release exe re-runs the full clean-profile gate | closed, gate target amended by DZ-D32 | DZ-Q5, DZ-D27, DZ-D32 |
| DZ-D30 | README and landing page state no affiliation with Docker Inc. or Portainer | closed, logo rule amended by DZ-D31 | DZ-D1, DZ-D27, DZ-D31 |
| DZ-D31 | Logo: an orca leaping out of a red ring over a two-tone sea; design rounds kept as a public gallery | closed | DZ-D30, DZ-D27 |
| DZ-D32 | Release exe built and published by a GitHub Actions workflow on tag push | closed | DZ-D29, DZ-D27, DZ-D12, DZ-D33 |
| DZ-D33 | The release workflow stamps the exe version from the git tag, on the runner only | closed | DZ-D32 |
| DZ-D34 | README: story-led, no UI screenshots, drawn figures, banner with the logo as the letter o | closed | DZ-D30, DZ-D31, DZ-D27 |
| DZ-D35 | Landing page: root `index.html`, rising-tide hero with a leaping pod, interactive body figures | closed | DZ-D34, DZ-D30, DZ-D31, DZ-D27, DZ-D1 |
| DZ-Q6 | Link the landing page from the README and the repository Website field | open | DZ-D35, DZ-D34 |

## Open questions

### DZ-Q6: Link the landing page from the README and the repository Website field?

The README does not mention the landing page and the repository `homepageUrl` is empty. A link written before Pages is switched on is dead, and the page's download button is dead until a release without a suffix exists (DZ-D34). Recommended: add both links after Pages is on and the first release exists. Decides: a README edit (item 6 of Phase 10) and a GitHub setting the user owns.

## Closed decisions

### DZ-D1: Docker Zero is the repository brand and the app name; audience is anyone on a business machine

Decision: name the repository and the app Docker Zero; target any user on a business machine with the same Docker Desktop licence problem and no admin rights.

Why: the user's answer; the pain is common on managed machines.

What breaks if ignored: tuning for one workstation hard-codes paths and names, and the public repo leaks identifiers.

Rejected: one-machine tool - too narrow for the stated audience. Colleagues-only - the user chose public reach.

### DZ-D2: The imported `portainer` pack and guide are reference only

Decision: `docs/baseline/portainer/` and `docs/guides/docker-engine-wsl/` (moved from `guides/` at the repository root) document the verified manual runtime and impose no rules here; this repository is treated as new and follows Wails3 and modern Go practice. The files may be moved.

Why: they were extracted from another repository whose `AGENTS.md` rules do not exist here.

What breaks if ignored: rules quoted from the old `AGENTS.md` (commit handling, guide template) constrain the app without a reason.

Rejected: none.

### DZ-D3: Tray icon with an Action Center style popup, plus an optional always-on-top corner box

Decision: a system tray icon opens a frameless popup anchored above the tray, styled like the Windows 11 Action Center; a setting also shows a small always-on-top box in the bottom-right corner.

Why: the user chose both forms.

What breaks if ignored: a permanent window only crowds the screen; a tray-only form hides status from users who want it visible.

Rejected: a single form only - the user asked for both.

### DZ-D4: Day 2 actions are system-level only; two stop modes; live status

Decision: the app offers Start Docker Zero (start the distro, ensure Portainer, open it), Open Portainer, Stop distro (`wsl --terminate <distro>`) and Shut down WSL (`wsl --shutdown`, all distros, with confirmation), plus live status of the distro and Portainer. Anything Portainer can do stays in Portainer.

Why: the user wants a GUI for the system-level actions that needed commands, and no duplicate of Portainer's job.

What breaks if ignored: duplicated container UI grows the app and drifts from Portainer.

Status poll: every 5 seconds while the popup or the corner box is visible, and briefly after an action. A hidden app stays idle.

Rejected: `docker stop portainer` as the only stop - keeps the distro's memory. Container lists in the app - Portainer owns them.

### DZ-D5: The Day 1 wizard covers every step; the WSL feature is not enabled by the app

Decision: the wizard runs, with live log and per-step status: detect WSL (stop with "ask IT to enable WSL" if absent), `wsl --install -d Ubuntu-24.04 --web-download --no-launch`, user `dev` and `wsl.conf` with systemd, Docker Engine from Docker's apt repository, the `GODEBUG=netdns=cgo` drop-in (PT-D3), the leaf certificate (PT-D19), the Portainer container with the certificate on `127.0.0.1:9443` with `unless-stopped` (PT-D2, PT-D7), the Portainer admin from a username and password the user types (DZ-D18), and `docker` wrappers for `cmd` (`docker.cmd`) and PowerShell (`docker.ps1`) on the user PATH (PT-D20). Distro `Ubuntu-24.04` and user `dev` are fixed.

Why: every step ran without elevation by hand; the user wants all of them and no agent help.

What breaks if ignored: a skipped step sends the user back to the terminal.

Rejected: putting the setup token on the clipboard for the web form - the user wants the app to collect username and password. Configurable distro and user names - the user chose fixed.

### DZ-D6: Settings content

Decision:

- Auto start: an `HKCU\Software\Microsoft\Windows\CurrentVersion\Run` value starts the app at sign-in, with a sub-option to start the distro too. It replaces the manual `WSL Docker keepalive.lnk`.
- Shortcuts: create or remove `Docker Zero` on the Desktop and in the Start Menu.
- Idle mode, three choices: app-held keepalive session (default; a hidden `sleep infinity` session while the app runs, PT-D5), `.wslconfig` `instanceIdleTimeout=-1` (account-wide, applied with a confirmed `wsl --shutdown`, PT-D16), or WSL default.
- Browser for Open Portainer: system default (default), Google Chrome by full path, or another chosen executable (PT-D17).
- Certificate: install, remove, regenerate. The `Cert:\CurrentUser\Root` import and removal show a Windows confirmation dialog, which is expected in a GUI app (PT-D21).

Why: the user accepted these proposals.

What breaks if ignored: the keepalive default avoids changing every WSL2 distro on the account; a `.wslconfig` change without `wsl --shutdown` does nothing.

Rejected: a Startup folder shortcut for auto start - `HKCU Run` is simpler to read and remove. Chrome-only - other users have other browsers.

### DZ-D7: Portainer `lts`, Update button, admin reset, uninstall in a danger zone

Decision: pin image `portainer/portainer-ce:lts`; an Update Portainer button pulls it and recreates the container on the same `portainer_data` volume and certificate; a Reset admin action per DZ-D18; an Uninstall all action in a Danger zone, with two confirmations and an offered `wsl --export` backup, removing everything the app created.

Why: the user's answers.

What breaks if ignored: `wsl --export` stops the distro and every container (PT-D21), so the backup step must say so; a leftover `instanceIdleTimeout=-1` or trusted root outlives the app (PT-D22).

Rejected: a tag picker - `lts` only.

### DZ-D8: Adopt an existing manual setup

Decision: on first run the app detects a setup made from the source guide (distro, container, certificate, `docker-wsl` folder, `WSL Docker keepalive.lnk`, `Portainer.lnk`, `.wslconfig` line) and takes it over, offering to migrate legacy shortcuts.

Why: the user and other early users already run the manual setup; reinstalling would destroy their data.

What breaks if ignored: duplicate auto start (old shortcut plus `HKCU Run`) and two launchers.

Rejected: fresh install only - forces an uninstall first.

### DZ-D9: Wails3 and a portable self-installing exe

Decision: Go with Wails3 (pinned to the installed `v3.0.0-beta.30` CLI and matching module), WebView2 frontend (Svelte, DZ-D11), shipped as one portable `.exe` that copies itself to `%LOCALAPPDATA%\Programs\docker-zero`.

Why: Wails3 has native tray and multi-window support that the tray popup plus corner box needs; the self-copy keeps shortcuts and `HKCU Run` valid when the download folder moves, the same pattern as `open-portainer.ps1`.

What breaks if ignored: an unpinned pre-release breaks builds on upgrade.

Rejected: Wails v2 - weaker tray-anchored windows and multi-window. NSIS per-user installer - more moving parts than a single exe.

### DZ-D10: Exe policy spike runs before app code

Decision: Phase 1 builds a minimal Wails3 tray app, copies it to `%LOCALAPPDATA%\Programs\docker-zero` and runs it on this business machine before any feature work.

Why: unsigned exes may be blocked by SmartScreen, AppLocker or Defender on managed machines; the current machine only proved that `RemoteSigned` PowerShell scripts run.

What breaks if ignored: a finished app that the target machines refuse to run.

Rejected: building first and testing policy at the end - the whole delivery form depends on it.

### DZ-D11: Svelte is the frontend framework

Closes DZ-Q1.

Decision: use Svelte 5 with TypeScript from the Wails3 `svelte` template (DZ-D12).

Why: the user asked for a light, fast app; Svelte compiles away its runtime (a few KB against about 35 KB for Vue 3, typical sizes, not measured on this app), and the popup is a handful of buttons and status values that need no large component library.

What breaks if ignored: a heavier framework adds memory to an app meant to sit resident in the tray.

Rejected: Vue 3 - larger runtime, its bigger UI ecosystem is not needed here. Vanilla TS and React - offered in the first interview and not picked.

### DZ-D12: The Wails3 `v3.0.0-beta.30` Svelte TypeScript template is named `svelte`; the unsigned tray exe runs from `%LOCALAPPDATA%\Programs`

Decision: scaffold with `wails3 init -t svelte` (Svelte + TypeScript + Vite); `svelte-js` is the JavaScript variant. Keep the portable self-copying exe (DZ-D9) as the delivery form.

Why: lesson - `svelte-ts` looked right from the v2 template naming and was written into DZ-D11 and Phase 2; `wails3 init -t svelte-ts` failed with `invalid template name` and `wails3 init -l` listed `svelte` as "Svelte + TypeScript + Vite". The Phase 1 spike (unsigned, built with `-H windowsgui`) started from `%LOCALAPPDATA%\Programs\docker-zero` on this business machine with no SmartScreen, AppLocker or antivirus prompt, and its frameless popup opened anchored at the bottom-right corner (screenshot observed).

What breaks if ignored: a scaffold step with `svelte-ts` fails before any file is created.

Rejected: NSIS per-user installer as a fallback delivery form - not needed while the unsigned exe runs.

### DZ-D13: Wails build assets live in `packaging/`; the scaffold is Windows only with pinned tools

Decision: keep the Wails build assets (`config.yml`, Taskfiles, icon, manifest, version info) in `packaging/` instead of the Wails default `build/`; build with `wails3 task build`, never `wails3 build`. Drop the Android, iOS, Linux, macOS, Docker, NSIS and MSIX assets and tasks, and the icon generation task. Pin `@wailsio/runtime` to `3.0.0-beta.30` (the template wrote `latest`, which resolves to `3.0.0-beta.34`) and run `golangci-lint` `v2.15.0` through `go run` from the Taskfile. Delete the template npm config file in `frontend/`.

Why: the agent session's permission policy denies Read and Edit inside any `build/` folder, so assets there cannot be maintained by an agent. `wails3 build` and `wails3 generate icons` hard-code `build/` paths (the icon command defaults `-macfilename` to `build/darwin/icon.icns` and failed with "The system cannot find the path specified"). The app targets Windows only (DZ-D9). The runtime must match the Go module, and `go run` with a version needs no global install. The template npm config file could not be read (a hook blocks that path) and produced an npm "Unknown project config" warning.

What breaks if ignored: `wails3 build` looks for `build/` and fails; an unpinned runtime drifts from the Go module.

Rejected: keeping `build/` and editing it through shell commands - works around the permission policy instead of respecting it. Installing `golangci-lint` globally - an extra setup step for every contributor.

### DZ-D14: The WSL runner reads state from `--list --quiet`, decodes both encodings, and the probe never starts a stopped distro

Decision:

- Distro state comes from `wsl.exe --list --quiet` (installed) and `--list --quiet --running`; a non-zero exit of the `--running` form means none running.
- Every `wsl.exe` child gets `WSL_UTF8=1`; output is still decoded as UTF-16LE when most odd bytes are zero, and BOMs are stripped.
- Children start with `CREATE_NO_WINDOW` and `HideWindow`.
- `Script` sends the script on stdin to `bash -s` with CRLF turned into LF.
- `Probe` runs Docker and Portainer checks only when the distro is already running, as `root`, then reads `https://127.0.0.1:9443/api/system/status` with certificate checks skipped.
- `go 1.26.0` in `go.mod` for `errors.AsType`.

Why: observed on this machine - `wsl.exe -l -v` wrote UTF-16LE without a BOM (`20 00 20 00 4e 00 ...`) and UTF-8 with `WSL_UTF8=1`; the verbose table's header and state words are localised on non-English Windows while names are not. A GUI-subsystem exe that starts a console program flashes a console window without the flag. Running anything inside a stopped distro boots it, which would defeat Stop. The probe must work before the leaf certificate exists; it is loopback only and sends nothing secret. A probe script piped from PowerShell (CRLF line ends) printed nothing for its `systemctl is-active docker` line while the LF-only lines after it worked; a trailing `\r` in an argument is the likely cause, not proven.

What breaks if ignored: wrong state on localised Windows; console flashes every 5 seconds; a status poll that restarts the distro the user just stopped.

Rejected: parsing `--list --verbose` - localised. A TLS-verified probe - fails until the certificate is installed. Probing as `dev` - depends on group membership that the wizard sets up later.

### DZ-D15: Day 2 widget mechanics

Decision:

- Start holds the distro with a hidden `wsl.exe --distribution Ubuntu-24.04 --exec /bin/sleep infinity` child for as long as the app runs (the keepalive idle mode of DZ-D6, built now because Start needs it); Stop distro and Shut down WSL release it first; app exit releases it.
- The app puts itself in a job object with `JOB_OBJECT_LIMIT_KILL_ON_JOB_CLOSE` before the first hold, so every child dies with the app. Superseded by DZ-D16: only the hold joins the job.
- Polling: every 5 s while the popup or the corner box is visible, every 1 s for 15 s after an action, never otherwise; one reading at launch. Hiding a window re-sends the last reading without probing.
- Shut down WSL asks for confirmation inside the popup, not with a native dialog.
- A second launch of the exe opens the popup (`SingleInstance.OnSecondInstanceLaunch`).
- Windows are told apart by `?view=popup` and `?view=corner`; the tray icon is drawn in Go and its dot follows the status level.

Why: observed - a force-killed build without the job object left `/bin/sleep infinity` running in the distro, which would keep it up forever; with the job object the hold ended. Job membership is inherited at process creation, so `wsl.exe` and the `wsl.exe` it starts are both covered without a race; assigning each child after `Start` would miss a grandchild created first. A native dialog takes focus and the popup hides on focus loss (`HideOnFocusLost`). The second-launch hook gives shortcuts and the Start Menu entry a visible effect and made UI Automation testing possible.

What breaks if ignored: the distro never idles out after a crash; the confirmation dialog closes the popup that asked for it; a hidden app polls `wsl.exe` forever.

Rejected: assigning only the hold's `wsl.exe` to a job after `Start` - races with its own child. A native confirmation dialog - hides the popup. Polling while hidden - wastes CPU in an app meant to sit idle in the tray.

### DZ-D16: Settings mechanics; only the hold joins the kill-on-close job

Decision:

- Only the hold joins the job object: `wsl.exe` starts with `CREATE_SUSPENDED`, is assigned, then resumed with `NtResumeProcess`, so its own child inherits the job with no race.
- Auto start is the `HKCU\...\Run` value `Docker Zero` = `"<install exe>"`, plus `--start-distro` for the sub-option, which starts the runtime quietly (no browser).
- Shortcuts are `Docker Zero.lnk` in the Desktop and Programs known folders, written through `WScript.Shell`.
- Before writing either, the app copies itself to `%LOCALAPPDATA%\Programs\docker-zero\docker-zero.exe` unless it already runs from there.
- Auto start, shortcuts, the `.wslconfig` line and the certificate are read from the machine each time the window loads; `%APPDATA%\docker-zero\settings.json` stores only idle mode, browser and corner box.
- Idle mode: keepalive holds a running distro at launch and on selection; WSL default and `.wslconfig` release it. Choosing `.wslconfig` adds `instanceIdleTimeout=-1` under `[general]`, any other mode removes it (a file left with only `[general]` is deleted); a pending change is applied only by a confirmed Apply now (`wsl --shutdown`).
- Browser: Chrome by full path from `ProgramFiles`, `ProgramFiles(x86)` or `LOCALAPPDATA` (PT-D17), or a chosen exe.
- Certificate: trust state read from `Cert:\CurrentUser\Root` through `CertFindCertificateInStore` by SHA-1 thumbprint (no dialog); trust and removal through `certutil -user -addstore|-delstore Root`, which shows Windows' confirmation dialog (PT-D21); regenerate rewrites the P-256 leaf in `/etc/portainer-certs` with the PT-D19 profile and restarts Portainer only when its container has `--sslcert`, else exit 3, reported as "run setup first".

Why: lesson - DZ-D15 put the whole app process in the job; a browser that Open Portainer starts is then a member too and would be killed when Docker Zero quits. Reading machine state avoids a stored flag that drifts when the user deletes a shortcut by hand. The Run value and the shortcuts must survive a cleaned Downloads folder (DZ-D9). `Import-Certificate` and `Remove-Item` on the root store fail without UI (PT-D19, PT-D21), and the dialog is the user's consent.

What breaks if ignored: closing the app closes the user's browser; a shortcut to a deleted download; a `.wslconfig` change that silently does nothing until a reboot.

Rejected: the whole process in the job (DZ-D15) - kills the browser. A Startup folder shortcut for auto start - `HKCU Run` is simpler to read and remove (DZ-D6). Storing auto start and shortcut flags in `settings.json` - drift. `Import-Certificate` - fails without UI.

### DZ-D17: Lesson: a phase is complete only when its written gate was observed; actions that stop live workloads or touch the trust store wait for the user

Decision: a phase stays `active` with "gate not met" until every check its roadmap gate names was observed; build, unit and throwaway-distro results do not substitute. Stop distro, Shut down WSL and certificate store changes on the real setup run only after the user approves them, and Windows' certificate dialogs are answered by the user.

Why: lesson - the Day 2 widget and Settings phases were first marked complete although Stop distro, Shut down WSL and the certificate Remove, Install and Regenerate had not run on the real distro and store; it looked reasonable because the same code paths had passed on the throwaway `dz-test` distro. The session's goal check rejected it, both phases went back to `active`, and both gates were met only after the user approved: Shut down WSL also ended the user's `--rm` `terraform-mcp-server` containers, and the regenerated certificate replaced thumbprint `270FF7DCB695691011F84F72E477E4D799FE11FE` with `2AEA72AAA4FB15117E3F8AF0C392AA51DEF54EBF`. The Shut down click was first refused by the agent tool's auto-mode classifier as interfering with workloads and ran once auto mode was off.

What breaks if ignored: a roadmap that reports a gate as met which nobody observed; live containers or trust anchors changed without consent.

Rejected: counting the throwaway-distro integration test as the gate - it does not show the real distro or store changing. Answering the certificate dialogs for the user through UI Automation - a trust decision the user owns.

### DZ-D18: The admin is created through the init API with an app-chosen setup token and reset with the helper and a password hash

Closes DZ-Q2.

Decision: Day 1 starts Portainer with `--setup-token <random value made by the app>`, waits for the status endpoint, then calls `POST /api/users/admin/init` with the typed username and password and the same value in the `X-Setup-Token` header; the token and the password stay in memory and out of logs and command lines. Reset admin stops the container, runs `portainer/helper-reset-password --password-hash <bcrypt hash>` on the `portainer_data` volume and starts the container again; the helper updates the existing admin whatever its username. Recreating `portainer_data` stays the last resort.

Why: the requirement (DZ-D5) is a typed username and password, and `--admin-password-file` can only create the user `admin`. Observed on throwaway distro `dz-test` (Alpine, Docker Engine 28.3.3, Portainer `lts` = `2.45.2`, fresh volume per run, published on `127.0.0.1` only):

- `--admin-password-file`: admin check `204`, login `admin` `200`, wrong password `422`, no setup-token line and no password in the container log, the admin kept after a restart without the flag, username fixed `admin`.
- init API with `--setup-token`: without the header `403` ("Provide the X-Setup-Token header"), with it `200`, login with a custom username `200`, neither token nor password found in the container log.
- helper with `--password-hash` (bcrypt from `htpasswd -nbB -C 10`): on the default `admin` and on a custom-username admin it reported the password updated for that user, the new password logged in `200`, the old one `422`, and no `admin` user appeared.

What breaks if ignored: a password on a command line (`--password`) or a token in the log; the helper run while the container is up (the official steps stop it first); a Day 1 flow that cannot honour the typed username.

Rejected: `--admin-password-file` - works but fixes the username to `admin`; kept as the fallback if the init API changes. `--no-setup-token` - the official docs limit it to networks the user fully controls, and supplying the token costs nothing. Reading the token from the container log - needs log parsing; supplying it avoids that. Helper `--password` - a plain password on the command line. Recreating `portainer_data` as the default reset - destroys Portainer settings.

Not verified: the third-party report of an authentication bypass on the init endpoints (CVE-2026-55761) was not checked against Portainer advisories, and the app calls init only on `127.0.0.1` with a token it chose; Docker Hub pull terms for organisations (PT-Q6; anonymous pulls of `portainer/portainer-ce:lts` and `portainer/helper-reset-password` succeeded).

### DZ-D19: Setup wizard mechanics: resumable steps, token through an env file, pinned certificate, no `docker.ps1`

Decision:

- The wizard runs eight steps in order (check WSL, install the distro, install Docker Engine, create the user, create the certificate, start Portainer, create the admin, add `docker.cmd`); each step checks what exists and reports `skipped`, so a run after a failure resumes and a run on a finished setup changes nothing.
- The distro is installed with `wsl --install --distribution Ubuntu-24.04 --web-download --no-launch`, plus `--name` and `--location` only when the name or folder differs from the default.
- The one-time setup token reaches Portainer through `PORTAINER_SETUP_TOKEN` in a root-only env file (`/run`, removed on exit), not through `--setup-token`; after the admin exists the container is recreated without it, so `docker inspect` shows no token.
- Admin API calls go to `127.0.0.1` over TLS verified against the leaf certificate read from the distro; the wizard never skips verification.
- Distro name, user, port, `docker.cmd` folder and PATH update are configuration; the production defaults are `Ubuntu-24.04`, `dev`, `9443`, and tests use a throwaway name with port `19443`.
- The wizard writes `docker.cmd` only, not `docker.ps1`; it overwrites `/etc/wsl.conf` with systemd and the default user.

Why: a failed apt or Docker Hub step must not force a reinstall; the flag form of the token would sit on a `docker run` line and in the container config for good (refines DZ-D18, which said the token stays off command lines); the status probe skips certificate checks because it sends nothing secret, while the admin call sends the password, so it pins the certificate; port `9443` belongs to the real Portainer, so a second distro in tests needs another port; PowerShell prefers a `.ps1` over a `.cmd` of the same name and a Restricted execution policy would then break `docker`, while `docker.cmd` alone was verified in PowerShell and `cmd` (PT-D20; the precedence is not tested here).

What breaks if ignored: a non-resumable wizard repeats a 100 MB download after any late failure; a token in the container config outlives its purpose; skipping certificate verification sends the admin password to whatever answers on the port; a test on port `9443` collides with the real setup; a `docker.ps1` can disable `docker` in PowerShell on a locked-down machine.

Rejected: `--setup-token` on the command line - persists in the config. `--no-setup-token` - see DZ-D18. `InsecureSkipVerify` for the admin call - password to an unverified peer. Writing `docker.ps1` as DZ-D5 listed - risk above, no gain over `docker.cmd`. Skipping the recreate after init - leaves the token behind.

Observed: `go test -tags integration ./internal/setup` with `DZ_SETUP_DISTRO` on a new distro took 169.8 s (first run 2 min 35 s): Ubuntu image registered, Docker Engine 29.9.2 from Docker's apt repository, the DNS drop-in, user in the `docker` group, the certificate served with `--sslcert`, Portainer 2.45.2, admin created, sign-in with the typed credentials `200`, a wrong password rejected, no setup token in the container config, a second run skipping the five steps that had work, and the distro unregistered afterwards. This also covers the flow on `Ubuntu-24.04` with systemd and the certificate flags, which DZ-D18 had left open.

### DZ-D20: Maintenance mechanics: update skips a current image, reset sends a hash, uninstall order and self-removal

Decision:

- Update pulls `portainer/portainer-ce:lts`, compares the container's image id with the tag, and recreates the container (same volume, certificate, no token) only when they differ.
- Reset admin hashes the password in the app with bcrypt cost 10 and sends only the hash; the helper script stops the container, runs `portainer/helper-reset-password --password-hash`, and starts the container again from an exit trap; the typed username is used only for a sign-in that proves the new password.
- One action at a time across setup, update, reset and uninstall.
- Uninstall order: optional `wsl --export` (refused if the file exists; a partial file is deleted on failure), remove the trusted certificate, unregister the distro, remove auto start, shortcuts, `docker.cmd` with its PATH entry and the `.wslconfig` line, remove settings; a failed step stops the run and a rerun skips finished steps.
- The `.wslconfig` line is removed only when the stored idle mode is `wslconfig`.
- The install folder is deleted by a detached `cmd.exe` started when the user quits after the uninstall, refused unless the folder is named `docker-zero`.

Why: the helper writes Portainer's database and the official steps stop the container first; the trap restarts it even when the helper fails; a password in a script breaks DZ-D18; the certificate thumbprint is only readable while the distro exists, so the trust step runs before the unregister and aborts it when Windows' dialog is declined; a running exe cannot delete its own folder; a `.wslconfig` line the user wrote is not the app's to remove (PT-D22).

What breaks if ignored: a restart on every update click; the helper run on a live container; a distro removed while its certificate stays trusted; an export over an existing backup; an uninstall that deletes a user-written `.wslconfig` line or a wrong folder.

Rejected: a tag picker (DZ-D7 keeps `lts`). Recreating `portainer_data` for reset - destroys settings. Deleting the install folder at the uninstall step - the running exe locks it. Continuing after a declined certificate dialog - the thumbprint is lost with the distro.

Observed: `TestMaintenance` on a throwaway distro passed in 577 s with every outcome in the roadmap's Phase 8 table.

### DZ-D21: Adoption mechanics: four replaceable items, read-only detection, replace before remove, one offer

Decision:

- Detection reads the distro state and, only while the distro runs, the container state and the certificate; it never starts a stopped distro and never changes anything.
- Four items can be replaced, each chosen by a checkbox: the Portainer shortcuts (matched by the `open-portainer.ps1` argument), the Startup keepalive shortcut (matched by `wsl.exe` and `--exec true`), the manual docker command folder with its PATH entry, the `.wslconfig` idle line.
- The replacement is made first and the legacy part removed after, so a failure leaves the old one; Docker Zero shortcuts that already exist are kept; the Startup shortcut becomes the `HKCU Run` value with `--start-distro`; the legacy folder loses only the files the guide created (launcher, icon, certificate copy, `docker.cmd`, `docker.ps1`) and then empty folders.
- Taking over the idle line sets the stored idle mode to `wslconfig` and drops the keepalive hold; the line itself is not edited.
- The app opens the setup window once per profile when it finds a setup (stored flag `adoptPrompted`).
- The distro, the container and the data volume are never touched by adoption.

Why: a reinstall destroys the data of early users; a shortcut is the user's, so only one whose arguments match the guide's is replaced; the idle line is the app's to remove at uninstall only once it manages it (DZ-D20); a window that reopens at every launch is noise.

What breaks if ignored: two sign-in starts and two launchers; a Portainer.lnk the user made for something else deleted; a legacy folder holding the user's own files emptied; a container restart during adoption.

Rejected: a fresh install only (DZ-D8). Recreating the container under the app's flags - risks the data and the user's settings. Matching shortcuts by name alone - removes unrelated links. Editing the `.wslconfig` line - it already holds the value the idle mode needs.

Observed: the run on this workstation in the roadmap's adoption table; the container id and start time were identical before and after. Defect found in the release run: with the distro stopped the container state is unread, so the card vanished with its result once the last item was replaced; the card now also stays while a result is shown.

### DZ-D22: Repository licence: MIT

Decision: license the repository under MIT, in `LICENSE` at the root; the copyright holder reads `Docker Zero contributors`.

Closes DZ-Q3.

Why: the user's answer; the release phase ships a licence file; MIT is the shortest permissive text; a contributors holder keeps usernames out of the public repository (DZ-D1).

What breaks if ignored: without a licence file the public code is all rights reserved and nobody may reuse it; a personal name in the holder line leaks an identifier.

Rejected: Apache-2.0 - adds an explicit patent grant and a NOTICE convention the project does not need.

### DZ-D23: Nothing at launch runs a command inside a stopped distro

Decision: every window is created hidden at launch and loads its page then, so code a page runs on load must not run `wsl.exe -d` against a stopped distro; the settings page reads the distro state first and, while it is not running, skips the certificate read and shows "Start Ubuntu-24.04 to read its certificate.", with Regenerate disabled.

Why: any `wsl.exe -d <distro>` command starts that distro; the settings page read `/etc/portainer-certs/portainer.crt` on load.

What breaks if ignored: a stopped distro starts about 2 s after every app launch, undoing Stop distro and the WSL default idle mode (observed on this workstation before the fix).

Rejected: reading the Windows-side copy in the install folder - it exists only after a trust and goes stale after a regenerate. Loading the settings page only when first shown - Wails creates the window and its page together.

### DZ-D24: Launch opens the setup window when nothing is installed; the take-over offer needs launchers to replace

Decision:

- At launch, while `Ubuntu-24.04` is not installed, the setup window opens by itself, at every launch.
- The once-per-profile take-over offer (`adoptPrompted`) fires only when detection lists at least one replaceable item; a running container alone does not count.
- The offer waits for the setup window's first `WindowRuntimeReady` event before showing it.

Supersedes the launch-offer rule of DZ-D21.

Why: a tray icon alone, often in the overflow area, tells a first-time user nothing; the app's own setup leaves a running container, which made the next launch offer to take over that same setup; Wails `v3.0.0-beta.30` `WebviewWindow.Focus` dereferences a nil WebView2 controller when called before the window exists.

What breaks if ignored: a first run that looks like nothing happened; an "Existing setup found / Nothing left to replace" window after every fresh setup; a crash at launch on a clean profile (panic in `edge.(*Chromium).GetController`, observed with a console build).

Rejected: offering on a running container (DZ-D21 as first built) - fires after the app's own setup. Setting `adoptPrompted` at the end of setup - leaves the same false offer on a profile whose settings were deleted. Showing without focus - the window can open behind others.

### DZ-D25: WebView2 data beside the settings, both folders removed after exit

Decision: set `WindowsOptions.WebviewUserDataPath` to `%APPDATA%\docker-zero\webview`; the uninstall data step deletes `settings.json` only; the detached `cmd.exe` started at quit deletes both the install folder and the settings folder, each refused unless named `docker-zero`.

Supersedes the settings-removal rule of DZ-D20 (`os.RemoveAll` of the settings folder during the run).

Why: by default Wails puts WebView2 data in `%APPDATA%\docker-zero.exe`, which uninstall never removed; WebView2 keeps its folder locked while the app runs, so it can only go after exit.

What breaks if ignored: an uninstall that leaves the browser cache and local storage of the app windows behind; deleting the settings folder during the run fails on the locked WebView2 files.

Rejected: WebView2 data in the install folder - a portable run from Downloads would create the install folder. Keeping the default folder and adding it to the removal - its name `docker-zero.exe` fails the folder-name guard, and loosening the guard widens what a bad path could delete.

Observed: the second live uninstall on this workstation removed the install folder, the settings folder with `webview`, the Run value, both shortcuts, the PATH entry and `.wslconfig`; `TestRemoveAfterExitDeletesFolders` passed.

### DZ-D26: Lesson: Windows line endings and a shared WebView2 folder broke checks during the release run

Decision: write repository files with LF (Python `write_bytes`, or `newline=""`), never `Path.write_text`, and never `git stash` in this working copy; end the `msedgewebview2.exe` processes of a debug build before launching another build on the same data folder.

Why: Python `write_text` on Windows writes CRLF; `core.autocrlf=true` checks files out with CRLF on `git stash pop`; `gofumpt` in `wails3 task lint` reports a CRLF Go file as unformatted. WebView2 refuses a user data folder held by a browser process started for a different exe name, and the app then exits without a window.

What breaks if ignored: lint fails on files whose content is correct; a release exe that silently exits looks like a product defect.

Rejected: nothing; these were process mistakes, each fixed by converting the files back to LF or by ending the stray processes.

Observed: `gofumpt` flagged `settings_service.go`, then `paths.go` and `host_windows.go` after a stash; the release exe exited twice while a console build named `dz-dbg.exe` still held `%APPDATA%\docker-zero\webview`.

### DZ-D27: The first release waits for the design phase

Decision: do not publish `v0.1.0` until Phase 10 (brand logo as app icon, UI upgrade, illustrated README, landing page) is done; the release exe is rebuilt from the commit that ends Phase 10.

Why: the user's instruction; Phase 10 replaces the app icon, the look of every window and the README screenshots.

What breaks if ignored: the first public download carries an icon and a UI about to change, and the README images stop matching the published exe.

Rejected: publish `v0.1.0` at `da09172` now and ship the design in a later version - the user declined it.

### DZ-D28: UI inventory: mocked frontend for every state, real exe for the main windows, committed

Decision: Phase 10 item 2 renders every window and state from the Svelte frontend in a browser with the Wails bindings replaced by per-state mocks, in light and dark; the built exe, driven by UI Automation, supplies reference screenshots of the main windows; the inventory page and its images are committed to the repository.

Closes DZ-Q4.

Why: several states exist only under real conditions (distro not installed, setup running, a failed pull, a stopped distro); mocks make each one repeatable for the before and after compare of item 4, and the real screenshots catch where the mock and WebView2 differ. Committed at the user's choice.

What breaks if ignored: real-exe-only capture needs an uninstall or a blocked Docker Hub per state and per redesign; mock-only capture hides a drift between the mock and the Go services.

Rejected: real exe only, and mocks only - the costs above. The inventory as a gitignored private Artifact, or a local file only - the user chose to commit it. Committed images must carry no account path, user name or machine name (public repository).

### DZ-D29: The rebuilt release exe re-runs the full clean-profile gate

Decision: before `v0.1.0` is published, the exe built from the commit that ends Phase 10 passes the Phase 9 gate again: Uninstall all, then setup from scratch through the window, on this account.

Closes DZ-Q5.

Why: Phase 10 changes the setup window and every other window the first gate exercised with the exe built at `da09172`.

What breaks if ignored: the published file is one whose setup flow was never run end to end.

Rejected: a smoke run on the existing setup (launch, each window, Start and Stop, Settings actions) - leaves the setup window flow unrun on the published exe. Cost accepted: about 8 min of setup plus an uninstall, and the Portainer admin password is random again afterwards.

Amended by DZ-D32: the exe the gate runs on is the one downloaded from the GitHub release the workflow created, not a local build.

### DZ-D30: README and landing page state no affiliation with Docker Inc. or Portainer

Decision: the README (Phase 10 item 5) and the landing page (item 6) carry one plain statement that Docker Zero is not affiliated with or endorsed by Docker Inc. or Portainer; the logo avoids the whale and both projects' marks (amended by DZ-D31: a whale as an animal is allowed, Docker's whale graphic is not).

Why: the name contains a Docker trademark and the app installs Docker Engine and Portainer CE, so a reader can take it for a vendor product.

What breaks if ignored: an implied endorsement on a public page.

Rejected: no statement - leaves the implication standing. Docker's trademark guidelines were not checked; the user checks them before the release.

### DZ-D31: Logo: an orca leaping out of a red ring over a two-tone sea; design rounds kept as a public gallery

Decision:

- The mark is a stylised orca (one arched silhouette, one white belly sweep, one eye patch) leaping from behind a red ring (the zero) towards the front left, its flukes behind the ring's lower right; inside the ring, below the orca, a light and a dark teal band of sea on a cream fill.
- Sources: `design/brand/logo.svg` (mark) and `design/brand/logo-tile.svg` (app icon on a cream tile); `packaging/appicon.png` at 1024 px and `packaging/windows/icon.ico` are rendered from them.
- Tray: the mark at 32 px, embedded as PNG in `internal/trayicon`; the ring is grey when stopped, amber when partly up, green when running (a red ring would read as a fault), and a red dot in the corner marks a failure.
- The seven design rounds (sheets, notes, choices) live in `design/brand/gallery/`, committed and public, linked from the landing page; the frontdoor design board of the logo is built there directly, never in `.frontdoor/`; the boards of other design work (README, landing page) use `.frontdoor/`, untracked (rule in `AGENTS.md`, DZ-D34).
- `board.py build` rewrites the page subtitle to "Local file, not committed"; after each build it is edited back to "Design rounds of the Docker Zero logo".

Amends DZ-D30: a whale as an animal is allowed; Docker's whale graphic (side view, containers on the back) is not.

Why: the user's brief (a whale passing a complete ring: Docker plus Zero) and choices over seven rounds: the sunset-and-sea idea and the orca, the orca travelling right to left from back to front, an arched orca in the style of the user's reference images, at most two flat beach elements so the mark stays clean at 16 px. The gallery is the user's request, to show visitors how the logo was reached.

What breaks if ignored: a mark redrawn from memory drifts from the chosen one; tray icons drawn in code stop matching the app icon; detail added inside the ring makes the 16 px tray icon unreadable.

Rejected: the blue tile with a white Z (the placeholder look); 17 other first-round concepts; gradient sunset scenes (round 2, too much detail); a front-view or fin-only orca (round 3); the first stylised orca (round 4, user disliked the vector); tray icons drawn in Go (cannot reproduce the mark); the gallery kept only in the ignored local board (lost to visitors); the board in `.frontdoor/` with a copy in the repository (two identical trees that drift once either is rebuilt); the gallery under `docs/brand/` (user moved it to `design/`, beside the UI inventory); stock reference images in the gallery (copyrighted, used only as a style brief, nothing traced).

### DZ-D32: The release exe is built and published by a GitHub Actions workflow on tag push

Decision: a workflow in `.github/workflows/` triggers on a pushed tag matching `v*`, builds `bin/docker-zero.exe` on `windows-latest` with `wails3 task build`, and creates the GitHub release for that tag with the exe and its SHA-256 attached; no locally built exe is published. Owned by Phase 9.

Amends DZ-D29: the clean-profile gate (Uninstall all, then setup from scratch through the window) runs on the exe downloaded from the release the workflow created, before the release is announced.

Why: the user wants a tag push to be the only step that makes a release appear; the published file is then reproducible from the tag and needs no hand upload.

What breaks if ignored: a gate run on a local exe says nothing about the published file, because the runner build has a different SHA-256; a hand-uploaded exe is not traceable to a commit.

Rejected: a local build published by hand - the user wants tag push to release automatically. A separate Phase 11 for the workflow - the user chose Phase 9. Test and lint on every pull request - not chosen; the workflow covers the tag push only. Lesson: the first workflow ran `wails3 task test` before the build; it looked fine because `frontend/dist` exists locally. On the clean runner `go test ./...` failed with `pattern all:frontend/dist: no matching files found` in `main.go` (`frontend/dist` is gitignored and created by the build), tag `v0.1.0-alpha` made no release, and the fix was to run the test step after the build.

Observed: tag `v0.1.1-alpha` ran every step green and created a pre-release whose exe matched its published SHA-256. Unverified: the unsigned runner-built exe behaves as the local one did on a managed machine (DZ-D12).

### DZ-D33: The release workflow stamps the exe version from the git tag, on the runner only

Decision: before `wails3 task build`, the workflow reads the tag, requires `vX.Y.Z` or `vX.Y.Z-suffix` (else fails), writes `X.Y.Z` to `fixed.file_version` and the tag without `v` (suffix kept) to `info.0409.ProductVersion` in `packaging/windows/info.json`, and never commits the change. The committed file holds `0.0.0` and `0.0.0-dev`, so a local build reads as a dev build.

Why: the tag is the one source of the version; `file_version` takes digits only, so a suffix such as `-alpha` can live only in `ProductVersion`.

What breaks if ignored: a hand-bumped file drifts from the tag (tag `v0.2.0`, exe reading `0.1.0`); a workflow that commits back needs write access to `main` and leaves automatic commits behind the tag.

Rejected: bump by hand before tagging - drifts. Commit the stamp back to the repository - branch protection and commit noise. Keep `0.1.0` in the committed file - a local build then looks like the release. Not stamped: `version` in `packaging/config.yml` and the manifest `assemblyIdentity` version, both still `0.1.0`; they do not feed the exe's file or product version.

### DZ-D34: README: story-led, no UI screenshots, drawn figures, banner with the logo as the letter o

Decision:

- The README is written story-first for a new visitor: the Docker Desktop licence wall, then pure WSL where every step and every later check is a command, then Docker Zero as the GUI that lets developers work with containers faster. Practical sections follow: what you get, get started, day to day, safe by design, questions, build from source.
- The README carries no screenshots of the app. The five PNGs in `docs/images/` are deleted. Drawn figures replace them, each a light and a dark WebP pair in `design/readme/asset/`: hero (banner); one banner per story paragraph (licence, commands, GUI); setup (the form and the eight steps); tray (the popup and the four ring colours); keep-the-distro-running (the three modes); safe by design (no UAC prompt, no Windows service, nothing under Program Files, certificate for localhost only). Seven navigation buttons and six icons for the `What you get` rows are small WebP pairs: the logo's red ring and two-tone sea as a badge around one line glyph. The `What you get` rows are not a table: on a phone GitHub squeezed the icon column to about 12 px. Each of the six rows is a paragraph: a 36 px icon inline before a bold text title, then the description (the user chose this over a drawn header image per row, after comparing both). Figure labels are the app's own labels, taken from `internal/setup`, the Svelte views and the README text; sources and the render commands are in `design/readme/source/` (`figures.html`, `icons.html`).
- Banner: the name set in Dela Gothic One with the logo as both letters o (Docker, Zero), the logo's two-tone sea along the bottom edge, the line "No admin prompt. No licence talk. No terminal marathon." (moved from a bold line under the banner, which wrapped and left one word alone) and, under it, the tagline "Docker and Portainer on Windows, set up for you and kept one click away." (the line says the benefit, the tagline says what it is and what the app does for the reader, the setup and the daily use; the earlier tagline "Docker Engine and Portainer on Windows. Zero admin rights." repeated the admin point already made by the name and the line); the type is chosen from six typographic directions of the second round on the README board.
- Layout: every figure sits directly under its section heading, never at the end of the section; the story has one banner above each of its three paragraphs; the navigation under the download button is seven image buttons on one row, each wrapped in an anchor, with the label as alt text; labels are one word where possible (Story, Features, Start, Daily, Safety, FAQ, Build) because the section names made the row about 930 px wide and it wrapped five and two on github.com, while the short labels make it about 610 px.
- Download button: the first thing under the one-line pitch, above the navigation. One image button "Download for Windows" (light and dark WebP) linking to `releases/latest/download/docker-zero.exe`, with no caption line under it (the user removed it; the WSL requirement stays in "You need", and the SHA-256 link moved into the Install step). The image carries no version number. The Install step points back at the button.
- The licence statement links Docker's own licence page ("Professional use in larger organizations") and tells the reader to check current terms with legal or procurement; nothing is stated about thresholds or prices.
- The README design board is the frontdoor default `.frontdoor/` (untracked). `design/brand/gallery/` holds the logo rounds only (amends the board rule recorded in DZ-D31).

Why: the user asked for a whole new README in a marketing tone with the story behind the app, and for the screenshots to be removed; in review the user found the download call to action missing and asked for it before anything else; the user asked for visual navigation buttons, figures right under the headings, one banner per story paragraph (a single banner after three paragraphs read as a wall of text), banners for the keep-running and safe-by-design sections, and icons in the `What you get` table; the UI is redesigned in Phase 10 items 3 and 4, so screenshots taken now would show a retired design. Drawn figures use only labels the code carries.

What breaks if ignored: screenshots of the old UI beside a new logo and a redesigned app; figure labels that drift from the app; a licence claim worded as fact instead of pointing at Docker's page; the download link breaks until a stable release exists, because `releases/latest` skips pre-releases (checked: the API answered 404 while the only release was the `v0.1.1-alpha` pre-release, which was then deleted with both test tags, so it answers 404 with no release at all), so the link starts working at the first tag without a suffix and the user accepted the dead link until then.

Rejected: keep the five screenshots and re-shoot them after item 4 (user asked to remove them). One three-panel story figure after the story text (the user asked for one banner per paragraph). Figures at the end of their section. Plain text navigation anchors (the user asked for buttons). Download button below the story or only in Install (the user wants it first). Download link to the releases list (works for pre-releases but costs a click), or pinned to the `v0.1.1-alpha` tag (goes stale on every release); the user chose the direct latest link. A second button for the releases page (diluted the one action; the Install step links the releases). A caption line under the button (removed in review). A table with an icon column for `What you get` (squeezed to about 12 px on a phone). A drawn header image per row with icon and title (the user chose the inline icon). Icons floated left of the text, or the whole row drawn as one image (misalign or cannot be copied). Navigation in two rows of four and three, or smaller buttons with the full section names (the text would be too small to read at the height that fits one row). A "pre-release" word on the caption (not chosen). The version number inside the image (stale at each release). First README banner round: a horizon scene, a sounding chart, a paper sun (the sun lost the zero ring). Other typefaces for the wordmark: wide geometric, contrast serif, lowercase mono (kept as sketches on the local board). Tagline "Docker Engine and Portainer on Windows, with no admin rights." (replaced by the Zero wording the user picked). Taglines not picked after the line moved into the banner: "Docker Engine and Portainer on Windows." (no new fact), "..., from the tray.", "Docker Engine on WSL2 with Portainer, from the Windows tray." (WSL2 is unfamiliar to a new reader), "Docker for Windows without Docker Desktop." (names Docker Desktop in the banner, closer to the trademark line), no tagline (the banner would not say what the product is), "One exe, one setup window: ...", "Docker Engine in WSL2, Portainer as the GUI, a tray icon to run it.", "A Windows tray app that sets up Docker Engine and Portainer." (first pick, replaced by the user with the set-up-and-one-click wording). Quote above the tagline was chosen over tagline above quote. Taglines not picked earlier: "Docker for Windows PCs where you are not an admin.", "Zero admin rights. Zero Docker Desktop.", "Docker Engine on WSL2 with Portainer, from the Windows tray.", "Docker and Portainer on Windows, no UAC prompt." Stating Docker's employee and revenue thresholds in the README (they change; the page is linked instead). Not checked: the figures against the redesigned UI (item 4), the README on github.com in light, dark and phone width, and the text fragment of the Docker link in a browser.

### DZ-D35: Landing page: root `index.html`, rising-tide hero with a leaping pod, interactive body figures

Decision:

- Placement: `index.html` at the repository root with an empty `.nojekyll`, served by Pages from `main` and `/ (root)`; the user switches Pages on. Social preview `design/landing/asset/og.png` (1200x630 PNG) is rendered from `design/landing/source/og.html`; the favicon is `design/brand/logo-tile.svg`; the design boards stay in `.frontdoor/` (rule in `AGENTS.md`).
- Links: GitHub links and the Open Graph URLs are absolute (`https://github.com/<owner>/docker-zero/...`, `https://<owner>.github.io/docker-zero/...`), because a page on Pages cannot reach github.com with a relative path. The public repository address is allowed there. The `AGENTS.md` rule against usernames means real usernames used inside the application (the user's reading), not the public repository address. Keep the address out of prose, alt text and comments.
- Action: one button "Download for Windows" linking `releases/latest/download/docker-zero.exe` (as the README, DZ-D34), repeated in the closing band; the latest release and its SHA-256 are linked in the Setup section.
- Voice: marketing pitch and storytelling in English; every claim from the README or the app's own labels (tray states Needs setup, Stopped, Starting, Running from `frontend/src/lib/status.svelte.ts`); the licence wall links Docker's page with the README's text fragment and states no threshold (DZ-D34); the no-affiliation statement stays in the footer (DZ-D30); the unsigned-exe SmartScreen prompt is stated.
- Hero: the headline and the download button over a wide sea of four drifting waves, a tray popup floating on the water. Start raises the tide in three steps while the popup rows turn on, then three orcas leap with splashes. Tapping the water, or a button, adds one leap. Once the tide is in, one orca leaps every 9 to 18 s at the side, only while the hero is on screen and the tab is visible. Under reduced motion the waves stop drifting and nothing leaps by itself, while Start and the tap still run at about a third of the duration. The orca is the logo's orca moved along a parabola behind the waves, so the water hides it below the surface.
- Body: three acts (The wall; The marathon, with a drag-to-compare slider between the manual commands and the Docker Zero window; The click), then Setup (eight-step demo and what you need), Day to day (four ring-state tabs and the four buttons), It stays up (the three idle modes as radio buttons), What you get, Safe by design, Questions, a closing download button; the two build commands carry Copy buttons.

Why: the user chose the rising tide of the first round's direction B as the start of the page and asked for more than the fidget alone, with orcas jumping out of the water; chose the horizon hero of the third round; asked for light repeating leaps so the page has life without blocking a tap on the water to make an orca jump; and chose the ring-colour tabs, the fidget cards and the slider of the second round for the body sections. Marketing pitch and storytelling is the user's answer on voice.

What breaks if ignored: a relative download link on Pages points into the Pages site, not the release; ambient leaps left on under reduced motion lock out Windows users who turned animation effects off; a root page without `.nojekyll` lets Jekyll convert the markdown under `docs/` into site pages; the download button is a 404 until the first tag without a suffix exists (DZ-D34).

Rejected: the first-round poster of the setup window as the whole page (the story becomes a row of cards). The dense card grid and the lifecycle timeline as whole pages (no hero moment and the story demoted; longest page with plain lists), kept as body figures instead. The hoop (the ring as a buoy, orcas threading it): left-to-right leaps only, a bare ring when stopped, lost behind the text on a phone, reads as a circus trick. The gauge (dusk to day, a four-mark tide gauge): a dark first impression on a cream brand, headline contrast dips when the text colour flips mid-transition, most code, gauge hidden under 700 px. Placement in `docs/` (the served folder would hold `docs/baseline/`). A Pages deploy through Actions (no workflow exists for it; a static file needs none). Relative links to GitHub (cannot reach it) and no GitHub links (no download). Leaps only on click (the user asked for life). Leaps that loop under reduced motion.

Lesson: the first build set the compare slider's Docker Zero window inside the dark band, so its text inherited cream from the band onto a cream card and the window showed ticks with no words; it looked right because the card's own background was set and its text colour was not. The section screenshots showed it; the structure check, run first without axe, did not name it. Fixed by giving the pane its own text colour.
