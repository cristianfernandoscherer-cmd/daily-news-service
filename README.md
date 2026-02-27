# Daily News Service

**Daily News Service** is a serverless application in Go that automates the curation of tech news. 
Every day at 7:00 AM, an AWS Lambda function is triggered to fetch the latest articles from sources like [Dev.to](https://dev.to/), AWS Blog, and GitHub.
Each news item is processed by Amazon Bedrock (Claude) to generate a concise summary, extract key points, and rank them by relevance.
The top 5 most important news are selected and registered in DynamoDB.

## 🚀 Architecture

The application is built with a clean, domain-driven architecture:

- **pkg/**: Infrastructure primitives (DynamoDB client, Bedrock/LLM client).
- **internal/news/**: Core domain logic (models, services for fetching, summarizing, and ranking).
- **cmd/lambdas/feed-parser/**: The singular entry point that orchestrates the entire flow.
- **iac_template.yaml**: AWS SAM template infrastructure as code.

## 🛠 Tech Stack

- **Go 1.22**
- **AWS SAM** (CloudFormation)
- **AWS Lambda** (Arm64)
- **Amazon DynamoDB** (NoSQL storage)
- **Amazon Bedrock** (Claude 3 Haiku for summaries)
- **LocalStack** (Local development & emulation)

## 🏗 Setup & Local Development

### Prerequisites

- Go 1.22
- Docker & Docker Compose
- AWS CLI & SAM CLI (optional for deployment)

### 1. Environment Setup

Copy the example environment file:
```bash
cp .env.example .env
```

### 2. Infrastructure (LocalStack)

Start LocalStack to emulate DynamoDB:
```bash
make infra-up
```

Create the articles table:
```bash
make migrate
```

### 3. Running Locally

You can run the full pipeline locally without deploying to AWS:
```bash
make local
```

Verify the results in the local DynamoDB:
```bash
aws --endpoint-url=http://localhost:4566 dynamodb scan --table-name articles
```

## 🚢 Deployment

O deploy é automatizado via **GitHub Actions** sempre que houver um push para a branch `master`.

### Configuração do GitHub Secrets
Para o deploy automático funcionar, adicione as seguintes Secrets no seu repositório GitHub:
*   `AWS_ACCESS_KEY_ID`: Sua chave de acesso AWS.
*   `AWS_SECRET_ACCESS_KEY`: Sua chave secreta AWS.

### Deploy Manual via SAM
Se preferir rodar manualmente da sua máquina:
```bash
# Ambientes disponíveis: dev, prod
sam build --template iac_template.yaml
sam deploy --config-env prod
```

## 📜 Available Commands

| Command | Description |
|---------|-------------|
| `make build` | Build the Lambda function via SAM |
| `make local` | Run the full pipeline locally |
| `make sam-local` | Run the Lambda function locally via SAM |
| `make infra-up` | Start LocalStack (DynamoDB emulation) |
| `make migrate` | Create DynamoDB tables in LocalStack |
| `make test` | Run all unit tests |
| `make deploy-prod` | Deploy to AWS Production environment |
| `make clean` | Remove build artifacts |
