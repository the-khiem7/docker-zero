---
baseline_schema: "2.0"
pack: "portainer"
document: "hallucination"
status: "active"
updated: "2026-10-15"
code_ref: "uncommitted"
---

# portainer hallucination journal

## Entry index

| ID | About | Status | Related |
|---|---|---|---|
| PT-Q6 | Docker Hub terms and pull limits | open, deferred by the user | PT-D8, PT-D23 |
| PT-D1 | Runtime: Docker Engine in WSL2 with Portainer CE | closed | PT-D2 |
| PT-D2 | Portainer exposure: localhost only, no engine TCP port | closed | PT-D1, PT-Q1 |
| PT-D3 | Lesson: `docker pull` DNS failure and the IPv6 misread | closed | PT-D1 |
| PT-D4 | Lesson: Portainer unreachable after `docker start` | closed | PT-D5, PT-D16 |
| PT-D5 | Keepalive: hidden `sleep infinity` session | closed, kept as fallback (PT-D16) | PT-D4, PT-D16 |
| PT-D6 | Distro creation and default user | closed | PT-D1 |
| PT-D7 | Portainer restart policy `unless-stopped` | closed | PT-D2, PT-D4 |
| PT-D8 | Origin of the work and the licence reading | closed, accepted by the user (PT-D23) | PT-D9, PT-Q5, PT-Q6, PT-D23 |
| PT-D9 | Candidate set, filters and the options dropped | closed | PT-D8, PT-D10, PT-Q7 |
| PT-D10 | Rancher Desktop trial and its evidence | superseded by PT-D1 | PT-D1, PT-D9 |
| PT-D11 | Docker Desktop removal and what it showed | closed | PT-D8, PT-D10 |
| PT-D12 | Guide and commit handling | closed | PT-D1, PT-D11 |
| PT-D13 | Lesson: the wrapper is not a starter, and the user's four-point model | closed | PT-D4, PT-D5, PT-Q1, PT-D16 |
| PT-D14 | Portainer setup token and admin password handling | closed | PT-D2, PT-D7 |
| PT-D15 | Lesson: quoting through `wsl.exe` from Windows PowerShell 5.1 | closed | PT-D3 |
| PT-D16 | Idle timeout off in `.wslconfig` plus a sign-in shortcut | closed, sign-in observed | PT-Q2, PT-D4, PT-D5, PT-D13 |
| PT-D17 | One-click launcher for Portainer in Chrome | closed | PT-D15, PT-D16, PT-Q8 |
| PT-D18 | Portainer local environment observed after the setup wizard | closed | PT-Q3, PT-D2 |
| PT-D19 | Locally trusted leaf certificate for Portainer | closed | PT-Q8, PT-D2, PT-D17 |
| PT-D20 | Windows-side `docker` wrapper `docker.cmd` | closed | PT-Q1, PT-D2, PT-D15 |
| PT-D21 | Lesson: `wsl --export` stops the distro; certificate store removal needs a dialog | closed | PT-Q4, PT-D19 |
| PT-D22 | Uninstall script covers every created file; no further testing | closed | PT-Q4, PT-D21, PT-D17 |
| PT-D23 | Licences of Docker Engine and Portainer CE accepted for business use; Business Edition not needed | closed | PT-Q5, PT-D8, PT-D2 |
| PT-D24 | Podman Desktop rejected as the fallback | closed | PT-Q7, PT-D1, PT-D9 |

## Open questions

### PT-Q6: What do Docker Hub terms and pull limits require here?

