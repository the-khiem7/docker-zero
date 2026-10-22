package ui

import (
	"context"
	"sync"

	"github.com/wailsapp/wails/v3/pkg/application"

	"github.com/the-khiem7/docker-zero/internal/setup"
)

// UninstallEvent carries a SetupView, which has the same shape for every
// stepped action, to the settings window.
const UninstallEvent = "uninstall"

func init() {
	application.RegisterEvent[SetupView](UninstallEvent)
}

// UninstallService is bound to the settings window.
type UninstallService struct {
	newWizard func() *setup.Wizard
	guard     *opGuard
	app       *application.App
	// emit publishes a view to the window; nil until attach.
	emit func(SetupView)
	// onDone runs after the uninstall, to refresh the status.
	onDone func()

	mu   sync.Mutex
	view SetupView
	w    *setup.Wizard
}

// NewUninstallService builds the service. newWizard makes a fresh wizard for each run.
func NewUninstallService(newWizard func() *setup.Wizard, guard *opGuard) *UninstallService {
	s := &UninstallService{newWizard: newWizard, guard: guard}
	s.view = s.fresh()
	return s
}

func (s *UninstallService) attach(app *application.App) {
	s.app = app
	s.emit = func(v SetupView) { app.Event.Emit(UninstallEvent, v) }
}

func (s *UninstallService) fresh() SetupView {
	v := SetupView{Steps: make([]StepView, len(setup.UninstallSteps)), Log: []string{}}
	for i, st := range setup.UninstallSteps {
		v.Steps[i] = StepView{ID: st.ID, Title: st.Title}
	}
	return v
}

// State returns the current uninstall state.
func (s *UninstallService) State() SetupView {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.snapshot()
}

func (s *UninstallService) snapshot() SetupView {
	v := s.view
	v.Steps = append([]StepView(nil), v.Steps...)
	v.Log = append([]string(nil), v.Log...)
	return v
}

// PickBackup asks where to save the backup and returns the path, or "".
func (s *UninstallService) PickBackup() (string, error) {
	return s.app.Dialog.SaveFile().
		SetMessage("Save a backup of Ubuntu-24.04").
		SetFilename("docker-zero-backup.tar").
		AddFilter("WSL backup", "*.tar").
		CanCreateDirectories(true).
		PromptForSingleSelection()
}

// Start removes everything Docker Zero created, after saving a backup to
// backupPath when it is not empty. The window has asked the user twice.
// Progress arrives as uninstall events.
func (s *UninstallService) Start(backupPath string) error {
	opts := setup.UninstallOptions{BackupPath: backupPath}
	if err := opts.Validate(); err != nil {
		return err
	}
	if err := s.guard.acquire("uninstall"); err != nil {
		return err
	}
	s.mu.Lock()
	s.view = s.fresh()
	s.view.Running = true
	s.w = s.newWizard()
	s.w.Emit = s.apply
	w := s.w
	v := s.snapshot()
	s.mu.Unlock()
	s.publish(v)

	go s.run(w, opts)
	return nil
}

func (s *UninstallService) run(w *setup.Wizard, opts setup.UninstallOptions) {
	ctx, cancel := context.WithTimeout(context.Background(), setupTimeout)
	defer cancel()
	err := w.Uninstall(ctx, opts)

	s.mu.Lock()
	s.view.Running = false
	s.view.Done = err == nil
	if err != nil {
		s.view.Error = err.Error()
	}
	v := s.snapshot()
	s.mu.Unlock()
	s.guard.release()
	s.publish(v)
	if s.onDone != nil {
		s.onDone()
	}
}

// Quit schedules the removal of the install folder and quits the app. The
// window offers it once the uninstall finished.
func (s *UninstallService) Quit() error {
	s.mu.Lock()
	w, done := s.w, s.view.Done
	s.mu.Unlock()
	if w == nil || !done {
		return nil
	}
	err := w.FinishUninstall()
	s.app.Quit()
	return err
}

// apply folds one wizard event into the view.
func (s *UninstallService) apply(e setup.Event) {
	s.mu.Lock()
	if e.Status == setup.StatusLog {
		s.view.Log = append(s.view.Log, e.Message)
		if over := len(s.view.Log) - maxLog; over > 0 {
			s.view.Log = s.view.Log[over:]
		}
	} else {
		for i := range s.view.Steps {
			if s.view.Steps[i].ID == e.Step {
				s.view.Steps[i].Status = e.Status
				s.view.Steps[i].Note = e.Message
			}
		}
	}
	v := s.snapshot()
	s.mu.Unlock()
	s.publish(v)
}

func (s *UninstallService) publish(v SetupView) {
	if s.emit != nil {
		s.emit(v)
	}
}
