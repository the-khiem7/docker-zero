# Docker Engine in WSL2 + Portainer

| Field | Value |
|---|---|
| Tool | [Docker Engine](https://docs.docker.com/engine/) (Apache 2.0) running in a WSL2 Linux distro, with [Portainer CE](https://www.portainer.io/) as the browser GUI - a Docker Desktop replacement with no Docker Desktop licence |
| Tested version | Docker Engine `29.9.0`, Compose `v5.6.0`, Portainer CE `2.45.2`, Ubuntu `24.04.5 LTS`, WSL `2.7.10.0` |
| Method | `wsl --install --web-download` for the distro, Docker's official apt repository for the engine, Portainer as a container |
| Admin needed | No, provided WSL2 is already enabled on the machine |
| Install location | `%LOCALAPPDATA%\wsl\<guid>\ext4.vhdx` (the distro disk) |
| Script | [`open-portainer.ps1`](open-portainer.ps1) - one-click launcher: starts the distro and Portainer if needed, then opens Chrome |

## Why the default installer needs admin

Docker Desktop is the vendor's Windows product. It is subject to a paid subscription in larger organizations (check the current Docker terms and ask your legal or procurement team), which is why many companies forbid it. Docker Engine itself is open source and free, so the route here runs the engine inside a WSL2 distro and skips Docker Desktop entirely.

WSL2 itself is a Windows feature. Enabling it needs administrator rights, so this guide treats an already enabled WSL2 as a prerequisite to request from IT. Installing a distro and everything inside it needs no elevation: the distro lives in the user profile, and `sudo` or root inside the distro is Linux root in the virtual machine, not a Windows UAC prompt.

## Rejected approaches

| Approach | Why rejected |
|---|---|
| Docker Desktop | Licence concern in an enterprise, the reason for this guide. |
| Rancher Desktop (MSI, version 1.24.0) | The installer is a dual-purpose MSI, but running it for the current user (`MSIINSTALLPERUSER=1 ALLUSERS=""`) still raises a UAC prompt. The Windows Installer log shows `Elevation required to install product`, and the package does not carry the flag that declares no elevation is needed. It needs IT to install it machine-wide. |
| Podman Desktop | Not tested. The vendor docs offer a per-user install scope, but state that enabling WSL needs administrator rights and do not say whether creating the Podman machine on an already enabled WSL needs them too. |
| Finch | CLI only, so it does not give a GUI. Its Windows MSI was not tested for elevation. |
| Docker Engine on TCP `2375` for a Windows `docker.exe` | The API has no authentication and is root-equivalent on the engine. Reachable by every local process once WSL forwards localhost. Not used. |

## Prerequisites

- A Windows standard user account.
- WSL2 enabled. Check with `wsl --status`: it must report `Default Version: 2`. If WSL is not installed, ask IT to enable it and stop here.
- Network access to `download.docker.com`, `registry-1.docker.io` and the Ubuntu WSL image host.

## Install

### 1. Install the Ubuntu distro

`--no-launch` skips the interactive first-run user prompt, so every later step is scripted.

```powershell
wsl --install -d Ubuntu-24.04 --web-download --no-launch
wsl -d Ubuntu-24.04 -u root -- bash -c "ps -p 1 -o comm="
```

The second command prints `systemd`. Ubuntu 24.04 images enable systemd by default through `/etc/wsl.conf`, which Docker needs to run as a service.

### 2. Install Docker Engine from Docker's apt repository

Open a root shell in the distro with `wsl -d Ubuntu-24.04 -u root`, then run the following (it is Docker's documented apt method for Ubuntu).

```bash
set -e
apt-get update
apt-get install -y ca-certificates curl

install -m 0755 -d /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/ubuntu/gpg -o /etc/apt/keyrings/docker.asc
chmod a+r /etc/apt/keyrings/docker.asc

. /etc/os-release
cat > /etc/apt/sources.list.d/docker.sources <<EOF
Types: deb
URIs: https://download.docker.com/linux/ubuntu
Suites: ${UBUNTU_CODENAME:-$VERSION_CODENAME}
Components: stable
Architectures: $(dpkg --print-architecture)
Signed-By: /etc/apt/keyrings/docker.asc
EOF

apt-get update
apt-get install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin
systemctl enable --now docker
docker --version
```

### 3. Work around the WSL DNS tunnel (still in the root shell)

The Go DNS resolver inside `dockerd` can fail against the WSL DNS proxy, so `docker pull` reports `lookup registry-1.docker.io on 10.255.255.254:53: server misbehaving` while `curl` in the same distro works. Force the system resolver for the engine only.

```bash
install -d /etc/systemd/system/docker.service.d
cat > /etc/systemd/system/docker.service.d/override.conf <<'EOF'
[Service]
Environment="GODEBUG=netdns=cgo"
EOF
systemctl daemon-reload
systemctl restart docker
```

### 4. Create a Linux user for daily work (still in the root shell)

Running as root works, but a non-root user in the `docker` group is the normal setup. The name `dev` is an example.

```bash
adduser --disabled-password --gecos "" dev
usermod -aG docker dev
cat > /etc/wsl.conf <<'EOF'
[boot]
systemd=true

[user]
default=dev
EOF
exit
```

Back in PowerShell, restart the distro so the default user takes effect.

```powershell
wsl --terminate Ubuntu-24.04
```

### 5. Start Portainer CE

Portainer is a container that talks to the engine through its socket. Bind it to `127.0.0.1` only, and use `unless-stopped` so a container you stopped by hand stays stopped after a restart.

```powershell
wsl -d Ubuntu-24.04 -- docker volume create portainer_data
wsl -d Ubuntu-24.04 -- docker run -d --name portainer --restart unless-stopped -p 127.0.0.1:9443:9443 -v /var/run/docker.sock:/var/run/docker.sock -v portainer_data:/data portainer/portainer-ce:lts
```

Open `https://localhost:9443` and accept the self-signed certificate warning. The first-run page asks for an admin username, a password of at least 12 characters and a one-time setup token. Portainer prints a new token to its log at every start, so use the last one.

```powershell
wsl -d Ubuntu-24.04 -- docker logs portainer 2>&1 | Select-String 'setup_token' | Select-Object -Last 1
```

Treat the token like a password: it lets whoever has it create the admin account. Create the admin account immediately, because Portainer locks the first-run setup after a few minutes (vendor behavior, untested). `docker restart portainer` re-opens it and prints a new token. The distro must stay running while you do this: apply the idle setting from Integrate first, or keep a `wsl` session open.

## Integrate

No PATH or `HKCU` change is needed. Run the engine through `wsl`, or open a shell with `wsl -d Ubuntu-24.04`. A Windows-side `docker.exe` is not part of this guide.

By default WSL does not start a distro at Windows sign-in, and it stops an idle distro shortly after the last `wsl.exe` session closes (the Microsoft default `instanceIdleTimeout` is 15 seconds), which also stops every container in it. Two per-user changes fix both. Neither needs admin rights.

### 1. Stop WSL from shutting the distro down when idle

Create `%USERPROFILE%\.wslconfig` if it does not exist. If it exists, add the `instanceIdleTimeout` line under an existing `[general]` section instead of overwriting the file. The setting applies to every WSL2 distro of your account, and Microsoft documents `-1` as disabling auto shutdown ([wsl-config](https://learn.microsoft.com/en-us/windows/wsl/wsl-config#:~:text=Set%20to%20%2D1%20to%20disable%20auto%20shutdown)).

```powershell
$cfg = Join-Path $env:USERPROFILE ".wslconfig"
if (Test-Path $cfg) { Get-Content $cfg; Write-Host "Exists: merge the setting by hand." }
else { Set-Content -Path $cfg -Value "[general]`r`ninstanceIdleTimeout=-1`r`n" -Encoding ASCII }
wsl --shutdown
```

`wsl --shutdown` applies the file and stops every running distro, including Portainer.

### 2. Start the distro at sign-in

A shortcut in the per-user Startup folder runs `wsl.exe` once at every sign-in. `conhost.exe --headless` hides the console window. With step 1 in place, the distro then stays up, and `--restart unless-stopped` brings Portainer back with the engine.

```powershell
$startup = [Environment]::GetFolderPath('Startup')
$ws = New-Object -ComObject WScript.Shell
$lnk = $ws.CreateShortcut((Join-Path $startup "WSL Docker keepalive.lnk"))
$lnk.TargetPath = Join-Path $env:SystemRoot "System32\conhost.exe"
$lnk.Arguments = "--headless wsl.exe -d Ubuntu-24.04 --exec true"
$lnk.WindowStyle = 7
$lnk.Description = "Starts the WSL2 distro (Docker Engine and Portainer) at sign-in"
$lnk.Save()
```

To start it now without signing out, run the shortcut. Starting from a fully stopped WSL takes about 30 seconds.

```powershell
Start-Process (Join-Path ([Environment]::GetFolderPath('Startup')) "WSL Docker keepalive.lnk")
```

Starting at a real sign-in is untested; the shortcut was only launched by hand.

Without step 1, replace `--exec true` in the shortcut with `--exec sleep infinity`. The hidden session then holds the distro open for as long as you stay signed in.

### 3. Open Portainer with one click

[`open-portainer.ps1`](open-portainer.ps1) checks `https://localhost:9443/`. If it answers, Chrome opens at once. If not, the script starts the distro, starts the `portainer` container, waits for the page and then opens Chrome. It needs no admin rights and prints each step. On failure it shows a message box: missing distro, missing container, or no answer within 90 seconds.

Chrome is started by its full path (`Program Files`, `Program Files (x86)`, then `%LOCALAPPDATA%`), not through the `App Paths` registry key. That key can point to another Chromium browser on the account. If Chrome is not found, the default browser opens instead.

Run it from the folder that holds the script. `-CreateShortcut` copies the script to `%LOCALAPPDATA%\Programs\docker-wsl` and creates `Portainer.lnk` on the Desktop and in the Start Menu, both pointing at that copy so moving the repository folder does not break them. The shortcut runs `powershell.exe -NoProfile -File <copy>` in a minimized window, with no execution-policy override, which works under the `RemoteSigned` policy of the current user.

```powershell
powershell.exe -NoProfile -File .\open-portainer.ps1 -CreateShortcut
```

Run it while Portainer is up. The script reads the Portainer logo from the running page, wraps it in `portainer.ico` beside the installed script, and uses it as the shortcut icon. The logo is fetched from your own instance and is not stored in this repository. If Portainer does not answer, the shortcut gets Chrome's icon instead; start Portainer and run the command again to switch it. The Desktop may keep showing the old icon for a moment until Windows refreshes its icon cache.

Pin the Start Menu entry to the taskbar if you want a button there. Run the script without `-CreateShortcut` to test it. Add `-NoBrowser` to prepare the environment without opening Chrome.

### 4. Trust a local certificate for Portainer (optional)

Portainer serves its own self-signed certificate, so Chrome shows "Not secure". This step replaces it with one self-signed leaf certificate that your account trusts. It is a leaf (`CA:FALSE`), not a private CA: nothing holds a key that could sign certificates for other hosts, and the trust covers this one certificate only. Skip the step if you accept the warning.

Importing the certificate adds a trust anchor to your Windows account (`Cert:\CurrentUser\Root`). Chrome, Edge and every other application of the account that uses the Windows store will trust it. Anyone who obtains the private key can impersonate `localhost` to those applications, so the key stays inside the distro, readable by root only, and never goes into a repository. A managed device may have security software or policy that reacts to a new root certificate. Check with your security team first if the device is managed. If the import is blocked, stop and do not work around it.

Open a root shell in the distro (`wsl -d Ubuntu-24.04 -u root`) and create the key and certificate. The certificate is valid for 1825 days (five years); lower `-days` for a shorter window. Subject, subject alternative names and key type are set in the file below.

```bash
d=/etc/portainer-certs
mkdir -p $d && chmod 755 $d
cat > $d/openssl.cnf <<'EOF'
[req]
distinguished_name = dn
prompt = no
x509_extensions = ext

[dn]
CN = portainer

[ext]
basicConstraints = critical,CA:FALSE
keyUsage = critical,digitalSignature
extendedKeyUsage = serverAuth
subjectAltName = DNS:localhost,IP:127.0.0.1,IP:::1
subjectKeyIdentifier = hash
EOF
openssl genpkey -algorithm EC -pkeyopt ec_paramgen_curve:P-256 -out $d/portainer.key
chmod 600 $d/portainer.key
openssl req -new -x509 -sha256 -days 1825 -key $d/portainer.key -config $d/openssl.cnf -out $d/portainer.crt
chmod 644 $d/portainer.crt
exit
```

Recreate the container with the certificate. The `portainer_data` volume keeps the admin account. The old container stays as `portainer-old` for rollback.

```powershell
wsl -d Ubuntu-24.04 -- docker stop portainer
wsl -d Ubuntu-24.04 -- docker rename portainer portainer-old
wsl -d Ubuntu-24.04 -- docker run -d --name portainer --restart unless-stopped -p 127.0.0.1:9443:9443 -v /var/run/docker.sock:/var/run/docker.sock -v portainer_data:/data -v /etc/portainer-certs:/certs:ro portainer/portainer-ce:lts --sslcert /certs/portainer.crt --sslkey /certs/portainer.key
```

Copy the public certificate to Windows and import it. Windows may ask you to confirm the import; accept it only if the certificate subject is `portainer`.

```powershell
$dest = Join-Path $env:LOCALAPPDATA "Programs\docker-wsl"
$null = New-Item -ItemType Directory -Force $dest
$wslDest = "/mnt/" + $dest.Substring(0, 1).ToLower() + $dest.Substring(2).Replace("\", "/")
wsl -d Ubuntu-24.04 -u root -- cp /etc/portainer-certs/portainer.crt "$wslDest/portainer.crt"
certutil.exe -user -addstore Root (Join-Path $dest "portainer.crt")
curl.exe -s -o NUL -w "%{http_code}" https://localhost:9443/
```

`curl.exe` without `-k` must print `200`. Log in to Portainer to confirm the admin account is intact, then remove the old container with `wsl -d Ubuntu-24.04 -- docker rm portainer-old`. To roll back instead, run `docker rm -f portainer`, then `docker rename portainer-old portainer` and `docker start portainer`.

The certificate expires on the date printed by `wsl -d Ubuntu-24.04 -u root -- openssl x509 -in /etc/portainer-certs/portainer.crt -noout -enddate`. Before then, remove the old certificate from the store (see Uninstall), delete the three files in `/etc/portainer-certs`, and repeat this step.

### 5. Call `docker` from PowerShell (optional)

`docker` is not a Windows program here, so PowerShell and `cmd` report that the command is not found. A one-line wrapper named `docker.cmd` in a folder on your user PATH forwards every call to the engine in the distro. It opens no network port and writes nothing outside your profile. Set the `DOCKER_WSL_DISTRO` environment variable to use a distro other than `Ubuntu-24.04`.

```powershell
$bin = Join-Path $env:LOCALAPPDATA "Programs\docker-wsl\bin"
$null = New-Item -ItemType Directory -Force $bin
$wrapper = @(
    '@echo off',
    'setlocal',
    'if not defined DOCKER_WSL_DISTRO set "DOCKER_WSL_DISTRO=Ubuntu-24.04"',
    'wsl.exe -d %DOCKER_WSL_DISTRO% -- docker %*',
    'exit /b %ERRORLEVEL%'
)
Set-Content -Path (Join-Path $bin "docker.cmd") -Value $wrapper -Encoding ASCII
$userPath = [Environment]::GetEnvironmentVariable("Path", "User")
if (($userPath -split ";") -notcontains $bin) {
    [Environment]::SetEnvironmentVariable("Path", ($userPath.TrimEnd(";") + ";" + $bin), "User")
}
```

Open a new terminal so it picks up the PATH change, then run `docker --version`. The wrapper returns the exit code of `docker`, and `docker compose` works through it.

The engine runs on Linux, so a path in a bind mount must be a Linux path. `-v "C:\work:/w"` fails with `invalid mode`. Write `/mnt/c/work`, or let the distro resolve the current directory:

```powershell
docker run --rm -v "$(wsl -d Ubuntu-24.04 -- pwd):/w" alpine:3.20 ls /w
```

## Verify

| Check | Command | Expected |
|---|---|---|
| Engine runs as a service | `wsl -d Ubuntu-24.04 -- systemctl is-active docker` | `active` |
| Non-root user can use Docker | `wsl -d Ubuntu-24.04 -- docker run --rm hello-world` | `Hello from Docker!` |
| Compose plugin | `wsl -d Ubuntu-24.04 -- docker compose version` | `Docker Compose version v5.6.0` |
| Portainer from Windows | `curl.exe -k -s -o NUL -w "%{http_code}" https://127.0.0.1:9443/` | `200` |
| Distro disk in user profile | `Get-ChildItem "$env:LOCALAPPDATA\wsl" -Recurse -Filter ext4.vhdx` | One `ext4.vhdx` file |
| Launcher prepares everything | `powershell.exe -NoProfile -File .\open-portainer.ps1 -NoBrowser` | `Portainer is ready.` and exit code `0` |
| Wrapper in a new terminal (step 5) | `docker --version` | `Docker version 29.9.0` or newer |
| Trusted certificate (step 4) | `curl.exe -s -o NUL -w "%{http_code}" https://localhost:9443/` | `200` without `-k` |
| Shortcut icon is Portainer's | `(New-Object -ComObject WScript.Shell).CreateShortcut((Join-Path ([Environment]::GetFolderPath('Desktop')) "Portainer.lnk")).IconLocation` | A path ending in `docker-wsl\portainer.ico,0` |
| Shortcuts exist | `Test-Path (Join-Path ([Environment]::GetFolderPath('Desktop')) "Portainer.lnk")` | `True` |

## Use

### Start and stop the environment

| Goal | Command | Effect |
|---|---|---|
| Start the engine and Portainer | `wsl -d Ubuntu-24.04 --exec true` | Boots the distro. systemd starts `docker`, and `--restart unless-stopped` brings Portainer back up. With `instanceIdleTimeout=-1` from Integrate it then stays up; without it the distro stops shortly after the command ends and Portainer stops with it. About 6 seconds from a stopped distro, about 30 seconds after `wsl --shutdown`. |
| Start it automatically | The Startup shortcut from Integrate | Runs the same command at every sign-in. |
| Open Portainer | The `Portainer` shortcut from Integrate | Opens Chrome at once when everything runs. Otherwise starts the distro and container first: about 6 seconds when only the container is stopped, about 14 seconds when the distro is stopped too. |
| Check state | `wsl -l -v` | `Running` or `Stopped` for `Ubuntu-24.04`. |
| Stop everything | `wsl --terminate Ubuntu-24.04` | Stops the distro and every container in it. Portainer returns at the next start. |
| Stop only Portainer | `wsl -d Ubuntu-24.04 -- docker stop portainer` | It stays stopped across restarts until you run `docker start portainer`. |

Containers you start with `docker run -d` do not come back by themselves after the distro stops. Add `--restart unless-stopped` to the ones that should.

### Common commands

```powershell
wsl -d Ubuntu-24.04 -- docker ps
wsl -d Ubuntu-24.04 -- docker run --rm -it alpine:3.20 sh
wsl -d Ubuntu-24.04 -- docker compose up -d
wsl -d Ubuntu-24.04 -- docker stop portainer
wsl -d Ubuntu-24.04 -- docker start portainer
```

Mounting the Docker socket into Portainer gives it root-equivalent control of the engine. Keep the port bound to `127.0.0.1`, set a strong admin password, and do not publish `9443` on other interfaces. On a shared Windows host, other local users can reach a port WSL forwards to localhost.

## Uninstall

`wsl --unregister` deletes the whole distro, including every image, container and volume in it. Export what you need first.

```powershell
$startup = Join-Path ([Environment]::GetFolderPath('Startup')) "WSL Docker keepalive.lnk"
if (Test-Path $startup) { Remove-Item $startup -Force }
foreach ($folder in "Desktop", "Programs") {
    $link = Join-Path ([Environment]::GetFolderPath($folder)) "Portainer.lnk"
    if (Test-Path $link) { Remove-Item $link -Force }
}
$launcher = Join-Path $env:LOCALAPPDATA "Programs\docker-wsl"
if (Test-Path $launcher) { Remove-Item $launcher -Recurse -Force }
certutil.exe -user -delstore Root portainer
$bin = Join-Path $env:LOCALAPPDATA "Programs\docker-wsl\bin"
$kept = ([Environment]::GetEnvironmentVariable("Path", "User") -split ";") | Where-Object { $_ -and $_ -ne $bin }
[Environment]::SetEnvironmentVariable("Path", ($kept -join ";"), "User")
wsl --export Ubuntu-24.04 "$env:USERPROFILE\ubuntu-24.04-backup.tar"
wsl --unregister Ubuntu-24.04
$startMenu = Join-Path $env:APPDATA "Microsoft\Windows\Start Menu\Programs\Ubuntu-24.04"
if (Test-Path $startMenu) { Remove-Item $startMenu -Recurse -Force }
$cfg = Join-Path $env:USERPROFILE ".wslconfig"
if (Test-Path $cfg) {
    $rest = @(Get-Content $cfg | Where-Object { $_ -notmatch '^\s*instanceIdleTimeout\s*=' })
    $content = @($rest | Where-Object { $_.Trim() -and $_.Trim() -ne "[general]" })
    if ($content.Count -eq 0) { Remove-Item $cfg -Force } else { Set-Content -Path $cfg -Value $rest -Encoding ASCII }
}
wsl --shutdown
```

The script removes, in order: the Startup shortcut, both `Portainer.lnk` files, the `docker-wsl` folder (the launcher script, the `bin` folder with `docker.cmd` and the copy of `portainer.crt`), the `portainer` certificate from your root store, the `docker-wsl\bin` entry of your user PATH, the distro with its disk and `Lxss` entry, a leftover Start Menu folder named after the distro, and the `instanceIdleTimeout` line of `%USERPROFILE%\.wslconfig` (the file itself goes if nothing else is in it). The setting applies to every WSL2 distro, so keep any other lines you added. The `ubuntu-24.04-backup.tar` export stays in your profile as your backup; delete it by hand when you no longer need it. Untested against the real distro: the Start Menu folder line and the `.wslconfig` lines (the `.wslconfig` logic was run on temporary copies only).

`wsl --export` stops the distro and every container in it; start it again with `wsl -d Ubuntu-24.04 --exec true` if you export without unregistering. The `certutil` line applies only if you did step 4 of Integrate: it removes the certificate in an interactive terminal. In a session that cannot show a dialog, `certutil` waits for an answer and `Remove-Item` on the store fails with "UI is not allowed". If you keep the distro, delete `/etc/portainer-certs` with `wsl -d Ubuntu-24.04 -u root -- rm -r /etc/portainer-certs` and recreate the container without the `--sslcert` and `--sslkey` flags. `wsl --unregister` removes the distro disk under `%LOCALAPPDATA%\wsl` and the `Lxss` entry in `HKCU`.

## Troubleshooting

| Symptom | Cause | Fix |
|---|---|---|
| `docker pull` fails with `lookup ... on 10.255.255.254:53: server misbehaving` | The Go resolver in `dockerd` does not work with the WSL DNS proxy | Apply step 3, then retry. |
| `dial tcp [2600:...]:443: connect: network is unreachable` | The distro has no IPv6 route and the pull tried an IPv6 address first | Apply step 3. If it persists, disabling IPv6 in the distro with `sysctl net.ipv6.conf.all.disable_ipv6=1` is a possible fallback (untested as a fix). |
| `permission denied` on `/var/run/docker.sock` | The Linux user is not in the `docker` group, or the distro was not restarted | Run step 4 again, then `wsl --terminate Ubuntu-24.04`. |
| `Cannot connect to the Docker daemon` | systemd is off or the service stopped | `wsl -d Ubuntu-24.04 -u root -- systemctl status docker`. |
| `https://localhost:9443` does not load | The container is stopped, or the first-run setup timed out | `wsl -d Ubuntu-24.04 -- docker restart portainer`. |
| `https://localhost:9443` stops working shortly after the last terminal closes | WSL stops an idle distro and every container with it | Apply step 1 of Integrate, then `wsl --shutdown` so the setting takes effect. Or keep one `wsl` session open. |
| `https://localhost:9443` is down right after sign-in | The distro has not started yet, or the Startup shortcut is missing | Wait about 30 seconds, check `wsl -l -v`, and apply step 2 of Integrate. |
| A container is gone after a reboot | Containers without a restart policy do not restart when the distro boots | Recreate it with `--restart unless-stopped`. |
| `wsl --install` asks for elevation or reports WSL is missing | WSL2 is not enabled on this machine | Stop and ask IT to enable WSL2. |
| The `Portainer` shortcut opens the wrong browser | The launcher found no Google Chrome in the three standard folders and used the default browser | Install Chrome, or pass the browser with `-BrowserPath` when creating the shortcut. |
| A message box says the container does not exist | The `portainer` container was removed or never created | Create it again with the Portainer step under Install. The `portainer_data` volume keeps the admin account. |
| A message box says Portainer did not answer within 90 seconds | The engine or the container is unhealthy | Run the launcher from a terminal to see which step stalls, then `wsl -d Ubuntu-24.04 -- docker logs portainer`. |
| Chrome shows a "Not secure" page | Portainer serves its own self-signed certificate, or Chrome remembers an earlier choice to proceed past the warning | Without a trusted certificate: choose Advanced, then proceed to localhost, or apply step 4 of Integrate. With one: open an incognito window to confirm the certificate is trusted, then close every Chrome window and reopen the page, which clears the remembered choice. |
| Chrome shows a certificate error again after a long time | The certificate expired | Renew it as described at the end of step 4 of Integrate. |
