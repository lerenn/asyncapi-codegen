.PHONY: check 
check:
	@dagger call check --dir .

.PHONY: examples
examples:
	@dagger call examples --dir .

.PHONY: generate
generate:
	@dagger call generate --dir .

.PHONY: lint
lint:
	@dagger call lint --dir .

.PHONY: publish
publish:
	@dagger call publish --dir .

.PHONY: test
test:
	@dagger call tests --dir .

.PHONY: help
help:
	@dagger functions