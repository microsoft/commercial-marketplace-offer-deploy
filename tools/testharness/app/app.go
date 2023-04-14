package app

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/labstack/echo"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/config"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/hosting"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/utils"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/api"
	"github.com/microsoft/commercial-marketplace-offer-deploy/sdk"
)

// TODO: this needs to go and pull from .env
var (
	location       = "eastus"
	resourceGroup  = "demo2"
	subscription   = "31e9f9a0-9fd2-4294-a0a3-0101246d9700"
	clientEndpoint = "https://dnsbobjac26.eastus.azurecontainer.io:443/api"
)

func GetRoutes(appConfig *config.AppConfig) hosting.Routes {

	return hosting.Routes{
		hosting.Route{
			Name:        "WebHookResponse",
			Method:      http.MethodPost,
			Path:        "/webhook",
			HandlerFunc: ReceiveEventNotification,
		},
		hosting.Route{
			Name:        "CreateDeployment",
			Method:      http.MethodGet,
			Path:        "/createdeployment",
			HandlerFunc: CreateDeployment,
		},
		hosting.Route{
			Name:        "StartDeployment",
			Method:      http.MethodGet,
			Path:        "/startdeployment/:deploymentId",
			HandlerFunc: StartDeployment,
		},
		hosting.Route{
			Name:        "CreateEventSubscription",
			Method:      http.MethodGet,
			Path:        "/createeventsubscription",
			HandlerFunc: CreateEventSubscription,
		},
	}
}

func getJsonAsMap(path string) map[string]interface{} {
	jsonMap, err := utils.ReadJson(path)
	if err != nil {
		log.Println(err)
	}
	return jsonMap
}

func ReceiveEventNotification(c echo.Context) error {
	log.Println("Event Web Hook received")
	var bodyJson any
	c.Bind(&bodyJson)
	return c.JSON(http.StatusOK, bodyJson)
}

func CreateEventSubscription(c echo.Context) error {
	//	templatePath := "./eventsubscription/mainTemplateBicep.json"
	//templateMap := getJsonAsMap(templatePath)

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Println(err)
	}

	client, err := sdk.NewClient(clientEndpoint, cred, nil)
	if err != nil {
		log.Panicln(err)
	}
	ctx := context.Background()

	subscriptionName := "webhook-1"
	apiKey := "1234"
	callbackclientEndpoint := "http://localhost:8280/webhook"

	request := api.CreateEventSubscriptionRequest{
		APIKey:   &apiKey,
		Callback: &callbackclientEndpoint,
		Name:     &subscriptionName,
	}

	res, err := client.CreateEventSubscription(ctx, request)
	if err != nil {
		log.Panicln(err)
	}

	return c.JSON(http.StatusOK, res)
}

func CreateDeployment(c echo.Context) error {
	log.Println("Inside CreateDeployment")
	templatePath := "./mainTemplateBicep.json"
	templateMap := getJsonAsMap(templatePath)
	log.Printf("The templateMap is %s", templateMap)

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Println(err)
	}
	log.Println("Got the credentials")

	log.Printf("Calling NewClient with endpoint %s", clientEndpoint)
	client, err := sdk.NewClient(clientEndpoint, cred, nil)
	if err != nil {
		log.Panicln(err)
	}
	log.Println("Got the client")
	ctx := context.Background()

	deploymentName := "TaggedDeployment"
	request := api.CreateDeployment{
		Name:           &deploymentName,
		Template:       templateMap,
		Location:       &location,
		ResourceGroup:  &resourceGroup,
		SubscriptionID: &subscription,
	}

	res, err := client.CreateDeployment(ctx, request)
	log.Printf("%v", res)
	if err != nil {
		log.Panicln(err)
	}

	return c.JSON(http.StatusOK, res)
}

func StartDeployment(c echo.Context) error {
	deploymentId, err := strconv.Atoi(c.Param("deploymentId"))

	if err != nil {
		log.Println(err)
	}

	paramsPath := "./parametersBicep.json"
	paramsMap := getJsonAsMap(paramsPath)

	cred, err := azidentity.NewDefaultAzureCredential(nil)

	if err != nil {
		log.Println(err)
	}

	client, err := sdk.NewClient(clientEndpoint, cred, nil)
	if err != nil {
		log.Println(err)
	}

	var ctx context.Context = context.Background()

	// TODO: properly construct the startdeployment params
	// create
	res, err := client.StartDeployment(ctx, int32(deploymentId), paramsMap)
	if err != nil {
		log.Println(err)
	}

	return c.JSON(http.StatusOK, res)
}

func BuildApp(configurationFilePath string) *hosting.App {
	builder := hosting.NewAppBuilder()

	appConfig := &config.AppConfig{}
	config.LoadConfiguration(configurationFilePath, nil, appConfig)
	builder.AddConfig(appConfig)

	builder.AddRoutes(func(options *hosting.RouteOptions) {
		routes := GetRoutes(appConfig)
		*options.Routes = routes
	})

	app := builder.Build(nil)
	return app
}
