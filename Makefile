SHELL := bash
.DEFAULT_GOAL := run

# App
name := http-server
version := 0.4.0
platforms := windows linux darwin
arch := amd64
entryPoint := main.go

# Registry
registry := docker.io
registryRepo := ezconnect

-include .makerc

init:
	@go install github.com/c9s/gomon@latest
	@make -s tidy

fmt:
	@go fmt ./...

lint:
	@golangci-lint run --fix ./...

# Build and exec instead of @go run $(entryPoint)
run:
	@go build -o dist/$(name) $(entryPoint)
	@dist/$(name) public

watch:
	@gomon main.go -- make -s run

test:
	@go test -v ./... -cover

test-clean:
	@go clean -testcache

build:
	@rm -rf dist
	@for p in $(platforms); do \
		echo dist/$(name)-$$p; \
		GOOS=$$p GOARCH=$(arch) go build -ldflags="-s -w" -o dist/$(name)-$$p $(entryPoint); \
		pushd dist > /dev/null; \
		tar -zcvf $(name)-$$p.tar.gz $(name)-$$p; \
		popd > /dev/null; \
	done

oci:
	@buildah bud -f ci/Dockerfile -t $(name):$(version)
	@buildah bud -f ci/Dockerfile-alpine -t $(name):$(version)-alpine
	@buildah bud -f ci/Dockerfile-hugo -t $(name):$(version)-hugo
ifneq ($(and $(registryUsername),$(registryPassword)),)
	@buildah login -u $(registryUsername) -p $(registryPassword) $(registry)
	buildah push $(name):$(version) $(registry)/$(registryRepo)/$(name):$(version)
	buildah push $(name):$(version)-alpine $(registry)/$(registryRepo)/$(name):$(version)-alpine
	buildah push $(name):$(version)-hugo $(registry)/$(registryRepo)/$(name):$(version)-hugo
endif
