package utils

import (
	"reflect"
)

func CheckKeys(map1, map2 interface{}) bool {
	map1Value := reflect.ValueOf(map1)
	map2Value := reflect.ValueOf(map2)

	if map1Value.Type().Kind() != reflect.Map || map2Value.Type().Kind() != reflect.Map {
		return false // Not maps
	}

	map1Keys := map1Value.MapKeys()
	map2Keys := map2Value.MapKeys()

	for _, key := range map1Keys {
		if !containsKey(map2Keys, key) {
			return false
		}
	}

	return true
}

func containsKey(keys []reflect.Value, target reflect.Value) bool {
	for _, key := range keys {
		if reflect.DeepEqual(key.Interface(), target.Interface()) {
			return true
		}
	}
	return false
}