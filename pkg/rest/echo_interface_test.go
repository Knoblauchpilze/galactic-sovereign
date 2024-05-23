package rest

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateEchoServer_DisablesBanner(t *testing.T) {
	assert := assert.New(t)

	e := createEchoServer()

	assert.True(e.HideBanner)
}

func TestCreateEchoServer_DisablesPort(t *testing.T) {
	assert := assert.New(t)

	e := createEchoServer()

	assert.True(e.HidePort)
}

func TestCreateEchoServer_AttachesLoggerWithPrefix(t *testing.T) {
	assert := assert.New(t)

	e := createEchoServer()

	assert.Equal("server", e.Logger.Prefix())
}