Images pulled so far: `hello-world`, `alpine:3.20`, `portainer/portainer-ce:lts`, all anonymously. Docker Hub terms and anonymous pull limits are a separate agreement from Docker Desktop's subscription. Read on 2026-10-15: "Unauthenticated and Docker Personal users are subject to a 6-hour pull rate limit" of "100 per IPv4 address or IPv6 /64 subnet" for unauthenticated users, and "A Docker pull includes both a version check and any download that occurs as a result of the pull" (https://docs.docker.com/docker-hub/usage/pulls/); workstations behind one shared address would share that quota (inference, not stated). The Subscription Service Agreement says the Services exclude "content available in a registry or repository", which falls under the content owner's terms (https://www.docker.com/legal/docker-subscription-service-agreement/). Not found: a Docker statement on whether an organisation may pull public images anonymously without a subscription (the Terms of Service URL tried returned 404). The user deferred this question.

## Closed decisions

### PT-D1: Docker Engine in a WSL2 distro with Portainer CE

Decision: run the engine in `Ubuntu-24.04` under WSL2 and use Portainer CE as the GUI.

Why: Docker Engine is Apache 2.0 and carries no Docker Desktop subscription; WSL2 was already enabled; every install step ran without elevation.

What breaks if ignored: a Docker Desktop install brings the licence problem back; a path that needs machine-scope install stops at a UAC prompt the account cannot answer.

Rejected:

- Docker Desktop - licence concern, the reason for the work.
- Rancher Desktop 1.24.0 - the MSI is dual-purpose (`ALLUSERS=2`) but `msiexec /i ... MSIINSTALLPERUSER=1 ALLUSERS=""` still logged `MSI_LUA: Elevation required to install product` and exited 1602 after the unanswered credential prompt; the MSI summary word count is 2, without the no-elevation bit; SHA-512 of the download matched the release file before it was deleted.
- Podman Desktop - not tested; vendor docs offer a per-user install scope but list enabling WSL as an admin task and do not say whether creating the machine on enabled WSL needs admin.
- Finch - command line only, no GUI; its MSI was not tested for elevation.

### PT-D2: Portainer bound to localhost, no TCP port on the engine

Decision: publish Portainer as `127.0.0.1:9443` only, mount `/var/run/docker.sock`, and never open `dockerd` on TCP `2375`.

Why: the socket gives Portainer root-equivalent control of the engine; the TCP API has no authentication and WSL forwards localhost to the Windows host.

What breaks if ignored: any local process or user on a shared host drives the engine as root.

Rejected: TCP `2375` for a Windows `docker.exe` - unauthenticated API. TCP with TLS - certificate management for no gain here.

### PT-D3: Lesson - `docker pull` DNS failure and the IPv6 misread

What was believed: the failure was IPv6. The first run printed `dial tcp [2600:...]:443: connect: network is unreachable`, and `ip -6 route` showed no default route.

What disproved it: after `net.ipv6.conf.all.disable_ipv6=1` the error changed to `lookup registry-1.docker.io on 10.255.255.254:53: server misbehaving`, the same error the very first pull had printed. `curl` and `getent` in the same distro resolved the name, so the system resolver worked and the Go resolver inside `dockerd` did not, against the WSL DNS tunnel.

Decision: set `Environment="GODEBUG=netdns=cgo"` for `docker.service` through a systemd drop-in. IPv6 disabling was reverted; `docker run alpine:3.20` then worked with IPv6 on.

What breaks if ignored: every `docker pull` fails while network checks from the shell pass.

Rejected: disabling IPv6 - not needed and hid the real error.

### PT-D4: Lesson - Portainer unreachable after `docker start`

What was believed: after `wsl -d Ubuntu-24.04 -- docker start portainer`, the container keeps running and `https://localhost:9443` loads. The user ran exactly that and got `ERR_CONNECTION_REFUSED`.

What disproved it: the command returns when the container starts, no `wsl.exe` session remains, and WSL stops the idle distro, killing the container. Observed: `wsl -l -v` showed `Stopped` after 90 idle seconds, and the container came back only because a later command restarted the distro. Microsoft documents `instanceIdleTimeout` default 15000 ms and `vmIdleTimeout` default 60000 ms. An earlier handover told the user to run `docker start` without saying a session must stay open, and an earlier estimate of "about a minute" was imprecise.

Decision: the environment needs a session that stays open (PT-D5). The guide states the idle shutdown before the start commands.

What breaks if ignored: Portainer disappears seconds after any start command, and the user reads it as a browser ban.

### PT-D5: Keepalive is a hidden `sleep infinity` session

Decision: run `conhost.exe --headless wsl.exe -d Ubuntu-24.04 --exec sleep infinity`.

Why: one attached `wsl` session stops the idle shutdown; `--headless` removes the console window. Verified from a cold start: after 100 seconds the distro was `Running` and Portainer returned `200`. A shortcut created from the same command line, launched from a temporary folder, started it.

What breaks if ignored: PT-D4 recurs after every reboot and idle period.

Rejected: a `.vbs` launcher - not tested, and Windows Script Host may be blocked by policy. A minimized console shortcut - leaves a taskbar button.

### PT-D6: Distro created with `--no-launch`, default user set in `wsl.conf`

Decision: `wsl --install -d Ubuntu-24.04 --web-download --no-launch`, then create user `dev`, add it to `docker`, and set `[user] default=dev` with `[boot] systemd=true` in `/etc/wsl.conf`.

Why: `--no-launch` skips the interactive first-run prompt, so every step is scripted; Ubuntu 24.04 ships with systemd enabled; a non-root user in the `docker` group is the normal setup.

What breaks if ignored: an interactive first-run prompt blocks an unattended install; running as root hides permission problems.

Rejected: `--web-download` omitted, which needs the Store - not tested.

### PT-D7: Portainer restart policy is `unless-stopped`

Decision: create the container with `--restart unless-stopped`.

Why: it comes back with the engine after every distro start, and a container stopped by hand stays stopped.

What breaks if ignored: `always` restarts a container the user stopped, restarting the first-run timer and token (a new `setup_token` is logged on every start, 8 seen); no policy leaves Portainer down after every boot.

Rejected: `always` - ignores a manual stop.

### PT-D8: Origin of the work and the licence reading

Decision: Docker Desktop is out; Docker Engine and the other open-source components are acceptable (accepted by the user, PT-D23).

Why: the user is not allowed to use Docker Desktop because of an enterprise licence concern, and is allowed vanilla Docker (docker CLI, WSL and similar); the request was a Windows Docker client with no licence problem. Apache 2.0 and zlib permit commercial use with no seat, revenue or headcount threshold, which is a different position from a product subscription.

Licences read from each project's own LICENSE file or repository page:

| Component | Licence | Evidence |
|---|---|---|
| Rancher Desktop | Apache 2.0 | [LICENSE](https://raw.githubusercontent.com/rancher-sandbox/rancher-desktop/main/LICENSE#:~:text=Version%202.0%2C%20January%202004) |
| Podman Desktop | Apache 2.0 | [LICENSE](https://raw.githubusercontent.com/podman-desktop/podman-desktop/main/LICENSE#:~:text=Version%202.0%2C%20January%202004) |
| Moby (`dockerd`) | Apache 2.0 | [LICENSE](https://raw.githubusercontent.com/moby/moby/master/LICENSE#:~:text=Version%202.0%2C%20January%202004) |
| Finch | Apache 2.0 | [repository page](https://github.com/runfinch/finch#:~:text=Apache%2D2.0%20license) |
| Portainer CE | zlib | [LICENSE](https://raw.githubusercontent.com/portainer/portainer/develop/LICENSE#:~:text=This%20software%20is%20provided) |

Not verified: a secondary source stated that Docker Desktop needs a paid subscription for companies above 250 employees or 10 million USD annual revenue; Docker's own terms were not read. The web searches returned comparison articles (Spacelift, AlternativeTo, SigNoz and similar), used as leads and never as licence evidence. Only LICENSE files were read; components a vendor installer or an image bundles were not.

What breaks if ignored: keeping the Docker Desktop product restores the licence breach; treating the LICENSE files as covering bundled third-party parts overstates what was checked. This is factual information, not legal advice (PT-Q5, PT-Q6).

Rejected: taking licences from comparison articles - they are SEO-driven and inconsistent.

### PT-D9: Candidate set, filters and the options dropped

Decision: shortlist Rancher Desktop, Podman Desktop and Docker Engine in a WSL2 distro, with Portainer as the GUI for the engine route.

Why: the user answered that a GUI is wanted and Kubernetes is not; the first reply offered four routes (Engine in a WSL2 distro with the Docker CLI, Rancher Desktop, Podman Desktop, Finch), and a second search added Portainer as a GUI for a plain engine.

What breaks if ignored: a CLI-only tool fails the GUI requirement; a Kubernetes-centred tool adds weight the user declined.

Rejected, with the one reason each lost:

- Finch - command line only; Windows MSI not tested for elevation.
- Colima and Lima - CLI oriented, macOS and Linux focus, no GUI, per comparison articles.
- OrbStack - macOS only, per the same articles.
- Container Desktop - the LICENSE file fetched from `container-desktop/container-desktop` carried an unrelated copyright holder, so the project's identity and licence were not established.
- Portainer alone - a UI over an engine, not a runtime.
- Docker Desktop alternatives lists as a method - see PT-D8.

### PT-D10: Rancher Desktop trial and its evidence

Decision (at the time): install Rancher Desktop 1.24.0 for the current user, container engine dockerd (moby), Kubernetes off, Privileged Service skipped, then delete the MSI.

Status: superseded by PT-D1; the trial ended without installing anything.

Why it was chosen: closest to the Docker Desktop experience, a GUI, and the user accepted that recommendation and the trial, and declined Kubernetes. Vendor docs said a non-admin user can run it, the installer offers a current-user scope, the Privileged Service is the admin-only part, and WSL must exist before install. The docs did not say whether the current-user scope still prompts for elevation. The docs page was read through a summarizer, so statements taken from it carry no text fragment.

Evidence gathered before running anything:

- Release `v1.24.0`; the only Windows assets were `Rancher.Desktop.Setup.1.24.0.msi` (670,932,992 bytes) and its `.sha512sum`; there was no zip build.
- SHA-512 of the download matched the published file: `6428f330cdf8dc9dbbc07bbde270046c9e86fb91e90392c0d25ac764e64eaa8f9d2b3e7faaffed669a48f7fd7445dadda18072dd51c905bf59df4787b70ed77e`.
- Read-only inspection of the MSI tables through the Windows Installer COM object: `ALLUSERS=2`, `MSIINSTALLPERUSER=0`, manufacturer SUSE, `WixUI_Mode=RD`; custom actions `SetALLUSERS` (sets `ALLUSERS` to 1), `SetMSIINSTALLPERUSER` (sets it to 0), `DetectWSL`, `UpdateWSL`, `ErrorWSLNotInstalled`; no `ServiceInstall` table; `APPLICATIONFOLDER` under `ProgramFiles64Folder`; launch conditions Windows 10 or later and no downgrade; summary word count 2.
- Run: `msiexec /i <msi> MSIINSTALLPERUSER=1 ALLUSERS="" /qb /l*v <log>`. The log shows `ALLUSERS` deleted, `MSIINSTALLPERUSER` set to 1, `AlwaysInstallElevated` 0 on machine and user, `MSI_LUA: Elevation required to install product, will prompt for credentials`, then `User cancelled installation` and exit 1602. Nothing was installed under `%LOCALAPPDATA%\Programs` or `Program Files`.

What breaks if ignored: the trial is repeated and stops at the same credential prompt, or the admin boundary is probed.

Rejected: asking IT to install the MSI machine-wide - the user declined; it remains the fallback only if the WSL route fails. Bypassing the prompt - against `AGENTS.md` hard rule 4.

Raw output, MSI `Property` table (all 19 rows, `Property | Value`):

```text
UpgradeCode | {1F717D5A-A55B-5FE2-9103-C0D74F7FBDE3}
WixUIRMOption | UseRM
MSIINSTALLPERUSER | 0
ALLUSERS | 2
ARPNOMODIFY | 1
ApplicationFolderName | Rancher Desktop
ARPPRODUCTICON | RancherDesktopIcon.exe
ARPURLINFOABOUT | https://rancherdesktop.io/
Manufacturer | SUSE
ProductCode | {E99D3E3A-149F-4F15-BC82-7904DBB0BA47}
ProductLanguage | 1033
ProductName | Rancher Desktop
ProductVersion | 1.24.0
DefaultUIFont | WixUI_Font_Normal
WixUI_Mode | RD
WIXUI_EXITDIALOGOPTIONALCHECKBOXTEXT | Run Rancher Desktop
WixShellExecTarget | [#mainExecutable]
ErrorDialog | ErrorDlg
SecureCustomProperties | ALLUSERS;MSIINSTALLPERUSER;RDRUNAFTERINSTALL;WIX_DOWNGRADE_DETECTED;WIX_UPGRADE_DETECTED
```

Raw output, other MSI tables (`Action | Type | Source | Target` for `CustomAction`):

```text
CustomAction
ErrorWSLNotInstalled | 19 |  | 100
SetALLUSERS | 51 | ALLUSERS | 1
SetMSIINSTALLPERUSER | 307 | MSIINSTALLPERUSER | 0
DetectWSL | 1 | CustomActionBinary | DetectWSL
UpdateWSL | 1025 | CustomActionBinary | UpdateWSL

LaunchCondition
Installed OR VersionNT >= 603 | Windows 10 and above is required
NOT WIX_DOWNGRADE_DETECTED | A newer version of [ProductName] is already installed.

ServiceInstall: table not present
Registry: table not present

Directory (first 25 rows only; the rest are bundled app folders such as locales, bin, lib)
APPLICATIONFOLDER | ProgramFiles64Folder | <short>|Rancher Desktop
(one row among them: <id> | <parent> | <short>|docker-cli-plugins)

Summary information
WordCount (property 15) = 2   -> bit 8 (no elevation needed) not set
Template  (property 7)  = x64;1033
```

Raw output, selected install-log lines (timestamps, thread ids, working directory and process id removed):

```text
Machine policy value 'AlwaysInstallElevated' is 0
User policy value 'AlwaysInstallElevated' is 0
Running product '{E99D3E3A-149F-4F15-BC82-7904DBB0BA47}' with user privileges: It's not assigned.
Command Line: MSIINSTALLPERUSER=1 ALLUSERS= [...] CLIENTUILEVEL=2 [...]
PROPERTY CHANGE: Deleting ALLUSERS property. Its current value is '2'.
PROPERTY CHANGE: Modifying MSIINSTALLPERUSER property. Its current value is '0'. Its new value: '1'.
MSI_LUA: Elevation required to install product, will prompt for credentials
Windows Installer installed the product. Product Name: Rancher Desktop. Product Version: 1.24.0. Product Language: 1033. Manufacturer: SUSE. Installation success or error status: 1602.
User cancelled installation.
```

### PT-D11: Docker Desktop removal and what it showed

Decision: remove Docker Desktop completely, its WSL distro and its data, before adding the replacement.

Why: the user instructed deleting Docker Desktop and its WSL distro and allowed deleting all of its data; a leftover Docker Desktop `docker` CLI and socket could conflict with the new engine; the licence forbids keeping it.

Observed first: Docker Desktop `4.85.0` (winget listed `4.95.0` as available) was a per-user install, with its uninstall key under `HKCU` and files under `%LOCALAPPDATA%\Programs\DockerDesktop`; WSL held one distro, `docker-desktop`, default and stopped; no `docker` was on the user PATH. The licence, not administrator rights, was the blocker for Docker Desktop.

Removal: `winget uninstall --id Docker.DockerDesktop --scope user` succeeded and took the distro with it; then `%LOCALAPPDATA%\Docker`, `%APPDATA%\Docker`, `%USERPROFILE%\.docker` and a leftover `Programs\DockerDesktop\tmp-delete` were deleted after listing file names only, with no file contents read, and the Start Menu shortcut and the `docker-desktop` Start Menu folder removed.

Lesson: a check for a Docker `Run` autostart key matched nothing because the name `docker-desktop` was a Start Menu folder, not a registry value; `Remove-ItemProperty` failed harmlessly and the folder was removed by path.

What breaks if ignored: two Docker CLIs and sockets compete; Docker Desktop stays installed against policy.

Rejected: keeping the data - the user said none was needed.

Raw output, before removal (the tools print WSL text as UTF-16; shown here as plain text, `<LOCALAPPDATA>` standing for the user's `%LOCALAPPDATA%`):

```text
> wsl --status
Default Distribution: docker-desktop
Default Version: 2
WSL1 is not supported with your current machine configuration.
Please enable the "Windows Subsystem for Linux" optional component to use WSL1.

> wsl --list --verbose
  NAME             STATE      VERSION
* docker-desktop   Stopped    2

> winget list --name "Docker"
Name           Id                   Version Available Source
Docker Desktop Docker.DockerDesktop 4.85.0  4.95.0    winget

HKCU uninstall entry: DisplayName "Docker Desktop", DisplayVersion 4.85.0,
  UninstallString "<LOCALAPPDATA>\Programs\DockerDesktop\Docker Desktop Installer.exe" "uninstall"
HKLM uninstall entries matching Docker, Rancher or Podman: none
HKCU uninstall entries matching Rancher or Podman: none
docker CLI on PATH: none
Existed: <LOCALAPPDATA>\Docker = True; Rancher Desktop folder = False; Program Files\Docker = False

> wsl --version
WSL version: 2.7.10.0
Kernel version: 6.18.33.2-2
WSLg version: 1.0.73.2
MSRDC version: 1.2.6676
Direct3D version: 1.611.1-81528511
DXCore version: 10.0.26100.1-240331-1435.ge-release
Windows version: 10.0.26300.9711
```

Raw output, removal and aftermath:

```text
> winget uninstall --id Docker.DockerDesktop --scope user --silent --disable-interactivity
Found Docker Desktop [Docker.DockerDesktop]
Starting package uninstall...
Successfully uninstalled
exit: 0

> wsl --list --verbose
Windows Subsystem for Linux has no installed distributions.

Leftovers listed by file count and size before deletion (names only, no contents read):
<LOCALAPPDATA>\Programs\DockerDesktop   4 files    8.4 MB   (a tmp-delete folder)
<LOCALAPPDATA>\Docker                  66 files   25.3 MB   (log, run, tasks, wsl, lock files, install logs)
%APPDATA%\Docker                       11 files    0.5 MB   (cagent, .trackid, analyticsmonitor.dat, features-overrides.json, locked-directories, marlin.dat, unleash-v2-docker-desktop.json)
%USERPROFILE%\.docker                  31 files    4.7 MB   (buildx, cli-plugins, contexts, desktop-build, mcp, models, modules, .token_seed)

> winget list --name Docker
No installed package found matching input criteria.
```

### PT-D12: Guide and commit handling

Decision: write `guides/docker-engine-wsl/README.md` for the chosen route, record Docker Desktop, Rancher Desktop, Podman Desktop, Finch and TCP `2375` under rejected approaches, and leave commits to the user.

Why: `AGENTS.md` rule 6 requires untested steps to be labelled, and rule 1 keeps a public repository free of identifiers; the user said they commit themselves. The template section "Why the default installer needs admin" is answered with the licence and the WSL feature, because Docker Desktop itself installed per user (PT-D11). An unrelated `AGENTS.md` edit in the working tree belongs to the user and is not touched.

What breaks if ignored: an agent commit overrides the user's own staging; an unlabelled step is read as verified.

Rejected: committing on the user's behalf.

### PT-D13: Lesson - the wrapper is not a starter, and the user's four-point model

The user stated their model of the developer experience:

- A `docker.cmd` wrapper is an alias that starts Docker faster with the command `docker`.
- `wsl -d Ubuntu-24.04 -- docker start portainer` stops right after the container runs, so the UI is unreachable.
- The fix is an infinite headless process that runs `wsl.exe`.
- A shortcut in the autostart folder runs that process when Windows starts.

Corrections: the wrapper only forwards the client command to `wsl -d Ubuntu-24.04 -- docker %*` and starts nothing that stays up; the container does not stop itself, WSL stops the whole idle distro; the second and third points are right apart from that; the Startup folder runs at sign-in, not at boot, and the shortcut had not been created. Idle timers are documented by Microsoft as [`instanceIdleTimeout`](https://learn.microsoft.com/en-us/windows/wsl/wsl-config#:~:text=Set%20to%20%2D1%20to%20disable%20auto%20shutdown) (15000 ms, `-1` disables) under `[general]` and `vmIdleTimeout` (60000 ms) under `[wsl2]` in `%USERPROFILE%\.wslconfig`.

What breaks if ignored: the wrapper is built to solve the idle shutdown and does not; the shutdown is read as a container fault.

Rejected: presenting `.wslconfig` as proven - it was untested at that point; PT-D16 records the later test.

### PT-D14: Portainer setup token and admin password handling

Decision: put the setup token on the user's clipboard and never print it; the user chooses the admin password and it is never requested.

Why: the token lets its holder create the admin account, and Portainer logs a new one at every start (8 had been logged); the first-run page needs an admin username, a password of at least 12 characters and that token, which the first guide draft omitted and the user hit.

What breaks if ignored: a token or password in chat or a transcript; a guide that stalls at the first-run page.

Rejected: reading the token aloud in the reply - it is a credential.

### PT-D15: Lesson - quoting through `wsl.exe` from Windows PowerShell 5.1

What was believed: `wsl -d Ubuntu-24.04 -u root -- bash -c "<script with $VAR and inner double quotes>"` runs the script as written.

What disproved it: the output stopped after `id -un`, and `os-release` values came back empty, so the quoting was lost on the way in. `wsl` output is also UTF-16, so unfiltered text shows spaced characters.

Decision: write LF scripts into the scratchpad and run them as `bash /mnt/c/<path>/script.sh`; pass only inner-quote-free strings with `-c`.

What breaks if ignored: a command appears to succeed while running a different, truncated command.

Rejected: `[Console]::OutputEncoding` tweaks - they did not change what the tool returned.

### PT-D16: Idle timeout off in `.wslconfig` plus a sign-in shortcut

Closes PT-Q2.

Decision: set `instanceIdleTimeout=-1` under `[general]` in `%USERPROFILE%\.wslconfig`, and add the Startup shortcut `WSL Docker keepalive.lnk` running `conhost.exe --headless wsl.exe -d Ubuntu-24.04 --exec true`.

Why: `.wslconfig` stops WSL from shutting down the idle distro, and the shortcut starts the distro, and Portainer through `unless-stopped` (PT-D7), at sign-in with no window. `.wslconfig` alone leaves Portainer down after sign-in until some `wsl` command runs; the shortcut alone needs a permanent hidden session (PT-D5). The user chose both after both were explained and asked for explicit consent, because each changes the machine.

Observed:

- `%USERPROFILE%\.wslconfig` did not exist and was created with `[general]` and `instanceIdleTimeout=-1`; `wsl --shutdown` applied it and stopped the then-running distro and Portainer.
- `wsl -d Ubuntu-24.04 --exec true` after a full shutdown returned in 31.5 seconds (about 6 seconds when only the distro was stopped); 100 seconds later with no session open the distro was `Running` and `https://localhost:9443` returned `200`.
- The Startup folder held two unrelated items before and still holds them; the shortcut was launched by hand from a terminated distro and after 45 seconds the distro was `Running` and Portainer returned `200`.

Observed after the entry was written: after a sign-out and sign-in with no terminal opened, `wsl -l -v` showed `Ubuntu-24.04` as `Running` and `https://localhost:9443` showed the Portainer login page with no warning (the user's screenshots; whether it was a sign-out or a restart was not stated). Not verified: memory held by the resident WSL VM. `.wslconfig` applies to every WSL2 distro on the account, not only this one.

What breaks if ignored: PT-D4 recurs after any idle period and after every sign-in.

Rejected: the `sleep infinity` shortcut alone (PT-D5) - needs a permanent hidden process; kept as the fallback if `.wslconfig` is removed. `.wslconfig` alone - Portainer is not up after sign-in. A `vmIdleTimeout` change - the VM stayed up without it.

### PT-D17: One-click launcher for Portainer in Chrome

Decision: ship `guides/docker-engine-wsl/open-portainer.ps1`; `-CreateShortcut` copies it to `%LOCALAPPDATA%\Programs\docker-wsl` and creates `Portainer.lnk` on the Desktop and in the Start Menu, pointing at the copy, minimized window, Portainer logo as icon (Chrome icon when Portainer does not answer).

Why: it checks the page, then starts the distro and the container only when needed, waits for the page, then opens Chrome; the user asked for exactly that and chose the Desktop plus Start Menu form, a minimized console with progress lines, and a self-installing script. The copy keeps the shortcut valid if the repository folder moves.

Facts that shaped it:

- The `HKCU` `App Paths` key for `chrome.exe` points to another Chromium-based browser, so the script starts Google Chrome by its full path (`Program Files`, `Program Files (x86)`, `%LOCALAPPDATA%`) and falls back to the default browser.
- Execution policy for the current user is `RemoteSigned`, so the shortcut runs `powershell.exe -NoProfile -File <copy>` with no override.
- Lesson: under `$ErrorActionPreference = "Stop"`, Windows PowerShell 5.1 turns native stderr into a terminating error; the script wraps native calls in a helper that relaxes the preference and returns exit code and text.

Observed, run from the repository folder with `-NoBrowser -NoDialog` unless noted:

- Everything running: page answered, exit 0, 2.7 seconds including PowerShell start.
- Container stopped: started, ready in 5.6 seconds.
- Distro terminated: distro and container started, ready in 14.5 seconds.
- Missing container, missing distro, and a URL that never answers: each printed a clear error and exited 1.
- With the browser: a Google Chrome window titled `Portainer` opened and no other Chromium process was running.
- Desktop shortcut launched from a terminated distro: finished in 13.2 seconds, distro `Running`, page `200`.
- The guide's Uninstall lines for the launcher removed both shortcuts and the script folder; the shortcuts were then recreated.

What breaks if ignored: launching by `Start-Process chrome` opens the wrong browser; a hidden window hides a 30-second wait; an execution-policy override widens what the shortcut may run.

Icon: the script reads the `apple-touch-icon` PNG (180 by 180) that the running Portainer page lists in its `<link>` tags, wraps it in a single-image `portainer.ico` beside the installed script and points both shortcuts at it. Observed: `-CreateShortcut` wrote a 1862-byte icon, both shortcuts report `portainer.ico,0`, and the shell rendered the Portainer "P" for the Desktop shortcut; `System.Drawing.Icon.ToBitmap` cannot decode a PNG-compressed icon, so the shell render was the check. Rejected for the icon: committing the logo to the repository - a public repo should not carry a vendor trademark file; the 16 by 16 inline `shortcut icon` of the page - too small for a desktop icon; keeping the Chrome icon - the user asked for the Portainer one.

Rejected: `Start-Process chrome` - wrong browser. `conhost --headless` for the launcher - no progress during a slow start. `-ExecutionPolicy Bypass` - not needed under `RemoteSigned`. A Chrome bookmark through a `portainer://` protocol - the user chose shortcuts.

### PT-D18: Portainer local environment observed after the setup wizard

Decision: the setup wizard needs no further step; Skip on the Edge Compute step leads to a working `local` environment. Closes PT-Q3.

Observed in the user's screenshot of the Portainer Home page (Community Edition 2.45.2 LTS, signed in as `admin`): one environment `local`, status `Up`, type Standalone, Docker `29.9.0`, socket `unix:///var/run/docker.sock`, 1 of 1 containers running, 1 node, 8 cores, 8.2 GB memory, group Unassigned. The Environment selector reads `None selected` until the user opens `local`. The login page shown earlier in the same session (browser tab titled Portainer, `Not secure` badge) is the launcher's output (PT-D17).

What breaks if ignored: none; the question only gated the guide's statement that Portainer manages the local engine, which now rests on an observation.

Rejected: none.

### PT-D19: Locally trusted leaf certificate for Portainer

Decision: Portainer serves one self-signed leaf certificate that the user's account trusts through `Cert:\CurrentUser\Root`; the key stays in the distro. Closes PT-Q8.

Certificate: `CN=portainer` only, subject alternative names `localhost`, `127.0.0.1`, `::1`, `CA:FALSE` critical, `keyUsage` critical `digitalSignature`, `extendedKeyUsage` `serverAuth`, ECDSA P-256 with SHA-256, five years (2026-10-15 to 2031-10-14), SHA-256 fingerprint `E5:6B:81:73:0E:4C:64:DF:98:1E:A2:9E:16:0D:4E:55:83:44:F3:29:36:D7:85:80:4F:FB:83:C2:57:55:4D:6F`, thumbprint `270FF7DCB695691011F84F72E477E4D799FE11FE`. Key and certificate in `/etc/portainer-certs` (key root `0600`, certificate `0644`); public `.crt` copied to `%LOCALAPPDATA%\Programs\docker-wsl\portainer.crt`. Container recreated with `--sslcert`, `--sslkey`, the read-only `/certs` mount and the same `portainer_data` volume; the old container was renamed `portainer-old`, then deleted by the user.

Why a leaf and not a private CA: a CA key could sign certificates for any host. Costs the user accepted: a new trust anchor affects every application of the account that uses the Windows store; anyone holding the key can impersonate `localhost` to that account's browsers; the certificate expires and the warning returns. The user chose to skip the security-team confirmation and accepts the risk of a security alert on a managed device.

Observed: `openssl s_client` shows the served fingerprint equals the file's; `curl.exe` without `-k` returns `200`; a Chrome incognito window opens `https://localhost:9443` with no "Not secure" badge and the admin login works. A normal Chrome window that had earlier proceeded past the warning kept the badge although its certificate dialog read "Certificate is valid"; closing every Chrome window and reopening the page from the `Portainer` shortcut removed the badge, which confirms the cause was the remembered proceed choice. `Import-Certificate` fails with "UI is not allowed in this operation" in a non-interactive session; `certutil -user -addstore Root` succeeded there. The five-year term was accepted by Chrome.

Guide: the certificate is an optional step 4 under Integrate, because the warning is only cosmetic for a localhost page and the trust-store import affects the whole account. The Troubleshooting row for "Not secure" stays permanently, because a reader may skip the step.

What breaks if ignored: a CA instead of a leaf widens the trust to any host; storing the key on the Windows side or in the repository exposes it; skipping the SAN entries gives a name mismatch because Chrome ignores `CN`; the certificate expires silently in 2031.

Rejected: a private CA - key could sign any host. `-k` or a bypass click - leaves the warning. A `portainer` hostname - not in the SAN, so `https://portainer:9443` still errors.

### PT-D20: Windows-side `docker` wrapper `docker.cmd`

Decision: a one-line `docker.cmd` in `%LOCALAPPDATA%\Programs\docker-wsl\bin`, on the user PATH, forwards every call to `wsl.exe -d %DOCKER_WSL_DISTRO% -- docker %*` (default distro `Ubuntu-24.04`). Closes PT-Q1.

Why: PowerShell and `cmd` have no `docker`, and the engine lives in the distro; a wrapper opens no port and writes only under the profile. The user chose it over doing nothing and over a Docker context over SSH (needs `sshd` in the distro and a key). A TCP `2375` listener stays rejected (PT-D2).

Observed in a new PowerShell process after the PATH change: `docker --version` and `docker ps` answer, `docker run --rm alpine:3.20 echo hi` prints `hi`, `docker nosuchcmd` returns exit code `1`, `docker version` works from `cmd.exe`, and a quoted `sh -c "echo a b"` arrives intact. The PATH change reaches only processes started afterwards.

What breaks if ignored: a bind mount written as `-v "C:\work:/w"` fails with `invalid mode` because the engine reads the Windows drive colon as a mount option; write `/mnt/c/work`, or `$(wsl -d Ubuntu-24.04 -- pwd)` for the current directory (both tested). Quoting through `wsl.exe` is still limited by PT-D15.

Rejected: no wrapper (user asked for it). Docker context over SSH - extra service and key for no gain.

### PT-D21: Lesson: `wsl --export` stops the distro; certificate store removal needs a dialog

Decision: the guide's Uninstall states that `wsl --export` stops the distro and every container, and removes the certificate with `certutil -user -delstore Root portainer`. Related to PT-Q4 (still open for the residue listed there).

Why it looked reasonable: export reads like a read-only copy, so the real `Ubuntu-24.04` was exported to seed a throwaway distro for the uninstall test. Disproved when `wsl -l -v` showed `Ubuntu-24.04` as `Stopped` and `curl.exe` to `https://localhost:9443/` returned `000`. Fix: `wsl -d Ubuntu-24.04 --exec true` brought it back; Portainer answered `200` and the `portainer` container was `Up`.

Observed on a throwaway copy imported as `Ubuntu-test` and on dummies: `rm -r /etc/portainer-certs` inside it, `wsl --export` then `wsl --unregister` (exit 0, vhdx and `Lxss` entry gone), user PATH entry removal logic, dummy Startup shortcut removal. `Remove-Item` on a root-store certificate fails with "The operation is on user root store and UI is not allowed"; `certutil -user -delstore Root` and an `X509Store.Remove` call both wait on a confirmation dialog that no one can answer in a non-interactive session. The throwaway certificate `CN=portainer-uninstall-test` therefore stayed in `Cert:\CurrentUser\Root` until the user ran `certutil -user -delstore Root portainer-uninstall-test` in an interactive terminal, which completed successfully.

What breaks if ignored: exporting the real distro as a "backup" silently stops Portainer and every other container; the guide's `Remove-Item` pipeline fails in any session without a dialog.

Rejected: `Remove-Item` on the store - fails without UI. Exporting the real distro to seed the test - stops it; a fresh `wsl --install` would not.

### PT-D22: Uninstall script covers every created file; no further testing

Decision: the guide's Uninstall block deletes every item the setup creates, and the user declined further testing of it. Closes PT-Q4.

Why: the user wants the script to be complete, not re-run. The block removes the Startup shortcut, both `Portainer.lnk` files, the `docker-wsl` folder (launcher script, `bin\docker.cmd`, `portainer.crt`), the `portainer` certificate from `Cert:\CurrentUser\Root`, the `docker-wsl\bin` user PATH entry, the distro, a leftover Start Menu folder named after the distro, and the `instanceIdleTimeout` line of `%USERPROFILE%\.wslconfig` (file deleted when nothing else is in it), then `wsl --shutdown`. The `ubuntu-24.04-backup.tar` export stays as the user's backup.

Evidence: run earlier on a throwaway distro copy and dummies, and by the user for the certificate removal (PT-D21, PT-D17). The `.wslconfig` logic was run on three temporary files (only the line: file deleted; extra `[wsl2]` content: kept with only the `instanceIdleTimeout` line removed; absent file: no error); the real `.wslconfig` was not touched. Not run against the real distro: the Start Menu folder line and the real `.wslconfig` edit; the guide labels both untested.

What breaks if ignored: a leftover `instanceIdleTimeout=-1` keeps every WSL2 distro on the account from idling out; a leftover trusted `portainer` root stays trusted after Portainer is gone; a leftover PATH entry points at a deleted folder.

Rejected: testing the real unregister - destroys the working environment. `Remove-Item` for the certificate - fails with "UI is not allowed" (PT-D21).

### PT-D23: Licences of Docker Engine and Portainer CE accepted for business use; Business Edition not needed

Decision: Docker Engine and Portainer CE stay in use in the business environment, and Portainer Business Edition is not bought, because the user reads Portainer's zlib licence as permitting business use. Closes PT-Q5.

Status of confirmation: the user gave this reading; no legal or procurement confirmation is recorded and none was requested. Still not found in primary sources: Portainer's image terms and trademark policy, a CE versus Business Edition feature table (RBAC, SSO, audit), CE node or user limits, Canonical's terms for redistributing the Ubuntu image, and any bundled component a vendor installer would have carried (none of whose bundled parts were read, PT-D8).

Primary-source reading on 2026-10-15 (quotations checked against the fetched page text; WebFetch summaries alone were not trusted):

- Moby, docker/cli, docker/compose and containerd LICENSE files are Apache 2.0. Docker's install page states "Apache License, Version 2.0. See LICENSE for the full license." (https://docs.docker.com/engine/install/ubuntu/).
- Docker's Desktop licence page says the licensing of Docker Engine "aren't changing" (https://docs.docker.com/subscription/desktop-license/). The Engine page ties the paid-subscription rule to "Commercial use of Docker Engine obtained via Docker Desktop" at companies over 250 employees or US$10 million revenue (https://docs.docker.com/engine/). No page states that Engine installed separately inside WSL2 is outside the Desktop rule; that reading is an inference.
- Portainer's LICENSE is zlib and says "Permission is granted to anyone to use this software for any purpose, including commercial applications" (https://raw.githubusercontent.com/portainer/portainer/develop/LICENSE).
- Portainer's production FAQ says "we do not recommend it, especially if you're using Portainer for business critical apps" (https://docs.portainer.io/faqs/getting-started/can-i-use-portainer-ce-in-a-production-environment). LTS is supported "up to 9 months" and CE support is through "community support channels" (https://docs.portainer.io/start/lifecycle).
- Portainer "no longer collects anonymous usage statistics" from 2.38.0 (https://docs.portainer.io/faqs/getting-started/what-information-does-portainer-collect); the running `2.45.2` is newer.
- Portainer CE mounts the Docker socket, which a security reviewer will flag (PT-D2).

What breaks if ignored: treating the user's reading as a legal confirmation; the Ubuntu image terms, Portainer image terms and the Docker Hub question (PT-Q6) remain unchecked, and an organisation audit may ask for them.

Rejected: Portainer Business Edition - not needed for the features in use. Waiting for legal confirmation before use - the user chose to proceed on the licence text.

### PT-D24: Podman Desktop rejected as the fallback

Decision: Podman Desktop is not the fallback if the WSL Docker route stops working. Closes PT-Q7.

Why: the user judges Podman's behaviour different from Docker's. Never tried; vendor docs offered a per-user install scope and listed enabling WSL or Hyper-V as an admin task without saying whether creating the WSL machine needs elevation. A Red Hat Jira story (`RUN-3245`) proposes a rootless Windows install; whether it shipped is unknown.

What breaks if ignored: none now; with no fallback, a failure of the WSL route leaves no Docker-compatible runtime until a new route is chosen.

Rejected: Podman Desktop - behaviour differs from Docker. Rancher Desktop stays rejected (PT-D10).
