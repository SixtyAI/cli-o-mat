APP_NAME?=omat

GOOS?=darwin
GOARCH?=arm64

.DEFAULT_GOAL := help
.PHONY: help

setup_ci: ## Install/update tools needed at CI-time.
	if [ $$(wc -l < tools.ci) -gt 0 ]; then \
		cat tools.ci | while read TOOL; do \
			echo "Installing $${TOOL}..."; \
			go install $${TOOL}; \
		done; \
	fi
	pip3.11 install chainjacking

setup_workstation: setup_ci ## Install/update tools needed at dev-time.
	if [ $$(wc -l < tools.dev) -gt 0 ]; then \
		cat tools.dev | while read TOOL; do \
			echo "Installing $${TOOL}..."; \
			go install $${TOOL}; \
		done; \
	fi

setup: setup_ci setup_workstation ## Install/update all tools.

clean: ## Clean up.
	go clean

version: ## Generate version string
	echo "$(shell cat .version) (SHA $(shell git rev-parse --short HEAD), $(shell date -u +%Y-%m-%dT%H:%M:%S%z))" > cmd/.version_string

generate: version ## Generate any code, as needed.
	GOOS=${GOOS} GOARCH=${GOARCH} go generate ./...

test: generate ## Run test suite.
	go test -race -coverprofile coverage.txt ./...

just_compile:
	GOOS=${GOOS} GOARCH=${GOARCH} go build -o "${APP_NAME}.${GOARCH}" main.go

compile: generate just_compile ## Compile code.

build: test generate compile ## Run tests, generate code, compile code.

build_amd64: version generate
	GOOS=darwin GOARCH=amd64 make just_compile

build_arm64: version generate
	GOOS=darwin GOARCH=arm64 make just_compile

dist: build_amd64 build_arm64 ## Compile for all architectures, and produce a fat binary.
	lipo -create -arch arm64 ${APP_NAME}.arm64 -arch x86_64 ${APP_NAME}.amd64 -output ${APP_NAME}

gci: ## Run GCI to normalize import ordering.
	gci write . --skip-generated -s 'standard,prefix(github.com/SixtyAI),default'

fmt: gci ## Tidy code.
	gofumpt -l -e -w $$(find . -name '*.go' | xargs grep -L 'Code generated by' | cut -d: -f1)

lint: clean ## Run Go linters, without auto-fixing.
	GOOS=${GOOS} GOARCH=${GOARCH} go vet ./...
	if [ $$(wc -l < tools.ci) -gt 0 ]; then \
		grep golang.org/x/tools/go/analysis/passes tools.ci | cut -d/ -f7 | while read TOOL; do \
			GOOS=${GOOS} GOARCH=${GOARCH} go vet -vettool=$$(which $$TOOL) ./...; \
		done \
	fi
	GOOS=${GOOS} GOARCH=${GOARCH} golangci-lint run --config ./.golangci.yml
	python3.11 -m chainjacking -gt $$GITHUB_TOKEN

vulncheck: ## Check for known vulnerabilities.
	govulncheck ./...

fix: ## Run transforms to simplify/correct known undesirable patterns of code.
	golangci-lint run --config ./.golangci.yml --fix

cloc: ## Run cloc.
	cloc .
	@echo

outdated: ## Check for outdated dependencies.
	go list -u -m all

update: ## Update dependencies.
	go get -t -u ./...
	go mod tidy

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "%-20s %s\n", $$1, $$2}'
