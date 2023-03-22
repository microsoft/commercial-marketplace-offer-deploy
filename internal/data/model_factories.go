package data

import (
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal"
)

func FromCreateDeployment(from *internal.CreateDeployment) *Deployment {
	//TODO: parse out template into the stages

	deployment := &Deployment{
		Name:     *from.Name,
		Status:   "New",
		Template: from.Template,
	}
	return deployment
}
