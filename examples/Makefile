.PHONY: clean
clean: ## Clean remains of run
	@docker-compose -f helloworld/docker-compose.yml down 
	@docker-compose -f ping/docker-compose.yml down 

.PHONY: run
run: ## Run the examples
	@docker-compose -f helloworld/docker-compose.yml run client
	@docker-compose -f helloworld/docker-compose.yml logs app
	@docker-compose -f ping/docker-compose.yml run client
	@docker-compose -f ping/docker-compose.yml logs server

.PHONY: help
help: ## Display this help message
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_\/-]+:.*?## / {printf "\033[34m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST) | \
		sort | \
		grep -v '#'
