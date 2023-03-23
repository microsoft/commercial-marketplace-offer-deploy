package data

import (
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/generated"
)

func FromCreateDeployment(from *generated.CreateDeployment) *Deployment {
	//TODO: parse out template into the stages

	deployment := &Deployment{
		Name:     *from.Name,
		Status:   "New",
		Template: from.Template.(map[string]interface{}),
	}
	return deployment
}
