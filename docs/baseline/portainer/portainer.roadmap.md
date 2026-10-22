---
baseline_schema: "2.0"
pack: "portainer"
document: "roadmap"
status: "active"
updated: "2026-10-15"
code_ref: "uncommitted"
---

# portainer roadmap

| Phase | Depends on | Status | Acceptance |
|---|---|---|---|
| Phase 1 | - | complete | Request understood, licences of the candidates read |
| Phase 2 | Phase 1 | complete | Route chosen with the user's answers |
| Phase 3 | Phase 2 | complete | Docker Desktop and its residue removed |
| Phase 4 | Phase 3 | abandoned | Rancher Desktop installed per user |
| Phase 5 | Phase 3 | complete | Engine runs in WSL2 and pulls images |
| Phase 6 | Phase 5 | complete | Portainer serves on localhost and the admin account exists |
| Phase 7 | Phase 6 | complete | Engine and Portainer up after sign-in with no terminal open |
| Phase 8 | Phase 6 | active | Guide matches what ran, checks clean, committed |
| Phase 9 | Phase 7 | complete | One-click shortcut opens Portainer in Chrome, starting what is stopped |
| Phase 10 | Phase 9 | complete | Portainer opens in Chrome with no certificate warning |
| Phase 11 | Phase 6 | complete | `docker` runs from PowerShell and `cmd` through a wrapper |

## Phase 1: Request and licence research

Status: complete. The user is not allowed Docker Desktop because of an enterprise licence concern and is allowed vanilla Docker (docker CLI, WSL and similar), and asked for a free-licence Docker client for Windows. Three web searches returned comparison articles, used as leads. LICENSE files read: Rancher Desktop, Podman Desktop, Moby, Portainer, and the Finch repository page. Results and the unverified claims: PT-D8 (origin and licence reading). The Docker Desktop threshold claim and the Container Desktop identity stayed unverified.

## Phase 2: Route selection

Status: complete. The first reply offered four routes (PT-D9). The user then answered: WSL2 check wanted; GUI wanted; Kubernetes not wanted; extra research on the licences of Rancher Desktop and Podman Desktop in enterprise use, and on options beyond them; yes to a Docker Desktop installation check and an installer per-user check; delete Docker Desktop and its distro; no data to keep; delete the downloaded MSI; try the install; write a guide. The recommended route was Rancher Desktop with dockerd and Kubernetes off (PT-D10).

## Phase 3: Remove Docker Desktop

Status: complete. Evidence: `winget uninstall --id Docker.DockerDesktop --scope user --silent` printed `Successfully uninstalled`; `winget list` finds no Docker; the `docker-desktop` WSL distro was gone afterwards; `%LOCALAPPDATA%\Programs\DockerDesktop`, `%LOCALAPPDATA%\Docker`, `%APPDATA%\Docker`, `%USERPROFILE%\.docker`, one Start Menu shortcut and one Start Menu folder were deleted; no Docker `Run` key existed. The state before removal and the lesson on the autostart check: PT-D11.

## Phase 4: Rancher Desktop trial

Status: abandoned. Evidence: the per-user `msiexec` run logged `MSI_LUA: Elevation required to install product, will prompt for credentials` and exited 1602 with nothing installed. The MSI and the install log were deleted afterwards. The user then reported that the installer required admin permission, and declined to ask IT for a machine-wide install. Full evidence and the MSI table readout: PT-D10.

## Phase 5: Engine in WSL2

Status: complete. Evidence: the user chose the WSL route; `wsl --install -d Ubuntu-24.04 --web-download --no-launch` ran with no UAC prompt; systemd was on by default; Docker Engine `29.9.0` and Compose `v5.6.0` installed from Docker's apt repository; user `dev` joined group `docker` and became the default user in `/etc/wsl.conf`; `systemctl is-active docker` returns `active`; `docker run --rm hello-world` and `alpine:3.20` ran as `dev`. Pulls first failed on DNS (PT-D3); distro and user choices: PT-D6.

## Phase 6: Portainer

Status: complete. Evidence: `curl -k` on `https://127.0.0.1:9443` returned `200` from inside the distro and from Windows. The first browser attempt failed with `ERR_CONNECTION_REFUSED` because the distro had stopped (PT-D4). After the keepalive was started the user opened the setup page, was given the setup token on the clipboard (PT-D14), created the admin account and reached the Edge Compute step; Skip was advised. The Home page then showed the `local` environment `Up`: PT-D18. Exposure and restart policy: PT-D2, PT-D7.

