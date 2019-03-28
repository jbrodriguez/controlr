package app

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"plugin/lib"
	"plugin/model"
	"plugin/services"

	"github.com/jbrodriguez/mlog"
	"github.com/jbrodriguez/pubsub"
	"github.com/vaughan0/go-ini"
)

const identCfg = "/boot/config/ident.cfg"

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

	mlog.Info("controlr %s starting ...", settings.Version)

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
	api := services.NewAPI(bus, settings, state)

	mlog.FatalIfError(core.Start())
	server.Start()
	api.Start()

	mlog.Info("Press Ctrl+C to stop ...")

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)
	mlog.Info("Received signal: (%s) ... shutting down the app now ...", <-c)

	api.Stop()
	server.Stop()
	core.Stop()

	if err := mlog.Stop(); err != nil {
		log.Printf("error stopping mlog: %s", err)
	}
}

func getUnraidInfo(apiDir, certDir string) (*model.State, error) {
	file, err := ini.LoadFile(filepath.Join(apiDir, "var.ini"))
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

	// use main ssl config file
	ident, err := ini.LoadFile(identCfg)
	if err != nil {
		return nil, err
	}

	var usessl, port, portssl string

	// if the key is missing, usessl, port and portssl are set to ""
	usessl, _ = ident.Get("", "USE_SSL")
	port, _ = ident.Get("", "PORT")
	portssl, _ = ident.Get("", "PORTSSL")

	// remove quotes from unRAID's ini file
	usessl = strings.Replace(usessl, "\"", "", -1)
	port = strings.Replace(port, "\"", "", -1)
	portssl = strings.Replace(portssl, "\"", "", -1)

	state.Cert = getCertificateName(certDir, state.Name)
	state.UseSelfCerts = false
	if state.Cert == "" {
		if err := provisionSelfCerts(certDir, state.Name); err != nil {
			return nil, err
		}

		state.UseSelfCerts = true
	}

	secure := state.Cert != ""

	// if usessl == "" this isn't a 6.4.x server, try to read emhttpPort; if that doesn't work
	// it will default to port 80
	// otherwise usessl has some value, the plugin will serve off http if the value is no, in any
	// other case, it will serve off https
	if usessl == "" {
		secure = false
		port, _ = file.Get("", "emhttpPort")
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

	return state, nil
}

func getCertificateName(certDir, name string) string {
	cert := "certificate_bundle.pem"

	exists, err := lib.Exists(filepath.Join(certDir, cert))
	if err != nil {
		mlog.Warning("unable to check for %s presence:(%s)", cert, err)
		return ""
	}

	if exists {
		mlog.Info("cert: found %s", cert)
		return cert
	}

	cert = name + "_unraid_bundle.pem"

	exists, err = lib.Exists(filepath.Join(certDir, cert))
	if err != nil {
		return ""
	}

	if exists {
		mlog.Info("cert: found %s", cert)
		return cert
	}

	return ""
}

func provisionSelfCerts(certDir, name string) error {
	certExists, err := lib.Exists(filepath.Join(certDir, "controlr_cert.pem"))
	if err != nil {
		mlog.Warning("unable to check for cert presence:(%s)", err)
		return err
	}

	keyExists, err := lib.Exists(filepath.Join(certDir, "controlr_key.pem"))
	if err != nil {
		mlog.Warning("unable to check for key presence:(%s)", err)
		return err
	}

	if certExists && keyExists {
		return nil
	}

	mlog.Info("cert: generating self certs")

	return lib.GenerateCerts(name, certDir)
}
