package app

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/labstack/echo/v4"

	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/utils"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/api"
	"github.com/microsoft/commercial-marketplace-offer-deploy/sdk"
)

// TODO: this needs to go and pull from .env
var (
	location      = "eastus"
	resourceGroup = "demo2"
	subscription  = "31e9f9a0-9fd2-4294-a0a3-0101246d9700"
	//clientEndpoint = "https://dnsbobjac67.eastus.azurecontainer.io:443/api"
	clientEndpoint = "http://localhost:8080"
	env            = loadEnvironmentVariables()
)

func getClientEndpoint() string {
	// no real need for viper here as we are just pulling 1 environment variable for the test harness
	endpoint := env.GetString("MODM_API_ENDPOINT")
	if len(endpoint) > 0 {
		return endpoint
	}
	return clientEndpoint
}

func getLocation() string {
	loc := env.GetString("MODM_DEPLOYMENT_LOCATION")
	if len(loc) > 0 {
		return loc
	}
	return location
}

func getSubscription() string {
	sub := env.GetString("MODM_SUBSCRIPTION")
	if len(sub) > 0 {
		return sub
	}
	return subscription
}

func getResourceGroup() string {
	rg := env.GetString("MODM_RESOURCE_GROUP")
	if len(rg) > 0 {
		return rg
	}
	return resourceGroup
}

func getTemplatePath() string {
	path := env.GetString("TEMPLATE_PATH")
	if len(path) > 0 {
		return path
	}
	return "./mainTemplateBicep.json"
}

func getParamsPath() string {
	templateParams := env.GetString("TEMPLATEPARAMS_PATH")
	if len(templateParams) > 0 {
		log.Printf("Found TEMPLATEPARAMS_PATH - %s", templateParams)
		return templateParams
	}
	return "./parametersBicep.json"
}

func getCallback() string {
	callback := env.GetString("CALLBACK_BASE_URL")
	if len(callback) > 0 {
		return callback
	}

	//TODO: use the value that's set on echo
	return "http://localhost:" + strconv.Itoa(8280)
}

func AddRoutes(e *echo.Echo) {
	e.GET("/", func(ctx echo.Context) error {
		return ctx.String(http.StatusOK, "Test Harness Up.")
	})
	e.GET("/createdeployment", CreateDeployment)
	e.GET("/startdeployment/:deploymentId", StartDeployment)
	e.GET("/createeventhook", CreateEventHook)
	e.GET("/dryrun/:deploymentId", DryRun)
	e.POST("/webhook", ReceiveEventNotification)
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

	json := c.JSON(http.StatusOK, bodyJson)
	log.Printf("ReceiveEventNotification response - %s", json)
	return json
}

func CreateEventHook(c echo.Context) error {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Println(err)
	}

	client, err := sdk.NewClient(getClientEndpoint(), cred, nil)
	if err != nil {
		log.Panicln(err)
	}
	ctx := context.Background()

	subscriptionName := "webhook-1"
	apiKey := "1234"
	callbackclientEndpoint := fmt.Sprintf("%s/webhook", getCallback())

	request := api.CreateEventHookRequest{
		APIKey:   &apiKey,
		Callback: &callbackclientEndpoint,
		Name:     &subscriptionName,
	}

	res, err := client.CreateEventHook(ctx, request)
	if err != nil {
		log.Panicln(err)
	}

	json := c.JSON(http.StatusOK, res)
	log.Printf("Create Event Hook response - %s", json)
	return json
}

func CreateDeployment(c echo.Context) error {
	location = getLocation()
	resourceGroup = getResourceGroup()
	subscription = getSubscription()

	log.Println("Inside CreateDeployment")
	templatePath := getTemplatePath()
	templateMap := getJsonAsMap(templatePath)
	log.Printf("The templateMap is %s", templateMap)

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Println(err)
	}
	log.Println("Got the credentials")

	log.Printf("Calling NewClient with endpoint %s", getClientEndpoint())
	client, err := sdk.NewClient(getClientEndpoint(), cred, nil)
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

	res, err := client.Create(ctx, request)
	log.Printf("%v", res)
	if err != nil {
		log.Panicln(err)
	}

	json := c.JSON(http.StatusOK, res)
	log.Printf("CreateDeployment response - %s", json)

	return json
}

func DryRun(c echo.Context) error {
	log.Println("Inside DryRun in the test harness")
	deploymentId, err := strconv.Atoi(c.Param("deploymentId"))

	if err != nil {
		log.Println(err)
	}

	paramsPath := getParamsPath()
	log.Printf("paramsPath - %v", paramsPath)
	paramsMap := getJsonAsMap(paramsPath)
	log.Printf("paramsMap - %v", paramsMap)

	cred, err := azidentity.NewDefaultAzureCredential(nil)

	if err != nil {
		log.Println(err)
	}

	client, err := sdk.NewClient(getClientEndpoint(), cred, nil)
	if err != nil {
		log.Println(err)
	}

	var ctx context.Context = context.Background()

	log.Printf("About to call DryRunDeployment - paramsMap: %s", paramsMap)
	res, err := client.DryRun(ctx, deploymentId, paramsMap)
	if err != nil {
		log.Println(err)
	}

	json := c.JSON(http.StatusOK, res)
	log.Printf("Dry run response - %s", json)

	return json
}

func StartDeployment(c echo.Context) error {
	deploymentId, err := strconv.Atoi(c.Param("deploymentId"))

	if err != nil {
		log.Println(err)
	}

	paramsPath := getParamsPath()
	paramsMap := getJsonAsMap(paramsPath)

	cred, err := azidentity.NewDefaultAzureCredential(nil)

	if err != nil {
		log.Println(err)
	}

	client, err := sdk.NewClient(getClientEndpoint(), cred, nil)
	if err != nil {
		log.Println(err)
	}

	var ctx context.Context = context.Background()

	// TODO: properly construct the startdeployment params
	// create
	res, err := client.Start(ctx, deploymentId, paramsMap, nil)
	if err != nil {
		log.Println(err)
	}

	return c.JSON(http.StatusOK, res)
}

func loadEnvironmentVariables() *viper.Viper {
	env := viper.New()
	env.AddConfigPath("./")
	env.SetConfigName(".env")
	env.SetConfigType("env")
	env.AutomaticEnv()

	err := env.ReadInConfig()
	if err != nil {
		log.Errorf("Error reading config file, %s", err)
	}
	return env
}
