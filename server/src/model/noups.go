package model

import "jbrodriguez/controlr/plugin/server/src/dto"

type NoUps struct {
	samples []dto.Sample
}

func NewNoUps() *NoUps {
	noups := &NoUps{
		samples: make([]dto.Sample, 0),
	}
	return noups
}

func (n *NoUps) GetStatus() []dto.Sample {
	return n.samples
}
