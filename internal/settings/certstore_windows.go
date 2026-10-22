//go:build windows

package settings

import (
	"encoding/hex"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

// IsTrusted reports whether the certificate with thumbprint is in the
// current user's Root store. Reading the store needs no dialog.
func IsTrusted(thumbprint string) (bool, error) {
	hash, err := hex.DecodeString(thumbprint)
	if err != nil {
		return false, err
	}
	name, err := windows.UTF16PtrFromString("Root")
	if err != nil {
		return false, err
	}
	store, err := windows.CertOpenStore(windows.CERT_STORE_PROV_SYSTEM, 0, 0,
		windows.CERT_SYSTEM_STORE_CURRENT_USER|windows.CERT_STORE_READONLY_FLAG,
		uintptr(unsafe.Pointer(name))) //nolint:gosec // the Win32 API takes a pointer
	if err != nil {
		return false, err
	}
	defer func() { _ = windows.CertCloseStore(store, 0) }()
	blob := windows.CryptHashBlob{Size: uint32(len(hash)), Data: &hash[0]} //nolint:gosec // a SHA-1 hash is 20 bytes
	ctx, err := windows.CertFindCertificateInStore(store, windows.X509_ASN_ENCODING|windows.PKCS_7_ASN_ENCODING,
		0, windows.CERT_FIND_SHA1_HASH, unsafe.Pointer(&blob), nil) //nolint:gosec // the Win32 API takes a pointer
	if errors.Is(err, windows.Errno(windows.CRYPT_E_NOT_FOUND)) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	_ = windows.CertFreeCertificateContext(ctx)
	return true, nil
}

// Trust adds the certificate to the current user's Root store. Windows asks
// the user to confirm with its own security dialog, which is expected; it is
// the only way to change that store as a standard user. The public PEM is
// kept at crtPath for reference; the key never leaves the distro.
func Trust(pemBytes []byte, crtPath string) error {
	if err := os.MkdirAll(filepath.Dir(crtPath), 0o700); err != nil {
		return err
	}
	if err := os.WriteFile(crtPath, pemBytes, 0o600); err != nil {
		return err
	}
	return certutil("-user", "-addstore", "Root", crtPath)
}

// Untrust removes the certificate with thumbprint from the current user's
// Root store, again behind Windows' confirmation dialog.
func Untrust(thumbprint string) error {
	return certutil("-user", "-delstore", "Root", thumbprint)
}

func certutil(args ...string) error {
	cmd := exec.Command("certutil.exe", args...) //nolint:gosec // fixed executable; args are a store name and a path or thumbprint this package built
	// Hide certutil's console; its confirmation dialog still shows.
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true, CreationFlags: 0x08000000}
	out, err := cmd.CombinedOutput()
	if err != nil {
		return errors.New("certutil " + args[1] + ": " + string(out))
	}
	return nil
}
