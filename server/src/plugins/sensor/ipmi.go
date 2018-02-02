/*
Ipmi parts based on code from dmacias
https://github.com/dmacias72/unRAID-plugins/blob/master/source/ipmi/usr/local/emhttp/plugins/ipmi/include/ipmi_helpers.php
Check LICENSE file in this folder
*/

package sensor

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"controlr/plugin/server/src/dto"
	"controlr/plugin/server/src/lib"

	"github.com/jbrodriguez/mlog"
	ini "github.com/vaughan0/go-ini"
)

const ipmiBinary = "/usr/sbin/ipmisensors"
const ipmiConfig = "/boot/config/plugins/ipmi/ipmi.cfg"

// IpmiSensor -
type IpmiSensor struct {
	sensors map[string]bool
}

// NewIpmiSensor -
func NewIpmiSensor() *IpmiSensor {
	lines, err := getConfig(ipmiConfig)
	if err != nil {
		mlog.Warning("ipmi:unable to parse config:(%s):(%s). No sensors available", ipmiConfig, err)
	}

	sensors := GetSensorIDs(lines)

	return &IpmiSensor{
		sensors: sensors,
	}
}

// GetReadings -
func (s *IpmiSensor) GetReadings(prefs dto.Prefs) []dto.Sample {
	args := []string{
		"--comma-separated-output",
		"--output-sensor-state",
		"--no-header-output",
		"--interpret-oem-data",
	}

	return s.Parse(prefs, s.sensors, lib.GetCmdOutput(ipmiBinary, args...))
}

// Parse -
func (s *IpmiSensor) Parse(prefs dto.Prefs, sensors map[string]bool, lines []string) []dto.Sample {
	samples := make([]dto.Sample, 0)

	for _, line := range lines {
		fields := strings.Split(line, ",")

		if _, ok := sensors[fields[0]]; !ok {
			continue
		}

		value := fields[4]
		unit := fields[5]

		// if temperature or fan, remove the precision (35.00 -> 35, 1250.00 -> 1250)
		if fields[2] == "Temperature" || fields[2] == "Fan" {
			index := strings.IndexByte(value, '.')
			if index > 0 {
				value = value[0:index]
			}
		}

		// if prefs is "F" and temperature, additionally convert appropriately
		if prefs.Unit == "F" && fields[2] == "Temperature" {
			fVal, _ := strconv.ParseFloat(value, 64)
			value = fmt.Sprintf("%d", lib.Round(9/5*fVal+32)) // probably an int(calculation) should suffice
			unit = "F"
		}

		sample := dto.Sample{Key: fields[1], Value: value, Unit: unit, Condition: "neutral"}
		samples = append(samples, sample)
	}

	return samples
}

// CheckIpmiPresence -
func CheckIpmiPresence() (bool, error) {
	exists, err := lib.Exists("/dev/ipmi0")
	if err != nil {
		return false, err
	}

	if exists {
		return true, nil
	}

	exists, err = lib.Exists("/dev/ipmi/0")
	if err != nil {
		return false, err
	}

	if exists {
		return true, nil
	}

	exists, err = lib.Exists("/dev/ipmidev/0")
	if err != nil {
		return false, err
	}

	if exists {
		return true, nil
	}

	return false, nil
}

// CheckNetworkEnabled -
func CheckNetworkEnabled() (bool, error) {
	file, err := ini.LoadFile(ipmiConfig)
	if err != nil {
		return false, err
	}

	net, _ := file.Get("", "NETWORK")
	net = strings.Replace(net, "\"", "", -1)

	return net == "enable", nil
}

// GetSensorIDs -
func GetSensorIDs(lines []string) map[string]bool {
	sensors := make(map[string]bool)

	for _, line := range lines {
		if strings.HasPrefix(line, "DISP_SENSOR") {
			fields := strings.Split(line, "=")           // get the value part
			value := fields[1][1 : len(fields[1])-1]     // remove those pesky quotes
			if n := strings.Index(value, "_"); n != -1 { // if dash is present remove the prefix
				value = value[n+1:]
			}
			sensors[value] = true
		}
	}

	return sensors
}

func getConfig(config string) ([]string, error) {
	b, err := ioutil.ReadFile(config)
	if err != nil {
		return make([]string, 0), err
	}

	return strings.Split(string(b), "\n"), nil
}
