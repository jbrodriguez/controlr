package lib

import (
	"io/ioutil"
	"path/filepath"
	"strings"

	"plugin/dto"

	"github.com/vaughan0/go-ini"
)

const nginx = "/var/run/nginx.origin"
const config = "/var/local/emhttp/var.ini"

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

func getOriginFromFile(address string) *dto.Origin {
	// use main ssl config file
	ident, err := ini.LoadFile(config)
	if err != nil {
		return nil
	}

	var usessl, portnossl, portssl, protocol, name, port string

	// if the key is missing, usessl, port and portssl are set to ""
	usessl, _ = ident.Get("", "USE_SSL")
	portnossl, _ = ident.Get("", "PORT")
	portssl, _ = ident.Get("", "PORTSSL")
	name, _ = ident.Get("", "NAME")

	// remove quotes from unRAID's ini file
	usessl = strings.Replace(usessl, "\"", "", -1)
	portnossl = strings.Replace(portnossl, "\"", "", -1)
	portssl = strings.Replace(portssl, "\"", "", -1)
	name = strings.Replace(name, "\"", "", -1)

	if usessl == "no" {
		protocol = "http"
		port = portnossl
	} else {
		protocol = "http"
		port = portssl
	}
	return &dto.Origin{Protocol: protocol, Host: name, Port: port, Name: name, Address: address}
}

func GetOrigin(apiDir string) *dto.Origin {
	exists, err := Exists(nginx)
	if err != nil {
		return nil
	}

	address, err := GetIPAddress(apiDir)
	if err != nil {
		return nil
	}

	if exists {
		data, err := ioutil.ReadFile(nginx)
		if err != nil {
			return nil
		}

		origin := string(data)
		origin = strings.Replace(origin, "\n", "", -1)

		params := GetParams(`(?P<protocol>^[^:]*)://(?P<name>[^\.]*)\.(?P<tld>[^:]*):(?P<port>.*)`, origin)

		return &dto.Origin{Protocol: params["protocol"], Host: params["name"], Port: params["port"], Name: params["name"], Address: address}
	} else {
		return getOriginFromFile(address)
	}
}
