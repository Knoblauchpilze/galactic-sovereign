package db

import (
	"strconv"

	"github.com/google/uuid"
)

func ToSliceInterface(ids []uuid.UUID) []interface{} {
	out := make([]interface{}, len(ids))

	for index, value := range ids {
		out[index] = value
	}

	return out
}

func GenerateInClauseForArgs(count int) string {
	var inClause string

	for i := 0; i < count; i++ {
		inClause += `$` + strconv.Itoa(i+1)
		if i < count-1 {
			inClause += `,`
		}
	}

	return inClause
}
