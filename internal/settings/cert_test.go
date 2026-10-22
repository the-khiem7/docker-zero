package settings

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"math/big"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/the-khiem7/docker-zero/internal/wsl"
)

type fakeScripter struct {
	out    string
	err    error
	script string
}

func (f *fakeScripter) Script(_ context.Context, _, _, script string) (string, error) {
	f.script = script
	return f.out, f.err
}

func selfSigned(t *testing.T) []byte {
	t.Helper()
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	tmpl := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject:      pkix.Name{CommonName: "portainer"},
		NotBefore:    time.Unix(0, 0),
		NotAfter:     time.Date(2031, 10, 14, 0, 0, 0, 0, time.UTC),
	}
	der, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, &key.PublicKey, key)
	if err != nil {
		t.Fatal(err)
	}
	return pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der})
}

func TestReadAndParseCert(t *testing.T) {
	pemBytes := selfSigned(t)
	got, err := ReadCertPEM(t.Context(), &fakeScripter{out: string(pemBytes)}, "d")
	if err != nil {
		t.Fatal(err)
	}
	thumb, notAfter, err := ParseCert(got)
	if err != nil {
		t.Fatal(err)
	}
	if len(thumb) != 40 || strings.ToUpper(thumb) != thumb {
		t.Fatalf("thumbprint %q", thumb)
	}
	if notAfter.Year() != 2031 {
		t.Fatalf("notAfter %v", notAfter)
	}
	if _, err := ReadCertPEM(t.Context(), &fakeScripter{out: ""}, "d"); !errors.Is(err, ErrNoCert) {
		t.Fatalf("empty: %v", err)
	}
}

func TestRegenerateScriptKeepsKeyInDistro(t *testing.T) {
	f := &fakeScripter{}
	if err := Regenerate(t.Context(), f, "d"); err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"ec_paramgen_curve:P-256", "-days 1825", "CA:FALSE", "IP:::1", "umask 077", "docker restart portainer"} {
		if !strings.Contains(f.script, want) {
			t.Errorf("script lacks %q", want)
		}
	}
	f.err = &wsl.CommandError{ExitCode: 3}
	if err := Regenerate(t.Context(), f, "d"); !errors.Is(err, ErrNotConfigured) {
		t.Fatalf("got %v", err)
	}
}

func TestChromeCandidatesAndFind(t *testing.T) {
	env := map[string]string{"ProgramFiles": `C:\PF`, "LOCALAPPDATA": `C:\L`}
	getenv := func(k string) string { return env[k] }
	if got := ChromeCandidates(getenv); len(got) != 2 {
		t.Fatalf("candidates %q", got)
	}
	dir := t.TempDir()
	exe := filepath.Join(dir, "Google", "Chrome", "Application", "chrome.exe")
	if err := os.MkdirAll(filepath.Dir(exe), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(exe, nil, 0o600); err != nil {
		t.Fatal(err)
	}
	env["LOCALAPPDATA"] = dir
	got, err := FindChrome(getenv, os.Stat)
	if err != nil || got != exe {
		t.Fatalf("got %q, %v", got, err)
	}
	if _, err := FindChrome(func(string) string { return "" }, os.Stat); !errors.Is(err, ErrChromeNotFound) {
		t.Fatalf("got %v", err)
	}
}

func TestRunCommand(t *testing.T) {
	if got := runCommand(`C:\p\docker-zero.exe`, true); got != `"C:\p\docker-zero.exe" --start-distro` {
		t.Fatalf("got %q", got)
	}
	if got := runCommand(`C:\p\docker-zero.exe`, false); got != `"C:\p\docker-zero.exe"` {
		t.Fatalf("got %q", got)
	}
}
