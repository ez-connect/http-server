# HTTP Server

A simple http server to serve static resource files. It's simple and hackable enough to be used for testing, local development, and learning.

## Installation

```
go install github.com/ez-connect/http-server
```

## Usage

```
http-server [-d <path/to/dir>] [-a <address>] [-p <port>] [-v]
```

- `-root` Which dir to serve? Defaults to `./public`
- `-host` host name, defaults to `localhost`
- `-port` a port to use, defaults to `8080`
- `-privates` protected dirs, defaults to `/private /protected`
- `-auth` Authentication URL
- `-redirect` authentication page, defaults to `/auth`
- `-v` Show the app version
