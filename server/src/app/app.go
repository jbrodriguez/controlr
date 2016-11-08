package app

import (
	"fmt"
	"github.com/jbrodriguez/mlog"
	"github.com/jbrodriguez/pubsub"
	"github.com/vaughan0/go-ini"
	"jbrodriguez/controlr/plugin/server/src/lib"
	"jbrodriguez/controlr/plugin/server/src/services"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
)

// App empty placeholder
type App struct {
}

// Setup app
func (a *App) Setup(version string) (*lib.Settings, error) {
	// look for controlr.conf at the following places
	// /boot/config/plugins/controlr/
	// <current dir>/controlr.conf
	// home := os.Getenv("HOME")

	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	locations := []string{
		"/boot/config/plugins/controlr",
		cwd,
	}

	settings, err := lib.NewSettings("controlr.conf", version, locations)

	return settings, err
}

// Run app
func (a *App) Run(settings *lib.Settings) {
	if settings.LogDir != "" {
		mlog.Start(mlog.LevelInfo, filepath.Join(settings.LogDir, "controlr.log"))
	} else {
		mlog.Start(mlog.LevelInfo, "")
	}

	mlog.Info("controlr v%s starting ...", settings.Version)

	var msg string
	if settings.Location == "" {
		msg = "No config file specified. Using app defaults ..."
	} else {
		msg = fmt.Sprintf("Using config file located at %s ...", settings.Location)
	}
	mlog.Info(msg)

	bus := pubsub.New(623)

	data := a.getUnraidInfo(settings.APIDir)
	if data == nil {
		mlog.Fatalf("Unable to retrieve unRAID info. Exiting now ...")
	}

	template := "http://127.0.0.1%s/"
	replacement := ""
	if data["port"] != "80" {
		replacement = ":" + data["port"]
	}
	data["backend"] = fmt.Sprintf(template, replacement)

	unraid := services.NewUnraid(bus, settings, data)
	server := services.NewServer(bus, settings, data)
	proxy := services.NewProxy(bus, settings, data)

	unraid.Start()
	server.Start()
	proxy.Start()

	mlog.Info("Press Ctrl+C to stop ...")

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)
	mlog.Info("Received signal: (%s) ... shutting down the app now ...", <-c)

	proxy.Stop()
	server.Stop()
	unraid.Stop()

	mlog.Stop()
}

func (a *App) getUnraidInfo(location string) map[string]string {
	var data map[string]string

	file, err := ini.LoadFile(filepath.Join(location, "var.ini"))
	if err != nil {
		return nil
	}

	data = make(map[string]string, 0)

	tmp, _ := file.Get("", "emhttpPort")
	data["port"] = strings.Replace(tmp, "\"", "", -1)

	tmp, _ = file.Get("", "GATEWAY")
	data["gateway"] = strings.Replace(tmp, "\"", "", -1)

	tmp, _ = file.Get("", "NAME")
	data["name"] = strings.Replace(tmp, "\"", "", -1)

	tmp, _ = file.Get("", "timeZone")
	data["timezone"] = strings.Replace(tmp, "\"", "", -1)

	tmp, _ = file.Get("", "version")
	data["version"] = strings.Replace(tmp, "\"", "", -1)

	return data
}
