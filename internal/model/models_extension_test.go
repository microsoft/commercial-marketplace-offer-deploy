package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func Test_Deployment_GetAzureDeploymentName(t *testing.T) {
	type test struct {
		input string
		want  string
	}

	tests := []test{
		{input: "Deployment name with spaces", want: "modm.1.Deployment-name-with-spaces"},
		{input: "test/slash", want: "modm.1.testslash"},
		{input: ",.~`{}|/<>[]rg-with-special-*&^%$#@!_+=.:'\"", want: "modm.1.rg-with-special"},
	}

	for _, tc := range tests {
		d := &Deployment{
			Model: gorm.Model{ID: 1},
			Name:  tc.input,
		}
		got := d.GetAzureDeploymentName()
		assert.Equal(t, tc.want, got)
	}
}
