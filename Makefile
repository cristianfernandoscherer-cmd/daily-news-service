BINARY_DIR := .aws-sam/build
ENV_FILE    := .env

# ─── Build ────────────────────────────────────────────────────────────────────
.PHONY: build
build:
	@echo "→ Building all Lambda functions..."
	sam build --template iac_template.yaml

.PHONY: build-FeedParserFunction
build-FeedParserFunction:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -tags lambda.norpc -o $(ARTIFACTS_DIR)/bootstrap ./cmd/lambdas/feed-parser/

# ─── Local Development ────────────────────────────────────────────────────────
.PHONY: local
local:
	@echo "→ Running local development mode..."
	@if [ -f $(ENV_FILE) ]; then \
		go run ./cmd/local/main.go || (echo "❌ Pipeline failed. Check if your $(ENV_FILE) has all required variables (see .env.example)"; exit 1); \
	else \
		echo "⚠️  No .env file found. Copy .env.example to .env first."; \
		exit 1; \
	fi

.PHONY: sam-local
sam-local: build
	@echo "→ Starting SAM local API (requires Docker + LocalStack)..."
	sam local invoke FeedParserFunction \
		--template .aws-sam/build/template.yaml \
		--env-vars env.local.json \
		--docker-network daily-news-service_default

# ─── Infrastructure: LocalStack ─────────────────────────────────────────────
.PHONY: infra-up
infra-up:
	@echo "→ Starting LocalStack..."
	docker-compose up -d localstack
	@echo "→ Waiting for LocalStack to be healthy..."
	@until curl -sf http://localhost:4566/_localstack/health | grep -qE '"dynamodb": "(running|available)"'; do \
		sleep 2; echo "  still waiting..."; done
	@echo "✅  LocalStack is ready."

.PHONY: infra-down
infra-down:
	docker-compose down -v

# ─── Database ────────────────────────────────────────────────────────────────
.PHONY: migrate
migrate:
	@echo "→ Running migrations..."
	go run ./.scripts/migrate/migrate.go

.PHONY: seed
seed:
	@echo "→ Seeding data..."
	go run ./.scripts/seed/seed.go

# ─── Tests ───────────────────────────────────────────────────────────────────
.PHONY: test
test:
	go test ./... -v -cover

.PHONY: test-coverage
test-coverage:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	@echo "→ Coverage report: coverage.html"

# ─── Code Quality ─────────────────────────────────────────────────────────────
.PHONY: lint
lint:
	golangci-lint run ./...

.PHONY: fmt
fmt:
	gofmt -s -w .
	goimports -w .

.PHONY: tidy
tidy:
	go mod tidy

# ─── Deploy ──────────────────────────────────────────────────────────────────
.PHONY: deploy-dev
deploy-dev: build
	sam deploy \
		--template-file .aws-sam/build/template.yaml \
		--stack-name daily-news-service-dev \
		--parameter-overrides Environment=dev \
		--capabilities CAPABILITY_IAM CAPABILITY_NAMED_IAM \
		--no-confirm-changeset \
		--resolve-s3

.PHONY: deploy-prod
deploy-prod: build
	sam deploy \
		--template-file .aws-sam/build/template.yaml \
		--stack-name daily-news-service-prod \
		--parameter-overrides Environment=prod \
		--capabilities CAPABILITY_IAM CAPABILITY_NAMED_IAM \
		--resolve-s3

# ─── Utilities ───────────────────────────────────────────────────────────────
.PHONY: clean
clean:
	rm -rf .aws-sam coverage.out coverage.html

.PHONY: validate
validate:
	sam validate --template iac_template.yaml --lint

.PHONY: logs
logs:
	sam logs --name FeedParserFunction --stack-name daily-news-service-dev --tail

.PHONY: help
help:
	@echo ""
	@echo "Daily News Service — Available Commands"
	@echo "────────────────────────────────────────"
	@echo "  make build          Build all Lambda functions via SAM"
	@echo "  make local          Run the full pipeline locally"
	@echo "  make infra-up       Start LocalStack (DynamoDB emulation)"
	@echo "  make infra-down     Stop LocalStack"
	@echo "  make migrate        Create DynamoDB tables in LocalStack"
	@echo "  make seed           Seed sample data"
	@echo "  make test           Run all tests"
	@echo "  make lint           Run golangci-lint"
	@echo "  make deploy-dev     Deploy to dev environment"
	@echo "  make deploy-prod    Deploy to production"
	@echo "  make validate       Validate SAM template"
	@echo "  make clean          Remove build artifacts"
	@echo ""
