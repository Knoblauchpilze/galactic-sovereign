package models

import "math/rand"

const (
	homeworldFields    = 163
	minFieldsForBeyond = 60
	maxFieldsForBeyond = 70
)

var (
	// https://board.en.ogame.gameforge.com/index.php?thread/790879-minimum-planet-size/
	minFields = map[int]int{
		0:  95,
		1:  97,
		2:  98,
		3:  123,
		4:  148,
		5:  148,
		6:  141,
		7:  163,
		8:  155,
		9:  151,
		10: 139,
		11: 134,
		12: 109,
		13: 81,
		14: 65,
	}

	maxFields = map[int]int{
		0:  108,
		1:  110,
		2:  139,
		3:  210,
		4:  215,
		5:  239,
		6:  242,
		7:  248,
		8:  243,
		9:  225,
		10: 205,
		11: 180,
		12: 121,
		13: 93,
		14: 74,
	}
)

type Coordinate struct {
	Galaxy      int
	SolarSystem int
	Position    int
}

func (c Coordinate) Fields(homeworld bool) int {
	if homeworld {
		return homeworldFields
	}

	min := fieldOrDefault(c.Position, minFields, minFieldsForBeyond)
	max := fieldOrDefault(c.Position, maxFields, maxFieldsForBeyond)

	return min + rand.Intn(max-min)
}

func fieldOrDefault(position int, table map[int]int, defaultFields int) int {
	value, ok := table[position]
	if ok {
		return value
	}

	return defaultFields
}
