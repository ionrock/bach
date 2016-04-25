.PHONY: clean install copy

SOURCEDIR=.
SOURCES := $(shell find $(SOURCEDIR) -path ./docker -prune -o -name '*.go')

VERSION=1.0.0
BUILD_TIME=`date +%FT%T%z`

LDFLAGS=
# LDFLAGS=-ldflags "-X github.com/ariejan/roll/core.Version=${VERSION} -X github.com/ariejan/roll/core.BuildTime=${BUILD_TIME}"

DOCKER_MACHINE_NAME=default
DOCKER_PREP=eval `docker-machine env $(DOCKER_MACHINE_NAME)`
DOCKER_IP=`docker-machine ip $(DOCKER_MACHINE_NAME)`


here: here.go
	go build ${LDFLAGS} here.go

clean:
	rm -f here
