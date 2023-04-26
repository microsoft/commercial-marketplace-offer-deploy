package app

import (
	"context"
	"net/http"
	"os"

	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/config"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/tasks"
	log "github.com/sirupsen/logrus"
)

type readinessTaskOptions struct {
	readinessFilePath string
	isReady           func() bool
	serviceUrl        string
	name              string
}

func newReadinessTask(appConfig *config.AppConfig, isReady func() bool) tasks.Task {
	options := readinessTaskOptions{
		readinessFilePath: appConfig.GetReadinessFilePath(),
		isReady:           isReady,
		serviceUrl:        appConfig.GetPublicBaseUrl(),
		name:              "Readiness Task",
	}

	return tasks.NewTask(options.name, func(ctx context.Context) error {
		statusCode := http.StatusNotFound

		for statusCode != http.StatusOK {
			response, err := http.Get(options.serviceUrl)

			if err != nil {
				continue
			}

			log.Infof("status code [%v] '%s'", response.StatusCode, options.serviceUrl)
			statusCode = response.StatusCode
		}
		// with a 200 OK, we're ready (we need the service url to be there so event grid can work)
		// this can be tested with a private URL that isn't reachable publicly, but event grid registration task will fail without a public url
		err := makeReady(&options)

		if err != nil {
			return err
		}
		log.Println("MODM is now ready.")

		return nil
	})
}

func makeReady(options *readinessTaskOptions) error {
	readiness, err := os.Create(options.readinessFilePath)

	if err != nil {
		return err
	}
	defer readiness.Close()

	return nil
}
