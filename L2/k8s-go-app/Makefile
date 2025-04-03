PROJECT := ~/sources/geekBrains/backend2/L2/k8s-go-app
GIT_COMMIT := $(shell git rev-parse HEAD)
VERSION := 4.5.6
APP_NAME := k8s-go-app
BIN_PATH := $(shell go env GOPATH)/bin

all:
	run

run:
	go install -ldflags="-X $(PROJECT)/version.Version=$(VERSION) -X $(PROJECT)/version.Commit=$(GIT_COMMIT)" $(APP_NAME)