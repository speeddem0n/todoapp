help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

build: ## Build web todo app.
	docker compose build todoapp

run: ## Start todo music app and postgres.
	docker compose up todoapp -d

down: ## Stop all services.
	docker compose down

test: ## Launch all unit tests.
	go test -v ./...

migarateup: ## Up migrations
	goose postgres -dir "scripts/migrations/migration" "host=localhost port=5439 user=postgres database=todoapp password=postgres sslmode=disable" up

migratedown: ## Down migrations
	goose postgres -dir "scripts/migrations/migration" "host=localhost port=5439 user=postgres database=todoapp password=postgres sslmode=disable" down

swag: ## Generate swagger docs
	swag init -d cmd/,internal/handler/,internal/models/ --output ./docs --parseInternal