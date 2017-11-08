package sensor

import (
	"controlr/plugin/server/src/dto"
	"controlr/plugin/server/src/lib"
)

type SensorKind int

const (
	NOSENSOR SensorKind = iota // does not exist
	SYSTEM              = iota // dynamix.system.temp
)

func IdentifySensor() (SensorKind, error) {
	exists, err := lib.Exists("/usr/local/emhttp/plugins/dynamix.system.temp")
	if err != nil {
		return NOSENSOR, err
	}

	if exists {
		return SYSTEM, nil
	}

	return NOSENSOR, nil
}

type Sensor interface {
	GetReadings(prefs dto.Prefs) []dto.Sample
}
