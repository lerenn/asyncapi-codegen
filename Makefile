include tools/make/help.mk

.PHONY: check
check: check-generation lint test ## Run all the checks locally

.PHONY: check-generation
check-generation: ## Check files are generated locally
	@sh ./scripts/check-generation.sh

.PHONY: clean
clean: local-env/stop ## Clean the project locally
	@rm -rf ./tmp/certs

.PHONY: generate
generate: ## Generate files locally
	@go generate ./...

.PHONY: lint
lint: ## Lint the code locally
	@go run github.com/golangci/golangci-lint/cmd/golangci-lint@v1.62.0 run ./...

# Dagger targets
.PHONY: dagger/check
dagger/check: ## Run all checks via Dagger
	@dagger call check --src-dir .

.PHONY: dagger/check-generation
dagger/check-generation: ## Check files are generated via Dagger
	@dagger call check-generation --src-dir .

.PHONY: dagger/develop
dagger/develop: ## Develop Dagger files
	@dagger develop

.PHONY: dagger/examples
dagger/examples: ## Run all examples via Dagger
	@dagger call examples --src-dir .

.PHONY: dagger/lint
dagger/lint: ## Lint the code via Dagger
	@dagger call lint --src-dir .

.PHONY: dagger/publish
dagger/publish: ## Publish with tag on git, docker hub, etc. via Dagger
	@dagger call publish --src-dir . --ssh-dir=~/.ssh

.PHONY: dagger/test
dagger/test: ## Perform tests via Dagger
	@dagger call test --src-dir .

.PHONY: dagger/all
dagger/all: dagger/check-generation dagger/lint dagger/test dagger/examples ## Run all CI checks via Dagger (same as GitHub workflow)

.PHONY: local-env/start
local-env/start: ## Start the local environment
	@go run ./tools/generate-certs
	@docker compose up -d

.PHONY: local-env/stop
local-env/stop: ## Stop the local environment
	@docker compose stop

.PHONY: local-env/teardown
local-env/teardown: ## Kill containers and delete volumes
	@docker compose down --volumes

.PHONY: publish
publish: dagger/publish ## Publish with tag on git, docker hub, etc. locally
	@git tag ${TAG} && git push origin ${TAG}

.PHONY: test
test: local-env/start local-env/wait ## Perform tests locally
	@go test ./...

.PHONY: local-env/wait
local-env/wait: ## Wait for services to be ready
	@echo "Waiting for Kafka services to be ready..."
	@for service in kafka kafka-tls kafka-tls-basic-auth; do \
		echo "Waiting for $$service..."; \
		timeout=180; \
		if [ "$$service" = "kafka" ]; then \
			port=9092; \
		elif [ "$$service" = "kafka-tls" ]; then \
			port=9094; \
		else \
			port=9096; \
		fi; \
		while [ $$timeout -gt 0 ]; do \
			if nc -z localhost $$port 2>/dev/null; then \
				if [ "$$service" = "kafka" ]; then \
					if docker compose exec -T $$service kafka-broker-api-versions --bootstrap-server localhost:$$port >/dev/null 2>&1 || \
					   docker compose exec -T $$service kafka-topics --bootstrap-server localhost:$$port --list >/dev/null 2>&1; then \
						echo "$$service is ready!"; \
						break; \
					fi; \
				else \
					echo "$$service is ready (port $$port is open)!"; \
					break; \
				fi; \
			fi; \
			if [ $$((timeout % 10)) -eq 0 ]; then \
				echo "  Waiting for $$service... ($$timeout seconds remaining)"; \
			fi; \
			sleep 2; \
			timeout=$$((timeout - 2)); \
		done; \
		if [ $$timeout -le 0 ]; then \
			echo "Timeout waiting for $$service to be ready"; \
			docker compose logs $$service | tail -20; \
			exit 1; \
		fi; \
	done
	@echo "All Kafka services are ready!"
	@echo "Creating SCRAM user for kafka-tls-basic-auth..."
	@timeout=60; \
	while [ $$timeout -gt 0 ]; do \
		if docker compose exec -T kafka-tls-basic-auth kafka-configs --bootstrap-server localhost:9096 --command-config /etc/kafka/client.properties --alter --add-config 'SCRAM-SHA-512=[password=password]' --entity-type users --entity-name user >/dev/null 2>&1; then \
			echo "SCRAM user created successfully!"; \
			break; \
		fi; \
		if [ $$((timeout % 10)) -eq 0 ]; then \
			echo "  Retrying SCRAM user creation... ($$timeout seconds remaining)"; \
		fi; \
		sleep 2; \
		timeout=$$((timeout - 2)); \
	done
	@sleep 5
