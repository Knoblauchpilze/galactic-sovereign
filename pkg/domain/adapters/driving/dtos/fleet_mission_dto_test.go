package dtos

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnit_FleetMissionDto_MarshalJSON(t *testing.T) {
	t.Run("marshals colonize mission", func(t *testing.T) {
		actual, err := json.Marshal(MissionColonize)
		require.NoError(t, err, "Actual err: %v", err)

		assert.JSONEq(t, `"colonize"`, string(actual))
	})

	t.Run("marshals transport mission", func(t *testing.T) {
		actual, err := json.Marshal(MissionTransport)
		require.NoError(t, err, "Actual err: %v", err)

		assert.JSONEq(t, `"transport"`, string(actual))
	})

	t.Run("rejects an invalid mission", func(t *testing.T) {
		_, err := json.Marshal(FleetMissionDto("attack"))

		assert.Error(t, err)
	})
}

func TestUnit_FleetMissionDto_UnmarshalJSON(t *testing.T) {
	t.Run("unmarshals colonize mission", func(t *testing.T) {
		var actual FleetMissionDto
		err := json.Unmarshal([]byte(`"colonize"`), &actual)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, MissionColonize, actual)
	})

	t.Run("unmarshals transport mission", func(t *testing.T) {
		var actual FleetMissionDto
		err := json.Unmarshal([]byte(`"transport"`), &actual)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, MissionTransport, actual)
	})

	t.Run("rejects an invalid mission", func(t *testing.T) {
		var actual FleetMissionDto
		err := json.Unmarshal([]byte(`"attack"`), &actual)

		assert.Error(t, err)
	})
}
