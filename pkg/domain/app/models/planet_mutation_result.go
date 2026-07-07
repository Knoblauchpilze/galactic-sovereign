package models

// MutationResult defines the result of a mutation. This result indicates
// whether the planet was deleted or not and if the planet is not deleted,
// the `Planet` field holds the latest version of the planet's data.
// The field should be ignored in case the `Deleted` boolean is true.
type PlanetMutationResult struct {
	Deleted bool
	Planet  Planet
}
