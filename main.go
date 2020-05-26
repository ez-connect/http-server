package main

import (
	"flag"
	"fmt"
	"log"
	"os/user"

	"github.com/labstack/echo/v4"
)

const (
	appName     = "Static server"
	appVersion  = "v0.0.1"
	defaultPort = 8080
)

func main() {
	/// Flag args
	root := flag.String("root", "", "The root dir")
	port := flag.Int("port", defaultPort, "Which port will be used to serve?")
	flag.Parse()

	if *root == "" {
		usr, err := user.Current()
		if err != nil {
			log.Fatal(err)
		}

		root = &usr.HomeDir
	}

	/// Server
	e := echo.New()
	e.Static("/", *root)

	addr := fmt.Sprintf(":%v", *port)
	fmt.Println(fmt.Sprintf("Static server is running: http://localhost:%v", *port))
	fmt.Println("Root dir:", *root)
	if err := e.Start(addr); err != nil {
		log.Fatal(err)
	}
}
