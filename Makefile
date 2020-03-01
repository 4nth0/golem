COMMIT=$(shell git rev-parse HEAD)
VERSION=$(shell git describe --tags --exact-match --always)
DATE=$(shell date +'%FT%TZ%z')

golem:
	CGO_ENABLED=0 go build -v -a -tags=netgo \
  	-ldflags '-s -w -extldflags "-static" -X main.Version=$(VERSION)'
