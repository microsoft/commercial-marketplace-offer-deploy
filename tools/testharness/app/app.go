package app

import (
	"context"
	"log"
	"net/http"

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
	location      = "eastus"
	resourceGroup = "demo2"
	subscription  = "31e9f9a0-9fd2-4294-a0a3-0101246d9700"
)

func GetRoutes(appConfig *config.AppConfig) hosting.Routes {

	return hosting.Routes{
		hosting.Route{
			Name:        "WebHookResponse",
			Method:      http.MethodGet,
			Path:        "/deploymentevent",
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
			Path:        "/startdeployment",
			HandlerFunc: StartDeployment,
		},
		hosting.Route{
			Name:        "CreateEventSubscription",
			Method:      http.MethodGet,
			Path:        "/ces",
			HandlerFunc: StartDeployment,
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
	return c.JSON(http.StatusOK, "Registered endpoint hit")
}

func CreateEventSubscription(c echo.Context) error {
	url := "http://localhost:8080/events/subscriptions"
//	templatePath := "./eventsubscription/mainTemplateBicep.json"
	//templateMap := getJsonAsMap(templatePath)

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Println(err)
	}

	client, err := sdk.NewClient(url, cred, nil)
	if err != nil {
		log.Panicln(err)
	}
	ctx := context.Background()

	subscriptionName := "TaggedSubscription"
	apiKey := "1234"
	callbackUrl := "http://localhost:8280/deploymentevent"
	request := api.CreateEventSubscriptionRequest{
		APIKey: &apiKey,
		Callback: &callbackUrl,
		Name: &subscriptionName,
	}

	res, err := client.CreateEventSubscription(ctx, request)
	if err != nil {
		log.Panicln(err)
	}

	return c.JSON(http.StatusOK, res)
}

func CreateDeployment(c echo.Context) error {
	log.Println("Inside CreateDeployment")
	url := "http://localhost:8080"
	templatePath := "./mainTemplateBicep.json"
	templateMap := getJsonAsMap(templatePath)
	log.Printf("The templateMap is %s", templateMap)

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Println(err)
	}
	log.Println("Got the credentials")

	client, err := sdk.NewClient(url, cred, nil)
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
	url := "http://localhost:8080"
	paramsPath := "./parametersBicep.json"
	paramsMap := getJsonAsMap(paramsPath)

	cred, err := azidentity.NewDefaultAzureCredential(nil)

	if err != nil {
		log.Println(err)
	}

	client, err := sdk.NewClient(url, cred, nil)
	if err != nil {
		log.Println(err)
	}

	var ctx context.Context = context.Background()

	// TODO: properly construct the startdeployment params
	// create
	res, err := client.StartDeployment(ctx, 1, paramsMap)
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
