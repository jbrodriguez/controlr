/* Copyright 2015, Bergware International.
 *
 * This program is free software; you can redistribute it and/or
 * modify it under the terms of the GNU General Public License version 2,
 * as published by the Free Software Foundation.
 *
 * The above copyright notice and this permission notice shall be included in
 * all copies or substantial portions of the Software.
 *
 * Plugin development contribution by gfjardim
 */

package sensor

import (
	"plugin/dto"
	"plugin/lib"
	"fmt"
	"strconv"
	"strings"
)

const sensorBinary string = "/usr/bin/sensors"

// SystemSensor -
type SystemSensor struct {
}

// NewSystemSensor -
func NewSystemSensor() *SystemSensor {
	return &SystemSensor{}
}

// GetReadings -
func (s *SystemSensor) GetReadings(prefs dto.Prefs) []dto.Sample {
	return s.Parse(prefs, lib.GetCmdOutput(sensorBinary, "-A"))
}

// Parse -
func (s *SystemSensor) Parse(prefs dto.Prefs, lines []string) []dto.Sample {
	samples := make([]dto.Sample, 0)

	for _, line := range lines {
		if strings.HasPrefix(line, "CPU Temp") {
			fields := strings.Fields(line)

			value := fields[2]

			index := strings.IndexByte(fields[2], '°')
			if index > 0 {
				strVal := fields[2][1 : index-1]
				fVal, _ := strconv.ParseFloat(strVal, 64)

				if prefs.Unit == "F" {
					value = fmt.Sprintf("%d", lib.Round(9/5*fVal+32))
				} else {
					value = fmt.Sprintf("%d", lib.Round(fVal))
				}
			}

			sample := dto.Sample{Key: "CPU", Value: value, Unit: prefs.Unit, Condition: "neutral"}

			samples = append(samples, sample)
		} else if strings.HasPrefix(line, "MB Temp") {
			fields := strings.Fields(line)

			value := fields[2]

			index := strings.IndexByte(fields[2], '°')
			if index > 0 {
				strVal := fields[2][1 : index-1]
				fVal, _ := strconv.ParseFloat(strVal, 64)

				if prefs.Unit == "F" {
					value = fmt.Sprintf("%d", lib.Round(9/5*fVal+32))
				} else {
					value = fmt.Sprintf("%d", lib.Round(fVal))
				}
			}

			sample := dto.Sample{Key: "BOARD", Value: value, Unit: prefs.Unit, Condition: "neutral"}

			samples = append(samples, sample)
		} else if strings.HasPrefix(line, "Array Fan") {
			fields := strings.Fields(line)
			sample := dto.Sample{Key: "FAN", Value: fields[2], Unit: "rpm", Condition: "neutral"}
			samples = append(samples, sample)
		}
	}

	return samples
}
