.PHONY: help
## help: prints this help message
help:
	@echo "Usage:"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

.PHONY: generate
## generate: runs `go generate`
generate:
	@go generate ./...

.PHONY: build
## build: builds server
build:
	@cd app &&\
	 go build -v -mod=vendor -o bin/stewart &&\
	 cp config.yml build/

.PHONY: vendor
## vendor: runs `go mod vendor`
vendor:
	@go mod vendor

.PHONY: test
## test: runs `go test`
test:
	@go test -mod=vendor ./...

.PHONY: lint
## lint: runs `golangci-lint`
lint:
	@golangci-lint run ./...

.PHONY: run
## run: runs app locally (don't forget to set all required environment variables)
run:
	@go run -v -mod=vendor cmd/bot/main.go --debug ${ARGS}

