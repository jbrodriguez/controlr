/* Copyright 2017, Juan B. Rodriguez
 * Copyright 2017, Derek Macias.
 * Copyright 2005-2016, Lime Technology
 * Copyright 2015, Dan Landon.
 * Copyright 2015, Bergware International.
 *
 * This program is free software; you can redistribute it and/or
 * modify it under the terms of the GNU General Public License version 2,
 * as published by the Free Software Foundation.
 *
 * The above copyright notice and this permission notice shall be included in
 * all copies or substantial portions of the Software.
 */
package model

import (
	"fmt"
	"jbrodriguez/controlr/plugin/server/src/dto"
	"strconv"
	"strings"

	"github.com/jbrodriguez/mlog"
	ini "github.com/vaughan0/go-ini"
)

const nutBinary string = "/usr/bin/upsc"
const nutConfig string = "/boot/config/plugins/nut/nut.cfg"

// nut
type Nut struct {
	legend  map[string]string
	address string
}

func NewNut() *Nut {
	nut := &Nut{}

	nut.legend = map[string]string{
		"OL":    "Online",
		"OB":    "On batt",
		"OL LB": "Online low batt",
		"OB LB": "Low batt",
	}

	var name, ip string
	var ok bool

	file, err := ini.LoadFile(nutConfig)
	if err != nil {
		name = "ups"
		ip = "127.0.0.1"
		mlog.Warning("nuts: unable to parse config: (%s). Using defaults ...", nutConfig)
	} else {
		name, ok = file.Get("", "NAME")
		if !ok {
			name = "ups"
			mlog.Warning("nuts: name not found. Defaulting to ups ...")
		}

		ip, ok = file.Get("", "IPADDR")
		if !ok {
			ip = "127.0.0.1"
			mlog.Warning("nuts: ip not found. Defaulting to 127.0.0.1 ...")
		}
	}

	nut.address = fmt.Sprintf("%s@%s", name, ip)

	return nut
}

func (n *Nut) GetStatus() []dto.Sample {
	return n.Parse(getData(nutBinary, n.address))
}

func (n *Nut) Parse(lines []string) []dto.Sample {
	samples := make([]dto.Sample, 0)

	var power float64
	var load float64

	for _, line := range lines {
		tokens := strings.Split(line, ":")

		key := strings.Trim(tokens[0], " ")
		val := strings.Trim(tokens[1], " ")

		switch key {
		case "ups.status":
			value, ok := n.legend[val]

			sample := dto.Sample{Key: "UPS STATUS"}

			if ok {
				sample.Value = value
				if strings.Index(strings.ToLower(value), "online") == 0 {
					sample.Condition = "green"
				} else {
					sample.Condition = "red"
				}
			} else {
				sample.Value = "Refreshing ..."
				sample.Condition = "orange"
			}

			samples = append(samples, sample)

			break

		case "battery.charge":
			sample := dto.Sample{Key: "UPS CHARGE", Value: val, Unit: "%"}

			charge, _ := strconv.ParseFloat(val, 64)
			if charge <= 10 {
				sample.Condition = "red"
			} else {
				sample.Condition = "green"
			}

			samples = append(samples, sample)
			break

		case "battery.runtime":
			left, _ := strconv.ParseFloat(val, 64)
			runtime := strconv.FormatFloat(left/60, 'f', 1, 64)

			sample := dto.Sample{Key: "UPS LEFT", Value: runtime, Unit: "m"}

			if left <= 5 {
				sample.Condition = "red"
			} else {
				sample.Condition = "green"
			}

			samples = append(samples, sample)
			break

		case "ups.realpower.nominal":
			power, _ = strconv.ParseFloat(val, 64)
			break

		case "ups.load":
			sample := dto.Sample{Key: "UPS LOAD", Value: val, Unit: "%"}

			load, _ = strconv.ParseFloat(val, 64)
			if load >= 90 {
				sample.Condition = "red"
			} else {
				sample.Condition = "green"
			}

			samples = append(samples, sample)
			break
		}
	}

	if power != 0 && load != 0 {
		value := power * load / 100
		watts := strconv.FormatFloat(value, 'f', 1, 64)

		sample := dto.Sample{Key: "UPS POWER", Value: watts, Unit: "w"}

		if load >= 90 {
			sample.Condition = "red"
		} else {
			sample.Condition = "green"
		}

		samples = append(samples, sample)
	}

	return samples
}
