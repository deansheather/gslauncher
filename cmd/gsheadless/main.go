package main

import (
	"log"
	"os"
	"runtime"

	"github.com/GrooveStats/gslauncher/internal/headless"
	"github.com/GrooveStats/gslauncher/internal/settings"
	"github.com/GrooveStats/gslauncher/internal/unlocks"
	"github.com/GrooveStats/gslauncher/internal/version"
)

func main() {
	log.Printf("GrooveStats Launcher (headless mode) %s (%s %s)", version.Formatted(), runtime.GOOS, runtime.GOARCH)

	err := settings.Load()
	if os.IsNotExist(err) {
		settingsPath, err := settings.SettingsFile()
		if err != nil {
			log.Fatalf("get settings path: %+v", err)
		}

		log.Printf("writing default settings file to %v", settingsPath)
		err = settings.Save()
		if err != nil {
			log.Fatalf("write settings: %+v", err)
		}
	} else if err != nil {
		log.Fatalf("failed to load settings: %+v", err)
	}

	unlockManager, err := unlocks.NewManager()
	if err != nil {
		log.Fatalf("failed to initialize downloader: %+v", err)
	}

	app := headless.NewApp(unlockManager)
	err = app.Run()
	if err != nil {
		log.Fatalf("run headless app: %+v", err)
	}
}
