.PHONY: clean install copy

SOURCEDIR=.
SOURCES := $(shell find $(SOURCEDIR) -path ./docker -prune -o -name '*.go')

BINDIR=bin

VERSION=1.0.0
BUILD_TIME=`date +%FT%T%z`

LDFLAGS=
# LDFLAGS=-ldflags "-X github.com/ariejan/roll/core.Version=${VERSION} -X github.com/ariejan/roll/core.BuildTime=${BUILD_TIME}"

DOCKER_MACHINE_NAME=default
DOCKER_PREP=eval `docker-machine env $(DOCKER_MACHINE_NAME)`
DOCKER_IP=`docker-machine ip $(DOCKER_MACHINE_NAME)`

GLIDE=$(GOPATH)/bin/glide

install: $(GLIDE)
	go install ./cmd/...

all: $(BINDIR)/present $(BINDIR)/we $(BINDIR)/toconfig $(BINDIR)/cluster

$(GLIDE):
	go get github.com/Masterminds/glide
	glide i

$(BINDIR)/present-$(GOOS)-$(GOARCH): $(SOURCES)
	go build -o $(BINDIR)/present-$(GOOS)-$(GOARCH) $(LDFLAGS) ./cmd/present/

$(BINDIR)/cluster-$(GOOS)-$(GOARCH): $(SOURCES)
	go build -o $(BINDIR)/cluster-$(GOOS)-$(GOARCH) $(LDFLAGS) ./cmd/cluster/

$(BINDIR)/we-$(GOOS)-$(GOARCH): $(SOURCES)
	go build -o $(BINDIR)/we-$(GOOS)-$(GOARCH) $(LDFLAGS) ./cmd/we/

$(BINDIR)/toconfig-$(GOOS)-$(GOARCH): $(SOURCES)
	go build -o $(BINDIR)/toconfig-$(GOOS)-$(GOARCH) $(LDFLAGS) ./cmd/toconfig/

$(BINDIR)/bach-$(GOOS)-$(GOARCH): $(SOURCES)
	go build -o $(BINDIR)/bach-$(GOOS)-$(GOARCH) $(LDFLAGS) ./cmd/bach/


clean:
	rm -f $(BINDIR)/*

present-example: present
	./present -s ./present.sh echo 'Hello World!'

we-example: we
	./we -e example_env.yml echo 'Hello World!'

toconfig-example: toconfig
	./toconfig -t example.conf.tmpl -c example.conf cat example.conf

test:
	go test


build-example:
	GOOS=linux GOARCH=amd64 go build -o example/cluster  ${LDFLAGS} ./cmd/cluster/
	GOOS=linux GOARCH=amd64 go build -o example/we       ${LDFLAGS} ./cmd/we/
	GOOS=linux GOARCH=amd64 go build -o example/toconfig ${LDFLAGS} ./cmd/toconfig/
	GOOS=linux GOARCH=amd64 go build -o example/bach     ${LDFLAGS} ./cmd/bach/

	docker-compose -f example/docker-compose.yml build

build-all: $(GLIDE) $(SOURCES)
	glide i
	for CLIAPP in we toconfig bach ; do \
	  for GOOS in linux darwin windows ; do \
	    for GOARCH in amd64 386 ; do \
	       GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o bin/we-$(GOOS)-$(GOARCH)  ${LDFLAGS} ./cmd/cluster/ ; \
	       GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o bin/toconfig-$(GOOS)-$(GOARCH)  ${LDFLAGS} ./cmd/cluster/ ; \
	       GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o bin/bach-$(GOOS)-$(GOARCH)  ${LDFLAGS} ./cmd/cluster/ ; \
	    done \
	  done \
	done
