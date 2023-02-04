build:
	go build ./...

docs:
	go generate ./main.go

install: build
	go install ./...

testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m

.PHONY: build docs install testacc
