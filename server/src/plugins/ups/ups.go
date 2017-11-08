package ups

import (
	"controlr/plugin/server/src/dto"
	"controlr/plugin/server/src/lib"
)

type UpsKind int

const (
	DNE UpsKind = iota // does not exist
	APC         = iota
	NUT         = iota
)

func IdentifyUps() (UpsKind, error) {
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

type Ups interface {
	GetStatus() []dto.Sample
}
