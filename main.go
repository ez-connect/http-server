package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
)

const (
	appName             = "HTTP Server"
	appVersion          = "v0.3.0"
	defaultHost         = "localhost"
	defaultPort         = 8080
	defaultRoot         = "./public"
	defaultPrivates     = "/private /protected"
	defaultRedirectPage = "/auth"
)

func main() {
	/// Flag args
	root := flag.String("root", defaultRoot, "Which dir to serve?")
	host := flag.String("host", defaultHost, "An address to use")
	port := flag.Int("port", defaultPort, "A port to use")
	privates := flag.String("privates", defaultPrivates, "Protected dirs")
	auth := flag.String("auth", "", "Authentication URL")
	redirectPage := flag.String("redirect", defaultRedirectPage, "Authentication page")
	v := flag.Bool("v", false, "Show the app version")
	flag.Parse()

	if *v {
		fmt.Println(appName)
		fmt.Println(appVersion)
		os.Exit(0)
	}

	/// Auth
	privateDirs := strings.Split(*privates, " ")

	/// Server
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		isAllowed := true
		for _, v := range privateDirs {
			if strings.Contains(r.URL.Path, v) {
				isAllowed = false
				client := &http.Client{}
				req, err := http.NewRequest(http.MethodPost, *auth, nil)
				if err != nil {
					break
				}

				// Token in cooke
				cookie, err := r.Cookie("token")
				if err != nil {
					break
				}

				token := cookie.Value
				log.Println(token)

				req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
				res, err := client.Do(req)
				if err == nil && res.StatusCode == http.StatusOK {
					isAllowed = true
				}

				break
			}
		}

		if isAllowed {
			filename := path.Join(*root, r.URL.Path)
			http.ServeFile(w, r, filename)
		} else {
			http.Redirect(w, r, *redirectPage, http.StatusSeeOther)
		}
	})

	addr := fmt.Sprintf("%s:%v", *host, *port)
	fmt.Println("Server", *root, "at", addr)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal(err)
	}
}
