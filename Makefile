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
GITHUB_TOKEN=$(shell aws secretsmanager get-secret-value --secret-id $(SECRET_NAME) --query 'SecretString' --output text --region $(AWS_REGION) | jq -r '.github_token')

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

# Local PostgreSQL development configuration
LOCAL_DB_HOST=localhost
LOCAL_DB_PORT=5432
LOCAL_DB_NAME=ahorro_transactions_local
LOCAL_DB_USER=postgres
LOCAL_DB_PASSWORD=local_password
LOCAL_POSTGRES_CONTAINER=ahorro-postgres-local
LOCAL_SSL_MODE=disable

# Main app arguments
APP_DIR=app
APP_BUILD_DIR=./build/service-handler
APP_LAMBDA_ZIP_BASE_NAME=$(SERVICE_NAME)-lambda
APP_LAMBDA_ZIP_NAME=$(APP_LAMBDA_ZIP_BASE_NAME).zip
APP_LAMBDA_HANDLER_ZIP=$(APP_BUILD_DIR)/$(APP_LAMBDA_ZIP_NAME)
APP_LAMBDA_BINARY=$(APP_BUILD_DIR)/bootstrap
APP_BINARY=$(APP_BUILD_DIR)/transactions_service

# S3 paths for different deployment types
BUILD_INFO_FILE=$(APP_DIR)/buildinfo/build-info.json
TIMESTAMP_FILE := .timestamp
TIMESTAMP := $(shell cat $(TIMESTAMP_FILE) 2>/dev/null || (date +build-%y%m%d-%H%M > $(TIMESTAMP_FILE) && cat $(TIMESTAMP_FILE)))
APP_LAMBDA_S3_BASE=s3://ahorro-artifacts/transactions
APP_LAMBDA_S3_PATH_LOCAL=$(APP_LAMBDA_S3_BASE)/$(INSTANCE_NAME)/$(APP_LAMBDA_ZIP_NAME)
APP_LAMBDA_S3_PATH_TIMESTAMP=$(APP_LAMBDA_S3_BASE)/$(TIMESTAMP)/$(APP_LAMBDA_ZIP_NAME)

# GitHub repository configuration
GITHUB_REPO=ahorro-transactions-service
GITHUB_TAG_NAME=$(TIMESTAMP)

# Schema generation
SCHEMA_TEMPLATE=schema/openapi.yml.tml
SCHEMA_OUTPUT=$(APP_DIR)/schema/openapi.yml

.PHONY: all build app-build-local app-build-lambda run package test clean deploy undeploy plan get-db-config get-db-endpoint get-db-port get-db-name show-db-config get-my-ip db-connect seed verify-seed pull-postgres deploy-public-custom drop-tables generate-schema generate-build-info db-start db-stop db-status db-get-identifier get-cognito-token show-cognito-config git-tag upload-and-tag local-db-start local-db-stop local-db-status local-db-create local-db-destroy local-db-connect local-drop-tables local-cleanup-port local-seed local-verify-seed local-run help

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
	@echo "  generate-build-info   - Generate build info with git metadata"
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
	@echo "  seed                  - Seed database with modular sample data (category_groups → categories → merchants → balances → transactions → transaction_entries)"
	@echo "  verify-seed           - Verify seed data integrity and relationships"
	@echo "  drop-tables           - Drop all tables (⚠️  DESTRUCTIVE)"
	@echo ""
	@echo "🏠 Local Development:"
	@echo "  local-db-start        - Start local PostgreSQL in Docker"
	@echo "  local-db-stop         - Stop local PostgreSQL container"
	@echo "  local-db-status       - Check local PostgreSQL container and database status"
	@echo "  local-db-create       - Create local database"
	@echo "  local-db-destroy      - Remove local PostgreSQL container and data"
	@echo "  local-db-connect      - Connect to local PostgreSQL database"
	@echo "  local-drop-tables     - Drop all tables in local database"
	@echo "  local-cleanup-port    - Clean up any processes using port 8080"
	@echo "  local-seed            - Seed local database with sample data"
	@echo "  local-verify-seed     - Verify local seed data"
	@echo "  local-run             - Run service locally with local database"
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
	@echo "  # AWS Development:"
	@echo "  make build && make deploy"
	@echo "  make get-cognito-token"
	@echo "  make db-quick-start && make seed"
	@echo ""
	@echo "  # Local Development:"
	@echo "  make local-db-start && make local-db-create && make local-seed && make local-run"

