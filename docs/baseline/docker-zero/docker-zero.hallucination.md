---
baseline_schema: "2.0"
pack: "docker-zero"
document: "hallucination"
status: "draft"
updated: "2026-10-22"
code_ref: "uncommitted"
---

# docker-zero hallucination journal

## Entry index

| ID | About | Status | Related |
|---|---|---|---|
| DZ-Q2 | Portainer admin init and reset mechanism | open | DZ-D5, DZ-D7 |
| DZ-D1 | Brand, audience and repository identity | closed | DZ-D2 |
| DZ-D2 | Imported `portainer` pack and guide are reference only | closed | DZ-D1, DZ-D8 |
| DZ-D3 | Widget form: tray popup plus optional corner box | closed | DZ-D4 |
| DZ-D4 | Day 2 actions: system-level only, two stop modes, live status | closed | DZ-D3 |
| DZ-D5 | Day 1 wizard covers every step; WSL feature not enabled by the app | closed | DZ-Q2, DZ-D6 |
| DZ-D6 | Settings: auto start, shortcuts, idle mode, browser, certificate | closed | DZ-D5, DZ-D7 |
| DZ-D7 | Portainer `lts`, Update button, admin reset, uninstall danger zone | closed | DZ-Q2, DZ-D6 |
| DZ-D8 | Adopt an existing manual setup | closed | DZ-D2 |
| DZ-D9 | Wails3 and portable self-installing exe | closed | DZ-D11, DZ-D10 |
| DZ-D10 | Exe policy spike runs before app code | closed | DZ-D9 |
| DZ-D11 | Frontend framework: Svelte | closed | DZ-Q1, DZ-D9, DZ-D12 |
| DZ-D12 | Wails3 Svelte template name and tray spike result | closed | DZ-D10, DZ-D11 |
| DZ-D13 | Scaffold shape: `packaging/` folder, pinned tools, Windows only | closed | DZ-D9 |
| DZ-D14 | How the WSL runner reads state and probes status | closed | DZ-D4 |
| DZ-D15 | Day 2 widget mechanics: hold, job object, polling, confirm, second launch | closed | DZ-D3, DZ-D4, DZ-D6 |

## Open questions

### DZ-Q2: How does the app create and reset the Portainer admin?

User requirement: the wizard asks for username and password and creates the admin; reset must exist, and if Portainer offers no account reset, the app resets Portainer itself. Candidates to verify against Portainer `lts` docs before Phase 6: the admin init API endpoint and whether it needs the setup token; the `portainer/helper-reset-password` image for a password reset; recreating `portainer_data` as the full reset (destroys Portainer settings). None verified.

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

Decision: the wizard runs, with live log and per-step status: detect WSL (stop with "ask IT to enable WSL" if absent), `wsl --install -d Ubuntu-24.04 --web-download --no-launch`, user `dev` and `wsl.conf` with systemd, Docker Engine from Docker's apt repository, the `GODEBUG=netdns=cgo` drop-in (PT-D3), the leaf certificate (PT-D19), the Portainer container with the certificate on `127.0.0.1:9443` with `unless-stopped` (PT-D2, PT-D7), the Portainer admin from a username and password the user types (DZ-Q2), and `docker` wrappers for `cmd` (`docker.cmd`) and PowerShell (`docker.ps1`) on the user PATH (PT-D20). Distro `Ubuntu-24.04` and user `dev` are fixed.

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

Decision: pin image `portainer/portainer-ce:lts`; an Update Portainer button pulls it and recreates the container on the same `portainer_data` volume and certificate; a Reset admin action per DZ-Q2; an Uninstall all action in a Danger zone, with two confirmations and an offered `wsl --export` backup, removing everything the app created.

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
- The app puts itself in a job object with `JOB_OBJECT_LIMIT_KILL_ON_JOB_CLOSE` before the first hold, so every child dies with the app.
- Polling: every 5 s while the popup or the corner box is visible, every 1 s for 15 s after an action, never otherwise; one reading at launch. Hiding a window re-sends the last reading without probing.
- Shut down WSL asks for confirmation inside the popup, not with a native dialog.
- A second launch of the exe opens the popup (`SingleInstance.OnSecondInstanceLaunch`).
- Windows are told apart by `?view=popup` and `?view=corner`; the tray icon is drawn in Go and its dot follows the status level.

Why: observed - a force-killed build without the job object left `/bin/sleep infinity` running in the distro, which would keep it up forever; with the job object the hold ended. Job membership is inherited at process creation, so `wsl.exe` and the `wsl.exe` it starts are both covered without a race; assigning each child after `Start` would miss a grandchild created first. A native dialog takes focus and the popup hides on focus loss (`HideOnFocusLost`). The second-launch hook gives shortcuts and the Start Menu entry a visible effect and made UI Automation testing possible.

What breaks if ignored: the distro never idles out after a crash; the confirmation dialog closes the popup that asked for it; a hidden app polls `wsl.exe` forever.

Rejected: assigning only the hold's `wsl.exe` to a job after `Start` - races with its own child. A native confirmation dialog - hides the popup. Polling while hidden - wastes CPU in an app meant to sit idle in the tray.
