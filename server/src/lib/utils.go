package lib

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/jessevdk/go-flags"
)

// Exists - check if File / Directory Exists
func Exists(path string) (bool, error) {
	_, err := os.Stat(path)

	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}

// SearchFile - look for file in the specified locations
func SearchFile(name string, locations []string) string {
	for _, location := range locations {
		if b, _ := Exists(filepath.Join(location, name)); b {
			return location
		}
	}

	return ""
}

type options struct {
	Port     string `short:"p" long:"port" description:"port to use"`
	Redirect bool   `short:"r" long:"redirect" description:"redirect http to https"`
}

// GetPort - define emhttp port from string
func GetPort(match []string) (error, bool, string) {
	if len(match) == 0 {
		// no match, i can't parse this
		return errors.New("Unable to parse emhttp flags"), false, "80"
	}

	allFlags := strings.Trim(match[1], " ")
	if allFlags == "&" {
		// emhttp & variant
		return nil, false, "80"
	}

	args := strings.Split(allFlags, " ")
	if len(args) <= 2 {
		// at the very least, I should have -p <port(s)>
		return errors.New("Invalid flags passed to emhttp"), false, "80"
	}

	opts := options{}

	_, err := flags.ParseArgs(&opts, args)
	if err != nil {
		// sent us incorrect flags
		return errors.New("Invalid flags passed to emhttp (parse)"), false, "80"
	}

	ports := strings.Split(opts.Port, ",")
	if len(ports) == 1 {
		// emhttp -p <m> & variant
		return nil, false, ports[0]
	}

	http := ports[0]
	https := ports[1]

	if opts.Redirect {
		if http != "" && https != "" {
			// emhttp -r -p <m>,<n> variant
			return nil, true, https
		}
	} else {
		if https != "" {
			// emhttp -p ,<n> variant and
			// emhttp -p <m>,<n> variant
			return nil, true, https
		}
	}

	// anything else is invalid
	return nil, false, "80"
}

func Get(client *http.Client, host, resource string) (string, error) {
	ep, err := url.Parse(host)
	if err != nil {
		return "", err
	}

	ep.Path = path.Join(ep.Path, resource)

	req, err := http.NewRequest("GET", ep.String(), nil)
	if err != nil {
		return "", err
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New(resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func Post(client *http.Client, host, resource string, args map[string]string) (string, error) {
	ep, err := url.Parse(host)
	if err != nil {
		return "", err
	}

	ep.Path = path.Join(ep.Path, resource)

	data := url.Values{}
	for k, v := range args {
		data.Set(k, v)
	}

	req, err := http.NewRequest("POST", ep.String(), strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	return string(resp.Status), nil
}

func Round(a float64) int {
	if a < 0 {
		return int(a - 0.5)
	}
	return int(a + 0.5)
}

func GetCmdOutput(command string, args string) []string {
	lines := make([]string, 0)

	if args != "" {
		ShellEx(command, func(line string) {
			lines = append(lines, line)
		}, args)
	} else {
		Shell(command, func(line string) {
			lines = append(lines, line)
		})
	}

	return lines
}
