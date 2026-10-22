package ui

import (
	"testing"
	"time"

	"github.com/the-khiem7/docker-zero/internal/status"
	"github.com/the-khiem7/docker-zero/internal/trayicon"
	"github.com/the-khiem7/docker-zero/internal/wsl"
)

func TestLevelOf(t *testing.T) {
	cases := []struct {
		s    status.Snapshot
		want trayicon.Level
	}{
		{status.Snapshot{}, trayicon.Failed},
		{status.Snapshot{WSLAvailable: true, Distro: wsl.StateNotInstalled}, trayicon.Failed},
		{status.Snapshot{WSLAvailable: true, Distro: wsl.StateStopped}, trayicon.Stopped},
		{status.Snapshot{WSLAvailable: true, Distro: wsl.StateRunning}, trayicon.Partial},
		{status.Snapshot{WSLAvailable: true, Distro: wsl.StateRunning, PortainerUp: true}, trayicon.Running},
	}
	for _, tc := range cases {
		if got := LevelOf(tc.s); got != tc.want {
			t.Errorf("LevelOf(%+v) = %d, want %d", tc.s, got, tc.want)
		}
	}
}

func TestDuePollsOnlyWhenVisibleOrBursting(t *testing.T) {
	s := newService(nil, nil)
	now := time.Now()
	s.lastAt = now.Add(-time.Hour)
	if s.due(now) {
		t.Fatal("hidden app polled")
	}
	s.visible[PopupWindow] = true
	if !s.due(now) {
		t.Fatal("visible popup not polled after an hour")
	}
	s.lastAt = now.Add(-2 * time.Second)
	if s.due(now) {
		t.Fatal("visible popup polled faster than every 5 s")
	}
	s.visible[PopupWindow] = false
	s.burstEnd = now.Add(time.Second)
	if !s.due(now) {
		t.Fatal("no poll during a burst")
	}
}