# Schema generation target
$(SCHEMA_OUTPUT): $(SCHEMA_TEMPLATE)
	@echo "Generating OpenAPI schema from template..."
	@mkdir -p $(dir $(SCHEMA_OUTPUT))
	@sed -e 's/$${INSTANCE_NAME}/$(INSTANCE_NAME)/g' \
	     -e 's/$${DOMAIN_NAME}/$(DOMAIN_NAME)/g' \
	     $(SCHEMA_TEMPLATE) > $(SCHEMA_OUTPUT)
	@echo "OpenAPI schema generated: $(SCHEMA_OUTPUT)"

generate-schema: $(SCHEMA_OUTPUT)

generate-build-info:
	@echo "Generating build info..."
	@mkdir -p $(dir $(BUILD_INFO_FILE))
	@eval $(shell ./scripts/github-meta.sh savak1990 ahorro-transactions-service main $(GITHUB_TOKEN)); \
	BUILD_TIME=$$(date -u +"%Y-%m-%dT%H:%M:%SZ"); \
	BUILD_USER=$$(whoami); \
	GO_VERSION=$$(go version | cut -d' ' -f3 2>/dev/null || echo "unknown"); \
	printf '{\n  "version": "%s",\n  "gitBranch": "%s",\n  "gitCommit": "%s",\n  "gitShort": "%s",\n  "buildTime": "%s",\n  "buildUser": "%s",\n  "goVersion": "%s",\n  "serviceName": "%s"\n}' \
		"$(TIMESTAMP)" "$$GIT_BRANCH" "$$GIT_COMMIT" "$$GIT_SHORT" "$$BUILD_TIME" "$$BUILD_USER" "$$GO_VERSION" "ahorro-transactions-service" \
		> $(BUILD_INFO_FILE)
	@echo "Build info generated: $(BUILD_INFO_FILE)"

# Build and package main app
$(APP_LAMBDA_BINARY): $(shell find $(APP_DIR) -type f -name '*.go') $(SCHEMA_OUTPUT) generate-build-info
	@echo "Building Lambda binary using Docker (ensures compatibility)..."
	@mkdir -p $(APP_BUILD_DIR)
	@docker run \
		-v $(PWD)/$(APP_DIR):/src \
		-v $(PWD)/$(APP_BUILD_DIR):/build \
		-v $(PWD)/.git:/src/.git \
		-w /src \
		golang:1.23-alpine \
		sh -c "apk add --no-cache git ca-certificates && \
		       CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
		       go build -ldflags='-s -w -extldflags=-static' -tags netgo -a \
		       -o /build/bootstrap main.go"

$(APP_BINARY): $(APP_DIR)/main.go $(SCHEMA_OUTPUT) $(BUILD_INFO_FILE)
	@mkdir -p $(APP_BUILD_DIR)
	cd $(APP_DIR) && go build -o ../$(APP_BINARY) main.go

$(APP_LAMBDA_HANDLER_ZIP): $(APP_LAMBDA_BINARY)
	@mkdir -p $(APP_BUILD_DIR)
	cd $(APP_BUILD_DIR) && zip $(APP_LAMBDA_ZIP_NAME) bootstrap

