package common

import "reflect"

// https://mangatmodi.medium.com/go-check-nil-interface-the-right-way-d142776edef1
func IsInterfaceNil(i interface{}) bool {
	if i == nil {
		return true
	}

	// Unlike what they say, `reflect.Array` makes the `IsNil` method panic, so we
	// removed it from here. A test verifies that the behavior works.
	switch reflect.TypeOf(i).Kind() {
	case reflect.Ptr, reflect.Map, reflect.Chan, reflect.Slice:
		return reflect.ValueOf(i).IsNil()
	}

	return false
}
