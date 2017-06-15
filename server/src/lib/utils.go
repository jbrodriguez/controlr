package lib

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jessevdk/go-flags"
)

const (
	bytec    = 1.0
	kilobyte = 1024 * bytec
	megabyte = 1024 * kilobyte
	gigabyte = 1024 * megabyte
	terabyte = 1024 * gigabyte
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

// ByteSize - convert to decimal notation
func ByteSize(bytes int64) string {
	unit := ""
	value := float32(bytes)

	switch {
	case bytes >= terabyte:
		unit = "T"
		value = value / terabyte
	case bytes >= gigabyte:
		unit = "G"
		value = value / gigabyte
	case bytes >= megabyte:
		unit = "M"
		value = value / megabyte
	case bytes >= kilobyte:
		unit = "K"
		value = value / kilobyte
	case bytes == 0:
		return "0"
	}

	stringValue := fmt.Sprintf("%.1f", value)
	stringValue = strings.TrimSuffix(stringValue, ".0")
	return fmt.Sprintf("%s%s", stringValue, unit)
}

// WriteLine - write line to file
func WriteLine(fullpath, line string) error {
	f, err := os.OpenFile(fullpath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(line + "\n")
	if err != nil {
		return err
	}

	return nil
}

// WriteLines - write multiple lines to file
func WriteLines(fullpath string, lines []string) error {
	f, err := os.OpenFile(fullpath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	for _, line := range lines {
		_, err = f.WriteString(line + "\n")
		if err != nil {
			return err
		}
	}

	return nil
}

// Round value
func Round(d, r time.Duration) time.Duration {
	if r <= 0 {
		return d
	}
	neg := d < 0
	if neg {
		d = -d
	}
	if m := d % r; m+m < r {
		d = d - m
	} else {
		d = d + r - m
	}
	if neg {
		return -d
	}
	return d
}

// Max - between two values
func Max(x, y int64) int64 {
	if x > y {
		return x
	}
	return y
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
			if http == "" {
				// emhttp -p ,<n> variant
				return nil, true, https
			} else {
				// emhttp -p <m>,<n> variant
				return nil, false, http
			}
		}
	}

	// anything else is invalid
	return nil, false, "80"
}

func publicKey(priv interface{}) interface{} {
	switch k := priv.(type) {
	case *rsa.PrivateKey:
		return &k.PublicKey
	case *ecdsa.PrivateKey:
		return &k.PublicKey
	default:
		return nil
	}
}

func pemBlockForKey(priv interface{}) *pem.Block {
	switch k := priv.(type) {
	case *rsa.PrivateKey:
		return &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(k)}
	case *ecdsa.PrivateKey:
		b, err := x509.MarshalECPrivateKey(k)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to marshal ECDSA private key: %v", err)
			os.Exit(2)
		}
		return &pem.Block{Type: "EC PRIVATE KEY", Bytes: b}
	default:
		return nil
	}
}

func GenerateCerts(name, location string) error {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}

	notBefore := time.Now()
	notAfter := notBefore.Add(10 * 365 * 24 * time.Hour) // 10 years

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return err
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Country:            []string{"PA"},
			Organization:       []string{"Apertoire"},
			OrganizationalUnit: []string{"ControlR"},
		},
		NotBefore: notBefore,
		NotAfter:  notAfter,

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	names := []string{"localhost", name, fmt.Sprintf("%s.local", name)}
	for _, n := range names {
		template.DNSNames = append(template.DNSNames, n)
	}

	template.EmailAddresses = []string{fmt.Sprintf("root@%s", name)}

	// template.DNSNames = append(template.DNSNames, "localhost")

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, publicKey(priv), priv)
	if err != nil {
		return err
	}

	certOut, err := os.Create(filepath.Join(location, "cert.pem"))
	if err != nil {
		return err
	}
	pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	certOut.Close()
	// log.Print("written cert.pem\n")

	keyOut, err := os.OpenFile(filepath.Join(location, "key.pem"), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	pem.Encode(keyOut, pemBlockForKey(priv))
	keyOut.Close()
	// log.Print("written key.pem\n")

	return nil
}
