package drivenadapters

import (
	"context"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/db"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/adapters/driven/database"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/adapters/driven/mappers"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	"github.com/google/uuid"
)

const (
	getSolarSystemQuery = `
SELECT
	ut.universe,
	$2::integer AS galaxy,
	$3::integer AS number,
	ut.orbits
FROM
	universe_topology AS ut
WHERE
	ut.universe = $1
	AND $2 >= 0
	AND $2 < ut.galaxies
	AND $3 >= 0
	AND $3 < ut.solar_systems
	`

	listSolarSystemPlanetsQuery = `
SELECT
	p.id,
	pl.name AS player_name,
	p.name,
	CASE
		WHEN h.planet IS NOT NULL THEN true
		ELSE false
	END AS homeworld,
	pc.position
FROM
	planet_coordinate AS pc
	INNER JOIN planet AS p ON pc.planet = p.id
	LEFT JOIN homeworld AS h ON h.planet = p.id
	INNER JOIN player AS pl ON pl.id = p.player
WHERE
	pc.universe = $1
	AND pc.galaxy = $2
	AND pc.solar_system = $3
ORDER BY
	pc.position ASC`
)

type SolarSystemRepository struct {
	conn database.Connection
}

func NewSolarSystemRepository(conn database.Connection) *SolarSystemRepository {
	return &SolarSystemRepository{
		conn: conn,
	}
}

// TODO: Add a test to verify that planet are ordered by position
func (r *SolarSystemRepository) GetSolarSystem(
	ctx context.Context,
	universe uuid.UUID,
	galaxy int,
	solarSystem int,
) (models.SolarSystem, error) {
	tx, err := r.conn.BeginTx(ctx)
	if err != nil {
		return models.SolarSystem{}, err
	}
	defer tx.Close(ctx)

	dbSolarSystem, err := db.QueryOneTx[mappers.DbSolarSystem](
		ctx,
		tx,
		getSolarSystemQuery,
		universe,
		galaxy,
		solarSystem,
	)
	if err != nil {
		return models.SolarSystem{}, parseDbError(err)
	}

	return loadSolarSystemDetails(ctx, tx, dbSolarSystem)
}

func loadSolarSystemDetails(
	ctx context.Context,
	tx database.Transaction,
	dbSolarSystem mappers.DbSolarSystem,
) (models.SolarSystem, error) {
	solarSystem := dbSolarSystem.ToDomain()

	var err error
	solarSystem.Planets, err = db.QueryAllTx[models.SolarSystemPlanet](
		ctx,
		tx,
		listSolarSystemPlanetsQuery,
		solarSystem.Universe,
		solarSystem.Galaxy,
		solarSystem.Number,
	)
	if err != nil {
		return solarSystem, err
	}

	return solarSystem, nil
}
