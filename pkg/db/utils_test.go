package db

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var defaultUuid = uuid.MustParse("08ce96a3-3430-48a8-a3b2-b1c987a207ca")

func TestToSliceInterface(t *testing.T) {
	assert := assert.New(t)

	testCases := map[string][]uuid.UUID{
		"nil":          nil,
		"empty":        {},
		"single-value": {defaultUuid},
		"multiple-values": {
			uuid.MustParse("7c9b9b54-f08d-4de5-bf59-b72ae8119bf7"),
			uuid.MustParse("bd1f54d9-b1cb-4003-81e7-f5b84f20a3e8"),
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			out := ToSliceInterface(testCase)

			assert.Equal(len(testCase), len(out))
			for id, actual := range out {
				assert.Equal(actual, testCase[id])
			}
		})
	}
}

type inClauseTestCase struct {
	count    int
	expected string
}

func TestGenerateInClauseForArgs(t *testing.T) {
	assert := assert.New(t)

	testCases := map[string]inClauseTestCase{
		"zero":     {0, ""},
		"one":      {1, "$1"},
		"multiple": {3, "$1,$2,$3"},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			actual := GenerateInClauseForArgs(testCase.count)
			assert.Equal(testCase.expected, actual)
		})
	}
}
