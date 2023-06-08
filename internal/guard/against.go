package guard

import (
	"reflect"

	log "github.com/sirupsen/logrus"
)

type GuardFunc[T any] func(i T)

func New[T any](guardAction GuardFunc[*T]) *T {
	instance := new(T)
	guardAction(instance)

	return instance
}

func AgainstNil(ptr interface{}) {
	if IsPointerNil(ptr) {
		panic("nil pointer")
	}
}

func IsPointerNil(ptr interface{}) bool {
	value := reflect.ValueOf(ptr)
	kind := value.Kind()

	if kind != reflect.Ptr {
		return false
	}

	v := value.IsNil()
	if v {
		callStack := GetCallstack()
		log.WithField("stacktrace", callStack).Errorf("%s is nil", value)
	}

	return value.IsNil()
}
