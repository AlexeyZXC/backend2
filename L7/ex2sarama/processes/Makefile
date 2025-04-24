#  make -f processes/Makefile build_container
VERSION := 1.0.8

build_container:
	cd ..
	docker build --build-arg VERSION=$(VERSION) -t processes:$(VERSION) -f processes/Dockerfile .