## Phase 7: Keep the environment running

Status: complete; the choice is made (PT-D16) and applied. Evidence: the idle shutdown was reproduced (PT-D4); the headless keepalive held the distro for 100 seconds as a first fix (PT-D5); the user chose `.wslconfig` plus a Startup shortcut; `%USERPROFILE%\.wslconfig` with `instanceIdleTimeout=-1` was created and applied with `wsl --shutdown`; with no session open the distro stayed `Running` for 100 seconds and Portainer returned `200`; the shortcut `WSL Docker keepalive.lnk` was created beside two unrelated items and, launched by hand from a terminated distro, brought both back within 45 seconds. A cold start after a full shutdown takes about 31 seconds. After a real sign-out and sign-in with no terminal opened, `wsl -l -v` showed `Ubuntu-24.04` as `Running` and the Portainer login page loaded with no warning (the user's screenshots). The user's model of the idle shutdown and its corrections: PT-D13.

## Phase 8: Guide

Status: active. `guides/docker-engine-wsl/README.md` and its Guides-table row were committed by the user in `eaa4839`. Uncommitted edits since: Portainer setup token and 12-character password, idle-shutdown wording, Start and stop table, Startup shortcut, uninstall additions, the one-click launcher section with `open-portainer.ps1`, and the optional certificate step with its Verify, Uninstall and Troubleshooting lines. The repo checks from `AGENTS.md` (secret scan, non-ASCII, identifiers, `.ps1` syntax) were clean on the working tree that includes the last guide edit and this pack. Handling rules: PT-D12.

## Phase 9: One-click launcher

Status: complete. `guides/docker-engine-wsl/open-portainer.ps1` exists and `Portainer.lnk` exists on the Desktop and in the Start Menu, pointing at the copy in `%LOCALAPPDATA%\Programs\docker-wsl`. Evidence: everything running 2.7 seconds; container stopped 5.6 seconds; distro terminated 14.5 seconds; the Desktop shortcut from a terminated distro 13.2 seconds with the page at `200`; three error paths exit 1 with clear messages; Google Chrome opened a tab titled Portainer. The guide gained the launcher section, Verify rows, Use row, Uninstall lines and Troubleshooting rows; the Uninstall lines for the launcher were run and the shortcuts recreated. Decisions and facts: PT-D17. The shortcut icon is the Portainer logo read from the running page and saved as `portainer.ico`; verified by the shell rendering it for the Desktop shortcut (PT-D17).

## Phase 10: Locally trusted certificate

Status: complete. Certificate generated, container recreated, certificate imported, `curl.exe` without `-k` returns `200`, a Chrome incognito window shows no warning and the admin login works; the guide has the steps, Verify row, Uninstall line and Troubleshooting rows. A normal Chrome window that had proceeded past the warning kept the badge until every Chrome window was closed and the page reopened. Fields, costs, rollback: PT-D19.

## Phase 11: Windows-side docker wrapper

Status: complete. `docker.cmd` and the user PATH entry exist; in a new PowerShell process `docker --version`, `docker ps`, `docker run --rm alpine:3.20 echo hi` and `docker nosuchcmd` (exit code `1`) behave as expected, and the guide has the step, a Verify row and the Uninstall lines. Decision and the bind-mount limit: PT-D20.

## Risks

| Risk | State |
|---|---|
| Portainer with the Docker socket is root-equivalent on the engine | Mitigated by localhost binding (PT-D2); other local users on a shared host can still reach the forwarded port, not verified |
| Components bundled in installers or images may carry other licences | The user accepted the licences on the Apache 2.0 and zlib texts without legal confirmation; image terms, trademark policy and the Ubuntu image position were not found (PT-D23) |
| Docker Hub terms and pull limits | Pull limits read; no statement found on organisational anonymous pulls; the user deferred the question (PT-Q6) |
| Uninstall steps | Block complete; Start Menu folder line and real `.wslconfig` edit untested (PT-D22) |
| No fallback runtime if the WSL route fails | Podman Desktop rejected (PT-D24) |

## Next action

1. Phase 8: the user commits the guide and pack changes, including `guides/docker-engine-wsl/open-portainer.ps1`. The guide's Integrate section must match PT-D16.
2. PT-Q6 (Docker Hub terms for organisational anonymous pulls) stays open until the user resumes it.
