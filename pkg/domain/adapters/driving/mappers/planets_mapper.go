package mappers

import (
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/adapters/driving/dtos"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
)

func ToPlanetResponse(planet models.Planet) dtos.PlanetDtoResponse {
	dto := dtos.PlanetDtoResponse{
		Id:        planet.Id,
		Player:    planet.Player,
		Name:      planet.Name,
		Homeworld: planet.Homeworld,
		Coordinate: dtos.CoordinateDtoResponse{
			Galaxy:      planet.Coordinate.Galaxy,
			SolarSystem: planet.Coordinate.SolarSystem,
			Position:    planet.Coordinate.Position,
		},
		CreatedAt:   planet.CreatedAt,
		UpdatedAt:   planet.UpdatedAt,
		Resources:   toPlanetResourcesResponse(planet.Resources),
		Storages:    toPlanetStoragesResponse(planet.Storages),
		Productions: toPlanetProductionsResponse(planet.Productions),
		Buildings:   toPlanetBuildingsResponse(planet.Buildings),
	}

	if planet.BuildingAction != nil {
		action := ToBuildingActionResponse(*planet.BuildingAction)
		dto.BuildingAction = &action
	}

	return dto
}

func toPlanetResourceResponse(
	resource models.PlanetResource,
) dtos.PlanetResourceDtoResponse {
	return dtos.PlanetResourceDtoResponse{
		Resource: resource.Resource,
		Amount:   resource.Amount,
	}
}

func toPlanetResourcesResponse(
	resources []models.PlanetResource,
) []dtos.PlanetResourceDtoResponse {
	out := make([]dtos.PlanetResourceDtoResponse, 0, len(resources))

	for _, r := range resources {
		dto := toPlanetResourceResponse(r)
		out = append(out, dto)
	}

	return out
}

func toPlanetStorageResponse(
	storage models.PlanetResourceStorage,
) dtos.PlanetResourceStorageDtoResponse {
	return dtos.PlanetResourceStorageDtoResponse{
		Resource: storage.Resource,
		Storage:  storage.Storage,
	}
}

func toPlanetStoragesResponse(
	storages []models.PlanetResourceStorage,
) []dtos.PlanetResourceStorageDtoResponse {
	out := make([]dtos.PlanetResourceStorageDtoResponse, 0, len(storages))

	for _, s := range storages {
		dto := toPlanetStorageResponse(s)
		out = append(out, dto)
	}

	return out
}

func toPlanetProductionResponse(
	production models.PlanetResourceProduction,
) dtos.PlanetResourceProductionDtoResponse {
	return dtos.PlanetResourceProductionDtoResponse{
		Building:   production.Building,
		Resource:   production.Resource,
		Production: production.Production,
	}
}

func toPlanetProductionsResponse(
	productions []models.PlanetResourceProduction,
) []dtos.PlanetResourceProductionDtoResponse {
	out := make([]dtos.PlanetResourceProductionDtoResponse, 0, len(productions))

	for _, p := range productions {
		dto := toPlanetProductionResponse(p)
		out = append(out, dto)
	}

	return out
}

func toPlanetBuildingResponse(
	building models.PlanetBuilding,
) dtos.PlanetBuildingDtoResponse {
	return dtos.PlanetBuildingDtoResponse{
		Building: building.Building,
		Level:    building.Level,
	}
}

func toPlanetBuildingsResponse(
	buildings []models.PlanetBuilding,
) []dtos.PlanetBuildingDtoResponse {
	out := make([]dtos.PlanetBuildingDtoResponse, 0, len(buildings))

	for _, b := range buildings {
		dto := toPlanetBuildingResponse(b)
		out = append(out, dto)
	}

	return out
}

func ToPlanetsResponse(planets []models.Planet) []dtos.PlanetDtoResponse {
	out := make([]dtos.PlanetDtoResponse, 0, len(planets))

	for _, p := range planets {
		dto := ToPlanetResponse(p)
		out = append(out, dto)
	}

	return out
}
