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

	LogDir string
	WebDir string
	APIDir string
	Port   string

	Location string
}

// NewSettings constructor
func NewSettings(name, version string, locations []string) (*Settings, error) {
	var logDir, webDir, apiDir, port string

	location := SearchFile(name, locations)

	flag.StringVar(&logDir, "logdir", "/boot/logs", "folder containing the log files")
	flag.StringVar(&webDir, "webdir", "", "folder containing the ui")
	flag.StringVar(&apiDir, "apidir", "/var/local/emhttp", "folders to look for api endpoints")
	flag.StringVar(&port, "port", "2378", "port to run the server")

	flag.Parse()

	s := &Settings{}
	s.Version = version
	s.LogDir = logDir
	s.WebDir = webDir
	s.APIDir = apiDir
	s.Port = port
	s.Location = location

	return s, nil
}
