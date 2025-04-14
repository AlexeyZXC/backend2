VERSION := 1.0.5
PIPOMASTER := 1

build_container:
	docker build --build-arg VERSION=$(VERSION) --build-arg PIPOMASTER=$(PIPOMASTER) -t pipo:$(VERSION) .