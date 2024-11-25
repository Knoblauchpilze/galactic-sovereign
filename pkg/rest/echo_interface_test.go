package rest

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnit_CreateEchoServer_DisablesBanner(t *testing.T) {
	assert := assert.New(t)

	e := createEchoServer()

	assert.True(e.HideBanner)
}

func TestUnit_CreateEchoServer_DisablesPort(t *testing.T) {
	assert := assert.New(t)

	e := createEchoServer()

	assert.True(e.HidePort)
}

func TestUnit_CreateEchoServer_AttachesLoggerWithPrefix(t *testing.T) {
	assert := assert.New(t)

	e := createEchoServer()

	assert.Equal("server", e.Logger.Prefix())
}
