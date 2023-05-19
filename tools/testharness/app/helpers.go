package app

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/microsoft/commercial-marketplace-offer-deploy/sdk"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func getClient(ctx context.Context) (*sdk.Client, error) {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Errorf("Error creating credential - %s", err)
		return nil, err
	}

	client, err := sdk.NewClient(getClientEndpoint(), cred, nil)
	if err != nil {
		log.Errorf("Error creating SDK client - %s", err)
		return nil, err
	}
	return client, nil
}

func getJsonAsMap(path string) map[string]interface{} {
	jsonMap, err := readJson(path)
	if err != nil {
		log.Println(err)
	}
	return jsonMap
}

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

func getTemplatePath(caseName string) string {
	path := env.GetString("TEMPLATE_PATH")
	if len(path) > 0 {
		log.Printf("Found TEMPLATE_PATH - %s", path)
	} else {
		path = "./templates"
	}
	return filepath.Join(path, caseName, "mainTemplate.json")
}

func getParamsPath(caseName string) string {
	templateParams := env.GetString("TEMPLATEPARAMS_PATH")
	if len(templateParams) > 0 {
		log.Printf("Found TEMPLATEPARAMS_PATH - %s", templateParams)
	} else {
		templateParams = "./templates"
	}
	return filepath.Join(templateParams, caseName, "parameters.json")
}

func getCallback() string {
	callback := env.GetString("CALLBACK_BASE_URL")
	if len(callback) > 0 {
		return callback
	}

	//TODO: use the value that's set on echo
	return "http://localhost:" + strconv.Itoa(8280)
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

func readJson(path string) (map[string]interface{}, error) {
	templateFile, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	template := make(map[string]interface{})
	if err := json.Unmarshal(templateFile, &template); err != nil {
		return nil, err
	}
	return template, nil
}
