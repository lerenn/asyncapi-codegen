.PHONY: all
all: generate lint examples test ## Run all the checks

.PHONY: clean
clean: clean ## Clean the project

.PHONY: lint
lint: local/lint ## Lint the code

.PHONY: test
test: local/test ## Perform tests

.PHONY: generate
generate: local/generate ## Generate files

.PHONY: publish
publish: local/publish ## Publish with tag on git, docker hub, etc.
