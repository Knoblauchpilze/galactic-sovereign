package controllers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/KnoblauchPilze/user-service/pkg/communication"
	"github.com/KnoblauchPilze/user-service/pkg/db"
	"github.com/KnoblauchPilze/user-service/pkg/errors"
	"github.com/KnoblauchPilze/user-service/pkg/persistence"
	"github.com/KnoblauchPilze/user-service/pkg/repositories"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

type mockContext struct {
	echo.Context

	params map[string]string

	status int
	data   interface{}
}

type mockUserRepository struct {
	repositories.UserRepository

	user persistence.User
	err  error

	createCalled int
	getCalled    int
	updateCalled int
	deleteCalled int
}

var defaultUuid = uuid.MustParse("08ce96a3-3430-48a8-a3b2-b1c987a207ca")
var defaultUser = persistence.User{
	Id:        defaultUuid,
	Email:     "e.mail@domain.com",
	Password:  "password",
	CreatedAt: time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC),
	UpdatedAt: time.Date(2009, 11, 17, 20, 34, 59, 651387237, time.UTC),
}
var defaultUserDto = communication.UserDtoResponse{
	Id:       defaultUuid,
	Email:    "e.mail@domain.com",
	Password: "password",

	CreatedAt: time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC),
}

func TestGetUser_WhenNoId_SetsStatusToBadRequest(t *testing.T) {
	assert := assert.New(t)

	mc := &mockContext{}
	mr := &mockUserRepository{}

	err := getUser(mc, mr)

	assert.Nil(err)
	assert.Equal(http.StatusBadRequest, mc.status)
	assert.Equal("Invalid id syntax", mc.data)
}

func TestGetUser_WhenIdSyntaxIsWrong_SetsStatusToBadRequest(t *testing.T) {
	assert := assert.New(t)

	mc := &mockContext{
		params: map[string]string{
			"id": "not-a-valid-id",
		},
	}
	mr := &mockUserRepository{}

	err := getUser(mc, mr)

	assert.Nil(err)
	assert.Equal(http.StatusBadRequest, mc.status)
	assert.Equal("Invalid id syntax", mc.data)
}

func TestGetUser_CallsRepositoryGet(t *testing.T) {
	assert := assert.New(t)

	mc := &mockContext{
		params: map[string]string{
			"id": defaultUuid.String(),
		},
	}
	mr := &mockUserRepository{}

	getUser(mc, mr)

	assert.Equal(1, mr.getCalled)
}

func TestGetUser_WhenRepositorySucceeds_SetsStatusToOk(t *testing.T) {
	assert := assert.New(t)

	mc := &mockContext{
		params: map[string]string{
			"id": defaultUuid.String(),
		},
	}
	mr := &mockUserRepository{}

	err := getUser(mc, mr)

	assert.Nil(err)
	assert.Equal(http.StatusOK, mc.status)
}

func TestGetUser_WhenRepositorySucceeds_ReturnsExpectedUser(t *testing.T) {
	assert := assert.New(t)

	mc := &mockContext{
		params: map[string]string{
			"id": defaultUuid.String(),
		},
	}
	mr := &mockUserRepository{
		user: defaultUser,
	}

	getUser(mc, mr)

	assert.Equal(defaultUserDto, mc.data)
}

func TestGetUser_WhenRepositoryFailsWithUnknownError_SetsStatusToInternalServerError(t *testing.T) {
	assert := assert.New(t)

	mc := &mockContext{
		params: map[string]string{
			"id": defaultUuid.String(),
		},
	}
	mr := &mockUserRepository{
		err: errDefault,
	}

	err := getUser(mc, mr)

	assert.Nil(err)
	assert.Equal(http.StatusInternalServerError, mc.status)
}

