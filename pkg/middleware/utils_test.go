package middleware

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPathFromRequest_NilRequest(t *testing.T) {
	assert := assert.New(t)

	actual := pathFromRequest(nil)

	assert.Equal("", actual)
}

func TestPathFromRequest_NoPath(t *testing.T) {
	assert := assert.New(t)
	req := http.Request{
		Host: "host",
	}

	actual := pathFromRequest(&req)

	assert.Equal("host", actual)
}

func TestPathFromRequest_WithPath(t *testing.T) {
	assert := assert.New(t)
	req := http.Request{
		Host: "host",
		URL: &url.URL{
			Path: "/path",
		},
	}

	actual := pathFromRequest(&req)

	assert.Equal("host/path", actual)
}

func TestFormatHttpStatusCode_2XX(t *testing.T) {
	assert := assert.New(t)

	actual := formatHttpStatusCode(http.StatusOK)
	assert.Equal("\x1b[1;32m200\x1b[0m", actual)

	actual = formatHttpStatusCode(http.StatusAccepted)
	assert.Equal("\x1b[1;32m202\x1b[0m", actual)
}

func TestFormatHttpStatusCode_3XX(t *testing.T) {
	assert := assert.New(t)

	actual := formatHttpStatusCode(http.StatusFound)
	assert.Equal("\x1b[1;36m302\x1b[0m", actual)

	actual = formatHttpStatusCode(http.StatusNotModified)
	assert.Equal("\x1b[1;36m304\x1b[0m", actual)
}

func TestFormatHttpStatusCode_4XX(t *testing.T) {
	assert := assert.New(t)

	actual := formatHttpStatusCode(http.StatusBadRequest)
	assert.Equal("\x1b[1;33m400\x1b[0m", actual)

	actual = formatHttpStatusCode(http.StatusForbidden)
	assert.Equal("\x1b[1;33m403\x1b[0m", actual)
}

func TestFormatHttpStatusCode_5XX(t *testing.T) {
	assert := assert.New(t)

	actual := formatHttpStatusCode(http.StatusInternalServerError)
	assert.Equal("\x1b[1;31m500\x1b[0m", actual)

	actual = formatHttpStatusCode(http.StatusBadGateway)
	assert.Equal("\x1b[1;31m502\x1b[0m", actual)
}

// func formatHttpStatusCode(status int) string {
// 	switch {
// 	case status >= 500:
// 		return logger.FormatWithColor(status, logger.Red)
// 	case status >= 400:
// 		return logger.FormatWithColor(status, logger.Yellow)
// 	case status >= 300:
// 		return logger.FormatWithColor(status, logger.Cyan)
// 	default:
// 		return logger.FormatWithColor(status, logger.Green)
// 	}
// }
