package rest

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSanitizePath(t *testing.T) {
	assert := assert.New(t)

	type testCase struct {
		in       string
		expected string
	}

	testCases := []testCase{
		{in: "", expected: "/"},
		{in: "/", expected: "/"},
		{in: "path", expected: "/path"},
		{in: "path/", expected: "/path"},
		{in: "/path", expected: "/path"},
		{in: "/path/", expected: "/path"},
		{in: "path/id", expected: "/path/id"},
		{in: "path/id/", expected: "/path/id"},
		{in: "/path/id", expected: "/path/id"},
		{in: "/path/id/", expected: "/path/id"},
	}

	for _, testCase := range testCases {
		t.Run("", func(t *testing.T) {
			actual := sanitizePath(testCase.in)

			assert.Equal(testCase.expected, actual)
		})
	}
}

func TestConcatenateEndpoints_AllEmpty(t *testing.T) {
	assert := assert.New(t)

	actual := concatenateEndpoints("", "")

	assert.Equal("/", actual)
}

func TestConcatenateEndpoints_EmptyEndpoint(t *testing.T) {
	assert := assert.New(t)

	actual := concatenateEndpoints("", "/some/path")

	assert.Equal("/some/path", actual)
}

func TestConcatenateEndpoints_EmptyPath(t *testing.T) {
	assert := assert.New(t)

	actual := concatenateEndpoints("/some/endpoint", "")

	assert.Equal("/some/endpoint", actual)
}

func TestConcatenateEndpoints_DoesNotGenerateDoubleSlash(t *testing.T) {
	assert := assert.New(t)

	actual := concatenateEndpoints("/some/endpoint", "/some/path")

	assert.Equal("/some/endpoint/some/path", actual)
}

func TestConcatenateEndpoints_ConcatenateCorrectly(t *testing.T) {
	assert := assert.New(t)

	actual := concatenateEndpoints("/some/endpoint", "some/path")

	assert.Equal("/some/endpoint/some/path", actual)
}
