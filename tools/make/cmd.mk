.PHONY: check
check: local/check ## Run all the checks

.PHONY: clean
clean: local/clean ## Clean the project

.PHONY: lint
lint: local/lint ## Lint the code

.PHONY: test
test: local/test ## Perform tests

.PHONY: generate
generate: local/generate ## Generate files

.PHONY: check-generation
check-generation: local/check-generation ## Check files are generated

.PHONY: publish
publish: local/publish ## Publish with tag on git, docker hub, etc.
