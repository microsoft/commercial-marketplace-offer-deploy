package app

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/config"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/tasks"
	log "github.com/sirupsen/logrus"
)

type readinessTaskOptions struct {
	readinessFilePath string
	signalReadiness   func()
	publicUrl         string
	name              string
}

func newReadinessTask(appConfig *config.AppConfig, signalReadiness func()) tasks.Task {
	options := readinessTaskOptions{
		readinessFilePath: appConfig.GetReadinessFilePath(),
		signalReadiness:   signalReadiness,
		publicUrl:         appConfig.GetPublicBaseUrl(),
		name:              "Readiness Task",
	}

	return tasks.NewTask(options.name, func(ctx context.Context) error {
		statusCode := http.StatusNotFound

		for statusCode != http.StatusOK {
			response, err := http.Get(options.publicUrl)

			if err != nil {
				continue
			}

			statusCode = response.StatusCode
		}
		// with a 200 OK, we're ready (we need the public address to be there so event grid can work)
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

	if options.signalReadiness == nil {
		return fmt.Errorf("the signalReadiness func is nil for task %s.", options.name)
	}

	options.signalReadiness()
	defer readiness.Close()

	return nil
}
