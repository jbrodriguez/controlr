package ups

import (
	"controlr/plugin/server/src/dto"
	"controlr/plugin/server/src/lib"
)

// Kind -
type Kind int

// DNE -
const (
	DNE Kind = iota // does not exist
	APC
	NUT
)

// Green -
const (
	Green  = "green"
	Red    = "red"
	Orange = "orange"
)

// IdentifyUps -
func IdentifyUps() (Kind, error) {
	exists, err := lib.Exists("/var/run/nut/upsmon.pid")
	if err != nil {
		return DNE, err
	}

	if exists {
		return NUT, nil
	}

	exists, err = lib.Exists("/var/run/apcupsd.pid")
	if err != nil {
		return DNE, err
	}

	if exists {
		return APC, nil
	}

	return DNE, nil
}

// Ups -
type Ups interface {
	GetStatus() []dto.Sample
}
