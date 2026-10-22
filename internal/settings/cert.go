package settings

import (
	"context"
	"crypto/sha1" //nolint:gosec // Windows identifies store certificates by SHA-1 thumbprint
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/the-khiem7/docker-zero/internal/target"
	"github.com/the-khiem7/docker-zero/internal/wsl"
)

// CertDir holds the Portainer leaf certificate and its key inside the
// distro. The key never leaves the distro: a copy on the Windows side would
// let anyone who reads it impersonate localhost to this account's browsers.
const CertDir = "/etc/portainer-certs"

// ErrNoCert is returned when the distro holds no certificate yet.
var ErrNoCert = errors.New("no Portainer certificate in the distro")

// ErrNotConfigured is returned when the Portainer container does not serve
// the certificate files, so regenerating them would change nothing.
var ErrNotConfigured = errors.New("the Portainer container does not use " + CertDir + "; run setup first")

// Cert describes the leaf certificate.
type Cert struct {
	Present    bool      `json:"present"`
	Thumbprint string    `json:"thumbprint,omitempty"`
	NotAfter   time.Time `json:"notAfter,omitzero"`
	Trusted    bool      `json:"trusted"`
}

// Scripter runs a bash script in the distro as a user.
type Scripter interface {
	Script(ctx context.Context, distro, user, script string) (string, error)
}

// ReadCertPEM returns the certificate PEM from distro.
func ReadCertPEM(ctx context.Context, w Scripter, distro string) ([]byte, error) {
	out, err := w.Script(ctx, distro, "root", "cat "+CertDir+"/portainer.crt 2>/dev/null || true\n")
	if err != nil {
		return nil, err
	}
	if !strings.Contains(out, "BEGIN CERTIFICATE") {
		return nil, ErrNoCert
	}
	return []byte(out), nil
}

// ParseCert returns the thumbprint (upper-case hex SHA-1 of the DER, as
// Windows shows it) and expiry of a PEM certificate.
func ParseCert(pemBytes []byte) (thumbprint string, notAfter time.Time, err error) {
	block, _ := pem.Decode(pemBytes)
	if block == nil || block.Type != "CERTIFICATE" {
		return "", time.Time{}, ErrNoCert
	}
	c, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return "", time.Time{}, err
	}
	sum := sha1.Sum(block.Bytes) //nolint:gosec // thumbprint, not a security check
	return strings.ToUpper(hex.EncodeToString(sum[:])), c.NotAfter, nil
}

// openSSLConfig is the leaf profile: CN=portainer, SANs localhost,
// 127.0.0.1 and ::1, not a CA, server authentication only. A private CA was
// rejected because its key could sign a certificate for any host.
const openSSLConfig = `[req]
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
`

// regenerateScript writes a new ECDSA P-256 key and a five-year leaf next
// to the old pair, swaps them in, and restarts Portainer when its container
// serves these files. Exit code 3 means the container is not configured.
var regenerateScript = `set -e
d=` + CertDir + `
mkdir -p "$d"
cat > "$d/openssl.cnf" <<'EOF'
` + openSSLConfig + `EOF
umask 077
openssl genpkey -algorithm EC -pkeyopt ec_paramgen_curve:P-256 -out "$d/portainer.key.new"
openssl req -new -x509 -sha256 -days 1825 -key "$d/portainer.key.new" -config "$d/openssl.cnf" -out "$d/portainer.crt.new"
chmod 600 "$d/portainer.key.new"
chmod 644 "$d/portainer.crt.new"
mv -f "$d/portainer.key.new" "$d/portainer.key"
mv -f "$d/portainer.crt.new" "$d/portainer.crt"
if docker inspect -f '{{json .Args}}' ` + target.PortainerContainer + ` 2>/dev/null | grep -q -- --sslcert; then
  docker restart ` + target.PortainerContainer + ` >/dev/null
else
  exit 3
fi
`

// Regenerate replaces the certificate and key in distro and restarts
// Portainer. The caller removes the old certificate from the Windows store
// first and offers to trust the new one afterwards.
func Regenerate(ctx context.Context, w Scripter, distro string) error {
	_, err := w.Script(ctx, distro, "root", regenerateScript)
	if cmdErr, ok := errors.AsType[*wsl.CommandError](err); ok && cmdErr.ExitCode == 3 {
		return ErrNotConfigured
	}
	if err != nil {
		return fmt.Errorf("regenerate certificate: %w", err)
	}
	return nil
}
