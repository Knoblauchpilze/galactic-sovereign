package usecases

import (
	"context"
	"time"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	domainerrors "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models/errors"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models/request"
	drivenports "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/ports/driven"
	domainservices "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/services"
)

type CreateShipActionUseCase struct {
	shipRepo      drivenports.ForFetchingShip
	planetMutator drivenports.ForMutatingPlanet
	clock         drivenports.ForFetchingTime
}

func NewCreateShipActionUseCase(
	shipRepo drivenports.ForFetchingShip,
	planetMutator drivenports.ForMutatingPlanet,
	clock drivenports.ForFetchingTime,
) *CreateShipActionUseCase {
	return &CreateShipActionUseCase{
		shipRepo:      shipRepo,
		planetMutator: planetMutator,
		clock:         clock,
	}
}

func (b *CreateShipActionUseCase) Create(
	ctx context.Context,
	req request.ShipActionCreationRequest,
) (models.ShipAction, error) {
	moment := b.clock.Now(ctx)

	ship, err := b.shipRepo.Get(ctx, req.Ship)
	if err != nil {
		if err == domainerrors.ErrNotFound {
			return models.ShipAction{}, domainerrors.ErrShipNotFound
		}

		return models.ShipAction{}, err
	}

	mutator := generateShipActionMutator(moment, ship, req.Count)
	result, err := b.planetMutator.Mutate(ctx, req.Planet, mutator)
	if err != nil {
		return models.ShipAction{}, err
	}
	if result.Deleted {
		return models.ShipAction{}, domainerrors.ErrNotFound
	}

	if len(result.Planet.ShipActions) == 0 {
		return models.ShipAction{}, domainerrors.ErrResourceCreationFailed
	}

	last := result.Planet.ShipActions[len(result.Planet.ShipActions)-1]
	return last, nil
}

func generateShipActionMutator(
	moment time.Time,
	ship models.Ship,
	count int,
) drivenports.PlanetMutator {
	return func(p *models.Planet) (bool, error) {
		err := domainservices.AdvancePlanetToTime(p, moment)
		if err != nil {
			return false, err
		}

		return false, p.AddShipAction(ship, count)
	}
}
