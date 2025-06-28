# Centralized Database Error Handling Implementation

## Summary

We've implemented a centralized error handling system that ensures **all database connection errors are properly converted to 503 Service Unavailable responses** instead of 500 Internal Server Errors.

## Key Components

### 1. Central Error Handler
```go
// handleServiceError is a centralized error handler for service layer errors
// Returns true if the error was handled (response was written), false if caller should continue
func (h *HandlerImpl) handleServiceError(w http.ResponseWriter, err error, operation string) bool
```

**Features:**
- ✅ **Database Connection Error Detection**: Converts connection errors to proper `DatabaseConnectionError` panics
- ✅ **Maintenance Error Handling**: Returns 503 for database maintenance
- ✅ **Timeout Error Handling**: Returns 503 for database timeouts  
- ✅ **Centralized Logging**: Consistent error logging with operation context
- ✅ **JSON Error Responses**: Uses `WriteJSONError` for consistent response format

### 2. Not Found Error Handler
```go
// handleNotFoundError handles "not found" errors with a specific pattern
func (h *HandlerImpl) handleNotFoundError(w http.ResponseWriter, err error, resourceType, resourceID string) bool
```

**Features:**
- ✅ **Pattern Matching**: Detects "resource not found: ID" patterns
- ✅ **Proper HTTP Status**: Returns 404 for not found resources
- ✅ **Consistent Messages**: Standardized "Resource not found" messages

## Updated Handlers

### Core Transaction Handlers
- ✅ `CreateTransaction` - Database connection errors → 503
- ✅ `ListTransactions` - Database connection errors → 503  
- ✅ `GetTransaction` - Database connection errors → 503 + Not found → 404
- ✅ `UpdateTransaction` - Database connection errors → 503
- ✅ `DeleteTransaction` - Database connection errors → 503

### Balance Handlers  
- ✅ `CreateBalance` - Database connection errors → 503
- ✅ `ListBalances` - Database connection errors → 503
- ✅ `GetBalance` - Database connection errors → 503 + Not found → 404
- ✅ `UpdateBalance` - Database connection errors → 503
- ✅ `DeleteBalance` - Database connection errors → 503

## Usage Pattern

### Before (Inconsistent)
```go
func (h *HandlerImpl) GetBalance(w http.ResponseWriter, r *http.Request) {
    balance, err := h.Service.GetBalance(r.Context(), balanceID)
    if err != nil {
        logrus.WithError(err).Error("GetBalance failed")
        // Manual error detection and conversion - INCONSISTENT
        if isDatabaseConnectionError(err) {
            panic(&aws.DatabaseConnectionError{...})
        }
        if err.Error() == fmt.Sprintf("balance not found: %s", balanceID) {
            http.Error(w, "Balance not found", http.StatusNotFound)
        } else {
            http.Error(w, err.Error(), http.StatusInternalServerError)
        }
        return
    }
    // ... success handling
}
```

### After (Centralized & Consistent)
```go
func (h *HandlerImpl) GetBalance(w http.ResponseWriter, r *http.Request) {
    balance, err := h.Service.GetBalance(r.Context(), balanceID)
    if err != nil {
        // Try to handle as "not found" error first
        if h.handleNotFoundError(w, err, "balance", balanceID) {
            return
        }
        // Handle all other errors (including database connection errors)
        h.handleServiceError(w, err, "GetBalance")
        return
    }
    // ... success handling
}
```

## Error Flow

### Database Connection Drop Scenario
```
1. User Request → Handler → Service → Repository → getDB()
2. getDB() → aws.GetGormDB() → Connection Fails → Panic(DatabaseConnectionError)
3. Repository → Service → Handler (panic propagates)
4. Handler.handleServiceError() → Detects connection error → Re-panic with proper type
5. Maintenance Middleware → Catches panic → Returns 503 JSON response
```

### Expected Response
```json
{
  "error": {
    "code": "DB_TIMEOUT",
    "message": "Database is temporarily unavailable, please retry in a few minutes"
  }
}
```

## Benefits

### 1. **Consistency**
- All handlers use the same error handling pattern
- Uniform response formats across all endpoints
- Consistent logging with operation context

### 2. **Maintainability**  
- Single place to update error handling logic
- Easy to add new error types
- Reduced code duplication

### 3. **Proper HTTP Status Codes**
- 503 Service Unavailable for database connection issues (not 500)
- 404 Not Found for missing resources (not 500)
- JSON error responses instead of plain text

### 4. **Debugging**
- Centralized logging with operation names
- Consistent error context
- Proper error type conversion for middleware

## Testing

With database stopped, all endpoints should now return:
- **Status**: `503 Service Unavailable`
- **Content-Type**: `application/json`
- **Body**: Proper JSON error with error code and message

Instead of the previous:
- **Status**: `500 Internal Server Error`
- **Content-Type**: `text/plain`
- **Body**: Raw error message

## Next Steps

1. **Deploy the updated code**
2. **Test with database stopped**: All endpoints should return 503 JSON responses
3. **Test recovery**: Start database and verify endpoints work normally
4. **Monitor logs**: Verify proper error logging and middleware handling
