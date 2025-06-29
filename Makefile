# Basic arguments
APP_NAME=ahorro
SERVICE_NAME=transactions
INSTANCE_NAME=$(shell whoami)
AWS_REGION=eu-west-1

FULL_NAME=$(APP_NAME)-$(SERVICE_NAME)-$(INSTANCE_NAME)

# Database credentials from AWS Secrets Manager
SECRET_NAME=$(APP_NAME)-app-secrets
DB_USERNAME=$(shell aws secretsmanager get-secret-value --secret-id $(SECRET_NAME) --query 'SecretString' --output text --region $(AWS_REGION) | jq -r '.transactions_db_username')
DB_PASSWORD=$(shell aws secretsmanager get-secret-value --secret-id $(SECRET_NAME) --query 'SecretString' --output text --region $(AWS_REGION) | jq -r '.transactions_db_password')
DOMAIN_NAME=$(shell aws secretsmanager get-secret-value --secret-id $(SECRET_NAME) --query 'SecretString' --output text --region $(AWS_REGION) | jq -r '.domain_name')

# Cognito configuration (fetched from AWS Cognito by name)
COGNITO_USER_POOL_NAME=ahorro-app-stable-user-pool
COGNITO_USER_POOL_CLIENT_NAME=ahorro-app-stable-client
COGNITO_USER_POOL_ID=$(shell aws cognito-idp list-user-pools --max-results 50 --region $(AWS_REGION) --query 'UserPools[?Name==`$(COGNITO_USER_POOL_NAME)`].Id' --output text)
COGNITO_CLIENT_ID=$(shell aws cognito-idp list-user-pool-clients --user-pool-id $$(aws cognito-idp list-user-pools --max-results 50 --region $(AWS_REGION) --query 'UserPools[?Name==`$(COGNITO_USER_POOL_NAME)`].Id' --output text) --region $(AWS_REGION) --query 'UserPoolClients[?ClientName==`$(COGNITO_USER_POOL_CLIENT_NAME)`].ClientId' --output text)

# Database configuration
DB_IDENTIFIER = $(APP_NAME)-$(SERVICE_NAME)-stable-db
DB_NAME = $(APP_NAME)_$(SERVICE_NAME)_stable_db
DB_ENDPOINT = $(shell aws rds describe-db-instances --db-instance-identifier $(DB_IDENTIFIER) --query 'DBInstances[0].Endpoint.Address' --output text --region $(AWS_REGION) 2>/dev/null || echo "")
DB_PORT = $(shell aws rds describe-db-instances --db-instance-identifier $(DB_IDENTIFIER) --query 'DBInstances[0].Endpoint.Port' --output text --region $(AWS_REGION) 2>/dev/null || echo "5432")

# Main app arguments
APP_DIR=app
APP_BUILD_DIR=./build/service-handler
APP_LAMBDA_ZIP_BASE_NAME=$(SERVICE_NAME)-lambda
APP_LAMBDA_ZIP_NAME=$(APP_LAMBDA_ZIP_BASE_NAME).zip
APP_LAMBDA_HANDLER_ZIP=$(APP_BUILD_DIR)/$(APP_LAMBDA_ZIP_NAME)
APP_LAMBDA_BINARY=$(APP_BUILD_DIR)/bootstrap
APP_BINARY=$(APP_BUILD_DIR)/transactions_service

# S3 paths for different deployment types
TIMESTAMP_FILE=$(APP_BUILD_DIR)/build-timestamp
TIMESTAMP=$(shell [ -f $(TIMESTAMP_FILE) ] && cat $(TIMESTAMP_FILE) || (mkdir -p $(APP_BUILD_DIR) && echo "build-$$(date +%y%m%d-%H%M)" | tee $(TIMESTAMP_FILE)))
APP_LAMBDA_S3_BASE=s3://ahorro-artifacts/transactions
APP_LAMBDA_S3_PATH_LOCAL=$(APP_LAMBDA_S3_BASE)/$(INSTANCE_NAME)/$(APP_LAMBDA_ZIP_NAME)
APP_LAMBDA_S3_PATH_TIMESTAMP=$(APP_LAMBDA_S3_BASE)/$(TIMESTAMP)/$(APP_LAMBDA_ZIP_NAME)

