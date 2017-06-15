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

	LogDir     string
	WebDir     string
	APIDir     string
	CertDir    string
	Port       string
	SPort      string
	ProxyPort  string
	SProxyPort string
	Dev        bool

	Location string
}

// NewSettings constructor
func NewSettings(name, version string, locations []string) (*Settings, error) {
	var logDir, webDir, apiDir, certDir, port, proxyPort, sport, sproxyPort string
	var dev bool

	location := SearchFile(name, locations)

	flag.StringVar(&logDir, "logdir", "/boot/logs", "folder containing the log files")
	flag.StringVar(&webDir, "webdir", "", "folder containing the ui")
	flag.StringVar(&apiDir, "apidir", "/var/local/emhttp", "folders to look for api endpoints")
	flag.StringVar(&certDir, "certdir", "/boot/config/plugins/controlr", "folders to look for https certs")
	flag.StringVar(&port, "port", "2378", "port to run the http server")
	flag.StringVar(&sport, "sport", "2379", "port to run the https server")
	flag.StringVar(&proxyPort, "proxyport", "2382", "port to run the http api endpoint")
	flag.StringVar(&sproxyPort, "sproxyport", "2383", "port to run the https api endpoint")
	flag.BoolVar(&dev, "dev", false, "work in dev mode for some features")

	flag.Parse()

	s := &Settings{}
	s.Version = version
	s.LogDir = logDir
	s.WebDir = webDir
	s.APIDir = apiDir
	s.CertDir = certDir
	s.Port = port
	s.ProxyPort = proxyPort
	s.SPort = sport
	s.SProxyPort = sproxyPort
	s.Location = location
	s.Dev = dev

	return s, nil
}
