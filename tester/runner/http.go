package main

import (
	"go.uber.org/zap"
	"reflect"
)

// returned true if all keys present in the subset have the same value in the superset
func IsSubset(subset, superset map[string][]string) bool {
	for key, expectedVal := range subset {
		if !reflect.DeepEqual(superset[key], expectedVal) {
			logger.Error("key is not equal in subset and superset",
				zap.String("key", key),
				zap.Strings("subset", expectedVal),
				zap.Strings("superset", superset[key]))
			return false
		}
	}
	return true
}
