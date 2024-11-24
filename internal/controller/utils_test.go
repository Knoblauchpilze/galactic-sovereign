package controller

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var defaultUuid = uuid.MustParse("08ce96a3-3430-48a8-a3b2-b1c987a207ca")

func TestMarshalNilToEmptySlice_WhenNil_ExpectMarshalToEmptySlice(t *testing.T) {
	assert := assert.New(t)

	var in []int

	actual, err := marshalNilToEmptySlice(in)

	assert.Nil(err)
	assert.Equal("[]", string(actual))
}

func TestMarshalNilToEmptySlice_WhenNotNil_ExpectMarshalCorrectData(t *testing.T) {
	assert := assert.New(t)

	in := []int{1, 2}

	actual, err := marshalNilToEmptySlice(in)

	assert.Nil(err)
	assert.Equal("[1,2]", string(actual))
}

var defaultKey = "my-key"

func TestFetchIdFromQueryParam_whenNoId_expectNotExistAndNoError(t *testing.T) {
	assert := assert.New(t)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	ctx, _ := generateTestEchoContextFromRequest(req)

	exists, _, err := fetchIdFromQueryParam(defaultKey, ctx)
	assert.False(exists)
	assert.Nil(err)
}

func TestFetchIdFromQueryParam_whenIdSetForOtherKey_expectNotExistAndNoError(t *testing.T) {
	assert := assert.New(t)

	req := generateRequestWithQueryParams("not-the-default-key", defaultUuid.String())
	ctx, _ := generateTestEchoContextFromRequest(req)

	exists, _, err := fetchIdFromQueryParam(defaultKey, ctx)
	assert.False(exists)
	assert.Nil(err)
}

func TestFetchIdFromQueryParam_whenIdSyntaxIsWrong_expectExistAndError(t *testing.T) {
	assert := assert.New(t)

	req := generateRequestWithQueryParams(defaultKey, "not-a-uuid")
	ctx, _ := generateTestEchoContextFromRequest(req)

	exists, _, err := fetchIdFromQueryParam(defaultKey, ctx)
	assert.True(exists)
	assert.Equal("invalid UUID length: 10", err.Error())
}

func TestFetchIdFromQueryParam_whenIdIsSet_expectExistCorrectIdAndNoError(t *testing.T) {
	assert := assert.New(t)

	req := generateRequestWithQueryParams(defaultKey, defaultUuid.String())
	ctx, _ := generateTestEchoContextFromRequest(req)

	exists, actual, err := fetchIdFromQueryParam(defaultKey, ctx)
	assert.True(exists)
	assert.Equal(defaultUuid, actual)
	assert.Nil(err)
}

func generateRequestWithQueryParams(key string, value string) *http.Request {
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	q := req.URL.Query()
	q.Add(key, value)

	req.URL.RawQuery = q.Encode()

	return req
}
