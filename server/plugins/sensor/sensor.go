/*
Ipmi parts based on code from dmacias
https://github.com/dmacias72/unRAID-plugins/blob/master/source/ipmi/usr/local/emhttp/plugins/ipmi/include/ipmi_helpers.php
Check LICENSE file in this folder
*/

package sensor

import (
	"plugin/dto"
	"plugin/lib"
)

// Kind -
type Kind int

// NOSENSOR -
const (
	NOSENSOR Kind = iota // does not exist
	SYSTEM               // dynamix.system.temp
	IPMI                 // ipmi
)

// IdentifySensor -
func IdentifySensor() (Kind, error) {
	ipmi, err := CheckIpmiPresence()
	if err != nil {
		return NOSENSOR, err
	}

	exists, err := lib.Exists("/usr/local/emhttp/plugins/ipmi")
	if err != nil {
		return NOSENSOR, err
	}

	network, err := CheckNetworkEnabled()
	if err != nil {
		return NOSENSOR, err
	}

	if exists && (ipmi || network) {
		return IPMI, nil
	}

	exists, err = lib.Exists("/usr/local/emhttp/plugins/dynamix.system.temp")
	if err != nil {
		return NOSENSOR, err
	}

	if exists {
		return SYSTEM, nil
	}

	return NOSENSOR, nil
}

// Sensor -
type Sensor interface {
	GetReadings(prefs dto.Prefs) []dto.Sample
}
