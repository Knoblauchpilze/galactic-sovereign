package request

import (
	"testing"
	"time"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUnit_FromBuildingActionCreationRequest(t *testing.T) {
	beforeConversion := time.Now()

	request := BuildingActionCreationRequest{
		Planet:   uuid.New(),
		Building: uuid.New(),
	}

	actual := FromBuildingActionCreationRequest(request)

	assert.Equal(t, request.Planet, actual.Planet)
	assert.Equal(t, request.Building, actual.Building)
	assert.Zero(t, actual.CurrentLevel)
	assert.Zero(t, actual.DesiredLevel)
	assert.True(t, actual.CreatedAt.After(beforeConversion))
	assert.Equal(t, time.Time{}, actual.CompletedAt)
	assert.Zero(t, actual.Version)
	assert.Equal(t, []models.BuildingActionCost{}, actual.Costs)
	assert.Equal(t, []models.BuildingActionResourceStorage{}, actual.Storages)
	assert.Equal(t, []models.BuildingActionResourceProduction{}, actual.Productions)
}
