package lib

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	byte     = 1.0
	kilobyte = 1024 * byte
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
