//go:build integration

package wsl

import (
	"os"
	"testing"
	"time"
)

// TestActionsReal drives Start, StartHold, Release and Terminate against a
// throwaway distro named by DZ_TEST_DISTRO, never the real one. Shutdown is
// not exercised because it stops every distro of the account.
//
//	DZ_TEST_DISTRO=dz-test go test -tags integration ./internal/wsl
func TestActionsReal(t *testing.T) {
	distro := os.Getenv("DZ_TEST_DISTRO")
	if distro == "" {
		t.Skip("DZ_TEST_DISTRO not set")
	}
	c := New()
	ctx := t.Context()

	want := func(s State) {
		t.Helper()
		got, err := c.DistroState(ctx, distro)
		if err != nil {
			t.Fatal(err)
		}
		if got != s {
			t.Fatalf("state %q, want %q", got, s)
		}
	}

	if err := c.Terminate(ctx, distro); err != nil {
		t.Fatal(err)
	}
	want(StateStopped)

	if err := c.Start(ctx, distro); err != nil {
		t.Fatal(err)
	}
	want(StateRunning)
	if err := c.Terminate(ctx, distro); err != nil {
		t.Fatal(err)
	}
	want(StateStopped)

	h, err := c.StartHold(distro)
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(3 * time.Second)
	want(StateRunning)
	if !h.Alive() {
		t.Fatal("hold ended early")
	}
	out, err := c.Script(ctx, distro, "root", "ps -o pid,comm | grep -c sleep || true\n")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("inside distro: %q", out)

	if err := c.Terminate(ctx, distro); err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Second)
	if h.Alive() {
		t.Fatal("hold survived terminate")
	}
	h.Release()
	want(StateStopped)

	// Releasing a live hold must end the sleep inside the distro too,
	// otherwise the distro would never idle out.
	h, err = c.StartHold(distro)
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(3 * time.Second)
	h.Release()
	time.Sleep(time.Second)
	out, err = c.Script(ctx, distro, "root", "ps -o comm | grep -c '^sleep' || true\n")
	if err != nil {
		t.Fatal(err)
	}
	if out != "0\n" {
		t.Fatalf("sleep processes after release: %q", out)
	}
	if err := c.Terminate(ctx, distro); err != nil {
		t.Fatal(err)
	}
}
