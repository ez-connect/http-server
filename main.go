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
	name    = "HTTP Server"
	version = "0.4.0"
)

func main() {
	/// Flag args
	var (
		root string
		fs   = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
		addr = fs.String("a", ":8080", "An address to use")

		// Authentication
		protected    = fs.String("protected", "", "Protected dirs, separated by a comma")
		auth         = fs.String("auth", "", "Authentication URL")
		exp          = fs.Int64("exp", 15*60, "Expire token")
		redirectPage = fs.String("redirect", "/auth", "Authentication page")

		// Log
		isVerbose = fs.Bool("v", false, "Verbose")
	)

	fs.Usage = usageFor(fs)
	// nolint:errcheck
	fs.Parse(os.Args[1:])

	// No args, print usage then exit
	if fs.NArg() < 1 {
		fs.Usage()
		os.Exit(1)
	}

	root = fs.Arg(0)

	/// Server
	http.HandleFunc("/", handle(root, *protected, *auth, *redirectPage, *exp, *isVerbose))

	fmt.Println(name, version)
	fmt.Println("Protected dirs:", *protected)
	fmt.Println("Authentication URL:", *auth)
	fmt.Println("Expired in:", *exp)
	fmt.Println("Redirect page:", *redirectPage)
	fmt.Println("Serve", root, "at:", *addr)

	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func handle(root, protected, auth, redirectPage string, exp int64, isVerbose bool) func(http.ResponseWriter, *http.Request) {
	var (
		protectedDirs = strings.Split(protected, ",")
		allowedToken  = map[string]int64{}
		debug         func(v ...interface{})
	)

	if isVerbose {
		debug = log.Println
	} else {
		debug = func(v ...interface{}) {}
	}

	return func(w http.ResponseWriter, r *http.Request) {
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
					isAllowed = true
					break
				}

				// Exp or
				debug(" Auth via:", auth)
				client := &http.Client{}
				req, err := http.NewRequest(http.MethodPost, auth, nil)
				if err != nil {
					break
				}

				req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
				res, err := client.Do(req)
				if err == nil && res.StatusCode == http.StatusOK {
					isAllowed = true

					// Add or update the token
					debug(" Add or update the token")
					allowedToken[token] = now + exp
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
			filename := path.Join(root, r.URL.Path)
			if _, err := os.Stat(filename); os.IsNotExist(err) {
				filename = path.Join(root, "404.html")
			}
			http.ServeFile(w, r, filename)
		} else {
			debug(" Not allowed, redirect to:", redirectPage)
			url := fmt.Sprintf("%s?url=%s", redirectPage, r.URL.Path)
			http.Redirect(w, r, url, http.StatusSeeOther)
		}
	}
}

// CLI help
func usageFor(fs *flag.FlagSet) func() {
	return func() {
		fmt.Println(name, fmt.Sprintf("v%s", version), "- Generate a gokit service")
		fmt.Println("USAGE:", "http-server", "[OPTIONS] path/to/your/service")
		fmt.Println("\nOPTIONS")
		fs.PrintDefaults()
	}
}
