package lib

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

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
func GetPort(match []string) (bool, string, error) {
	if len(match) == 0 {
		// no match, i can't parse this
		return false, "80", errors.New("unable to parse emhttp flags")
	}

	allFlags := strings.Trim(match[1], " ")
	if allFlags == "&" {
		// emhttp & variant
		return false, "80", nil
	}

	args := strings.Split(allFlags, " ")
	if len(args) <= 2 {
		// at the very least, I should have -p <port(s)>
		return false, "80", errors.New("invalid flags passed to emhttp")
	}

	opts := options{}

	_, err := flags.ParseArgs(&opts, args)
	if err != nil {
		// sent us incorrect flags
		return false, "80", errors.New("invalid flags passed to emhttp (parse)")
	}

	ports := strings.Split(opts.Port, ",")
	if len(ports) == 1 {
		// emhttp -p <m> & variant
		return false, ports[0], nil
	}

	http := ports[0]
	https := ports[1]

	if opts.Redirect {
		if http != "" && https != "" {
			// emhttp -r -p <m>,<n> variant
			return true, https, nil
		}
	} else {
		if https != "" {
			// emhttp -p ,<n> variant and
			// emhttp -p <m>,<n> variant
			return true, https, nil
		}
	}

	// anything else is invalid
	return false, "80", nil
}

// Get -
func Get(client *http.Client, host, resource string) (string, error) {
	ep, err := url.Parse(host)
	if err != nil {
		return "", err
	}

	ep.Path = filepath.Join(ep.Path, resource)

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

// Post -
func Post(client *http.Client, host, resource string, args map[string]string) (string, error) {
	ep, err := url.Parse(host)
	if err != nil {
		return "", err
	}

	ep.Path = filepath.Join(ep.Path, resource)

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

	return resp.Status, nil
}

// Round -
func Round(a float64) int {
	if a < 0 {
		return int(a - 0.5)
	}
	return int(a + 0.5)
}

// GetCmdOutput -
func GetCmdOutput(command string, args ...string) []string {
	lines := make([]string, 0)

	if len(args) > 0 {
		ShellEx(command, func(line string) {
			lines = append(lines, line)
		}, args...)
	} else {
		Shell(command, func(line string) {
			lines = append(lines, line)
		})
	}

	return lines
}

// GenerateCerts - ControlR auto generated TLS cert
func GenerateCerts(name, location string) error {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return err
	}

	template := &x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Country:            []string{"PA"},
			Organization:       []string{"Apertoire"},
			OrganizationalUnit: []string{"ControlR"},
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().AddDate(10, 0, 0),

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	names := []string{"localhost", name, fmt.Sprintf("%s.local", name)}
	template.DNSNames = append(template.DNSNames, names...)

	template.EmailAddresses = []string{fmt.Sprintf("root@%s", name)}

	cert, err := x509.CreateCertificate(rand.Reader, template, template, &priv.PublicKey, priv)
	if err != nil {
		return fmt.Errorf("x590.CreateCertificate: %s", err)
	}

	privDER, err := x509.MarshalPKCS8PrivateKey(priv)
	if err != nil {
		return fmt.Errorf("x509.MarshalPKCS8PrivateKey: %s", err)
	}

	err = ioutil.WriteFile(filepath.Join(location, "controlr_key.pem"), pem.EncodeToMemory(
		&pem.Block{Type: "PRIVATE KEY", Bytes: privDER}), 0600)
	if err != nil {
		return fmt.Errorf("key.WriteFile: %s", err)
	}

	err = ioutil.WriteFile(filepath.Join(location, "controlr_cert.pem"), pem.EncodeToMemory(
		&pem.Block{Type: "CERTIFICATE", Bytes: cert}), 0644)
	if err != nil {
		return fmt.Errorf("certificate.WriteFile: %s", err)
	}

	return nil
}
