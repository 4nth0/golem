COMMIT=$(shell git rev-parse HEAD)
VERSION=$(shell git describe --tags --always)
DATE=$(shell date +'%FT%TZ%z')

lint:
	docker run --rm -v $(pwd):/app -w /app golangci/golangci-lint:v1.49.0 golangci-lint run -v

build:
	CGO_ENABLED=0 go build -v -a -tags=netgo \
  	-ldflags '-s -w -extldflags "-static" -X main.Version=$(VERSION)'
