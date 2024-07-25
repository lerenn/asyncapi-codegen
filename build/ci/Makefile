include ../../tools/make/help.mk

SRC_DIR := ../..

.PHONY: check 
check: ## Run all the checks before commit on Dagger
	@dagger call check --src-dir $(SRC_DIR)

.PHONY: run-examples
run-examples: ## Run all the examples on Dagger
	@dagger call run-examples --src-dir $(SRC_DIR)

.PHONY: check-generation
check-generation: ## Check files are generated on Dagger
	@dagger call check-generation --src-dir $(SRC_DIR)

.PHONY: lint
lint: ## Lint the code on Dagger
	@dagger call lint --src-dir $(SRC_DIR)

.PHONY: publish   
publish: ## Publish with tag on git, docker hub, etc. on Dagger
	@dagger call publish --src-dir $(SRC_DIR) --ssh-dir=~/.ssh

.PHONY: test
test: ## Perform tests on Dagger
	@dagger call tests --src-dir $(SRC_DIR)