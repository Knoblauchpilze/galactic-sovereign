package models

import "math"

type BuildingShipSpeedup struct {
	Scaling  ScalingMode
	Base     float64
	Progress float64
}

func (s *BuildingShipSpeedup) FactorAt(level int) float64 {
	switch s.Scaling {
	case LinearScaling:
		return s.Base + s.Progress*float64(level)
	case GeometricScaling:
		return s.Base + math.Pow(s.Progress, float64(level))
	}

	// Voluntarily ignore unknown values and return 1.
	// Adding a scaling mode is a big enough operation that it
	// should not be rushed and will surface the update to do
	// here.
	return 1.0
}
