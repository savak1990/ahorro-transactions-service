# Basic arguments
APP_NAME=ahorro
SERVICE_NAME=transactions
INSTANCE_NAME=$(shell whoami)

FULL_NAME=$(APP_NAME)-$(SERVICE_NAME)-$(INSTANCE_NAME)
CATEGORIES_DB_TABLE_NAME=$(FULL_NAME)-categories-db
TRANSACTIONS_DB_TABLE_NAME=$(FULL_NAME)-transactions-db

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

.PHONY: all build run package test clean deploy undeploy plan

# Build and package main app
$(APP_LAMBDA_BINARY): $(shell find $(APP_DIR) -type f -name '*.go')
	@mkdir -p $(APP_BUILD_DIR)
	cd $(APP_DIR) && GOOS=linux GOARCH=amd64 go build -o ../$(APP_LAMBDA_BINARY) main.go

$(APP_BINARY): $(APP_DIR)/main.go
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

test: $(APP_BINARY)
	cd $(APP_DIR) && go test ./...

run: app-build-local
	CATEGORIES_DB_TABLE_NAME=$(CATEGORIES_DB_TABLE_NAME) TRANSACTIONS_DB_TABLE_NAME=$(TRANSACTIONS_DB_TABLE_NAME) ./$(APP_BINARY)

package: $(APP_LAMBDA_HANDLER_ZIP)

package-timestamp: $(APP_LAMBDA_HANDLER_ZIP_TIMESTAMP)

upload: $(APP_LAMBDA_HANDLER_ZIP)
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

deploy:
	cd deploy && \
	terraform init && \
	terraform apply -auto-approve \
		-var="app_name=$(APP_NAME)" \
		-var="service_name=$(SERVICE_NAME)" \
		-var="env=$(INSTANCE_NAME)"

undeploy:
	cd deploy && \
	terraform init && \
	terraform destroy -auto-approve \
		-var="app_name=$(APP_NAME)" \
		-var="service_name=$(SERVICE_NAME)" \
		-var="env=$(INSTANCE_NAME)"

show-api-url:
	@cd deploy && terraform output -raw api_url

# Clean up build artifacts and virtual environments

clean:
	rm -rf ./build
