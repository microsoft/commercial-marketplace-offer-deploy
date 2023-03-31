package handlers

import (
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm/internal/resource"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/azure/azure-sdk-for-go/services/eventgrid/2018-01-01/eventgrid"
	"github.com/labstack/echo"
	internal "github.com/microsoft/commercial-marketplace-offer-deploy/internal/azure/eventgrid"
	"gorm.io/gorm"
)

// this API handler is the webook endpoint that receives event grid events
func EventGridWebHookHandler(c echo.Context, db *gorm.DB) error {
	validationResult := eventGridSubscriptionValidation(c, db)

	if validationResult != nil {
		return validationResult
	}

	messages := []eventgrid.Event{}
	err := c.Bind(&messages)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	for _, message := range messages {
		resourceId, err := getResourceId(c, message)

		if err != nil {
			continue
		}
		resourcesClient, err := armresources.NewClient(resourceId.SubscriptionID, testsuite.cred, testsuite.options)
	}
	return nil
}

func eventGridSubscriptionValidation(c echo.Context, db *gorm.DB) error {
	webhookValidator := internal.NewWebHookValidationEventHandler(c.Bind)
	result := webhookValidator.Handle(c.Request())

	if result.Error != nil {
		return echo.NewHTTPError(http.StatusBadRequest, result.Error.Error())
	}

	if result.Handled {
		return c.JSON(http.StatusOK, &result.Response)
	}

	return nil
}

func getResourceId(c echo.Context, message eventgrid.Event) (internal.ResourceEventData, error) {
	data := message.Data.(internal.ResourceEventData)

	messages := []eventgrid.Event{}
	err := c.Bind(&messages)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	resourceId, err := resource.ParseResourceID(ResourceURI)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return resourceId, nil
}

// func getResource(resourceId resource.ResourceID) (resource.ResourceID, error)) {
// 	resourcesClient, err := armresources.NewClient(testsuite.subscriptionID, testsuite.cred, testsuite.options)
// }
