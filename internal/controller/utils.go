package controller

import "encoding/json"

func marshalNilToEmptySlice[T any](in []T) ([]byte, error) {
	toMarshal := make([]T, 0)
	if in != nil {
		toMarshal = in
	}

	return json.Marshal(toMarshal)
}