# Combined build and package targets
build: new-timestamp $(APP_BINARY) $(APP_LAMBDA_BINARY)

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
	SSL_MODE=require \
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
	@if [ ! -f "scripts/seed_database.sh" ]; then \
		echo "Error: scripts/seed_database.sh not found!"; \
		exit 1; \
	fi
	@echo "Running modular seed scripts..."
	docker run --pull=missing -v "$(PWD):/app" -w /app \
		-e DB_HOST=$(shell $(MAKE) -s get-db-endpoint) \
		-e DB_PORT=$(shell $(MAKE) -s get-db-port) \
		-e DB_USER=$(DB_USERNAME) \
		-e DB_PASSWORD=$(DB_PASSWORD) \
		-e DB_NAME=$(shell $(MAKE) -s get-db-name) \
		--entrypoint=/bin/bash \
		postgres:15-alpine \
		/app/scripts/seed_database.sh
	@echo "Database seeding completed!"

verify-seed:
	@echo "Verifying seeded data in PostgreSQL database..."
	@echo "Host: $(shell $(MAKE) -s get-db-endpoint)"
	@echo "Database: $(shell $(MAKE) -s get-db-name)"
	@echo "Username: $(DB_USERNAME)"
	@echo ""
	@if [ ! -f "scripts/verify_seed_data.sh" ]; then \
		echo "Error: scripts/verify_seed_data.sh not found!"; \
		exit 1; \
	fi
	@echo "Running seed data verification..."
	docker run --pull=missing -v "$(PWD):/app" -w /app \
		-e DB_HOST=$(shell $(MAKE) -s get-db-endpoint) \
		-e DB_PORT=$(shell $(MAKE) -s get-db-port) \
		-e DB_USER=$(DB_USERNAME) \
		-e DB_PASSWORD=$(DB_PASSWORD) \
		-e DB_NAME=$(shell $(MAKE) -s get-db-name) \
		--entrypoint=/bin/bash \
		postgres:15-alpine \
		/app/scripts/verify_seed_data.sh
	@echo "Seed data verification completed!"

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

# Clean up build artifacts and virtual environments

clean:
	rm -rf ./build ./app/schema/openapi.yml ./app/buildinfo/build-info.json .timestamp

# ================================
# Local Development Setup
# ================================

# Start local PostgreSQL container
local-db-start:
	@echo "Starting local PostgreSQL container..."
	@if docker ps -a --format "table {{.Names}}" | grep -q "^$(LOCAL_POSTGRES_CONTAINER)$$"; then \
		if docker ps --format "table {{.Names}}" | grep -q "^$(LOCAL_POSTGRES_CONTAINER)$$"; then \
			echo "PostgreSQL container is already running"; \
		else \
			echo "Starting existing PostgreSQL container..."; \
			docker start $(LOCAL_POSTGRES_CONTAINER); \
		fi; \
	else \
		echo "Creating new PostgreSQL container..."; \
		docker run -d \
			--name $(LOCAL_POSTGRES_CONTAINER) \
			-e POSTGRES_DB=$(LOCAL_DB_NAME) \
			-e POSTGRES_USER=$(LOCAL_DB_USER) \
			-e POSTGRES_PASSWORD=$(LOCAL_DB_PASSWORD) \
			-p $(LOCAL_DB_PORT):5432 \
			-v ahorro-postgres-data:/var/lib/postgresql/data \
			--restart unless-stopped \
			postgres:15-alpine; \
	fi
	@echo "Waiting for PostgreSQL to be ready..."
	@for i in 1 2 3 4 5 6 7 8 9 10; do \
		if docker exec $(LOCAL_POSTGRES_CONTAINER) pg_isready -U $(LOCAL_DB_USER) -d $(LOCAL_DB_NAME) >/dev/null 2>&1; then \
			echo "Local PostgreSQL is ready at localhost:$(LOCAL_DB_PORT)"; \
			exit 0; \
		fi; \
		echo "Waiting... ($$i/10)"; \
		sleep 2; \
	done; \
	echo "Failed to connect to PostgreSQL after 20 seconds"; \
	exit 1

