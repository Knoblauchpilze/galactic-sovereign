package communication

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/KnoblauchPilze/user-service/pkg/persistence"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var defaultActionId = uuid.MustParse("91336067-9884-4280-bb37-411124561e73")
var defaultBuildignAction = persistence.BuildingAction{
	Id:           defaultActionId,
	Planet:       defaultPlanetId,
	Building:     defaultBuildingId,
	CurrentLevel: 37,
	DesiredLevel: 38,
	CreatedAt:    someTime,
	CompletedAt:  someOtherTime,
}
var defaultBuildingActionDtoResponse = BuildingActionDtoResponse{
	Id:           defaultActionId,
	Planet:       defaultPlanetId,
	Building:     defaultBuildingId,
	CurrentLevel: 37,
	DesiredLevel: 38,
	CreatedAt:    someTime,
	CompletedAt:  someOtherTime,
}

func TestBuildingActionDtoRequest_MarshalsToCamelCase(t *testing.T) {
	assert := assert.New(t)

	dto := BuildingActionDtoRequest{
		Planet:   defaultPlanetId,
		Building: defaultBuildingId,
	}

	out, err := json.Marshal(dto)

	assert.Nil(err)
	expectedJson := `
	{
		"planet": "65801b9b-84e6-411d-805f-2eb89587c5a7",
		"building": "461ba465-86e6-4234-94b8-fc8fab03fa74"
	}`
	assert.JSONEq(expectedJson, string(out))
}

func TestFromBuildingActionDtoRequest(t *testing.T) {
	assert := assert.New(t)

	beforeConversion := time.Now()

	dto := BuildingActionDtoRequest{
		Planet:   defaultPlanetId,
		Building: defaultBuildingId,
	}

	actual := FromBuildingActionDtoRequest(dto)

	assert.Nil(uuid.Validate(actual.Id.String()))
	assert.Equal(defaultPlanetId, actual.Planet)
	assert.Equal(defaultBuildingId, actual.Building)
	assert.Equal(0, actual.CurrentLevel)
	assert.Equal(0, actual.DesiredLevel)
	assert.True(actual.CreatedAt.After(beforeConversion))
	assert.True(actual.CompletedAt.IsZero())
}

func TestToBuildingActionDtoResponse(t *testing.T) {
	assert := assert.New(t)

	actual := ToBuildingActionDtoResponse(defaultBuildignAction)

	assert.Equal(defaultActionId, actual.Id)
	assert.Equal(defaultPlanetId, actual.Planet)
	assert.Equal(defaultBuildingId, actual.Building)
	assert.Equal(37, actual.CurrentLevel)
	assert.Equal(38, actual.DesiredLevel)
	assert.Equal(someTime, actual.CreatedAt)
	assert.Equal(someOtherTime, actual.CompletedAt)
}

func TestBuildingActionDtoResponse_MarshalsToCamelCase(t *testing.T) {
	assert := assert.New(t)

	out, err := json.Marshal(defaultBuildingActionDtoResponse)

	assert.Nil(err)
	expectedJson := `
	{
		"id": "91336067-9884-4280-bb37-411124561e73",
		"planet": "65801b9b-84e6-411d-805f-2eb89587c5a7",
		"building": "461ba465-86e6-4234-94b8-fc8fab03fa74",
		"currentLevel": 37,
		"desiredLevel": 38,
		"createdAt": "2024-05-05T20:50:18.651387237Z",
		"completedAt": "2024-07-28T10:30:02.651387236Z"
	}`
	assert.JSONEq(expectedJson, string(out))
}
