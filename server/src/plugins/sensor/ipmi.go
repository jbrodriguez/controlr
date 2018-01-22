/*
Ipmi parts based on code from dmacias
https://github.com/dmacias72/unRAID-plugins/blob/master/source/ipmi/usr/local/emhttp/plugins/ipmi/include/ipmi_helpers.php
Check LICENSE file in this folder
*/

package sensor

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"

	"controlr/plugin/server/src/dto"
	"controlr/plugin/server/src/lib"

	"github.com/jbrodriguez/mlog"
	ini "github.com/vaughan0/go-ini"
)

const ipmiBinary = "/usr/sbin/ipmi-sensors"
const ipmiConfig string = "/boot/config/plugins/ipmi/ipmi.cfg"

// IpmiSensor -
type IpmiSensor struct {
	network  bool
	ipaddr   string
	user     string
	password string
}

// NewIpmiSensor -
func NewIpmiSensor() *IpmiSensor {

	var ipaddr, user, password string
	var network bool

	file, err := ini.LoadFile(ipmiConfig)
	if err != nil {
		mlog.Warning("ipmi:unable to parse config:(%s). Using defaults ...", ipmiConfig)
	} else {
		net, ok := file.Get("", "NETWORK")
		if ok {
			net = strings.Replace(net, "\"", "", -1)
		}

		network = net == "enable"

		if network {
			ipaddr, _ = file.Get("", "IPADDR")
			ipaddr = strings.Replace(ipaddr, "\"", "", -1)

			user, _ = file.Get("", "USER")
			user = strings.Replace(user, "\"", "", -1)

			pwd, _ := file.Get("", "PASSWORD")
			pwd = strings.Replace(pwd, "\"", "", -1)

			data, err := base64.StdEncoding.DecodeString(pwd)
			if err != nil {
				mlog.Warning("ipmi:unable to decode pwd:(%s)", err)
				password = ""
			} else {
				password = string(data)
			}
		}
	}

	return &IpmiSensor{
		network:  network,
		ipaddr:   ipaddr,
		user:     user,
		password: password,
	}

}

// GetReadings -
func (s *IpmiSensor) GetReadings(prefs dto.Prefs) []dto.Sample {
	args := make([]string, 0)

	args = append(args, "--comma-separated-output")
	args = append(args, "--output-sensor-state")
	args = append(args, "--no-header-output")
	args = append(args, "--interpret-oem-data")

	if s.network {
		args = append(args, "--always-prefix")
		args = append(args, "-h "+s.ipaddr)
		args = append(args, "-u "+s.user)
		args = append(args, "-p "+s.password)
		args = append(args, "--session-timeout=5000")
		args = append(args, "--retransmission-timeout=1000")
	}

	mlog.Warning("ipmi:args:(%v)", args)

	return make([]dto.Sample, 0)
	// return s.Parse(prefs, lib.GetCmdOutput(ipmiBinary, args...))
}

// Parse -
func (s *IpmiSensor) Parse(prefs dto.Prefs, lines []string) []dto.Sample {
	samples := make([]dto.Sample, 0)

	for _, line := range lines {
		if strings.Contains(line, "CPU Temp") {
			fields := strings.Split(line, ",")

			value := fields[4]

			index := strings.IndexByte(fields[4], '.')
			if index > 0 {
				strVal := fields[4][0:index]
				fVal, _ := strconv.ParseFloat(strVal, 64)

				if prefs.Unit == "F" {
					value = fmt.Sprintf("%d", lib.Round(9/5*fVal+32))
				} else {
					value = fmt.Sprintf("%d", lib.Round(fVal))
				}
			}

			sample := dto.Sample{Key: "CPU", Value: value, Unit: prefs.Unit, Condition: "neutral"}

			samples = append(samples, sample)
		} else if strings.Contains(line, "System Temp") {
			fields := strings.Split(line, ",")

			value := fields[4]

			index := strings.IndexByte(fields[4], '.')
			if index > 0 {
				strVal := fields[4][0:index]
				fVal, _ := strconv.ParseFloat(strVal, 64)

				if prefs.Unit == "F" {
					value = fmt.Sprintf("%d", lib.Round(9/5*fVal+32))
				} else {
					value = fmt.Sprintf("%d", lib.Round(fVal))
				}
			}

			sample := dto.Sample{Key: "BOARD", Value: value, Unit: prefs.Unit, Condition: "neutral"}

			samples = append(samples, sample)
		} else if strings.Contains(line, "FAN1") {
			fields := strings.Split(line, ",")

			value := fields[4]
			index := strings.IndexByte(fields[4], '.')
			if index > 0 {
				value = fields[4][0:index]
			}

			sample := dto.Sample{Key: "FAN", Value: value, Unit: "rpm", Condition: "neutral"}
			samples = append(samples, sample)
		}
	}

	return samples
}
