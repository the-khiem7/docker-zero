//go:build integration

package settings

import (
	"errors"
	"os"
	"testing"

	"github.com/the-khiem7/docker-zero/internal/target"
	"github.com/the-khiem7/docker-zero/internal/wsl"
)

// TestCertReadOnlyReal reads the real distro's certificate and looks it up
// in the current user's Root store. It changes nothing.
func TestCertReadOnlyReal(t *testing.T) {
	pemBytes, err := ReadCertPEM(t.Context(), wsl.New(), target.Distro)
	if errors.Is(err, ErrNoCert) {
		t.Skip("no certificate in the real distro")
	}
	if err != nil {
		t.Fatal(err)
	}
	thumb, notAfter, err := ParseCert(pemBytes)
	if err != nil {
		t.Fatal(err)
	}
	trusted, err := IsTrusted(thumb)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("thumbprint %s, expires %s, trusted %v", thumb, notAfter.Format("2006-01-02"), trusted)
}

// TestRegenerateThrowaway runs the regeneration script in the throwaway
// distro named by DZ_TEST_DISTRO. It has no Portainer container, so the
// script must create the files and then report ErrNotConfigured.
func TestRegenerateThrowaway(t *testing.T) {
	distro := os.Getenv("DZ_TEST_DISTRO")
	if distro == "" {
		t.Skip("DZ_TEST_DISTRO not set")
	}
	c := wsl.New()
	if err := Regenerate(t.Context(), c, distro); !errors.Is(err, ErrNotConfigured) {
		t.Fatalf("got %v, want ErrNotConfigured", err)
	}
	pemBytes, err := ReadCertPEM(t.Context(), c, distro)
	if err != nil {
		t.Fatal(err)
	}
	if _, _, err := ParseCert(pemBytes); err != nil {
		t.Fatal(err)
	}
	perms, err := c.Script(t.Context(), distro, "root", "stat -c '%a %n' "+CertDir+"/portainer.key "+CertDir+"/portainer.crt\n")
	if err != nil {
		t.Fatal(err)
	}
	want := "600 " + CertDir + "/portainer.key\n644 " + CertDir + "/portainer.crt\n"
	if perms != want {
		t.Fatalf("perms %q, want %q", perms, want)
	}
}
