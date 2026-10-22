package setup

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/the-khiem7/docker-zero/internal/settings"
)

// api talks to the Portainer API over HTTPS on the loopback address. The
// connection is verified against the leaf certificate the wizard created, so
// the admin password only goes to the server that holds its key.
type api struct {
	base string
	http *http.Client
}

func newAPI(base string, certPEM []byte) (*api, error) {
	pool := x509.NewCertPool()
	if !pool.AppendCertsFromPEM(certPEM) {
		return nil, settings.ErrNoCert
	}
	return &api{
		base: base,
		http: &http.Client{
			Timeout: 15 * time.Second,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{RootCAs: pool, MinVersion: tls.VersionTLS12},
				// Loopback only: a corporate proxy must not see the password.
				Proxy: nil,
			},
		},
	}, nil
}

// do sends a request and returns the status code and the server's message,
// which never repeats the request body.
func (a *api) do(ctx context.Context, method, path string, header map[string]string, body any) (int, string, error) {
	var rd io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return 0, "", err
		}
		rd = bytes.NewReader(b)
	}
	req, err := http.NewRequestWithContext(ctx, method, a.base+path, rd)
	if err != nil {
		return 0, "", err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	for k, v := range header {
		req.Header.Set(k, v)
	}
	resp, err := a.http.Do(req)
	if err != nil {
		return 0, "", err
	}
	defer func() { _ = resp.Body.Close() }()
	raw, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<16))
	var msg struct {
		Message string `json:"message"`
	}
	_ = json.Unmarshal(raw, &msg)
	return resp.StatusCode, msg.Message, nil
}

func (a *api) ready(ctx context.Context) bool {
	code, _, err := a.do(ctx, http.MethodGet, "/system/status", nil, nil)
	return err == nil && code == http.StatusOK
}

// adminExists asks whether the first-run admin account was created.
func (a *api) adminExists(ctx context.Context) (bool, error) {
	code, msg, err := a.do(ctx, http.MethodGet, "/users/admin/check", nil, nil)
	switch {
	case err != nil:
		return false, err
	case code == http.StatusNoContent:
		return true, nil
	case code == http.StatusNotFound:
		return false, nil
	}
	return false, fmt.Errorf("admin check answered HTTP %d %s", code, msg)
}

func (a *api) initAdmin(ctx context.Context, token string, cr Credentials) error {
	code, msg, err := a.do(ctx, http.MethodPost, "/users/admin/init",
		map[string]string{"X-Setup-Token": token},
		map[string]string{"Username": cr.Username, "Password": cr.Password})
	if err != nil {
		return err
	}
	if code != http.StatusOK && code != http.StatusCreated {
		return fmt.Errorf("admin account refused by Portainer: HTTP %d %s", code, msg)
	}
	return nil
}

func (a *api) login(ctx context.Context, cr Credentials) error {
	code, msg, err := a.do(ctx, http.MethodPost, "/auth", nil,
		map[string]string{"username": cr.Username, "password": cr.Password})
	if err != nil {
		return err
	}
	if code != http.StatusOK {
		return fmt.Errorf("sign-in with the new admin answered HTTP %d %s", code, msg)
	}
	return nil
}

// client returns the API client, reading the certificate from the distro
// the first time.
func (w *Wizard) client(ctx context.Context) (*api, error) {
	if w.api != nil {
		return w.api, nil
	}
	pem, err := settings.ReadCertPEM(ctx, w.WSL, w.Config.Distro)
	if err != nil {
		return nil, err
	}
	base := w.APIBase
	if base == "" {
		base = fmt.Sprintf("https://127.0.0.1:%d/api", w.Config.Port)
	}
	w.api, err = newAPI(base, pem)
	return w.api, err
}

// waitReady polls the status endpoint until Portainer answers.
func (w *Wizard) waitReady(ctx context.Context) error {
	a, err := w.client(ctx)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeoutCause(ctx, w.ReadyTimeout, errors.New("no answer from Portainer in time"))
	defer cancel()
	for {
		if a.ready(ctx) {
			return nil
		}
		select {
		case <-ctx.Done():
			return context.Cause(ctx)
		case <-time.After(w.Wait):
		}
	}
}

// quiet runs a root script without logging its output.
func (w *Wizard) quiet(ctx context.Context, script string) (string, error) {
	out, err := w.WSL.Script(ctx, w.Config.Distro, "root", script)
	return strings.TrimSpace(out), err
}

func (w *Wizard) stepPortainer(ctx context.Context, _ Credentials) (bool, error) {
	state, err := w.quiet(ctx, portainerStateScript)
	if err != nil {
		return false, err
	}
	switch state {
	case "running":
	case "missing":
		w.logf("Pulling and starting %s", portainerImage)
		if w.token, err = w.NewToken(); err != nil {
			return false, err
		}
		if _, err := w.root(ctx, portainerRunScript(w.Config.Port, w.token, false)); err != nil {
			return false, err
		}
	default:
		w.logf("Starting the stopped container")
		if _, err := w.root(ctx, "docker start portainer >/dev/null\n"); err != nil {
			return false, err
		}
	}
	if err := w.waitReady(ctx); err != nil {
		return false, err
	}
	return state == "running", nil
}

func (w *Wizard) stepAdmin(ctx context.Context, cr Credentials) (bool, error) {
	a, err := w.client(ctx)
	if err != nil {
		return false, err
	}
	exists, err := a.adminExists(ctx)
	if err != nil {
		return false, err
	}
	if exists {
		w.logf("Portainer already has an admin account; its password is unchanged")
		return true, nil
	}
	if w.token == "" {
		// The container predates this run, so its token is unknown: restart
		// it with a token of our own, which also reopens the first-run window.
		if w.token, err = w.NewToken(); err != nil {
			return false, err
		}
		if err := w.recreate(ctx, w.token); err != nil {
			return false, err
		}
	}
	if err := a.initAdmin(ctx, w.token, cr); err != nil {
		return false, err
	}
	w.token = ""
	// Recreate without the token so it is not left in the container config.
	if err := w.recreate(ctx, ""); err != nil {
		return false, err
	}
	if err := a.login(ctx, cr); err != nil {
		return false, err
	}
	w.logf("Signed in as the new admin")
	return false, nil
}

// recreate replaces the container, keeping the data volume, and waits until
// Portainer answers again.
func (w *Wizard) recreate(ctx context.Context, token string) error {
	if _, err := w.root(ctx, portainerRunScript(w.Config.Port, token, true)); err != nil {
		return err
	}
	return w.waitReady(ctx)
}
