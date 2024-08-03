package service

import (
	"context"
	"testing"

	"github.com/KnoblauchPilze/user-service/pkg/communication"
	"github.com/KnoblauchPilze/user-service/pkg/db"
	"github.com/KnoblauchPilze/user-service/pkg/persistence"
	"github.com/KnoblauchPilze/user-service/pkg/repositories"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

var defaultResourceId = uuid.MustParse("3e0aaf91-8b81-403f-b967-8bdba748594d")
var defaultResourceName = "my-resource"

var defaultResource = persistence.Resource{
	Id:   defaultResourceId,
	Name: defaultResourceName,

	CreatedAt: testDate,
	UpdatedAt: testDate,
}

func Test_ResourceService(t *testing.T) {
	s := ServiceTestSuite{
		generateRepositoriesMock:      generateValidResourceRepositoryMock,
		generateErrorRepositoriesMock: generateErrorResourceRepositoryMock,

		repositoryInteractionTestCases: map[string]repositoryInteractionTestCase{
			"list": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewResourceService(pool, repos)
					_, err := s.List(ctx)
					return err
				},

				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertResourceRepoIsAMock(repos, assert)

					assert.Equal(1, m.listCalled)
				},
			},
			"list_repositoryFails": {
				generateRepositoriesMock: generateErrorResourceRepositoryMock,
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewResourceService(pool, repos)
					_, err := s.List(ctx)
					return err
				},
				expectedError: errDefault,
			},
		},

		returnTestCases: map[string]returnTestCase{
			"list": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) interface{} {
					s := NewResourceService(pool, repos)
					out, _ := s.List(ctx)
					return out
				},

				expectedContent: []communication.ResourceDtoResponse{
					{
						Id:   defaultResource.Id,
						Name: defaultResource.Name,

						CreatedAt: defaultResource.CreatedAt,
					},
				},
			},
		},
	}

	suite.Run(t, &s)
}

func generateValidResourceRepositoryMock() repositories.Repositories {
	return repositories.Repositories{
		Resource: &mockResourceRepository{
			resource: defaultResource,
		},
	}
}

func generateErrorResourceRepositoryMock() repositories.Repositories {
	return repositories.Repositories{
		Resource: &mockResourceRepository{
			err: errDefault,
		},
	}
}

func assertResourceRepoIsAMock(repos repositories.Repositories, assert *require.Assertions) *mockResourceRepository {
	m, ok := repos.Resource.(*mockResourceRepository)
	if !ok {
		assert.Fail("Provided resource repository is not a mock")
	}
	return m
}
