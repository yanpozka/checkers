BINARY=api

VERSION=0.1-$(shell git rev-parse HEAD)
BUILD_TIME=$(shell date +%FT%T%z)

LDFLAGS=-ldflags '-extldflags "-lm -lstdc++ -static" -X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME}'


build: # build
	make clean
	CGO_ENABLED=0 go build -v -tags netgo ${LDFLAGS} -o ${BINARY}

clean:
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi

install:
	go install
