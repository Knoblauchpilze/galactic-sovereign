package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsInterfaceNil_Nil(t *testing.T) {
	assert := assert.New(t)

	assert.NotPanics(func() {
		IsInterfaceNil(nil)
	})
	assert.True(IsInterfaceNil(nil))
}

func TestIsInterfaceNil_Value(t *testing.T) {
	assert := assert.New(t)

	assert.NotPanics(func() {
		IsInterfaceNil(1)
	})
	assert.False(IsInterfaceNil(1))
}

func TestIsInterfaceNil_Ptr(t *testing.T) {
	assert := assert.New(t)

	var out *int

	assert.NotPanics(func() {
		IsInterfaceNil(out)
	})
	assert.True(IsInterfaceNil(out))

	out = new(int)
	assert.False(IsInterfaceNil(out))
}

func TestIsInterfaceNil_Map(t *testing.T) {
	assert := assert.New(t)

	var out map[string]int

	assert.NotPanics(func() {
		IsInterfaceNil(out)
	})
	assert.True(IsInterfaceNil(out))

	out = make(map[string]int)
	assert.False(IsInterfaceNil(out))
}

func TestIsInterfaceNil_Array(t *testing.T) {
	assert := assert.New(t)

	var out [1]int

	assert.NotPanics(func() {
		IsInterfaceNil(out)
	})
	assert.False(IsInterfaceNil(out))
}

func TestIsInterfaceNil_Chan(t *testing.T) {
	assert := assert.New(t)

	var out chan string

	assert.NotPanics(func() {
		IsInterfaceNil(out)
	})
	assert.True(IsInterfaceNil(out))

	out = make(chan string)
	assert.False(IsInterfaceNil(out))
}

func TestIsInterfaceNil_Slice(t *testing.T) {
	assert := assert.New(t)

	var out []int

	assert.NotPanics(func() {
		IsInterfaceNil(out)
	})
	assert.True(IsInterfaceNil(out))

	out = make([]int, 0)
	assert.False(IsInterfaceNil(out))
}
