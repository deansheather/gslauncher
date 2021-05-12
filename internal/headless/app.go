package headless

import (
	"fmt"

	"github.com/GrooveStats/gslauncher/internal/session"
	"github.com/GrooveStats/gslauncher/internal/unlocks"
)

// App is a headless version of the GUI package. It launches Stepmania
// immediately, and returns when it ends.
type App struct {
	unlockManager *unlocks.Manager
}

func NewApp(unlockManager *unlocks.Manager) *App {
	app := &App{
		unlockManager: unlockManager,
	}

	return app
}

func (app *App) Run() error {
	return app.launchSM()
}

func (app *App) launchSM() error {
	session, err := session.Launch(app.unlockManager)
	if err != nil {
		return fmt.Errorf("launch session: %w", err)
	}

	session.Wait()
	return nil
}
