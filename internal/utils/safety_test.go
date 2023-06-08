package utils

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

type testNil struct {
	foo string
}

func Test_IsNil(t *testing.T) {
	var nilTest *testNil
	if nilTest == nil {
		t.Log("nilTest is nil")
	}
	assert.True(t, IsPointerNil(nilTest))
}