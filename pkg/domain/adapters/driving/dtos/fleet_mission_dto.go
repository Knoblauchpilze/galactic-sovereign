package dtos

import (
	"encoding/json"
	"fmt"
)

type FleetMissionDto string

const (
	MissionColonize FleetMissionDto = "Colonize"
)

func (m FleetMissionDto) MarshalJSON() ([]byte, error) {
	if !m.valid() {
		return nil, fmt.Errorf("invalid fleet mission %q", m)
	}

	return json.Marshal(string(m))
}

func (m *FleetMissionDto) UnmarshalJSON(data []byte) error {
	var value string
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	mission := FleetMissionDto(value)
	if !mission.valid() {
		return fmt.Errorf("invalid fleet mission %q", value)
	}

	*m = mission
	return nil
}

func (m FleetMissionDto) valid() bool {
	return m == MissionColonize
}
