package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const (
	appName    = "HTTP Server"
	appVersion = "v0.1.0"
	defaultDir = "./public"
)

func main() {
	/// Flag args
	d := flag.String("d", defaultDir, "Which dir to serve?")
	a := flag.String("a", "localhost", "An address to use")
	p := flag.Int("p", 5000, "A port to use")
	v := flag.Bool("v", false, "Show the app version")
	flag.Parse()

	if *v {
		fmt.Println(appName)
		fmt.Println(appVersion)
		os.Exit(0)
	}

	// Default path
	if *d == defaultDir {
		// If not exist, use the working dir
		if _, err := os.Stat(defaultDir); os.IsNotExist(err) {
			workingDir := "./"
			d = &workingDir
		}
	}

	/// Server
	e := echo.New()
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "[${status}] ${method} ${uri} ${remote_ip} ${error}\n",
	}))
	e.Static("/", *d)

	fmt.Println("Server at:", *d)
	if err := e.Start(fmt.Sprintf("%s:%v", *a, *p)); err != nil {
		log.Fatal(err)
	}
}