# Stop local PostgreSQL container
local-db-stop:
	@echo "Stopping local PostgreSQL container..."
	@if docker ps --format "table {{.Names}}" | grep -q "^$(LOCAL_POSTGRES_CONTAINER)$$"; then \
		docker stop $(LOCAL_POSTGRES_CONTAINER); \
		echo "PostgreSQL container stopped"; \
	else \
		echo "PostgreSQL container is not running"; \
	fi

# Check local PostgreSQL container status
local-db-status:
	@echo "Checking local PostgreSQL container status..."
	@echo "Container: $(LOCAL_POSTGRES_CONTAINER)"
	@echo "Expected Database: $(LOCAL_DB_NAME)"
	@echo "Expected Port: $(LOCAL_DB_PORT)"
	@echo ""
	@if docker ps -a --format "table {{.Names}}" | grep -q "^$(LOCAL_POSTGRES_CONTAINER)$$"; then \
		if docker ps --format "table {{.Names}}" | grep -q "^$(LOCAL_POSTGRES_CONTAINER)$$"; then \
			echo "✅ Container Status: RUNNING"; \
			if docker exec $(LOCAL_POSTGRES_CONTAINER) pg_isready -U $(LOCAL_DB_USER) -d $(LOCAL_DB_NAME) >/dev/null 2>&1; then \
				echo "✅ Database Status: READY"; \
				echo "✅ Connection: SUCCESS"; \
				echo ""; \
				echo "Database Details:"; \
				docker exec $(LOCAL_POSTGRES_CONTAINER) psql -U $(LOCAL_DB_USER) -d $(LOCAL_DB_NAME) -c "SELECT version();" -t | head -1 | xargs echo "PostgreSQL Version:"; \
				docker exec $(LOCAL_POSTGRES_CONTAINER) psql -U $(LOCAL_DB_USER) -d $(LOCAL_DB_NAME) -c "SELECT current_database(), current_user;" -t | xargs echo "Connected as:"; \
				echo ""; \
				echo "Tables:"; \
				docker exec $(LOCAL_POSTGRES_CONTAINER) psql -U $(LOCAL_DB_USER) -d $(LOCAL_DB_NAME) -c "\dt" 2>/dev/null || echo "No tables found (database may need migration)"; \
			else \
				echo "❌ Database Status: NOT READY"; \
				echo "❌ Connection: FAILED"; \
			fi; \
		else \
			echo "⏸️  Container Status: STOPPED"; \
			echo "❌ Database Status: NOT ACCESSIBLE"; \
		fi; \
	else \
		echo "❌ Container Status: NOT EXISTS"; \
		echo "❌ Database Status: NOT AVAILABLE"; \
		echo ""; \
		echo "💡 Run 'make local-db-start' to create and start the container"; \
	fi

# Create/initialize local database with schema
local-db-create: local-db-start
	@echo "Creating database schema in local PostgreSQL..."
	@echo "Database: $(LOCAL_DB_NAME)"
	@echo "Host: $(LOCAL_DB_HOST):$(LOCAL_DB_PORT)"
	@echo "User: $(LOCAL_DB_USER)"
	@echo ""
	@echo "The application will auto-migrate the schema on first connection."
	@echo "Local database is ready for use!"

# Remove local PostgreSQL container and all data
local-db-destroy:
	@echo "WARNING: This will DESTROY the local PostgreSQL container and ALL DATA!"
	@echo "Container: $(LOCAL_POSTGRES_CONTAINER)"
	@echo "Volume: ahorro-postgres-data"
	@echo ""
	@echo -n "Are you sure you want to continue? (y/N): "; \
	read CONFIRM; \
	if [ "$$CONFIRM" != "y" ] && [ "$$CONFIRM" != "Y" ]; then \
		echo "Operation cancelled."; \
		exit 1; \
	fi
	@echo "Removing PostgreSQL container..."
	@docker rm -f $(LOCAL_POSTGRES_CONTAINER) 2>/dev/null || echo "Container not found"
	@echo "Removing PostgreSQL data volume..."
	@docker volume rm ahorro-postgres-data 2>/dev/null || echo "Volume not found"
	@echo "Local PostgreSQL container and data destroyed"

