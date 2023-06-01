package mapper

import (
	"github.com/microsoft/commercial-marketplace-offer-deploy/sdk"
)

type GetHealthCheckResponseMapper struct {

}

func (m *GetHealthCheckResponseMapper) Map(isHealthy bool) sdk.GetHealthResponse {
	return sdk.GetHealthResponse{
		IsHealthy: &isHealthy,
	}
}