# GitHub repository configuration
GITHUB_REPO=ahorro-transactions-service
GITHUB_TAG_NAME=$(TIMESTAMP)

# Schema generation
SCHEMA_TEMPLATE=schema/openapi.yml.tml
SCHEMA_OUTPUT=$(APP_DIR)/schema/openapi.yml

.PHONY: all build app-build-local app-build-lambda run package test clean deploy undeploy plan get-db-config get-db-endpoint get-db-port get-db-name show-db-config get-my-ip db-connect seed pull-postgres deploy-public-custom drop-tables generate-schema db-start db-stop db-status db-get-identifier get-cognito-token show-cognito-config git-tag upload-and-tag help

# Default target
all: build

# Help target
help:
	@echo "Ahorro Transactions Service - Available Makefile targets:"
	@echo ""
	@echo "📦 Build & Package:"
	@echo "  build                 - Build both local and Lambda binaries"
	@echo "  app-build-local       - Build local binary"
	@echo "  app-build-lambda      - Build Lambda binary (using Docker)"
	@echo "  package               - Create Lambda deployment package"
	@echo "  package-timestamp     - Create timestamped Lambda package"
	@echo "  generate-schema       - Generate OpenAPI schema from template"
	@echo ""
	@echo "🧪 Testing & Running:"
	@echo "  test                  - Run Go tests"
	@echo "  run                   - Run service locally"
	@echo "  clean                 - Clean build artifacts"
	@echo ""
	@echo "🚀 Deployment:"
	@echo "  deploy                - Deploy infrastructure and service"
	@echo "  undeploy              - Destroy infrastructure"
	@echo "  plan                  - Show Terraform plan"
	@echo "  upload                - Upload Lambda package to S3 (s3://ahorro-artifacts/transactions/\$INSTANCE_NAME/)"
	@echo "  upload-timestamp      - Upload timestamped package to S3 (s3://ahorro-artifacts/transactions/\$TIMESTAMP/)"
	@echo "  upload-and-tag        - Upload timestamped package AND create git tag (RECOMMENDED for CI/CD)"
	@echo "  git-tag               - Create and push git tag with timestamp (\$TIMESTAMP)"
	@echo ""
	@echo "🗄️  Database Operations:"
	@echo "  db-connect            - Connect to PostgreSQL database"
	@echo "  db-status             - Check database status"
	@echo "  db-start              - Start database instance"
	@echo "  db-stop               - Stop database instance (saves costs)"
	@echo "  db-quick-start        - Start DB and wait until available"
	@echo "  db-quick-stop         - Stop DB and wait until stopped"
	@echo "  seed                  - Seed database with sample data"
	@echo "  drop-tables           - Drop all tables (⚠️  DESTRUCTIVE)"
	@echo ""
	@echo "🔧 Configuration & Utilities:"
	@echo "  show-db-config        - Show database configuration"
	@echo "  show-api-url          - Show deployed API URL"
	@echo "  get-my-ip             - Get your public IP address"
	@echo "  pull-postgres         - Update PostgreSQL Docker image"
	@echo ""
	@echo "🔐 Authentication:"
	@echo "  get-cognito-token     - Get Cognito IdToken for API access"
	@echo "  show-cognito-config   - Show Cognito configuration"
	@echo ""
	@echo "💡 Examples:"
	@echo "  make build && make deploy"
	@echo "  make get-cognito-token"
	@echo "  make db-quick-start && make seed"

# Schema generation target
$(SCHEMA_OUTPUT): $(SCHEMA_TEMPLATE)
	@echo "Generating OpenAPI schema from template..."
	@mkdir -p $(dir $(SCHEMA_OUTPUT))
	@sed -e 's/$${INSTANCE_NAME}/$(INSTANCE_NAME)/g' \
	     -e 's/$${DOMAIN_NAME}/$(DOMAIN_NAME)/g' \
	     $(SCHEMA_TEMPLATE) > $(SCHEMA_OUTPUT)
	@echo "OpenAPI schema generated: $(SCHEMA_OUTPUT)"

generate-schema: $(SCHEMA_OUTPUT)

