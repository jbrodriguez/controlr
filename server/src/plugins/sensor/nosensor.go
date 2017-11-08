package sensor

import "controlr/plugin/server/src/dto"

type NoSensor struct {
	samples []dto.Sample
}

func NewNoSensor() *NoSensor {
	nosensor := &NoSensor{
		samples: make([]dto.Sample, 0),
	}
	return nosensor
}

func (n *NoSensor) GetReadings(prefs dto.Prefs) []dto.Sample {
	return n.samples
}