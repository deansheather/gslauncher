package main

import (
	"bytes"
	"flag"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/GrooveStats/gslauncher/internal/gui"
	"github.com/GrooveStats/gslauncher/internal/headless"
	"github.com/GrooveStats/gslauncher/internal/settings"
	"github.com/GrooveStats/gslauncher/internal/unlocks"
	"github.com/GrooveStats/gslauncher/internal/version"
)

const MAX_LOG_SIZE = 1024 * 1024 // 1 MiB

func redirectLog() {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		log.Print("failed to get cache directory: ", err)
		return
	}

	filename := filepath.Join(cacheDir, "groovestats-launcher", "log.txt")

	old, err := os.ReadFile(filename)
	if err == nil {
		if len(old) > MAX_LOG_SIZE {
			old = old[len(old)-MAX_LOG_SIZE:]
			idx := bytes.IndexByte(old, byte('\n'))
			old = old[idx+1:]
		}
	}

	logFile, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Print("failed to open log file: ", err)
		return
	}

	if old != nil {
		logFile.Write(old)
		logFile.WriteString("-----\n")
	}

	if settings.Get().Debug {
		log.SetOutput(io.MultiWriter(os.Stderr, logFile))
	} else {
		log.SetOutput(logFile)
	}
}

func main() {
	autolaunch := flag.Bool("autolaunch", false, "automatically launch StepMania")
	headlessMode := flag.Bool("headless", false, "launch without GUI window (implies --autolaunch)")
	flag.Parse()
	if *headlessMode {
		*autolaunch = true
	}

	modeStr := "headless"
	if !*headlessMode {
		modeStr = "gui"
		redirectLog()
	}
	log.Printf("GrooveStats Launcher (%s mode) %s (%s %s)", modeStr, version.Formatted(), runtime.GOOS, runtime.GOARCH)

	err := settings.Load()
	if os.IsNotExist(err) {
		settingsPath, err := settings.SettingsFile()
		if err != nil {
			log.Fatalf("get settings path: %+v", err)
		}

		settings.DetectSM()

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

	if *headlessMode {
		app := headless.NewApp(unlockManager)
		err = app.Run()
	} else {
		app := gui.NewApp(unlockManager, *autolaunch)
		err = app.Run()
	}
	if err != nil {
		log.Fatalf("run %v app: %+v", modeStr, err)
	}
}
