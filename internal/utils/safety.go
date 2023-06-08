package utils

import (
	"reflect"
	log "github.com/sirupsen/logrus"
)

func IsPointerNil(ptr interface{}) bool {
	value := reflect.ValueOf(ptr)
	kind := value.Kind()

	if kind != reflect.Ptr {
		return false
	}

	v := value.IsNil()
	if v {
		log.Tracef("%s is nil", value)
		callStack := GetCallstack()
		log.Trace(callStack)
	}

	return value.IsNil()
}