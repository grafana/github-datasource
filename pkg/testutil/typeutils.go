package testutil

import (
	"reflect"
	"testing"
)

// TypesAreEqual compares the types a and b. If they are not equal, then false is returned. If they are equal, then true is returned.
func TypesAreEqual(a interface{}, b interface{}) bool {
	return reflect.TypeOf(a) == reflect.TypeOf(b)
}

// EnsureTypeEquality uses the test object and fails the test if the types are not equal
func EnsureTypeEquality(t *testing.T, actual interface{}, expected interface{}) {
	if !TypesAreEqual(actual, expected) {
		t.Errorf("Types are not equal. Expected '%s', received '%s", reflect.TypeOf(actual).String(), reflect.TypeOf(expected).String())
	}
}

// GetTestQueryFunction returns a function that satisfies the TestQuery function in the TestClient object
func GetTestQueryFunction(expected interface{}) func(*testing.T, interface{}) {
	return func(t *testing.T, actual interface{}) {
		EnsureTypeEquality(t, actual, expected)
	}
}
