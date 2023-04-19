package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/microsoft/commercial-marketplace-offer-deploy/tools/testharness/app"
//	testharness "github.com/microsoft/commercial-marketplace-offer-deploy/tools/testharness/app"
)

var (
	//configurationFilePath string = "."
	port                  int    = 8080
)

func main() {
	formattedPort := ":" + strconv.Itoa(port)
	log.Printf("Server starting on %s", formattedPort)

	e := echo.New()
	app.AddRoutes(e)
	
	if err := e.Start(":8080"); err != http.ErrServerClosed {
		log.Fatal(err)
	  }
	// app := testharness.BuildApp(configurationFilePath)

	// options := &hosting.AppStartOptions{
	// 	Port: &port,
	// 	WebServer: true,
	// }
	// log.Fatal(app.Start(options))
}
