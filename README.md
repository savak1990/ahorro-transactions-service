# Ahorro Transactions Service

This project is a Go application that provides a REST API for managing user transactions, balances, and per-user category prioritization, backed by AWS Aurora PostgreSQL and DynamoDB. It includes infrastructure-as-code (Terraform), functional tests (pytest), and supports easy local and CI testing with isolated databases.

---

## Features

- RESTful API for creating, updating, and retrieving transactions and balances
- Per-user category prioritization with DynamoDB
- Aurora PostgreSQL for transactional data (transactions, balances, outbox)
- Infrastructure managed with Terraform modules
- Functional end-to-end tests using Python and pytest
- Easy local development and test setup via Makefile
- Event-driven architecture with outbox pattern and Kafka integration

---

## Prerequisites

- Go 1.20 or later
- Python 3.8+ (for functional tests)
- AWS CLI configured with appropriate credentials
- Terraform 1.3+ installed
- [Optional] Docker (if you want to use DynamoDB Local for development)

---

## Installation

1. **Clone the repository:**
    ```bash
    git clone https://github.com/your-username/ahorro-transactions-service.git
    cd ahorro-transactions-service
    ```

2. **Install Go dependencies:**
    ```bash
    go mod tidy
    ```

3. **Install Python dependencies (for tests):**
    ```bash
    python3 -m venv .venv
    source .venv/bin/activate
    pip install -r test/requirements.txt
    ```

4. **Install Terraform dependencies:**
    ```bash
    cd terraform
    terraform init
    cd ..
    ```

---

## Usage

### Run the Application

1. **Deploy Aurora PostgreSQL and DynamoDB tables (using Terraform):**
    ```bash
    make deploy
    ```

2. **Run the Go server locally:**
    ```bash
    make app-run
    ```
    The server will start on `localhost:8080`.

3. **API Endpoints:**
    - `POST   /transactions` — Create a new transaction (expense, income, or movement)
    - `GET    /transactions` — List transactions for the user (supports filtering, sorting, pagination)
    - `GET    /transactions/{transaction_id}` — Get transaction details
    - `PUT    /transactions/{transaction_id}` — Update transaction
    - `DELETE /transactions/{transaction_id}` — Delete transaction (soft delete)
    - `GET    /balances/{user_id}` — List all balances for the user
    - `POST   /balances/{user_id}` — Create new balance for the user
    - `PUT    /balances/{balance_id}` — Update balance information
    - `DELETE /balances/{balance_id}` — Delete balance information
    - `GET    /categories/{user_id}` — List categories for the user, sorted by personalized score
    - `GET    /health` — Health check
    - `GET    /info` — Info endpoint

---

### Running Functional Tests

Functional tests will:
- Deploy dedicated test Aurora PostgreSQL and DynamoDB tables
- Start the Go server with the test databases
- Run all Python/pytest tests
- Clean up the test tables after

To run all functional tests:
```bash
make functional-test
```

---

### Clean Up

To destroy all infrastructure created by Terraform:
```bash
make undeploy
```

To remove build artifacts:
```bash
make clean
```

---

## Configuration

- The PostgreSQL and DynamoDB table names are configurable via environment variables (`DB_TABLE_NAME`, etc.).
- The Makefile and Terraform use `DB_TABLE_NAME` and `DB_TABLE_TEST_NAME` for production and test tables, respectively.
- You can override these variables when running `make`:
    ```bash
    make app-run DB_TABLE_NAME=my-custom-table
    ```

---

## Architecture Overview

- **Aurora PostgreSQL** stores transactions, balances, and the outbox table for event publishing.
- **DynamoDB** stores per-user categories and prioritization scores.
- **Outbox Pattern:** All transaction writes also insert an event into the outbox table. A background Lambda reads the outbox and publishes events to Kafka (MSK), then deletes successfully sent events.
- **API Gateway + Lambda** provide the REST API endpoints.
- **Kafka (MSK)** is used for event-driven integration with other services.

---

## Contributing

Contributions are welcome! Please open issues or submit pull requests.

---

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

---

## Contact

For questions or feedback, please contact [your-email@example.com].

# Ahorro Transactions Service API Endpoints

Base URL (local):
```
http://localhost:8080
```

## Health & Info
- `GET /health` — Health check
- `GET /info` — Service info

## Transactions APIs
- `POST /transactions` — Create a new transaction
- `GET /transactions` — List transactions for the user (supports filtering, sorting, pagination)
  - Query params: `user_id`, `type`, `category`, `sort_by`, `order`, `count`, `startKey`
- `GET /transactions/{transaction_id}` — Get transaction details
- `PUT /transactions/{transaction_id}` — Update transaction
- `DELETE /transactions/{transaction_id}` — Delete transaction

## Categories APIs
- `GET /categories/{user_id}` — List categories for the user (sorted by personalized score, paginated)

---

### Example Requests

#### Create Transaction
```http
POST http://localhost:8080/transactions
Content-Type: application/json

{
  "user_id": "user1",
  "group_id": "group1",
  "type": "expense",
  "amount": 100.50,
  "balance_id": "balance1",
  "category": "Groceries",
  "description": "Weekly groceries",
  "transacted_at": "2024-06-19T12:00:00Z"
}
```

#### List Transactions
```http
GET http://localhost:8080/transactions?user_id=user1&count=10&sort_by=transacted_at&order=desc
```

#### Get Transaction
```http
GET http://localhost:8080/transactions/{transaction_id}
```

#### Update Transaction
```http
PUT http://localhost:8080/transactions/{transaction_id}
Content-Type: application/json
{
  "user_id": "user1",
  "group_id": "group1",
  "type": "expense",
  "amount": 120.00,
  "balance_id": "balance1",
  "category": "Groceries",
  "description": "Updated groceries",
  "transacted_at": "2024-06-19T12:00:00Z"
}
```

#### Delete Transaction
```http
DELETE http://localhost:8080/transactions/{transaction_id}
```

#### List Categories
```http
GET http://localhost:8080/categories/user1
```

---

See also: [ARCHITECTURE.md](../docs/architecture.md) for full API and data model details.