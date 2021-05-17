package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"path"
	"strings"
	"time"
)

const (
	appName    = "HTTP Server"
	appVersion = "v0.3.0"
)

func main() {
	/// Flag args
	root := flag.String("root", "./", "Which dir to serve?")
	host := flag.String("host", "localhost", "An address to use")
	port := flag.Int("port", 8080, "A port to use")
	auth := flag.String("auth", "", "Authentication URL")
	exp := flag.Int64("exp", 15*60, "Expire token")
	redirectPage := flag.String("redirect", "/auth", "Authentication page")
	v := flag.Bool("v", false, "Verbose")
	protected := flag.String("protected", "", "Protected dirs")
	flag.Parse()

	var debug func(v ...interface{})
	if *v {
		debug = log.Println
	} else {
		debug = func(v ...interface{}) {}
	}

	/// Auth
	allowedToken := map[string]int64{}
	protectedDirs := strings.Fields(*protected)

	/// Server
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		isAllowed := true
		debug(r.Method, r.URL.Path)
		for _, v := range protectedDirs {
			if strings.Contains(r.URL.Path, v) {
				debug("Protected:", r.URL.Path, v)
				isAllowed = false

				// Token in cooke
				cookie, err := r.Cookie("token")
				if err != nil {
					debug(" Token not found")
					break
				}

				token := cookie.Value
				debug(" Token:", token)
				if token == "" {
					break
				}

				// Checking the token
				now := time.Now().Unix()
				value, hasToken := allowedToken[token]
				// Already in `allowedToken` and not exp
				if hasToken && value > now {
					debug(" Token exist, not expired")
					hasToken = true
					isAllowed = true
					break
				}

				// Exp or
				debug(" Auth via:", *auth)
				client := &http.Client{}
				req, err := http.NewRequest(http.MethodPost, *auth, nil)
				if err != nil {
					break
				}

				req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
				res, err := client.Do(req)
				if err == nil && res.StatusCode == http.StatusOK {
					isAllowed = true

					// Add or update the token
					debug(" Add or update the token")
					allowedToken[token] = now + *exp
				} else {
					// Delete the invalid one
					debug(" Remove invalid token")
					delete(allowedToken, token)
				}

				break
			}
		}

		if isAllowed {
			debug(" Allowed")
			filename := path.Join(*root, r.URL.Path)
			http.ServeFile(w, r, filename)
		} else {
			debug(" Not allowed, redirect to:", *redirectPage)
			http.Redirect(w, r, *redirectPage, http.StatusSeeOther)
		}
	})

	addr := fmt.Sprintf("%s:%v", *host, *port)

	fmt.Println(fmt.Sprintf("%s %s", appName, appVersion))
	fmt.Println("Protected dirs:", protectedDirs)
	fmt.Println("Authentication URL:", *auth)
	fmt.Println("Expired in:", *exp)
	fmt.Println("Redirect page:", *redirectPage)
	fmt.Println("Serve", *root, "at:", addr)

	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal(err)
	}
}
