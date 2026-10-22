package setup

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"github.com/the-khiem7/docker-zero/internal/settings"
	"github.com/the-khiem7/docker-zero/internal/target"
)

// resetImage is Portainer's own tool for replacing an admin password.
const resetImage = "portainer/helper-reset-password"

// bcryptCost matches the hashes Portainer writes itself.
const bcryptCost = 10

// ErrNoPortainer is returned by a maintenance action when the container does
// not exist yet.
var ErrNoPortainer = errors.New("there is no Portainer container yet; run setup first")

// ErrPasswordTooLong is returned for a password bcrypt cannot hash.
var ErrPasswordTooLong = errors.New("the password is longer than 72 bytes, the most Portainer can store")

// hashPattern is the shape of a bcrypt hash. The hash is placed inside a
// shell script, so nothing outside this alphabet may reach it.
var hashPattern = regexp.MustCompile(`^\$2[aby]\$[0-9]{2}\$[./A-Za-z0-9]{53}$`)

// currentScript prints "current" when the container already runs the image
// the tag points at now, and "newer" otherwise.
const currentScript = `if [ "$(docker inspect -f '{{.Image}}' ` + target.PortainerContainer + ` 2>/dev/null)" = "$(docker image inspect -f '{{.Id}}' ` + portainerImage + ` 2>/dev/null)" ]; then echo current; else echo newer; fi
`

// resetScript stops Portainer, replaces the admin password with the hash on
// the data volume, and starts Portainer again whatever the helper did. The
// official procedure stops the container first because the helper writes to
// Portainer's database.
func resetScript(hash string) string {
	return `set -e
trap 'docker start ` + target.PortainerContainer + ` >/dev/null 2>&1 || true' EXIT
docker stop ` + target.PortainerContainer + ` >/dev/null
docker run --rm -v portainer_data:/data ` + resetImage + ` --password-hash '` + hash + `'
`
}

// prepare validates the configuration for a maintenance action.
func (w *Wizard) prepare() error {
	w.Config = w.Config.withDefaults()
	if err := w.Config.validate(); err != nil {
		return err
	}
	w.api = nil
	w.cur = ""
	return nil
}

// UpdatePortainer pulls the pinned image and, when it changed, recreates the
// container on the same data volume and certificate. It reports whether the
// container was replaced. A failed pull leaves the running container alone.
func (w *Wizard) UpdatePortainer(ctx context.Context) (bool, error) {
	if err := w.prepare(); err != nil {
		return false, err
	}
	state, err := w.quiet(ctx, portainerStateScript)
	if err != nil {
		return false, err
	}
	if state == "missing" {
		return false, ErrNoPortainer
	}
	w.logf("Pulling %s", portainerImage)
	if _, err := w.root(ctx, "docker pull "+portainerImage+"\n"); err != nil {
		return false, err
	}
	got, err := w.quiet(ctx, currentScript)
	if err != nil {
		return false, err
	}
	if got == "current" {
		w.logf("Portainer already runs the latest image")
		return false, nil
	}
	if err := settings.Generate(ctx, w.WSL, w.Config.Distro); err != nil {
		return false, err
	}
	w.logf("Recreating the container on the same data volume")
	if err := w.recreate(ctx, ""); err != nil {
		return false, err
	}
	return true, nil
}

// ResetAdmin sets a new password for the Portainer admin account. The helper
// changes the existing administrator whatever its name, so the username only
// serves the sign-in that proves the new password works. The password stays in
// this process: only its bcrypt hash reaches the distro.
func (w *Wizard) ResetAdmin(ctx context.Context, cr Credentials) error {
	if err := w.prepare(); err != nil {
		return err
	}
	if err := cr.Validate(); err != nil {
		return err
	}
	raw, err := bcrypt.GenerateFromPassword([]byte(cr.Password), bcryptCost)
	if errors.Is(err, bcrypt.ErrPasswordTooLong) {
		return ErrPasswordTooLong
	}
	if err != nil {
		return err
	}
	hash := string(raw)
	if !hashPattern.MatchString(hash) {
		return errors.New("unexpected password hash format")
	}
	state, err := w.quiet(ctx, portainerStateScript)
	if err != nil {
		return err
	}
	if state == "missing" {
		return ErrNoPortainer
	}
	w.logf("Replacing the admin password")
	if _, err := w.quiet(ctx, resetScript(hash)); err != nil {
		return redact(err, hash)
	}
	if err := w.waitReady(ctx); err != nil {
		return err
	}
	a, err := w.client(ctx)
	if err != nil {
		return err
	}
	if err := a.login(ctx, cr); err != nil {
		return fmt.Errorf("the password was replaced, but signing in as %q failed; check the username: %w", cr.Username, err)
	}
	w.logf("Signed in with the new password")
	return nil
}

// redact keeps a hash out of an error message.
func redact(err error, secret string) error {
	if err == nil || !strings.Contains(err.Error(), secret) {
		return err
	}
	return errors.New(strings.ReplaceAll(err.Error(), secret, "<hash>"))
}
