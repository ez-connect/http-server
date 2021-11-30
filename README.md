# HTTP Server

A simple http server to serve static resource files. It's simple and hackable enough to be used for testing, local development, and learning.

## Installation

```
go install github.com/ez-connect/http-server
```

Or download from [releases](https://github.com/ez-connect/http-server/releases) section

## Usage

```
http-server [OPTIONS] path/to/dir

OPTIONS
  -a string
        An address to use (default ":8080")
  -auth string
        Authentication URL
  -exp int
        Expire token (default 900)
  -protected string
        Protected dirs
  -redirect string
        Authentication page (default "/auth")
  -v    Verbose
```

## Example

Serve the current dir at `8080`

```bash
http-server .
```

Serve the current dir, `/document` is protected, and `http://localhost:1337/auth` is the authentication endpoint

```bash
http-server -auth http://localhost:1337/auth -protected /document .
```

Serve `./public` dir, `/document` & `/blog` are protected, `http://localhost:1337/auth` is the authentication endpoint, and `/login` is your login page.

```bash
http-server -auth http://localhost:1337/auth -redirect /login -protected /document /blog ./public
```

> AUTHENTICATION
> - Use cookie-based authentication for your static page
> - Use token-based authentication for your authentication endpoin
