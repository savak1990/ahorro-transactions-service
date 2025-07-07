# Copilot Instructions for Ahorro Transactions Service

## Fresh Start Workflow

1. **Build and Start the Application Locally**  
   Run the following command to build the application and start it locally:
   ```bash
   make local-full-start
   ```

2. **Automigrate and Create Database Tables**  
   In a separate terminal, call the `/health` endpoint to trigger GORM automigration and create database tables:
   ```bash
   curl -X GET http://localhost:8080/health
   ```

3. **Seed the Database with Test Data**  
   Populate the database tables with test data:
   ```bash
   make seed
   ```

---

## Regenerating Seed Data

1. **Drop All Tables**  
   Reset the database by dropping all tables:
   ```bash
   make local-drop-tables
   ```

2. **Stop the Service**  
   Stop the running service (e.g., Ctrl+C in the terminal where the service is running).

3. **Follow Fresh Start Instructions**  
   Use the steps in the "Fresh Start Workflow" section to rebuild and reseed the database.

---

## Restarting the Application

To restart the application without resetting the database or seed data, simply run:
```bash
make local-full-start
```

---

## Testing the Service

- Use `curl` commands to test the service endpoints.  
- Verify the expected behavior and data integrity manually.  
- Example: Testing the `/transactions` endpoint:
  ```bash
  curl -X GET http://localhost:8080/transactions
  ```

---

## Notes for Copilot

- **Avoid Running Tests:**  
  Do not attempt to run tests, as the project has limited test coverage. Focus on ensuring the project compiles successfully.

- **Avoid Creating New Files for Tests:**
  Do not create new files for tests. Instead, use the existing structure and curl requests to test functionality.

- **Don't Create Summary Files:**  
  Do not create summary files for test results. Use the provided scripts and curl commands to verify functionality.

- **Compilation Check:**  
  Always ensure the project compiles without errors when making changes.

- **Endpoint Testing:**  
  Use `curl` commands to test endpoints and verify the required functionality.

- **Database Auto-Migration:**  
  The `/health` endpoint triggers GORM auto-migration. Always call this endpoint after starting the service to ensure the database schema is up-to-date.

- **Local Development Workflow:**  
  Follow the `Makefile` targets for local development. Key commands include:
  - `make local-full-stop` to stop the service and drop tables.
  - `make local-full-start` for a complete setup.
  - `curl -X GET http://localhost:8080/health` to create tables
  - `make seed` to populate the database with test data.

---

## Suggestions for Improvement

1. **Document Key Endpoints:**  
   Include a list of frequently used API endpoints and example `curl` commands for testing.

2. **Error Handling:**  
   Add instructions for handling common errors, such as database connection issues or missing environment variables.

3. **Environment Variables:**  
   Document any required environment variables or AWS credentials needed for local development.

4. **Build Artifacts:**  
   Mention the location of build artifacts and how to clean them if needed:
   ```bash
   make clean
   ```

5. **Database Connection:**  
   Include a command to connect to the local database for manual inspection:
   ```bash
   make local-db-connect
   ```
