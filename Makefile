build:
	go build ./...

docs:
	go generate ./main.go

install: build
	go install ./...

.PHONY: build docs install
