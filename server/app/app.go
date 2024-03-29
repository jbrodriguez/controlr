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

	mlog.Info("state(%+v)", state)

	core := services.NewCore(bus, settings, state)
	api := services.NewAPI(bus, settings, state)

	mlog.FatalIfError(core.Start())
	api.Start()

	mlog.Info("Press Ctrl+C to stop ...")

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)
	mlog.Info("Received signal: (%s) ... shutting down the app now ...", <-c)

	api.Stop()
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

	return getNetworkInfo(state, apiDir, certDir, file)
}

func getNetworkInfo(state *model.State, apiDir, certDir string, file ini.File) (*model.State, error) {
	state.Cert = getCertificateName(certDir, state.Name)
	state.UseSelfCerts = false
	if state.Cert == "" {
		if err := provisionSelfCerts(certDir, state.Name); err != nil {
			return nil, err
		}

		state.UseSelfCerts = true
	}

	origin := lib.GetOrigin(apiDir)

	state.Secure = strings.Contains(origin.Protocol, "https")
	state.Origin = *origin
	state.Host = fmt.Sprintf("%s://%s", origin.Protocol, origin.Host)

	return state, nil
}

func getIPAddress(apiDir string) (string, error) {
	network, err := ini.LoadFile(filepath.Join(apiDir, "network.ini"))
	if err != nil {
		return "", err
	}

	var tmp string
	tmp, _ = network.Get("eth0", "IPADDR:0")
	ipaddress := strings.Replace(tmp, "\"", "", -1)

	return ipaddress, nil
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
