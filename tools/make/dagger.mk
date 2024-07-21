DAGGER_CMD := dagger -vvv

.PHONY: dagger/check 
dagger/check: ## Run all the checks before commit on Dagger
	@$(DAGGER_CMD) call check --src-dir .

.PHONY: dagger/examples
dagger/examples: ## Run all the examples on Dagger
	@$(DAGGER_CMD) call examples --src-dir .

.PHONY: dagger/check-generation
dagger/check-generation: ## Check files are generated on Dagger
	@$(DAGGER_CMD) call check-generation --src-dir .

.PHONY: dagger/lint
dagger/lint: ## Lint the code on Dagger
	@$(DAGGER_CMD) call lint --src-dir .

.PHONY: dagger/publish   
dagger/publish: ## Publish with tag on git, docker hub, etc. on Dagger
	@$(DAGGER_CMD) call publish --src-dir . --ssh-dir=~/.ssh

.PHONY: dagger/test
dagger/test: ## Perform tests on Dagger
	@$(DAGGER_CMD) call tests --src-dir .

.PHONY: dagger/help
dagger/help: ## Display Dagger module help message
	@$(DAGGER_CMD) functions