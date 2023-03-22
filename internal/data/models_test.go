package data

import (
	"log"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDeploymentTemplateMarshaling(t *testing.T) {
	template := make(map[string]interface{})

	attr := make(map[string]interface{})
	attr["field"] = "test"
	template["mapField"] = attr

	model := &Deployment{Template: template}

	db := NewDatabase(&DatabaseOptions{UseInMemory: true}).Instance()
	db.Save(model)

	var result *Deployment
	db.First(&result)

	log.Print(result.Template)
	require.NotNil(t, result)
	require.Equal(t, "test", result.Template["mapField"].(map[string]interface{})["field"])
}