func TestGetUser_WhenRepositoryFailsWithNoMatchingRows_SetsStatusToNotFound(t *testing.T) {
	assert := assert.New(t)

	mc := &mockContext{
		params: map[string]string{
			"id": defaultUuid.String(),
		},
	}
	mr := &mockUserRepository{
		err: errors.NewCode(db.NoMatchingSqlRows),
	}

	err := getUser(mc, mr)

	assert.Nil(err)
	assert.Equal(http.StatusNotFound, mc.status)
}

func TestDeleteUser_WhenNoId_SetsStatusToBadRequest(t *testing.T) {
	assert := assert.New(t)

	mc := &mockContext{}
	mr := &mockUserRepository{}

	err := deleteUser(mc, mr)

	assert.Nil(err)
	assert.Equal(http.StatusBadRequest, mc.status)
	assert.Equal("Invalid id syntax", mc.data)
}

func TestDeleteUser_WhenIdSyntaxIsWrong_SetsStatusToBadRequest(t *testing.T) {
	assert := assert.New(t)

	mc := &mockContext{
		params: map[string]string{
			"id": "not-a-valid-id",
		},
	}
	mr := &mockUserRepository{}

	err := deleteUser(mc, mr)

	assert.Nil(err)
	assert.Equal(http.StatusBadRequest, mc.status)
	assert.Equal("Invalid id syntax", mc.data)
}

func TestDeleteUser_CallsRepositoryDelete(t *testing.T) {
	assert := assert.New(t)

	mc := &mockContext{
		params: map[string]string{
			"id": defaultUuid.String(),
		},
	}
	mr := &mockUserRepository{}

	deleteUser(mc, mr)

	assert.Equal(1, mr.deleteCalled)
}

func TestDeleteUser_WhenRepositorySucceeds_SetsStatusToNoContent(t *testing.T) {
	assert := assert.New(t)

	mc := &mockContext{
		params: map[string]string{
			"id": defaultUuid.String(),
		},
	}
	mr := &mockUserRepository{}

	err := deleteUser(mc, mr)

	assert.Nil(err)
	assert.Equal(http.StatusNoContent, mc.status)
}

func TestDeleteUser_WhenRepositoryFailsWithUnknownError_SetsStatusToInternalServerError(t *testing.T) {
	assert := assert.New(t)

	mc := &mockContext{
		params: map[string]string{
			"id": defaultUuid.String(),
		},
	}
	mr := &mockUserRepository{
		err: errDefault,
	}

	err := deleteUser(mc, mr)

	assert.Nil(err)
	assert.Equal(http.StatusInternalServerError, mc.status)
}

func TestDeleteUser_WhenRepositoryFailsWithNoMatchingRows_SetsStatusToNotFound(t *testing.T) {
	assert := assert.New(t)

	mc := &mockContext{
		params: map[string]string{
			"id": defaultUuid.String(),
		},
	}
	mr := &mockUserRepository{
		err: errors.NewCode(db.NoMatchingSqlRows),
	}

	err := deleteUser(mc, mr)

	assert.Nil(err)
	assert.Equal(http.StatusNotFound, mc.status)
}

func (m *mockContext) Request() *http.Request {
	return httptest.NewRequest(http.MethodGet, "http://localhost:3000", nil)
}

func (m *mockContext) Param(key string) string {
	if m.params == nil {
		return ""
	}

	if value, ok := m.params[key]; ok {
		return value
	}

	return ""
}

func (m *mockContext) JSON(status int, message interface{}) error {
	m.status = status
	m.data = message
	return nil
}

func (m *mockContext) NoContent(status int) error {
	m.status = status
	return nil
}

func (m *mockUserRepository) Create(ctx context.Context, user persistence.User) error {
	m.createCalled++
	return m.err
}

func (m *mockUserRepository) Get(ctx context.Context, id uuid.UUID) (persistence.User, error) {
	m.getCalled++
	return m.user, m.err
}

func (m *mockUserRepository) Update(ctx context.Context, user persistence.User) (persistence.User, error) {
	m.updateCalled++
	return m.user, m.err
}

func (m *mockUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	m.deleteCalled++
	return m.err
}
