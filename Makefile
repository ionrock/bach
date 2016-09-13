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

all: present we toconfig cluster

install:
	go install ./cmd/...

present: $(SOURCES)
	go build -o present ${LDFLAGS} ./cmd/present/

cluster: $(SOURCES)
	go build -o cluster ${LDFLAGS} ./cmd/cluster/


clean:
	rm -f present
	rm -f we

present-example: present
	./present -s ./present.sh echo 'Hello World!'


we: $(SOURCES)
	go build -o we ${LDFLAGS} ./cmd/we/

we-example: we
	./we -e example_env.yml echo 'Hello World!'

toconfig: $(SOURCES)
	go build -o toconfig ${LDFLAGS} ./cmd/toconfig/

toconfig-example: toconfig
	./toconfig -t example.conf.tmpl -c example.conf cat example.conf

test:
	go test
