package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"time"
)

const (
	appName    = "HTTP Server"
	appVersion = "0.3.3"
)

func main() {
	/// Flag args
	f := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	root := f.String("root", "./", "Which dir to serve?")
	host := f.String("host", "", "An address to use")
	port := f.Int("port", 8080, "A port to use")
	auth := f.String("auth", "", "Authentication URL")
	exp := f.Int64("exp", 15*60, "Expire token")
	redirectPage := f.String("redirect", "/auth", "Authentication page")
	v := f.Bool("v", false, "Verbose")
	protected := f.String("protected", "", "Protected dirs")
	f.Parse(os.Args[1:])

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

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Expose-Headers", "Authorization")

		if isAllowed {
			debug(" Allowed")
			filename := path.Join(*root, r.URL.Path)
			if _, err := os.Stat(filename); os.IsNotExist(err) {
				http.Redirect(w, r, "/404.html", http.StatusSeeOther)
			} else {
				http.ServeFile(w, r, filename)
			}
		} else {
			debug(" Not allowed, redirect to:", *redirectPage)
			url := fmt.Sprintf("%s?url=%s", *redirectPage, r.URL.Path)
			http.Redirect(w, r, url, http.StatusSeeOther)
		}
	})

	addr := fmt.Sprintf("%s:%v", *host, *port)

	fmt.Println(fmt.Sprintf("%s v%s", appName, appVersion))
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
