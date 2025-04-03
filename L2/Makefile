PROJECT := ~/sources/geekBrains/backend2/L2
GIT_COMMIT := $(shell git rev-parse HEAD)
VERSION := 4.5.6
APP_NAME := k8s-go-app
BIN_PATH := $(shell go env GOPATH)/bin

all:
	run

run:
	rm -f build/* && \
	go build -o ./build/$(APP_NAME) -ldflags "-X $(APP_NAME)/version.Version=$(VERSION)" $(PROJECT)/cmd/main.go && \
	cp ./build/$(APP_NAME) $(BIN_PATH) && \
	$(BIN_PATH)/$(APP_NAME)