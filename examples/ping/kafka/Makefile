EXAMPLE_NAME := ping-kafka
DOCKER_COMPOSE := docker-compose -p $(EXAMPLE_NAME)

.PHONY: clean
clean: ## Clean remains of run
	@$(DOCKER_COMPOSE) down

.PHONY: run
run: ## Run the example
	@$(MAKE) -C ../../.. brokers/up
	@$(DOCKER_COMPOSE) run user
	@$(DOCKER_COMPOSE) logs app

.PHONY: help
help: ## Display this help message
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_\/-]+:.*?## / {printf "\033[34m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST) | \
		sort | \
		grep -v '#'
