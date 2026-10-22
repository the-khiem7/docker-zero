package ui

import (
	"context"
	"sync"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"

	"github.com/the-khiem7/docker-zero/internal/setup"
)

// SetupEvent carries a SetupView to the setup window.
const SetupEvent = "setup"

func init() {
	application.RegisterEvent[SetupView](SetupEvent)
}

const (
	// setupTimeout bounds a whole run; the image download dominates it.
	setupTimeout = 45 * time.Minute
	// maxLog is the number of log lines kept for a window opened mid-run.
	maxLog = 400
)

// StepView is one wizard step as the window renders it.
type StepView struct {
	ID     setup.StepID `json:"id"`
	Title  string       `json:"title"`
	Status setup.Status `json:"status"`
	Note   string       `json:"note,omitempty"`
}

// SetupView is the whole wizard state. The window rebuilds itself from it,
// so reopening the window mid-run shows the current progress.
type SetupView struct {
	Steps   []StepView `json:"steps"`
	Log     []string   `json:"log"`
	Running bool       `json:"running"`
	Done    bool       `json:"done"`
	Error   string     `json:"error,omitempty"`
}

// SetupService is bound to the setup window.
type SetupService struct {
	newWizard func() *setup.Wizard
	// emit publishes a view to the window; nil until attach.
	emit func(SetupView)
	// onDone runs after a successful setup, to refresh the status.
	onDone func()

	mu   sync.Mutex
	view SetupView
}

// NewSetupService builds the service. newWizard makes a fresh wizard for each run.
func NewSetupService(newWizard func() *setup.Wizard) *SetupService {
	s := &SetupService{newWizard: newWizard}
	s.view = s.fresh()
	return s
}

func (s *SetupService) attach(app *application.App) {
	s.emit = func(v SetupView) { app.Event.Emit(SetupEvent, v) }
}

func (s *SetupService) fresh() SetupView {
	v := SetupView{Steps: make([]StepView, len(setup.Steps)), Log: []string{}}
	for i, st := range setup.Steps {
		v.Steps[i] = StepView{ID: st.ID, Title: st.Title}
	}
	return v
}

// State returns the current wizard state.
func (s *SetupService) State() SetupView {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.snapshot()
}

func (s *SetupService) snapshot() SetupView {
	v := s.view
	v.Steps = append([]StepView(nil), v.Steps...)
	v.Log = append([]string(nil), v.Log...)
	return v
}

// Start checks the admin credentials and starts the setup in the background.
// Progress arrives as setup events. The credentials are held only for the run.
func (s *SetupService) Start(username, secret string) error {
	cr := setup.Credentials{Username: username, Password: secret}
	if err := cr.Validate(); err != nil {
		return err
	}
	s.mu.Lock()
	if s.view.Running {
		s.mu.Unlock()
		return setup.ErrBusy
	}
	s.view = s.fresh()
	s.view.Running = true
	v := s.snapshot()
	s.mu.Unlock()
	s.publish(v)

	go s.run(cr)
	return nil
}

func (s *SetupService) run(cr setup.Credentials) {
	ctx, cancel := context.WithTimeout(context.Background(), setupTimeout)
	defer cancel()
	w := s.newWizard()
	w.Emit = s.apply
	err := w.Run(ctx, cr)

	s.mu.Lock()
	s.view.Running = false
	s.view.Done = err == nil
	if err != nil {
		s.view.Error = err.Error()
	}
	v := s.snapshot()
	s.mu.Unlock()
	s.publish(v)
	if err == nil && s.onDone != nil {
		s.onDone()
	}
}

// apply folds one wizard event into the view.
func (s *SetupService) apply(e setup.Event) {
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

func (s *SetupService) publish(v SetupView) {
	if s.emit != nil {
		s.emit(v)
	}
}
