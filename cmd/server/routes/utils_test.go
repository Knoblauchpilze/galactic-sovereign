package routes

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSanitizePath_TrimPrefix_EmptyRoute(t *testing.T) {
	assert := assert.New(t)

	actual := sanitizePath("", trimPrefix)

	assert.Equal("", actual)
}

func TestSanitizePath_TrimPrefix_RemovesLeadingSlash(t *testing.T) {
	assert := assert.New(t)

	actual := sanitizePath("/route", trimPrefix)

	assert.Equal("route", actual)
}

func TestSanitizePath_TrimPrefix_RemovesTrailingSlash(t *testing.T) {
	assert := assert.New(t)

	actual := sanitizePath("/route/", trimPrefix)

	assert.Equal("route", actual)
}

func TestSanitizePath_TrimPrefix_PreservesPath(t *testing.T) {
	assert := assert.New(t)

	actual := sanitizePath("route/with/path/", trimPrefix)

	assert.Equal("route/with/path", actual)
}

func TestSanitizePath_AddPrefix_EmptyRoute(t *testing.T) {
	assert := assert.New(t)

	actual := sanitizePath("", addPrefix)

	assert.Equal("/", actual)
}

func TestSanitizePath_AddPrefix_AddLeadingSlash(t *testing.T) {
	assert := assert.New(t)

	actual := sanitizePath("route", addPrefix)

	assert.Equal("/route", actual)
}

func TestSanitizePath_AddPrefix_DoesNotDuplicateLeadingSlash(t *testing.T) {
	assert := assert.New(t)

	actual := sanitizePath("/route", addPrefix)

	assert.Equal("/route", actual)
}

func TestSanitizePath_AddPrefix_RemovesTrailingSlash(t *testing.T) {
	assert := assert.New(t)

	actual := sanitizePath("/route/", addPrefix)

	assert.Equal("/route", actual)
}

func TestSanitizePath_AddPrefix_PreservesPath(t *testing.T) {
	assert := assert.New(t)

	actual := sanitizePath("route/with/path/", addPrefix)

	assert.Equal("/route/with/path", actual)
}

func TestConcatenateEndpoints_AllEmpty(t *testing.T) {
	assert := assert.New(t)

	actual := ConcatenateEndpoints("", "")

	assert.Equal("/", actual)
}

func TestConcatenateEndpoints_EmptyEndpoint(t *testing.T) {
	assert := assert.New(t)

	actual := ConcatenateEndpoints("", "/some/path")

	assert.Equal("/some/path", actual)
}

func TestConcatenateEndpoints_EmptyPath(t *testing.T) {
	assert := assert.New(t)

	actual := ConcatenateEndpoints("/some/endpoint", "")

	assert.Equal("/some/endpoint", actual)
}

func TestConcatenateEndpoints_DoesNotGenerateDoubleSlash(t *testing.T) {
	assert := assert.New(t)

	actual := ConcatenateEndpoints("/some/endpoint", "/some/path")

	assert.Equal("/some/endpoint/some/path", actual)
}

func TestConcatenateEndpoints_ConcatenateCorrectly(t *testing.T) {
	assert := assert.New(t)

	actual := ConcatenateEndpoints("/some/endpoint", "some/path")

	assert.Equal("/some/endpoint/some/path", actual)
}
