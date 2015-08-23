VERSION=0.0.1
TARGETS_NOVENDOR=$(shell glide novendor)

all: build

bundle:
	glide install

fmt:
	@echo $(TARGETS_NOVENDOR) | xargs go fmt

test:
	go test -cover $(TARGETS_NOVENDOR)

build:
	go build -o bin/goflake-server -ldflags "-X main.version=${VERSION}" .

build-all:
	GOOS=linux GOARCH=amd64 go build -o bin/linux/amd64/goflake-server-${VERSION}/goflake-server -ldflags "-X main.version=${VERSION}" .
	GOOS=darwin GOARCH=amd64 go build -o bin/darwin/amd64/goflake-server-${VERSION}/goflake-server -ldflags "-X main.version=${VERSION}" .

dist: build-all
	cd bin/linux/amd64 && tar zcvf gflake-server-linux-amd64-${VERSION}.tar.gz goflake-server-${VERSION}
	cd bin/darwin/amd64 && tar zcvf gflake-server-darwin-amd64-${VERSION}.tar.gz goflake-server-${VERSION}

clean:
	rm -rf bin/*
