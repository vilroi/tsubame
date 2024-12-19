all: check build

check: 
	go vet ./...

build:
	CGO_ENABLED=0 go build -o tsubame *.go

clean: tsubame
	rm tsubame

.PHONY: check clean build
