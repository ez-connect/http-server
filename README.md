# HTTP Server

A simple http server to serve static resource files. It's simple and hackable enough to be used for testing, local development, and learning.

## Installation

```
go install github.com/ez-connect/http-server
```

Or download from [releases](https://github.com/ez-connect/http-server/releases) section

## Usage

```
http-server [-root <path/to/dir>] [-host <host_name>] [-port <port_number>] [-redirect <path/to/redirect/page>] [-v] [protected <path/to/protect/dir...>]
```

- `-root` Which dir to serve? Defaults to `./`
- `-host` host name, defaults to empty
- `-port` a port to use, defaults to `8080`
- `-auth` Authentication URL, defaults to empty
- `-redirect` authentication page, defaults to `/auth`
- `-v` Verbose ouput
- `-protected` protected dirs, defaults to empty

## Example

Serve the current dir at `8080`

```bash
http-server
```

Serve the current dir, `/document` is protected, and `http://localhost:1337/auth` is the authentication endpoint

```bash
http-server -auth http://localhost:1337/auth -protected /document
```

Serve `./public` dir, `/document` & `/blog` are protected, `http://localhost:1337/auth` is the authentication endpoint, and `/login` is your login page.

```bash
http-server -root ./public -auth http://localhost:1337/auth -redirect /login -protected /document /blog
```

> AUTHENTICATION
> - Use cookie-based authentication for your static page
> - Use token-based authentication for your authentication endpoin
