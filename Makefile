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

# Main app arguments
APP_DIR=app
APP_BUILD_DIR=./build/service-handler
APP_LAMBDA_ZIP_BASE_NAME=$(SERVICE_NAME)-lambda
APP_LAMBDA_ZIP_NAME=$(APP_LAMBDA_ZIP_BASE_NAME).zip
APP_LAMBDA_ZIP_TIMESTAMP_NAME=$(APP_LAMBDA_ZIP_BASE_NAME)-$(shell date +%y%m%d-%H%M).zip
APP_LAMBDA_HANDLER_ZIP=$(APP_BUILD_DIR)/$(APP_LAMBDA_ZIP_NAME)
APP_LAMBDA_HANDLER_ZIP_TIMESTAMP=$(APP_BUILD_DIR)/$(APP_LAMBDA_ZIP_TIMESTAMP_NAME)
APP_LAMBDA_BINARY=$(APP_BUILD_DIR)/bootstrap
APP_BINARY=$(APP_BUILD_DIR)/transactions_service
APP_LAMBDA_S3_PATH=s3://ahorro-artifacts/transactions

# Schema generation
SCHEMA_TEMPLATE=schema/openapi.yml.tml
SCHEMA_OUTPUT=$(APP_DIR)/schema/openapi.yml

.PHONY: all build app-build-local app-build-lambda run package test clean deploy undeploy plan get-db-config get-db-endpoint get-db-port get-db-name show-db-config get-my-ip connect-db seed pull-postgres deploy-public-custom drop-tables generate-schema

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
	@mkdir -p $(APP_BUILD_DIR)
	cd $(APP_DIR) && GOOS=linux GOARCH=amd64 go build -o ../$(APP_LAMBDA_BINARY) main.go

$(APP_BINARY): $(APP_DIR)/main.go $(SCHEMA_OUTPUT)
	@mkdir -p $(APP_BUILD_DIR)
	cd $(APP_DIR) && go build -o ../$(APP_BINARY) main.go

$(APP_LAMBDA_HANDLER_ZIP): $(APP_LAMBDA_BINARY)
	@mkdir -p $(APP_BUILD_DIR)
	cd $(APP_BUILD_DIR) && zip $(APP_LAMBDA_ZIP_NAME) bootstrap

$(APP_LAMBDA_HANDLER_ZIP_TIMESTAMP): $(APP_LAMBDA_BINARY)
	@mkdir -p $(APP_BUILD_DIR)
	cd $(APP_BUILD_DIR) && zip $(APP_LAMBDA_ZIP_TIMESTAMP_NAME) bootstrap

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
	DB_USERNAME="$(DB_USERNAME)" \
	DB_PASSWORD="$(DB_PASSWORD)" \
	./$(APP_BINARY)

package: $(APP_LAMBDA_HANDLER_ZIP)

package-timestamp: $(APP_LAMBDA_HANDLER_ZIP_TIMESTAMP)

upload: $(APP_LAMBDA_HANDLER_ZIP)
	aws s3 rm $(APP_LAMBDA_S3_PATH)/$(APP_LAMBDA_ZIP_NAME) --quiet || true
	aws s3 cp $(APP_LAMBDA_HANDLER_ZIP) $(APP_LAMBDA_S3_PATH)/$(APP_LAMBDA_ZIP_NAME)

upload-timestamp: $(APP_LAMBDA_HANDLER_ZIP_TIMESTAMP)
	aws s3 cp $(APP_LAMBDA_HANDLER_ZIP_TIMESTAMP) $(APP_LAMBDA_S3_PATH)/$(APP_LAMBDA_ZIP_TIMESTAMP_NAME)

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
deploy-public:
	@echo "WARNING: Enabling public database access for development!"
	@IP=$$(curl -4 -s ifconfig.me || curl -s ipv4.icanhazip.com || dig +short myip.opendns.com @resolver1.opendns.com); \
	echo "Your IPv4: $$IP"; \
	echo "Using CIDR block: $$IP/32"; \
	cd deploy && \
	terraform apply -auto-approve \
		-var="app_name=$(APP_NAME)" \
		-var="service_name=$(SERVICE_NAME)" \
		-var="env=$(INSTANCE_NAME)" \
		-var="enable_db_public_access=true" \
		-var="my_ip_cidr=$$IP/32"

