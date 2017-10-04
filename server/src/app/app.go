package app

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"jbrodriguez/controlr/plugin/server/src/lib"
	"jbrodriguez/controlr/plugin/server/src/model"
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

	//
	state, err := getUnraidInfo(settings.APIDir, settings.CertDir)
	if err != nil {
		mlog.Fatalf("Unable to retrieve unRAID info (%s). Exiting now ...", err)
	}

	// mlog.Info("Connections to emhttp via %s:%s ...", data["protocol"], data["port"])

	core := services.NewCore(bus, settings, state)
	server := services.NewServer(bus, settings, state)
	api := services.NewApi(bus, settings, state)

	core.Start()
	server.Start()
	api.Start()

	mlog.Info("Press Ctrl+C to stop ...")

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)
	mlog.Info("Received signal: (%s) ... shutting down the app now ...", <-c)

	api.Stop()
	server.Stop()
	core.Stop()

	mlog.Stop()
}

func getUnraidInfo(apiDir, certDir string) (*model.State, error) {
	file, err := ini.LoadFile(filepath.Join(apiDir, "var.ini"))
	if err != nil {
		return nil, err
	}

	secure, err := lib.Exists(filepath.Join(certDir, "certificate_bundle.pem"))
	if err != nil {
		return nil, err
	}

	state := &model.State{}

	var tmp string
	var ok bool

	tmp, _ = file.Get("", "NAME")
	state.Name = strings.Replace(tmp, "\"", "", -1)

	tmp, _ = file.Get("", "timeZone")
	state.Timezone = strings.Replace(tmp, "\"", "", -1)

	tmp, _ = file.Get("", "version")
	state.Version = strings.Replace(tmp, "\"", "", -1)

	tmp, ok = file.Get("", "csrf_token")
	if !ok {
		state.CsrfToken = ""
	} else {
		state.CsrfToken = strings.Replace(tmp, "\"", "", -1)
	}

	var usessl, port, portssl string

	// if the key is missing, usessl, port and portssl are set to ""
	usessl, ok = file.Get("", "USE_SSL")
	port, ok = file.Get("", "PORT")
	portssl, ok = file.Get("", "PORTSSL")

	// remove quotes from unRAID's ini file
	usessl = strings.Replace(usessl, "\"", "", -1)
	port = strings.Replace(port, "\"", "", -1)
	portssl = strings.Replace(portssl, "\"", "", -1)

	// if usessl == "" this isn't a 6.4.x server, try to read emhttpPort; if that doesn't work
	// it will default to port 80
	// otherwise usessl has some value, the plugin will serve off http if the value is no, in any
	// other case, it will serve off https
	if usessl == "" {
		secure = false
		port, ok = file.Get("", "emhttpPort")
		port = strings.Replace(port, "\"", "", -1)
	} else if usessl == "no" {
		secure = false
	}

	state.Secure = secure

	if secure {
		if portssl == "" || portssl == "443" {
			portssl = ""
		} else {
			portssl = ":" + portssl
		}

		state.Host = fmt.Sprintf("https://127.0.0.1%s/", portssl)
	} else {
		if port == "" || port == "80" {
			port = ""
		} else {
			port = ":" + port
		}

		state.Host = fmt.Sprintf("http://127.0.0.1%s/", port)
	}

	//

	// cat := exec.Command("cat", "/boot/config/go")
	// grep := exec.Command("grep", "^/usr/local/sbin/emhttp")

	// // Run the pipeline
	// output, stderr, err := lib.Pipeline(cat, grep)
	// if err != nil {
	// 	mlog.Warning("Failed to run commands to get emhttp port from config: %s\n", err)
	// }

	// // Print the stderr, if any
	// if len(stderr) > 0 {
	// 	mlog.Warning("Error while reading config (stderr): %s\n", stderr)
	// }

	// re := regexp.MustCompile(emhttpRe)
	// args := re.FindStringSubmatch(strings.Trim(string(output), "\n\r"))

	// err, secure, port := lib.GetPort(args)
	// if err != nil {
	// 	mlog.Warning("Unable to get emhttp port (using defaults now): %s", err)
	// }

	// if secure {
	// 	data["protocol"] = "https"
	// } else {
	// 	data["protocol"] = "http"
	// }
	// data["port"] = port

	return state, nil
}
