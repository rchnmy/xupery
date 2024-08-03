APP=xupery
REP=github.com/rchnmy/$(APP)

.PHONY: *

all: build

init:
	@go mod init $(REP)
	@go mod tidy

update:
	@go get -u ./...

build: export GOOS=linux
build: export GOARCH=amd64
build: export CGO_ENABLED=0
build: export GOMEMLIMIT=800MiB
build:
	@go build -trimpath -ldflags="-s -w" -o $(APP) ./cmd

clean:
	@-go clean
	@-rm -rf go.mod go.sum pgdata/ $(APP)
