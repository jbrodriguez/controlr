package app

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"regexp"
	"strings"
	"syscall"

	"jbrodriguez/controlr/plugin/server/src/lib"
	"jbrodriguez/controlr/plugin/server/src/services"

	"github.com/jbrodriguez/mlog"
	"github.com/jbrodriguez/pubsub"
	"github.com/vaughan0/go-ini"
)

const emhttpRe = `.*?emhttp(.*)$`

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

	mlog.Info("Connections to emhttp via %s:%s ...", data["protocol"], data["port"])

	exists, err := lib.Exists(filepath.Join(settings.CertDir, "cert.pem"))
	if err != nil {
		mlog.Warning("Unable to check for certs presence: %s", err)
	}

	if !exists {
		mlog.Info("No certs are available, generating ...")
		err := lib.GenerateCerts(data["name"], settings.CertDir)
		if err != nil {
			mlog.Warning("Unable to generate certs: %s", err)
		}
	}

	template := "%s://127.0.0.1%s/"
	port := ""
	if (data["protocol"] == "http" && data["port"] != "80") || (data["protocol"] == "https" && data["port"] != "443") {
		port = ":" + data["port"]
	}
	data["backend"] = fmt.Sprintf(template, data["protocol"], port)

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

	tmp, _ := file.Get("", "NAME")
	data["name"] = strings.Replace(tmp, "\"", "", -1)

	tmp, _ = file.Get("", "timeZone")
	data["timezone"] = strings.Replace(tmp, "\"", "", -1)

	tmp, _ = file.Get("", "version")
	data["version"] = strings.Replace(tmp, "\"", "", -1)

	token, ok := file.Get("", "csrf_token")
	if !ok {
		data["csrf_token"] = ""
	} else {
		data["csrf_token"] = strings.Replace(token, "\"", "", -1)
	}

	cat := exec.Command("cat", "/boot/config/go")
	grep := exec.Command("grep", "^/usr/local/sbin/emhttp")

	// Run the pipeline
	output, stderr, err := lib.Pipeline(cat, grep)
	if err != nil {
		mlog.Warning("Failed to run commands to get emhttp port from config: %s\n", err)
	}

	// Print the stderr, if any
	if len(stderr) > 0 {
		mlog.Warning("Error while reading config (stderr): %s\n", stderr)
	}

	re := regexp.MustCompile(emhttpRe)
	args := re.FindStringSubmatch(strings.Trim(string(output), "\n\r"))

	err, secure, port := lib.GetPort(args)
	if err != nil {
		mlog.Warning("Unable to get emhttp port (using defaults now): %s", err)
	}

	if secure {
		data["protocol"] = "https"
	} else {
		data["protocol"] = "http"
	}
	data["port"] = port

	return data
}
