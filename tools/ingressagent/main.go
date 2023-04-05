package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/labstack/echo"
	apiserver "github.com/microsoft/commercial-marketplace-offer-deploy/cmd/apiserver/app"
	operator "github.com/microsoft/commercial-marketplace-offer-deploy/cmd/operator/app"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/hosting"
	"github.com/spf13/viper"
	"golang.ngrok.com/ngrok"
	"golang.ngrok.com/ngrok/config"
)

func main() {
	if err := run(context.Background()); err != nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context) error {
	tunnel, err := ngrok.Listen(ctx,
		config.HTTPEndpoint(),
		ngrok.WithAuthtoken(viper.GetString("NGROK_AUTHTOKEN")),
	)
	if err != nil {
		return err
	}
	defer tunnel.Close()

	log.Println("tunnel created:", tunnel.URL())

	var appName string
	flag.StringVar(&appName, "app", "", "App name to tunnel to")
	port := flag.Int("port", 8080, "Port to tunnel to. default: 8080")
	flag.Parse()

	return getApp(appName).Start(*port, func(e *echo.Echo) {
		e.Listener = tunnel
	})
}

func getApp(appName string) *hosting.App {
	configPath := getExecutionDirectory()
	fmt.Println("app name: ", appName)

	switch appName {
	case "operator":
		return operator.BuildApp(configPath)
	case "apiserver":
		return apiserver.BuildApp(configPath)
	default:
		log.Fatal("invalid app name")
	}
	return nil
}

func getExecutionDirectory() string {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	path := filepath.Dir(ex)
	return path
}
