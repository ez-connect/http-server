GOPATH:=$(shell go env GOPATH)

name := http-server
buildDir := dist
platforms := windows linux darwin
arch := amd64
entryPoint := main.go

.PHONY: proto build

proto:
	@protoc proto/sample.proto --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative

run:
	@go run $(entryPoint)

build-debug:
	@go build -o $(buildDir)/$(name) $(entryPoint)

build:
	@for p in $(platforms); do \
		echo $(buildDir)/$(name)-$$p; \
		GOOS=$$p GOARCH=$(arch) go build -ldflags="-s -w" -o $(buildDir)/$(name)-$$p $(entryPoint); \
		pushd $(buildDir); \
		tar -zcvf $(name)-$$p.tar.gz $(name)-$$p; \
		popd; \
	done

test:
	@go test -v ./... -cover

docker:
	@docker build -t sample:latest .
