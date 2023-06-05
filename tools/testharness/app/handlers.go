package app

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/microsoft/commercial-marketplace-offer-deploy/sdk"
	log "github.com/sirupsen/logrus"
)

func SetCase(c echo.Context) error {
	caseName = c.Param("caseName")
	return c.String(http.StatusOK, fmt.Sprintf("Case set to %s", caseName))
}

func Cancel(c echo.Context) error {
	deploymentId, err := strconv.Atoi(c.Param("deploymentId"))
	if err != nil {
		log.Println(err)
	}

	location = getLocation()
	resourceGroup = getResourceGroup()
	subscription = getSubscription()

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

	res, err := client.Cancel(ctx, deploymentId)
	if err != nil {
		log.Println(err)
	}

	return c.JSON(http.StatusOK, res)

}

func ReceiveEventHook(c echo.Context) error {
	log.Print("Event Hook Received")
	reader := c.Request().Body
	return c.Stream(http.StatusOK, "application/json", reader)
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

	hookName := "webhook-1"
	apiKey := "1234"
	callbackclientEndpoint := fmt.Sprintf("%s/webhook", getCallback())

	request := sdk.CreateEventHookRequest{
		APIKey:   &apiKey,
		Callback: &callbackclientEndpoint,
		Name:     &hookName,
	}

	res, err := client.CreateEventHook(ctx, request)
	if err != nil {
		log.Panicln(err)
	}

	json := c.JSON(http.StatusOK, res)
	log.Printf("Create Event Hook response - %s", json)
	return json
}

func Redeploy(c echo.Context) error {
	deploymentId, err := strconv.Atoi(c.Param("deploymentId"))
	if err != nil {
		log.Println(err)
	}
	stageName := c.Param("stageName")
	var stageId = uuid.Nil

	for _, v := range deployment.Stages {
		if strings.EqualFold(*v.DeploymentName, stageName) {
			stageId = uuid.MustParse(*v.ID)
		}
	}

	if err != nil {
		log.Println(err)
	}

	location = getLocation()
	resourceGroup = getResourceGroup()
	subscription = getSubscription()

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

	retryOptions := &sdk.RetryOptions{
		StageId: stageId,
	}

	resp, err := client.Retry(ctx, deploymentId, retryOptions)
	if err != nil {
		log.Println(err)
	}

	json := c.JSON(http.StatusOK, resp)
	log.Printf("Redeploy response - %s", json)

	return json
}

func HealthStatus(c echo.Context) error {
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
	res, err := client.HealthStatus(ctx)
	if err != nil {
		log.Println(err)
	}
	json := c.JSON(http.StatusOK, res)
	log.Printf("HealthStatus response - %s", json)

	return json
}

func CreateDeployment(c echo.Context) error {
	location = getLocation()
	resourceGroup = getResourceGroup()
	subscription = getSubscription()

	log.Println("Inside CreateDeployment")
	templatePath := getTemplatePath(caseName)

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
	request := sdk.CreateDeployment{
		Name:           &deploymentName,
		Template:       templateMap,
		Location:       &location,
		ResourceGroup:  &resourceGroup,
		SubscriptionID: &subscription,
	}

	res, err := client.Create(ctx, request)
	deployment = res
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

	paramsPath := getParamsPath(caseName)
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
	options := &sdk.DryRunOptions{
		Retries: 1,
	}
	res, err := client.DryRun(ctx, deploymentId, paramsMap, options)
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

	paramsPath := getParamsPath(caseName)
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

func GetStatus(c echo.Context) error {
	ctx := c.Request().Context()
	deploymentId, _ := strconv.Atoi(c.Param("deploymentId"))
	operationName := sdk.OperationType(c.Param("operationName"))

	client, _ := getClient(ctx)

	status, _ := client.GetDeploymentOperationStatus(ctx, deploymentId, operationName)

	return c.JSON(http.StatusOK, status)
}
