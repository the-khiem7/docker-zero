package wsl

import (
	"context"
	"encoding/binary"
	"errors"
	"io"
	"slices"
	"strings"
	"testing"
	"unicode/utf16"
)

func utf16le(s string) []byte {
	u := utf16.Encode([]rune(s))
	b := make([]byte, 2*len(u))
	for i, v := range u {
		binary.LittleEndian.PutUint16(b[2*i:], v)
	}
	return b
}

func TestDecodeOutput(t *testing.T) {
	const text = "Ubuntu-24.04\r\ndocker-desktop\r\n"
	cases := map[string][]byte{
		"utf16le":     utf16le(text),
		"utf16le bom": append([]byte{0xff, 0xfe}, utf16le(text)...),
		"utf8":        []byte(text),
		"utf8 bom":    append([]byte{0xef, 0xbb, 0xbf}, text...),
	}
	for name, in := range cases {
		t.Run(name, func(t *testing.T) {
			if got := decodeOutput(in); got != text {
				t.Fatalf("got %q, want %q", got, text)
			}
		})
	}
}

func TestParseNames(t *testing.T) {
	got := parseNames("Ubuntu-24.04\r\n\r\n  docker-desktop \r\n")
	if want := []string{"Ubuntu-24.04", "docker-desktop"}; !slices.Equal(got, want) {
		t.Fatalf("got %q, want %q", got, want)
	}
}

type call struct {
	args  []string
	stdin string
}

// fake answers each wsl.exe call from a table keyed by the joined args.
type fake struct {
	calls   []call
	answers map[string]Result
	err     error
}

func (f *fake) run(_ context.Context, _ string, args []string, stdin io.Reader) (Result, error) {
	c := call{args: args}
	if stdin != nil {
		b, _ := io.ReadAll(stdin)
		c.stdin = string(b)
	}
	f.calls = append(f.calls, c)
	if f.err != nil {
		return Result{}, f.err
	}
	return f.answers[strings.Join(args, " ")], nil
}

func TestDistroState(t *testing.T) {
	cases := []struct {
		name    string
		all     string
		running Result
		want    State
	}{
		{"running", "Ubuntu-24.04\r\n", Result{Stdout: utf16le("Ubuntu-24.04\r\n")}, StateRunning},
		{"stopped", "Ubuntu-24.04\r\nother\r\n", Result{Stdout: []byte("other\n")}, StateStopped},
		{"stopped, none running exits non-zero", "Ubuntu-24.04\r\n", Result{ExitCode: 1, Stdout: utf16le("There are no running distributions.")}, StateStopped},
		{"not installed", "other\r\n", Result{}, StateNotInstalled},
		{"case-insensitive", "ubuntu-24.04\r\n", Result{Stdout: []byte("UBUNTU-24.04\n")}, StateRunning},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			f := &fake{answers: map[string]Result{
				"--list --quiet":           {Stdout: utf16le(tc.all)},
				"--list --quiet --running": tc.running,
			}}
			got, err := NewWithRunner(f.run).DistroState(t.Context(), "Ubuntu-24.04")
			if err != nil {
				t.Fatal(err)
			}
			if got != tc.want {
				t.Fatalf("got %q, want %q", got, tc.want)
			}
		})
	}
}

func TestScriptSendsLFOnStdin(t *testing.T) {
	f := &fake{answers: map[string]Result{
		"--distribution Ubuntu-24.04 --user root --exec /bin/bash -s": {Stdout: []byte("active\n")},
	}}
	out, err := NewWithRunner(f.run).Script(t.Context(), "Ubuntu-24.04", "root", "echo a\r\necho b\r\n")
	if err != nil {
		t.Fatal(err)
	}
	if out != "active\n" {
		t.Fatalf("stdout %q", out)
	}
	if got := f.calls[0].stdin; got != "echo a\necho b\n" {
		t.Fatalf("stdin %q", got)
	}
}

func TestCommandErrorCarriesDecodedMessage(t *testing.T) {
	f := &fake{answers: map[string]Result{
		"--terminate nope": {ExitCode: -1, Stdout: utf16le("There is no distribution with the supplied name.\r\n")},
	}}
	err := NewWithRunner(f.run).Terminate(t.Context(), "nope")
	cmdErr, ok := errors.AsType[*CommandError](err)
	if !ok {
		t.Fatalf("got %v, want *CommandError", err)
	}
	if cmdErr.Message != "There is no distribution with the supplied name." {
		t.Fatalf("message %q", cmdErr.Message)
	}
}

func TestAvailableWrapsSentinel(t *testing.T) {
	f := &fake{err: ErrWSLUnavailable}
	if err := NewWithRunner(f.run).Available(t.Context()); !errors.Is(err, ErrWSLUnavailable) {
		t.Fatalf("got %v", err)
	}
}

func TestStopActionsArgs(t *testing.T) {
	f := &fake{answers: map[string]Result{}}
	c := NewWithRunner(f.run)
	if err := c.Start(t.Context(), "Ubuntu-24.04"); err != nil {
		t.Fatal(err)
	}
	if err := c.Terminate(t.Context(), "Ubuntu-24.04"); err != nil {
		t.Fatal(err)
	}
	if err := c.Shutdown(t.Context()); err != nil {
		t.Fatal(err)
	}
	want := [][]string{
		{"--distribution", "Ubuntu-24.04", "--exec", "/bin/true"},
		{"--terminate", "Ubuntu-24.04"},
		{"--shutdown"},
	}
	for i, w := range want {
		if !slices.Equal(f.calls[i].args, w) {
			t.Fatalf("call %d: got %q, want %q", i, f.calls[i].args, w)
		}
	}
}

func TestExportArgs(t *testing.T) {
	f := &fake{answers: map[string]Result{}}
	c := NewWithRunner(f.run)
	if err := c.Export(t.Context(), "dz-throwaway", `C:\b\backup.tar`); err != nil {
		t.Fatal(err)
	}
	if want := []string{"--export", "dz-throwaway", `C:\b\backup.tar`}; !slices.Equal(f.calls[0].args, want) {
		t.Fatalf("got %q, want %q", f.calls[0].args, want)
	}
}

func TestInstallAndUnregisterArgs(t *testing.T) {
	f := &fake{answers: map[string]Result{}}
	c := NewWithRunner(f.run)
	if err := c.Install(t.Context(), "Ubuntu-24.04", "ubuntu-24.04", ""); err != nil {
		t.Fatal(err)
	}
	if err := c.Install(t.Context(), "Ubuntu-24.04", "dz-throwaway", `C:\d`); err != nil {
		t.Fatal(err)
	}
	if err := c.Unregister(t.Context(), "dz-throwaway"); err != nil {
		t.Fatal(err)
	}
	want := [][]string{
		{"--install", "--distribution", "Ubuntu-24.04", "--web-download", "--no-launch"},
		{"--install", "--distribution", "Ubuntu-24.04", "--web-download", "--no-launch", "--name", "dz-throwaway", "--location", `C:\d`},
		{"--unregister", "dz-throwaway"},
	}
	for i, w := range want {
		if !slices.Equal(f.calls[i].args, w) {
			t.Fatalf("call %d: got %q, want %q", i, f.calls[i].args, w)
		}
	}
}
