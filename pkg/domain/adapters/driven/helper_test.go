package drivenadapters

import (
	"math"
	"math/rand/v2"
	"os"
	"testing"
	"time"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/db"
	eassert "github.com/Knoblauchpilze/easy-assert/assert"
	integrationdb "github.com/Knoblauchpilze/galactic-sovereign/pkg/testing/integrationdb"
	"github.com/stretchr/testify/assert"
)

var (
	someTime      = time.Date(2024, time.November, 29, 17, 53, 29, 0, time.UTC)
	someOtherTime = time.Date(2026, time.June, 1, 8, 20, 15, 0, time.UTC)

	sharedDbContainer = &integrationdb.Suite{}
)

func TestMain(m *testing.M) {
	code := m.Run()
	sharedDbContainer.Teardown()
	os.Exit(code)
}

func newTestConnection(t *testing.T) db.Connection {
	t.Helper()
	return sharedDbContainer.NewTestConnection(t)
}

func randFloat(t *testing.T, min float64, max float64, precision int) float64 {
	t.Helper()

	scale := math.Pow(10, float64(precision))
	minScaled := int64(math.Ceil(min * scale))
	maxScaled := int64(math.Floor(max * scale))
	if minScaled > maxScaled {
		t.Fatalf(
			"no representable values in range [%f, %f] for precision=%d",
			min,
			max,
			precision,
		)
	}

	valueScaled := minScaled + rand.Int64N(maxScaled-minScaled+1)
	return float64(valueScaled) / scale
}

func assertEqualIgnoringFields[T any](
	t *testing.T,
	actual T,
	expected T,
	ignoredFields ...string,
) {
	t.Helper()

	equal := eassert.EqualsIgnoringFields(actual, expected, ignoredFields...)
	assert.True(t, equal, "Expected actual=%+v and expected=%+v to be equal", actual, expected)
}

func assertContainsIgnoringFields[T any](
	t *testing.T,
	collection []T,
	expected T,
	ignoredFields ...string,
) {
	t.Helper()

	equal := eassert.ContainsIgnoringFields(collection, expected, ignoredFields...)
	assert.True(t, equal, "Expected collection=%+v to contain expected=%+v", collection, expected)
}