# Build and package main app
$(APP_LAMBDA_BINARY): $(shell find $(APP_DIR) -type f -name '*.go') $(SCHEMA_OUTPUT)
	@echo "Building Lambda binary using Docker (ensures compatibility)..."
	@mkdir -p $(APP_BUILD_DIR)
	@docker run \
		-v $(PWD)/$(APP_DIR):/src \
		-v $(PWD)/$(APP_BUILD_DIR):/build \
		-w /src \
		golang:1.23-alpine \
		sh -c "apk add --no-cache git ca-certificates && \
		       CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
		       go build -ldflags='-s -w -extldflags=-static' -tags netgo -a \
		       -o /build/bootstrap main.go"

$(APP_BINARY): $(APP_DIR)/main.go $(SCHEMA_OUTPUT)
	@mkdir -p $(APP_BUILD_DIR)
	cd $(APP_DIR) && go build -o ../$(APP_BINARY) main.go

$(APP_LAMBDA_HANDLER_ZIP): $(APP_LAMBDA_BINARY)
	@mkdir -p $(APP_BUILD_DIR)
	cd $(APP_BUILD_DIR) && zip $(APP_LAMBDA_ZIP_NAME) bootstrap

# Combined build and package targets
build: $(APP_BINARY) $(APP_LAMBDA_BINARY)

app-build-local: $(APP_BINARY)

app-build-lambda: $(APP_LAMBDA_BINARY)

test: $(APP_BINARY)
	cd $(APP_DIR) && go test ./...

run: app-build-local get-db-config
	DB_HOST=$(shell $(MAKE) -s get-db-endpoint) \
	DB_PORT=$(shell $(MAKE) -s get-db-port) \
	DB_NAME=$(shell $(MAKE) -s get-db-name) \
	DB_USER=$(DB_USERNAME) \
	DB_PASSWORD=$(DB_PASSWORD) \
	LOG_LEVEL=$(LOG_LEVEL) \
	./$(APP_BINARY)

package: $(APP_LAMBDA_HANDLER_ZIP)

package-timestamp: $(APP_LAMBDA_HANDLER_ZIP)

upload: $(APP_LAMBDA_HANDLER_ZIP)
	@echo "Uploading Lambda package to: $(APP_LAMBDA_S3_PATH_LOCAL)"
	aws s3 rm $(APP_LAMBDA_S3_PATH_LOCAL) --quiet || true
	aws s3 cp $(APP_LAMBDA_HANDLER_ZIP) $(APP_LAMBDA_S3_PATH_LOCAL)

upload-timestamp: $(APP_LAMBDA_HANDLER_ZIP)
	@echo "Uploading timestamped Lambda package to: $(APP_LAMBDA_S3_PATH_TIMESTAMP)"
	aws s3 cp $(APP_LAMBDA_HANDLER_ZIP) $(APP_LAMBDA_S3_PATH_TIMESTAMP)

upload-and-tag: $(APP_LAMBDA_HANDLER_ZIP)
	@echo "=== Deploying with timestamp: $(TIMESTAMP) ==="
	@echo "Uploading timestamped Lambda package to: $(APP_LAMBDA_S3_PATH_TIMESTAMP)"
	aws s3 cp $(APP_LAMBDA_HANDLER_ZIP) $(APP_LAMBDA_S3_PATH_TIMESTAMP)
	@echo "Creating Git tag: $(GITHUB_TAG_NAME)"
	@./scripts/git-tag.sh "$(GITHUB_TAG_NAME)" "$(GITHUB_REPO)" "$(AWS_REGION)" "$(SECRET_NAME)"
	@echo "=== Deployment completed successfully! ==="

git-tag:
	@echo "Creating Git tag: $(GITHUB_TAG_NAME)"
	@./scripts/git-tag.sh "$(GITHUB_TAG_NAME)" "$(GITHUB_REPO)" "$(AWS_REGION)" "$(SECRET_NAME)"

# Terraform deployment helpers

plan:
	cd deploy && \
	terraform init && \
	terraform plan \
		-var="app_name=$(APP_NAME)" \
		-var="service_name=$(SERVICE_NAME)" \
		-var="env=$(INSTANCE_NAME)"

refresh:
	cd deploy && \
	terraform init && \
	terraform refresh \
		-var="app_name=$(APP_NAME)" \
		-var="service_name=$(SERVICE_NAME)" \
		-var="env=$(INSTANCE_NAME)"

