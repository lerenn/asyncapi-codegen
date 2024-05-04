.PHONY: dagger/check 
dagger/check: ## Run all the checks before commit on Dagger
	@dagger call check --dir .

.PHONY: dagger/examples
dagger/examples: ## Run all the examples on Dagger
	@dagger call examples --dir .

.PHONY: dagger/check-generation
dagger/check-generation: ## Check files are generated on Dagger
	@dagger call check-generation --dir .

.PHONY: dagger/lint
dagger/lint: ## Lint the code on Dagger
	@dagger call lint --dir .

.PHONY: dagger/publish
dagger/publish: ## Publish with tag on git, docker hub, etc. on Dagger
	@dagger call publish --dir . --tag ${TAG}

.PHONY: dagger/test
dagger/test: ## Perform tests on Dagger
	@dagger call tests --dir .

.PHONY: dagger/help
dagger/help: ## Display Dagger module help message
	@dagger functions