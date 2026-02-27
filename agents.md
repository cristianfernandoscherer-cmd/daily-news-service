# AI Agent Guidelines - Daily News Service

This document provides essential information for AI agents interacting with this project. Follow these guidelines to ensure consistency and architectural integrity.

## 🏗 Project Architecture

We follow a **Clean Architecture / Hexagonal** approach to keep the domain logic isolated from infrastructure concerns.

- **`cmd/lambdas/`**: Entry points for AWS Lambda functions. Should be minimal, orchestrating dependencies.
- **`internal/`**: Core application logic.
    - **`news/models/`**: Domain entities (e.g., `Article`).
    - **`news/service/`**: Business logic / Use cases.
    - **`news/handler/`**: Adapters for Lambda triggers.
    - **`news/repository/`**: Interfaces for data persistence.
    - **`infra/database/`**: Implementations of repositories (e.g., DynamoDB).
- **`pkg/`**: Shared infrastructure primitives and utilities.
    - **`llm/`**: LLM client interfaces and implementations (Bedrock).
    - **`database/`**: Database client wrappers.

## 🏷 Naming Conventions & Ambiguity

To avoid confusion between different layers:

- **LLM Client**: Use descriptive names like `ResumeArticle` or `CallLLM`. Avoid generic names like `Summarize` at this level.
- **Services**: Use business-oriented names. For example, use `ResumeArticles` when adding AI-generated summaries and scores to articles.
- **Avoid Ambiguity**: Ensure method names clearly indicate their scope (e.g., `UpdateNewsHandler.Handle` vs `FetcherService.Fetch`).

## 💉 Dependency Injection

- We use a custom `Container` located in `internal/di/container.go`.
- All dependencies should be wired in `di.Build`.
- Avoid global state; inject dependencies into handlers and services via constructors.

## 🛠 Tech Stack & Tools

- **Language**: Go 1.22+
- **Infrastructure**: AWS SAM (CloudFormation), DynamoDB, AWS Lambda.
- **LLM**: Amazon Bedrock (Claude 3 Haiku).
- **Local Dev**: LocalStack for DynamoDB emulation. Use `Makefile` commands for common tasks (`make infra-up`, `make migrate`, `make local`).

## 🧠 Strategic Intent

- **Goal**: Automate tech news curation for software engineers.
- **Flow**: Fetch (RSS/API) -> Enrich (LLM) -> Rank (Logic) -> Store (DynamoDB).
- **Skills**: "Skills" are specialized services that add value to the curated news (e.g., Daily Briefing narrations).

## 🚀 Running Locally

To run the project in a development environment, follow these steps:

1. **Setup Environment**: `cp .env.example .env`
2. **Start Infrastructure**: `make infra-up` (Starts LocalStack for DynamoDB).
3. **Run Migrations**: `make migrate` (Creates the `articles` table).
4. **Execute Pipeline**: `make local` (Runs the Go application locally).

Verification:
```bash
aws --endpoint-url=http://localhost:4566 dynamodb scan --table-name articles
```

## 🚀 Execution Rules

1. **Testability**: Always write unit tests with mocks for external dependencies (DB, LLM).
2. **Logging**: Use `log/slog` for structured logging.
3. **Configuration**: Use `internal/config` for environment variables.
4. **Error Handling**: Wrap errors with context (`fmt.Errorf("context: %w", err)`).
