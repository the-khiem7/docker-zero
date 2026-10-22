---
baseline_schema: "2.0"
pack: "portainer"
document: "introduction"
status: "active"
updated: "2026-10-15"
code_ref: "uncommitted"
---

# portainer - Docker Engine in WSL2 with a Portainer GUI

## Scope

Replace Docker Desktop on one Windows workstation with Docker Engine running in a WSL2 distro and Portainer CE as the GUI, using only standard-user rights. In scope: the workstation setup, the guide `guides/docker-engine-wsl/README.md`, its row in the `README.md` Guides table. Out of scope: Kubernetes, a Windows-side `docker.exe` (a `docker.cmd` wrapper covers the need, PT-D20), Edge Compute, production hosting.

## Origin

- The user is not allowed to use Docker Desktop because of an enterprise licence concern, and is allowed vanilla Docker (docker CLI, WSL and similar). The request was for a free-licence Docker client for Windows (PT-D8).
- Constraints given along the way: GUI wanted, Kubernetes not wanted, delete Docker Desktop and its distro with all its data, delete the downloaded installer, do not ask IT, write a guide, the user commits (PT-D9, PT-D11, PT-D12).
- Starting state observed: Docker Desktop `4.85.0` installed per user, WSL2 enabled with one stopped default distro `docker-desktop`, no other distro, no `docker` on the user PATH (PT-D11).
- Path taken, in order: licence research (Phase 1), route chosen (Phase 2), Docker Desktop removed (Phase 3), Rancher Desktop per-user trial stopped at an elevation prompt (Phase 4), pivot to Docker Engine in WSL2 (Phase 5), Portainer (Phase 6), idle-shutdown lesson and keepalive (Phase 7), guide (Phase 8).

## Options considered

| Option | Result | Detail |
|---|---|---|
| Docker Desktop | Disallowed, removed | PT-D8, PT-D11 |
| Rancher Desktop 1.24.0 | Rejected, MSI needs elevation | PT-D10, PT-D1 |
| Podman Desktop | Rejected as the fallback, behaviour differs from Docker | PT-D24, PT-D1 |
| Finch | Not chosen, CLI only | PT-D9 |
| Colima, Lima, OrbStack, Container Desktop | Dropped | PT-D9 |
| Docker Engine in WSL2 | Chosen | PT-D1 |
| Portainer CE as the GUI | Chosen | PT-D1, PT-D2 |

## Current truth

Inspected state: working tree on top of HEAD `0d9ec66`. The guide was committed there; later edits to it are uncommitted. Machine state below was observed running on the workstation.

| Item | State |
|---|---|
| Docker Desktop | Removed with `winget uninstall` at user scope; its WSL distro, `%LOCALAPPDATA%` and `%APPDATA%` folders, `%USERPROFILE%\.docker` and Start Menu entries deleted |
| Rancher Desktop | Not installed; the MSI and the install log were deleted (PT-D10) |
| WSL | `2.7.10.0`, kernel `6.18.33.2-2`, default version 2; enabled before this work |
| Distro | `Ubuntu-24.04` (24.04.5 LTS), systemd on, disk under `%LOCALAPPDATA%\wsl\<guid>\ext4.vhdx`, default Linux user `dev` in group `docker` |
| Docker Engine | `29.9.0`, Compose `v5.6.0`, from Docker's apt repository, enabled as a systemd service |
| Engine override | `/etc/systemd/system/docker.service.d/override.conf` sets `GODEBUG=netdns=cgo` (PT-D3) |
| Portainer CE | `2.45.2` (image tag `lts`), container `portainer`, `127.0.0.1:9443`, volume `portainer_data`, restart policy `unless-stopped`, Docker socket mounted |
| Portainer first run | Admin account created and Skip chosen on Edge Compute; Home shows environment `local` as `Up`, Standalone, Docker `29.9.0`, 1 of 1 containers running (PT-D18) |
| Idle behavior | WSL's default stops the distro after the last `wsl.exe` session ends, taking every container with it (PT-D4); `%USERPROFILE%\.wslconfig` now sets `instanceIdleTimeout=-1` and the distro stayed up 100 seconds with no session (PT-D16) |
| Sign-in start | Startup shortcut `WSL Docker keepalive.lnk` runs `conhost.exe --headless wsl.exe -d Ubuntu-24.04 --exec true`; launched by hand it brought the distro and Portainer up; after a real sign-in with no terminal open the distro was `Running` and the Portainer login page loaded (PT-D16) |
| Keepalive session | The `sleep infinity` session (PT-D5) is a fallback and is not running |
| One-click launcher | `Portainer.lnk` on the Desktop and in the Start Menu runs `%LOCALAPPDATA%\Programs\docker-wsl\open-portainer.ps1` in a minimized window with the Portainer logo as icon; it starts the distro and container if needed and opens Google Chrome by full path (PT-D17) |
| Certificate | Portainer serves a locally trusted leaf certificate (`CN=portainer`, ECDSA P-256, valid to 2031-10-14), imported into `Cert:CurrentUserRoot`; a Chrome incognito window shows no warning (PT-D19) |
| Windows-side `docker` | `docker.cmd` in `%LOCALAPPDATA%\Programs\docker-wsl\bin` on the user PATH forwards to the distro; bind mounts need Linux paths (PT-D20) |
| Guide | Written, indexed in `README.md`, repo checks clean on the working tree; the Uninstall block deletes every created file, the Start Menu folder and real `.wslconfig` lines untested (PT-D22) |

Topology, outside to inside: Windows host, browser or `wsl.exe` client, WSL2 virtual machine, `Ubuntu-24.04`, `dockerd`, containers. Windows reaches Portainer through WSL localhost forwarding on `127.0.0.1:9443`.

## Target

- Docker Engine and Portainer start without an open terminal after sign-in, with no admin prompt.
- The guide documents only steps run as a standard user, and labels the rest untested.

## Constraints

| Constraint | Reason |
|---|---|
| Standard user, no UAC prompt, no writes to `C:\Program Files`, `HKLM`, services | `AGENTS.md` hard rule 3; the account cannot elevate |
| WSL2 must already be enabled | Enabling the feature needs admin; the guide stops and asks IT otherwise |
| Docker Desktop not allowed | Enterprise licence concern |
| Do not ask IT for a machine-wide install | The user declined it for now (PT-D10) |
| The user commits; the agent does not stage or commit | The user's instruction (PT-D12) |
| No credential or token is printed in chat or written to a file | Secrets handling (PT-D14) |
| The certificate's private key stays in the distro; only the public `.crt` reaches Windows | A key in the repository or on the Windows side lets anyone impersonate `localhost` to the account's browsers (PT-D19) |
| Public repository: no usernames, company names, absolute user paths, profile names | `AGENTS.md` hard rule 1; this pack is committed with the repository |
