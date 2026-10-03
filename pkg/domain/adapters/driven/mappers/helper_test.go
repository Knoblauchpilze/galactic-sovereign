package mappers

import (
	"os"
	"testing"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/adapters/driven/database"
	integrationdb "github.com/Knoblauchpilze/galactic-sovereign/pkg/testing/integrationdb"
)

var (
	sharedDbContainer = &integrationdb.Suite{}
)

func TestMain(m *testing.M) {
	code := m.Run()
	sharedDbContainer.Teardown()
	os.Exit(code)
}

func newTestConnection(t *testing.T) database.Connection {
	t.Helper()
	return sharedDbContainer.NewTestConnection(t)
}
