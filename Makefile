.PHONY: dagger/check 
dagger/check:
	@dagger call check --dir .

.PHONY: dagger/examples
dagger/examples:
	@dagger call examples --dir .

.PHONY: dagger/generate
dagger/generate:
	@dagger call generate --dir .

.PHONY: dagger/lint
dagger/lint:
	@dagger call lint --dir .

.PHONY: dagger/publish
dagger/publish:
	@dagger call publish --dir .

.PHONY: dagger/test
dagger/test:
	@dagger call tests --dir .

.PHONY: all
all: generate lint examples test ## Run all the checks

.PHONY: clean
clean: dev/down ## Clean the project
	@rm -rf ./tmp/certs

.PHONY: dev/up
dev/up: ## Start the development environment
	@go run ./tools/generate-certs
	@docker-compose up -d

.PHONY: dev/down
dev/down: ## Stop the development environment
	@docker-compose down

.PHONY: lint
lint: ## Lint the code
	@go run github.com/golangci/golangci-lint/cmd/golangci-lint@v1.55 run ./...

.PHONY: test
test: dev/up ## Perform tests
	@go test ./...

.PHONY: generate
generate: ## Generate files
	@go generate ./...

.PHONY: publish
publish: dagger/publish ## Publish with tag on git, docker hub, etc.
	@git tag ${TAG} && git push origin ${TAG}

.PHONY: help
help:
	@dagger functions