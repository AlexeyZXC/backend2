PROJECT := /home/user/sources/geekBrains/backend2/L2/k8s-go-app

GIT_COMMIT := $(shell git rev-parse HEAD)

VERSION := 4.5.6
APP_NAME := k8s-go-app
USERNAME := user1

BIN_PATH := $(shell go env GOPATH)/bin

all:
	run

run:
	go install -ldflags="-X $(PROJECT)/version.Version=$(VERSION) -X $(PROJECT)/version.Commit=$(GIT_COMMIT)" $(APP_NAME)

build_container:
	docker build --build-arg=GIT_COMMIT=$(GIT_COMMIT) --build-arg=VERSION=$(VERSION) --build-arg=PROJECT=$(PROJECT) -t $(APP_NAME):$(VERSION) .

apply_deploy:
	kubectl apply -f deployment.yaml

apply_service:
	kubectl apply -f service.yaml

apply_ingress:
	kubectl apply -f ingress2.yaml	