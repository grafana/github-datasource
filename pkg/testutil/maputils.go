package testutil

import "testing"

// EnsureKeyIsSet ensures that a single key is set in the map (m)
func EnsureKeyIsSet(t *testing.T, m map[string]interface{}, key string) {
	if _, ok := m[key]; !ok {
		t.Errorf("key '%s' is not in map", key)
	}
}

// EnsureKeysAreSet ensures that all of the provided keys are set in the map (m)
func EnsureKeysAreSet(t *testing.T, m map[string]interface{}, keys ...string) {
	for _, v := range keys {
		EnsureKeyIsSet(t, m, v)
	}
}

// GetTestVariablesFunction provides a function that satisfies the TestVariables function of a TestClient
func GetTestVariablesFunction(keys ...string) func(*testing.T, map[string]interface{}) {
	return func(t *testing.T, m map[string]interface{}) {
		EnsureKeysAreSet(t, m, keys...)
	}
}
