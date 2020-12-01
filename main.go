package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

const (
	appName    = "HTTP Server"
	appVersion = "v0.2.0"
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
	fmt.Println("Server at:", *d)
	addr := fmt.Sprintf("%s:%v", *a, *p)
	fs := http.FileServer(http.Dir(*d))
	err := http.ListenAndServe(addr, fs)
	if err != nil {
		log.Fatal(err)
	}
}
