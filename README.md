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

- `-d` Which dir to serve? `./public` if it exist or the working dir
- `-a` an address to use, defaults to `localhost`
- `-p` a port to use, defaults to `5000`
- `-v` Show the app version