# Use this only for development purposes
deploy:
	@echo "Deploying the service with $(DB_IDENTIFIER)..."
	cd deploy && \
	terraform init && \
	terraform apply -auto-approve \
		-var="app_name=$(APP_NAME)" \
		-var="service_name=$(SERVICE_NAME)" \
		-var="env=$(INSTANCE_NAME)"

undeploy:
	@echo "Undeploying the service with $(DB_IDENTIFIER)..."
	cd deploy && \
	terraform init && \
	terraform destroy -auto-approve \
		-var="app_name=$(APP_NAME)" \
		-var="service_name=$(SERVICE_NAME)" \
		-var="env=$(INSTANCE_NAME)"

show-api-url:
	@cd deploy && terraform output -raw api_url

# Database configuration helpers
get-db-config:
	@echo "Fetching database configuration from AWS RDS..."
	@if [ -z "$(DB_ENDPOINT)" ]; then \
		echo "Error: Unable to get database endpoint. Check if DB instance '$(DB_IDENTIFIER)' exists."; \
		exit 1; \
	fi

get-db-endpoint:
	@echo "$(DB_ENDPOINT)"

get-db-port:
	@echo "$(DB_PORT)"

get-db-name:
	@echo "$(DB_NAME)"

show-db-config: get-db-config
	@echo "Database Endpoint: $(shell $(MAKE) -s get-db-endpoint)"
	@echo "Database Port: $(shell $(MAKE) -s get-db-port)"
	@echo "Database Name: $(shell $(MAKE) -s get-db-name)"
	@echo "Database Username: $(DB_USERNAME)"
	@echo "Database Password: [HIDDEN]"

# Cognito authentication helper
get-cognito-token:
	@read -p "Enter username: " USERNAME; \
	read -s -p "Enter password: " PASSWORD; \
	echo ""; \
	aws cognito-idp initiate-auth \
		--auth-flow USER_PASSWORD_AUTH \
		--client-id $(COGNITO_CLIENT_ID) \
		--auth-parameters USERNAME=$$USERNAME,PASSWORD=$$PASSWORD \
		--region $(AWS_REGION) \
		--query 'AuthenticationResult.IdToken' \
		--output text 2>/dev/null || echo "Authentication failed"

show-cognito-config:
	@echo "Cognito Configuration:"
	@echo "User Pool Name: $(COGNITO_USER_POOL_NAME)"
	@echo "User Pool ID: $(COGNITO_USER_POOL_ID)"
	@echo "Client Name: $(COGNITO_USER_POOL_CLIENT_NAME)"
	@echo "Client ID: $(COGNITO_CLIENT_ID)"

# Public database access helpers (SECURITY WARNING: Only for development!)
get-my-ip:
	@echo "Your current public IPv4 address:"
	@IP=$$(curl -4 -s ifconfig.me || curl -s ipv4.icanhazip.com || dig +short myip.opendns.com @resolver1.opendns.com); \
	echo "IPv4: $$IP (use $$IP/32 for CIDR)"

seed:
	@echo "Seeding PostgreSQL database with sample data..."
	@echo "Host: $(shell $(MAKE) -s get-db-endpoint)"
	@echo "Database: $(shell $(MAKE) -s get-db-name)"
	@echo "Username: $(DB_USERNAME)"
	@echo ""
	@if [ ! -f "sql/seed_data.sql" ]; then \
		echo "Error: sql/seed_data.sql not found!"; \
		exit 1; \
	fi
	@echo "Running seed script..."
	docker run --pull=missing -v "$(PWD):/app" -w /app \
		-e PGPASSWORD="$(DB_PASSWORD)" \
		postgres:15-alpine psql \
		--host=$(shell $(MAKE) -s get-db-endpoint) \
		--port=$(shell $(MAKE) -s get-db-port) \
		--username=$(DB_USERNAME) \
		--dbname=$(shell $(MAKE) -s get-db-name) \
		--file=sql/seed_data.sql
	@echo "Database seeding completed!"