# Drop all tables in local database
local-drop-tables: local-db-start
	@echo "WARNING: This will DROP ALL TABLES in the local database!"
	@echo "Container: $(LOCAL_POSTGRES_CONTAINER)"
	@echo "Database: $(LOCAL_DB_NAME)"
	@echo ""
	@echo -n "Are you sure you want to continue? (y/N): "; \
	read CONFIRM; \
	if [ "$$CONFIRM" != "y" ] && [ "$$CONFIRM" != "Y" ]; then \
		echo "Operation cancelled."; \
		exit 1; \
	fi
	@echo "Dropping all tables in local database..."
	docker exec $(LOCAL_POSTGRES_CONTAINER) psql -U $(LOCAL_DB_USER) -d $(LOCAL_DB_NAME) -c \
		"DROP TABLE IF EXISTS transaction_entry CASCADE; DROP TABLE IF EXISTS transaction CASCADE; DROP TABLE IF EXISTS balance CASCADE; DROP TABLE IF EXISTS merchant CASCADE; DROP TABLE IF EXISTS category CASCADE; DROP TABLE IF EXISTS category_group CASCADE;"
	@echo "All tables dropped successfully from local database!"

# Seed local database with sample data
local-seed: local-db-create
	@echo "Seeding local PostgreSQL database with sample data..."
	@echo "Host: $(LOCAL_DB_HOST):$(LOCAL_DB_PORT)"
	@echo "Database: $(LOCAL_DB_NAME)"
	@echo "Username: $(LOCAL_DB_USER)"
	@echo ""
	@if [ ! -f "scripts/seed_database.sh" ]; then \
		echo "Error: scripts/seed_database.sh not found!"; \
		exit 1; \
	fi
	@echo "Running modular seed scripts..."
	docker run --rm --pull=missing -v "$(PWD):/app" -w /app \
		--network host \
		-e DB_HOST=$(LOCAL_DB_HOST) \
		-e DB_PORT=$(LOCAL_DB_PORT) \
		-e DB_USER=$(LOCAL_DB_USER) \
		-e DB_PASSWORD=$(LOCAL_DB_PASSWORD) \
		-e DB_NAME=$(LOCAL_DB_NAME) \
		--entrypoint=/bin/bash \
		postgres:15-alpine \
		/app/scripts/seed_database.sh
	@echo "Local database seeding completed!"

# Verify local seed data
local-verify-seed: local-db-create
	@echo "Verifying seeded data in local PostgreSQL database..."
	@echo "Host: $(LOCAL_DB_HOST):$(LOCAL_DB_PORT)"
	@echo "Database: $(LOCAL_DB_NAME)"
	@echo "Username: $(LOCAL_DB_USER)"
	@echo ""
	@if [ ! -f "scripts/verify_seed_data.sh" ]; then \
		echo "Error: scripts/verify_seed_data.sh not found!"; \
		exit 1; \
	fi
	@echo "Running seed data verification..."
	docker run --rm --pull=missing -v "$(PWD):/app" -w /app \
		--network host \
		-e DB_HOST=$(LOCAL_DB_HOST) \
		-e DB_PORT=$(LOCAL_DB_PORT) \
		-e DB_USER=$(LOCAL_DB_USER) \
		-e DB_PASSWORD=$(LOCAL_DB_PASSWORD) \
		-e DB_NAME=$(LOCAL_DB_NAME) \
		--entrypoint=/bin/bash \
		postgres:15-alpine \
		/app/scripts/verify_seed_data.sh
	@echo "Local seed data verification completed!"

# Clean up any processes using port 8080
local-cleanup-port:
	@echo "Cleaning up any processes using port 8080..."
	@kill $(shell lsof -ti:8080) 2>/dev/null || true
	@echo "Port 8080 cleanup completed"

