PKG := "github.com/hapoon/kiku"
PKG_LIST := $(shell go list ${PKG}/...)
VERSION := 0.1.0

tag:
	@git tag ${VERSION}
	@git tag|grep -v ^v

.DEFAULT_GOAL := check
init:
	@go install golang.org/x/lint/golint@latest

check: lint vet race ## Check project

lint: ## Lint the files
	@golint -set_exit_status ${PKG_LIST}

vet: ## Vet the files
	@go vet ${PKG_LIST}

test: ## Run tests
	@go test -short ${PKG_LIST}

cover: ## Run tests and get coverage
	@go test -cover ${PKG_LIST}

race: ## Run tests with data race detector
	@go test -race ${PKG_LIST}

benchmark: ## Run benchmarks
	@go test -run="-" -bench=".*" ${PKG_LIST}

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' ${MAKEFILE_LIST} | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