drop-tables:
	@echo "WARNING: This will DROP ALL TABLES in the database!"
	@echo "Host: $(shell $(MAKE) -s get-db-endpoint)"
	@echo "Database: $(shell $(MAKE) -s get-db-name)"
	@echo "Username: $(DB_USERNAME)"
	@echo ""
	@echo -n "Are you sure you want to continue? (y/N): "; \
	read CONFIRM; \
	if [ "$$CONFIRM" != "y" ] && [ "$$CONFIRM" != "Y" ]; then \
		echo "Operation cancelled."; \
		exit 1; \
	fi
	@echo "Dropping all tables..."
	docker run --pull=missing \
		-e PGPASSWORD="$(DB_PASSWORD)" \
		postgres:15-alpine psql \
		--host=$(shell $(MAKE) -s get-db-endpoint) \
		--port=$(shell $(MAKE) -s get-db-port) \
		--username=$(DB_USERNAME) \
		--dbname=$(shell $(MAKE) -s get-db-name) \
		--command="DROP TABLE IF EXISTS transaction_entry CASCADE; DROP TABLE IF EXISTS transaction CASCADE; DROP TABLE IF EXISTS balance CASCADE; DROP TABLE IF EXISTS merchant CASCADE; DROP TABLE IF EXISTS category CASCADE;"
	@echo "All tables dropped successfully!"

# Docker image management
pull-postgres:
	@echo "Pulling latest PostgreSQL Docker image..."
	docker pull postgres:15-alpine
	@echo "PostgreSQL Docker image updated!"

# Clean up build artifacts and virtual environments

clean:
	rm -rf ./build ./app/schema/openapi.yml

# Database operations
db-get-identifier:
	@echo "$(DB_IDENTIFIER)"

db-connect:
	@echo "Connecting to PostgreSQL database..."
	@echo "Host: $(shell $(MAKE) -s get-db-endpoint)"
	@echo "Port: $(shell $(MAKE) -s get-db-port)"
	@echo "Database: $(shell $(MAKE) -s get-db-name)"
	@echo "Username: $(DB_USERNAME)"
	@echo ""
	docker run -it --pull=missing postgres:15-alpine psql \
		"postgres://$(DB_USERNAME):$(DB_PASSWORD)@$(shell $(MAKE) -s get-db-endpoint):$(shell $(MAKE) -s get-db-port)/$(shell $(MAKE) -s get-db-name)?sslmode=require"

db-status:
	@echo "Checking database status..."
	@DB_ID=$(shell $(MAKE) -s db-get-identifier); \
	STATUS=$$(aws rds describe-db-instances --db-instance-identifier "$$DB_ID" --query 'DBInstances[0].DBInstanceStatus' --output text --region $(AWS_REGION)); \
	echo "Database Status: $$STATUS"

db-start:
	@echo "Starting database instance..."
	@DB_ID=$(shell $(MAKE) -s db-get-identifier); \
	echo "Starting RDS instance: $$DB_ID"; \
	aws rds start-db-instance --db-instance-identifier "$$DB_ID" --region $(AWS_REGION); \
	echo "Database start initiated. Use 'make db-status' to check progress."

db-stop:
	@echo "Stopping database instance to save costs..."
	@DB_ID=$(shell $(MAKE) -s db-get-identifier); \
	echo "Stopping RDS instance: $$DB_ID"; \
	aws rds stop-db-instance --db-instance-identifier "$$DB_ID" --region $(AWS_REGION); \
	echo "Database stop initiated. This will save costs while stopped."

db-wait-available:
	@echo "Waiting for database to become available..."
	@DB_ID=$(shell $(MAKE) -s db-get-identifier); \
	echo "Monitoring RDS instance: $$DB_ID"; \
	aws rds wait db-instance-available --db-instance-identifier "$$DB_ID" --region $(AWS_REGION); \
	echo "Database is now available!"

db-wait-stopped:
	@echo "Waiting for database to stop completely..."
	@DB_ID=$(shell $(MAKE) -s db-get-identifier); \
	echo "Monitoring RDS instance: $$DB_ID"; \
	aws rds wait db-instance-stopped --db-instance-identifier "$$DB_ID" --region $(AWS_REGION); \
	echo "Database is now stopped!"

# Combined commands for convenience
db-quick-start: db-start db-wait-available
	@echo "Database is ready for use!"

db-quick-stop: db-stop db-wait-stopped
	@echo "Database is now stopped and saving costs!"