# Run the service locally with local database
local-run: app-build-local local-db-create local-cleanup-port
	@echo "Running service locally with local PostgreSQL database..."
	@echo "Database: $(LOCAL_DB_HOST):$(LOCAL_DB_PORT)/$(LOCAL_DB_NAME)"
	@echo "Service will be available at: http://localhost:8080"
	@echo ""
	DB_HOST=$(LOCAL_DB_HOST) \
	DB_PORT=$(LOCAL_DB_PORT) \
	DB_NAME=$(LOCAL_DB_NAME) \
	DB_USER=$(LOCAL_DB_USER) \
	DB_PASSWORD=$(LOCAL_DB_PASSWORD) \
	SSL_MODE=$(LOCAL_SSL_MODE) \
	LOG_LEVEL=$(LOG_LEVEL) \
	./$(APP_BINARY)

# Connect to local PostgreSQL database
local-db-connect: local-db-start
	@echo "Connecting to local PostgreSQL database..."
	@echo "Host: $(LOCAL_DB_HOST):$(LOCAL_DB_PORT)"
	@echo "Database: $(LOCAL_DB_NAME)"
	@echo "Username: $(LOCAL_DB_USER)"
	@echo ""
	docker run -it --rm --pull=missing \
		--network host \
		postgres:15-alpine psql \
		"postgres://$(LOCAL_DB_USER):$(LOCAL_DB_PASSWORD)@$(LOCAL_DB_HOST):$(LOCAL_DB_PORT)/$(LOCAL_DB_NAME)?sslmode=disable"

# At the top
TIMESTAMP_FILE := .timestamp
# Always generate a new timestamp when requested
new-timestamp:
	@date +build-%y%m%d-%H%M > $(TIMESTAMP_FILE)

TIMESTAMP := $(shell cat $(TIMESTAMP_FILE) 2>/dev/null || (date +build-%y%m%d-%H%M > $(TIMESTAMP_FILE) && cat $(TIMESTAMP_FILE)))

$(TIMESTAMP_FILE):
	@date +build-%y%m%d-%H%M > $(TIMESTAMP_FILE)

# Make build-info depend on .timestamp
$(BUILD_INFO_FILE): $(TIMESTAMP_FILE)
	@echo "Generating build info..."
	@mkdir -p $(dir $(BUILD_INFO_FILE))
	@if [ -n "$$GITHUB_TOKEN" ]; then \
	  eval $(shell ./scripts/github-meta.sh savak1990 ahorro-transactions-service main $$GITHUB_TOKEN); \
	else \
	  GIT_BRANCH=$$(git rev-parse --abbrev-ref HEAD 2>/dev/null || echo "unknown"); \
	  GIT_COMMIT=$$(git rev-parse HEAD 2>/dev/null || echo "unknown"); \
	  GIT_SHORT=$$(git rev-parse --short HEAD 2>/dev/null || echo "unknown"); \
	fi; \
	BUILD_TIME=$$(date -u +"%Y-%m-%dT%H:%M:%SZ"); \
	BUILD_USER=$$(whoami); \
	GO_VERSION=$$(go version | cut -d' ' -f3 2>/dev/null || echo "unknown"); \
	printf '{\n  "version": "%s",\n  "gitBranch": "%s",\n  "gitCommit": "%s",\n  "gitShort": "%s",\n  "buildTime": "%s",\n  "buildUser": "%s",\n  "goVersion": "%s",\n  "serviceName": "%s"\n}' \
		"$(TIMESTAMP)" "$$GIT_BRANCH" "$$GIT_COMMIT" "$$GIT_SHORT" "$$BUILD_TIME" "$$BUILD_USER" "$$GO_VERSION" "ahorro-transactions-service" \
		> $(BUILD_INFO_FILE)
	@echo "Build info generated: $(BUILD_INFO_FILE)"