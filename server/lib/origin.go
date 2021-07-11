package lib

import (
	"io/ioutil"
	"path/filepath"
	"strings"

	"plugin/dto"

	"github.com/vaughan0/go-ini"
)

const nginx = "/var/run/nginx.origin"

func GetIPAddress(apiDir string) (string, error) {
	network, err := ini.LoadFile(filepath.Join(apiDir, "network.ini"))
	if err != nil {
		return "", err
	}

	var tmp string
	tmp, _ = network.Get("eth0", "IPADDR:0")
	ipaddress := strings.Replace(tmp, "\"", "", -1)

	return ipaddress, nil
}

func GetOrigin(apiDir string) *dto.Origin {
	exists, err := Exists(nginx)
	if err != nil {
		return nil
	}

	if !exists {
		return nil
	}
	data, err := ioutil.ReadFile(nginx)
	if err != nil {
		return nil
	}

	origin := string(data)
	origin = strings.Replace(origin, "\n", "", -1)

	address, err := GetIPAddress(apiDir)
	if err != nil {
		return nil
	}

	params := GetParams(`(?P<protocol>^[^:]*)://(?P<host>[^:]*):(?P<port>.*)`, origin)

	return &dto.Origin{Protocol: params["protocol"], Host: params["host"], Port: params["port"], Address: address}
}
