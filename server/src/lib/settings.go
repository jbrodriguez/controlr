package lib

import (
	"github.com/namsral/flag"
)

// Config template
type Config struct {
	Version string
}

// Settings template
type Settings struct {
	Config

	LogDir  string
	WebDir  string
	APIDir  string
	CertDir string
	Port    string
	ApiPort string
	ShowUps bool
	Dev     bool

	Location string
}

// NewSettings constructor
func NewSettings(name, version string, locations []string) (*Settings, error) {
	var logDir, webDir, apiDir, certDir, port, apiPort string
	var showUps, dev bool

	location := SearchFile(name, locations)

	flag.StringVar(&logDir, "logdir", "/boot/logs", "folder containing the log files")
	flag.StringVar(&webDir, "webdir", "", "folder containing the ui")
	flag.StringVar(&apiDir, "apidir", "/var/local/emhttp", "folders to look for api endpoints")
	flag.StringVar(&certDir, "certdir", "/boot/config/ssl/certs", "folders to look for https certs")
	flag.StringVar(&port, "port", "2378", "port to run the http server")
	flag.StringVar(&apiPort, "apiport", "2382", "port to run the http api endpoint")
	flag.BoolVar(&showUps, "showups", false, "whether to provide ups status or not")
	flag.BoolVar(&dev, "dev", false, "work in dev mode for some features")

	flag.Parse()

	s := &Settings{}
	s.Version = version
	s.LogDir = logDir
	s.WebDir = webDir
	s.APIDir = apiDir
	s.CertDir = certDir
	s.Port = ":" + port
	s.ApiPort = ":" + apiPort
	s.Location = location
	s.ShowUps = showUps
	s.Dev = dev

	return s, nil
}
