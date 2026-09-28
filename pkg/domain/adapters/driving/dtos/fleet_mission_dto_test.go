package dtos

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnit_FleetMissionDto_MarshalJSON(t *testing.T) {
	t.Run("marshals a valid mission", func(t *testing.T) {
		actual, err := json.Marshal(MissionColonize)
		require.NoError(t, err, "Actual err: %v", err)

		assert.JSONEq(t, `"Colonize"`, string(actual))
	})

	t.Run("rejects an invalid mission", func(t *testing.T) {
		_, err := json.Marshal(FleetMissionDto("Attack"))

		assert.Error(t, err)
	})
}

func TestUnit_FleetMissionDto_UnmarshalJSON(t *testing.T) {
	t.Run("unmarshals a valid mission", func(t *testing.T) {
		var actual FleetMissionDto
		err := json.Unmarshal([]byte(`"Colonize"`), &actual)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, MissionColonize, actual)
	})

	t.Run("rejects an invalid mission", func(t *testing.T) {
		var actual FleetMissionDto
		err := json.Unmarshal([]byte(`"Attack"`), &actual)

		assert.Error(t, err)
	})
}
