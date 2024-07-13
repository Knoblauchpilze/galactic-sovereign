package controller

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
