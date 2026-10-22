package setup

import (
	"crypto/rand"
	"encoding/hex"
	"strconv"
	"strings"

	"github.com/the-khiem7/docker-zero/internal/settings"
	"github.com/the-khiem7/docker-zero/internal/target"
)

// markerPresent is printed by a script that found its work already done.
const markerPresent = "present"

// enableSystemdScript turns on systemd for the next boot without dropping
// settings that wsl.conf already holds.
const enableSystemdScript = `set -e
if [ -f /etc/wsl.conf ] && grep -q '^\[boot\]' /etc/wsl.conf; then
  sed -i '/^\[boot\]/a systemd=true' /etc/wsl.conf
else
  printf '[boot]\nsystemd=true\n' >> /etc/wsl.conf
fi
`

// installDockerScript installs Docker Engine from Docker's apt repository,
// the documented method for Ubuntu, unless the service already runs.
const installDockerScript = `set -e
if command -v docker >/dev/null 2>&1 && systemctl is-active --quiet docker; then echo ` + markerPresent + `; exit 0; fi
export DEBIAN_FRONTEND=noninteractive
apt-get update -q
apt-get install -y -q ca-certificates curl
install -m 0755 -d /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/ubuntu/gpg -o /etc/apt/keyrings/docker.asc
chmod a+r /etc/apt/keyrings/docker.asc
codename=$(grep '^VERSION_CODENAME=' /etc/os-release | cut -d= -f2)
cat > /etc/apt/sources.list.d/docker.sources <<EOF
Types: deb
URIs: https://download.docker.com/linux/ubuntu
Suites: ${codename}
Components: stable
Architectures: $(dpkg --print-architecture)
Signed-By: /etc/apt/keyrings/docker.asc
EOF
apt-get update -q
apt-get install -y -q docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin
systemctl enable --now docker
docker --version
`

// dnsScript makes dockerd use the system resolver. Its built-in Go resolver
// fails against the WSL DNS tunnel, so every pull would report a DNS error
// while curl in the same distro works.
const dnsScript = `set -e
f=/etc/systemd/system/docker.service.d/override.conf
want='[Service]
Environment="GODEBUG=netdns=cgo"'
if [ "$(cat "$f" 2>/dev/null)" = "$want" ]; then exit 0; fi
install -d /etc/systemd/system/docker.service.d
printf '%s\n' "$want" > "$f"
systemctl daemon-reload
systemctl restart docker
n=0; until docker info >/dev/null 2>&1 || [ $n -ge 30 ]; do sleep 1; n=$((n+1)); done
docker info >/dev/null
echo "Docker uses the system DNS resolver"
`

// userScript creates the default Linux user in the docker group and points
// wsl.conf at it. The user name was validated against a safe alphabet.
func userScript(user string) string {
	return `set -e
if id -u ` + user + ` >/dev/null 2>&1 && id -nG ` + user + ` | grep -qw docker && grep -q '^default=` + user + `$' /etc/wsl.conf 2>/dev/null; then echo ` + markerPresent + `; exit 0; fi
id -u ` + user + ` >/dev/null 2>&1 || adduser --disabled-password --gecos "" ` + user + `
usermod -aG docker ` + user + `
cat > /etc/wsl.conf <<'EOF'
[boot]
systemd=true

[user]
default=` + user + `
EOF
`
}

// portainerRunScript creates the Portainer container. With a token it passes
// the one-time setup token through a root-only env file that is removed at
// once, so the value is not on a command line or in the shell history.
// replace removes an existing container first; the volume is kept.
func portainerRunScript(port int, token string, replace bool) string {
	var b strings.Builder
	b.WriteString("set -e\numask 077\n")
	b.WriteString("docker volume create portainer_data >/dev/null\n")
	if replace {
		b.WriteString("docker rm -f " + target.PortainerContainer + " >/dev/null 2>&1 || true\n")
	}
	envFlag := ""
	if token != "" {
		b.WriteString("trap 'rm -f /run/dz-portainer.env' EXIT\n")
		b.WriteString("printf 'PORTAINER_SETUP_TOKEN=%s\\n' '" + token + "' > /run/dz-portainer.env\n")
		envFlag = " --env-file /run/dz-portainer.env"
	}
	b.WriteString("docker run -d --name " + target.PortainerContainer + " --restart unless-stopped" +
		" -p 127.0.0.1:" + strconv.Itoa(port) + ":9443" +
		" -v /var/run/docker.sock:/var/run/docker.sock -v portainer_data:/data" +
		" -v " + settings.CertDir + ":/certs:ro" + envFlag + " " + portainerImage +
		" --sslcert /certs/portainer.crt --sslkey /certs/portainer.key >/dev/null\n")
	return b.String()
}

// portainerStateScript prints the container state, or "missing".
const portainerStateScript = `docker inspect -f '{{.State.Status}}' ` + target.PortainerContainer + ` 2>/dev/null || echo missing
`

// dockerCmd is the wrapper that forwards docker to the engine in the distro.
// The distro is overridable per shell through DOCKER_WSL_DISTRO.
func dockerCmd(distro string) string {
	lines := []string{
		"@echo off",
		"setlocal",
		`if not defined DOCKER_WSL_DISTRO set "DOCKER_WSL_DISTRO=` + distro + `"`,
		"wsl.exe -d %DOCKER_WSL_DISTRO% -- docker %*",
		"exit /b %ERRORLEVEL%",
		"",
	}
	return strings.Join(lines, "\r\n")
}

// newToken returns a 256-bit one-time setup token as hex.
func newToken() (string, error) {
	var b [32]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", err
	}
	return hex.EncodeToString(b[:]), nil
}
