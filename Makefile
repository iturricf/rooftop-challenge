.PHONY: all clean test build

all: test clean build

test:
	cd src && go test ./...

build:
	cd src && go mod download && go mod vendor && \
	go build -o ../bin/rtchallenge main/main.go

clean:
	@rm -rf bin
