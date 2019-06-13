help: ## This help.
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help

migrations-up: ## sync the migrations
	migrate -source file://migrations/ -database postgres://postgres:db_pass@db/posts_service?sslmode=disable up