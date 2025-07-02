# Ahorro Transactions Service

[![CodeBuild Status](https://codebuild.eu-west-1.amazonaws.com/badges?projectName=ahorro-transactions-service-build)](https://eu-west-1.console.aws.amazon.com/codesuite/codebuild/projects/ahorro-transactions-service-build/history)
[![CodePipeline Status](https://aws-codepipeline-badges.s3.amazonaws.com/eu-west-1/transactions/badge.svg)](https://eu-west-1.console.aws.amazon.com/codesuite/codepipeline/pipelines/transactions/view)

A Go microservice for managing financial transactions, balances, and categories with comprehensive local development support and AWS deployment capabilities.

## Features

- **RESTful API** for transactions, balances, categories, and merchants
- **PostgreSQL** for transactional data with auto-migration
- **UUID prefixing system** for easy entity identification
- **Local development** with Docker PostgreSQL
- **AWS deployment** with Lambda, RDS, and API Gateway
- **Terraform IaC** for infrastructure management
- **Comprehensive testing** and seeding scripts

---

## Quick Start Guide

### Prerequisites

- **Go 1.23+**
- **Docker** (for local PostgreSQL)
- **AWS CLI** (configured with credentials)
- **Terraform 1.3+**
- **Make** (for running commands)

### Installation

```bash
git clone <repository-url>
cd ahorro-transactions-service
go mod tidy
```

---

## Development Workflows

### 🏠 1. Local Development Flow

**Complete local development with Docker PostgreSQL - no AWS dependencies.**

#### Setup Local Environment

```bash
# 1. Start local PostgreSQL container
make local-db-start

# 2. Create database schema (auto-migration)
make local-db-create

# 3. Seed with sample data
make local-seed

# 4. Verify seed data
make local-verify-seed

# 5. Run service locally
make local-run
```

**Service will be available at: http://localhost:8080**

#### Working with Local Database

```bash
# Connect to local database
make local-db-connect

# Check database status
make local-db-status

# Reset all tables (destructive)
make local-drop-tables

# Re-seed after table reset
make local-seed
```

#### Stop Local Development

```bash
# Stop the service (Ctrl+C in terminal where local-run is running)

# Stop PostgreSQL container (keeps data)
make local-db-stop

# OR completely remove container and data
make local-db-destroy
```

#### Local Development Commands Summary

| Command | Purpose |
|---------|---------|
| `make local-db-start` | Start PostgreSQL container |
| `make local-db-create` | Initialize database schema |
| `make local-seed` | Add sample data |
| `make local-run` | Run service on port 8080 |
| `make local-db-connect` | Interactive database connection |
| `make local-db-status` | Check container/database status |
| `make local-db-stop` | Stop container (preserve data) |
| `make local-db-destroy` | Remove container and all data |

---

### ☁️ 2. AWS Remote Database Development

**Local service connecting to AWS RDS PostgreSQL.**

#### Setup Remote Database Connection

```bash
# 1. Ensure AWS credentials are configured
aws configure

# 2. Check database configuration
make show-db-config

# 3. Start remote database (if stopped)
make db-quick-start

# 4. Seed remote database
make seed

# 5. Verify remote seed data
make verify-seed

# 6. Run service locally with remote DB
make run
```

**Service connects to AWS RDS and runs locally on port 8080**

#### Working with Remote Database

```bash
# Connect to remote database
make db-connect

# Check database status
make db-status

# Reset remote tables (⚠️ DESTRUCTIVE)
make drop-tables

# Get database endpoint info
make get-db-endpoint
make get-db-port
make get-db-name
```

#### Cleanup Remote Development

```bash
# Stop remote database (saves costs)
make db-quick-stop

# Check if stopped
make db-status
```

#### Remote Database Commands Summary

| Command | Purpose |
|---------|---------|
| `make show-db-config` | Display RDS connection info |
| `make db-quick-start` | Start RDS and wait until ready |
| `make seed` | Seed remote database |
| `make run` | Run service with remote DB |
| `make db-connect` | Connect to remote database |
| `make db-quick-stop` | Stop RDS (cost savings) |
| `make drop-tables` | Reset remote tables |

---

### 🚀 3. AWS Lambda Deployment

**Full serverless deployment with Lambda, API Gateway, and RDS.**

#### Build and Deploy to Lambda

```bash
# 1. Build Lambda package
make build

# 2. Upload to S3 and deploy infrastructure
make deploy

# 3. Get API URL
make show-api-url

# 4. Seed production database
make seed

# 5. Get authentication token (if needed)
make get-cognito-token
```

#### Lambda Deployment Commands

| Command | Purpose |
|---------|---------|
| `make build` | Build local and Lambda binaries |
| `make deploy` | Deploy full infrastructure |
| `make show-api-url` | Get deployed API endpoint |
| `make undeploy` | Destroy all infrastructure |
| `make upload-and-tag` | Upload with Git tag (CI/CD) |

#### Production Management

```bash
# Check infrastructure plan
make plan

# Upload new version
make upload-and-tag

# Redeploy with new code
make deploy

# Monitor costs by stopping DB when not needed
make db-quick-stop
make db-quick-start  # when ready to use again
```

---

## API Reference

### Base URLs

- **Local Development:** `http://localhost:8080`
- **AWS Deployment:** Use `make show-api-url` to get endpoint

### Core Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/health` | Health check |
| `GET` | `/info` | Service information |
| `POST` | `/transactions` | Create transaction |
| `GET` | `/transactions` | List transactions |
| `GET` | `/transactions/{id}` | Get transaction |
| `PUT` | `/transactions/{id}` | Update transaction |
| `DELETE` | `/transactions/{id}` | Delete transaction |
| `GET` | `/balances/{user_id}` | List user balances |
| `POST` | `/balances` | Create balance |
| `PUT` | `/balances/{id}` | Update balance |
| `DELETE` | `/balances/{id}` | Delete balance |
| `GET` | `/categories` | List categories |
| `GET` | `/merchants` | List merchants |

### Example Request

```bash
# Create a transaction (local)
curl -X POST http://localhost:8080/transactions \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "99bb2200-0011-2233-4455-667788990011",
    "group_id": "88aa1100-0011-2233-4455-667788990011",
    "type": "expense",
    "amount": 25.50,
    "balance_id": "ba001111-1111-1111-1111-111111111111",
    "category_id": "ca001111-1111-1111-1111-111111111111",
    "merchant_id": "4e001111-1111-1111-1111-111111111111",
    "description": "Coffee and pastry",
    "transacted_at": "2025-07-02T10:30:00Z"
  }'
```

---

## Database Schema

The service uses a PostgreSQL database with the following main entities:

- **category_group** - Logical groupings of categories
- **category** - Transaction categories (Food, Transport, etc.)
- **merchant** - Businesses/vendors
- **balance** - User accounts/wallets
- **transaction** - Financial transactions
- **transaction_entry** - Detailed transaction line items

### UUID Prefixing System

All entities use prefixed UUIDs for easy identification:

| Entity | Prefix | Example |
|--------|--------|---------|
| Balance | `ba` | `ba001111-1111-1111-1111-111111111111` |
| Category | `ca` | `ca001111-1111-1111-1111-111111111111` |
| CategoryGroup | `c9` | `c9001111-1111-1111-1111-111111111111` |
| Merchant | `4e` | `4e001111-1111-1111-1111-111111111111` |
| Transaction | `7a` | `7a001111-1111-1111-1111-111111111111` |
| TransactionEntry | `7e` | `7e001111-1111-1111-1111-111111111111` |

---

## Development Tips

### Useful Commands

```bash
# View all available commands
make help

# Clean build artifacts
make clean

# Run tests
make test

# Check your public IP (for database access)
make get-my-ip

# Pull latest PostgreSQL image
make pull-postgres
```

### Troubleshooting

**Port 8080 already in use:**
```bash
make local-cleanup-port  # Kills processes on port 8080
```

**Database connection issues:**
```bash
make local-db-status     # Check local container status
make show-db-config      # Check remote database config
```

**Seed data foreign key errors:**
```bash
make local-drop-tables   # Reset local tables
make local-seed          # Re-seed in correct order
```

---

## Architecture

- **Language:** Go 1.23 with Gin HTTP framework
- **Database:** PostgreSQL with GORM ORM
- **Infrastructure:** Terraform (AWS Lambda, RDS, API Gateway)
- **Authentication:** AWS Cognito (production)
- **Development:** Docker PostgreSQL container
- **Build:** Docker-based Lambda binary compilation

---

## Contributing

1. Use local development flow for feature development
2. Test with both local and remote databases
3. Ensure all tests pass: `make test`
4. Follow the UUID prefixing convention
5. Update seed data if schema changes

---

## License

MIT License - see LICENSE file for details.