deploy-public-custom:
	@echo "WARNING: Enabling public database access for development!"
	@IP=$$(curl -4 -s ifconfig.me || curl -s ipv4.icanhazip.com || dig +short myip.opendns.com @resolver1.opendns.com); \
	echo "Your current IPv4: $$IP"; \
	echo -n "Enter CIDR block (e.g., $$IP/32 or 0.0.0.0/0): "; \
	read CIDR_BLOCK; \
	if [ -z "$$CIDR_BLOCK" ]; then \
		echo "Error: CIDR block cannot be empty"; \
		exit 1; \
	fi; \
	echo "Using CIDR block: $$CIDR_BLOCK"; \
	cd deploy && \
	terraform apply -auto-approve \
		-var="app_name=$(APP_NAME)" \
		-var="service_name=$(SERVICE_NAME)" \
		-var="env=$(INSTANCE_NAME)" \
		-var="enable_db_public_access=true" \
		-var="my_ip_cidr=$$CIDR_BLOCK"

deploy-private:
	@echo "Disabling public database access..."
	cd deploy && \
	terraform apply -auto-approve \
		-var="app_name=$(APP_NAME)" \
		-var="service_name=$(SERVICE_NAME)" \
		-var="env=$(INSTANCE_NAME)" \
		-var="enable_db_public_access=false"

undeploy:
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
	@echo "Fetching database configuration from Terraform..."

get-db-endpoint:
	@cd deploy && terraform output -raw db_endpoint

get-db-port:
	@cd deploy && terraform output -raw db_port

get-db-name:
	@cd deploy && terraform output -raw db_name

show-db-config: get-db-config
	@echo "Database Endpoint: $(shell $(MAKE) -s get-db-endpoint)"
	@echo "Database Port: $(shell $(MAKE) -s get-db-port)"
	@echo "Database Name: $(shell $(MAKE) -s get-db-name)"
	@echo "Database Username: $(DB_USERNAME)"
	@echo "Database Password: [HIDDEN]"

# Public database access helpers (SECURITY WARNING: Only for development!)
get-my-ip:
	@echo "Your current public IPv4 address:"
	@IP=$$(curl -4 -s ifconfig.me || curl -s ipv4.icanhazip.com || dig +short myip.opendns.com @resolver1.opendns.com); \
	echo "IPv4: $$IP (use $$IP/32 for CIDR)"

connect-db:
	@echo "Connecting to PostgreSQL database..."
	@echo "Host: $(shell $(MAKE) -s get-db-endpoint)"
	@echo "Port: $(shell $(MAKE) -s get-db-port)"
	@echo "Database: $(shell $(MAKE) -s get-db-name)"
	@echo "Username: $(DB_USERNAME)"
	@echo ""
	docker run -it --pull=missing postgres:15-alpine psql \
		"postgres://$(DB_USERNAME):$(DB_PASSWORD)@$(shell $(MAKE) -s get-db-endpoint):$(shell $(MAKE) -s get-db-port)/$(shell $(MAKE) -s get-db-name)?sslmode=require"

seed:
	@echo "Seeding PostgreSQL database with sample data..."
	@echo "Host: $(shell $(MAKE) -s get-db-endpoint)"
	@echo "Database: $(shell $(MAKE) -s get-db-name)"
	@echo "Username: $(DB_USERNAME)"
	@echo ""
	@if [ ! -f "seed/seed_data.sql" ]; then \
		echo "Error: seed/seed_data.sql not found!"; \
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
		--file=seed/seed_data.sql
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
