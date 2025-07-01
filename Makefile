.PHONY: help test build check-links

help: ## Show this help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

test: ## Run tests with race detection
	go run golang.org/x/tools/cmd/goimports@latest -local $$(go list -m) -w .
	go test -race ./...

build: ## Build all binaries
	go build ./cmd/knife
	go build ./cmd/cutter
	go build ./cmd/hagane
	go build ./cmd/objls
	go build ./cmd/typels

check-links: ## Check for broken links in documentation
	./script/check-links.sh