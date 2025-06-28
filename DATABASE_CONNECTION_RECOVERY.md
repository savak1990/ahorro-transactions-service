# Database Connection Recovery Guide

## How Database Connection Works in Lambda

### Connection Lifecycle in Single Lambda Instance

```
Lambda Cold Start → Lazy Init → First DB Request → Connection Success/Failure → Handle Subsequent Requests → Connection Recovery
```

## Current Implementation Details

### 1. Lazy Initialization (Fast Lambda Start)
- Lambda starts **quickly** without connecting to database
- Database connection is attempted **only when first request needs DB access**
- Health check (`/health`) returns `"not_configured"` without triggering connection

### 2. Connection Recovery Mechanism

#### Before (Problematic - sync.Once)
```go
sync.Once → First Failure → All subsequent requests fail until Lambda dies
```

#### After (Recovery-Enabled - Mutex + Retry Logic)
```go
RWMutex → Connection Health Check → Retry with Rate Limiting → Fresh Connection
```

## Request Scenarios

### Scenario 1: Database Stopped → Started During Lambda Life

```
┌─────────────┬──────────────────┬─────────────────┬────────────────┐
│ Time        │ Database Status  │ Request         │ Response       │
├─────────────┼──────────────────┼─────────────────┼────────────────┤
│ T0          │ Stopped          │ GET /health     │ 200 "not_conf" │
│ T1          │ Stopped          │ GET /transactions│ 503 (fast!)   │
│ T2          │ Stopped          │ GET /transactions│ 503 (rate limited) │
│ T3 (+30s)   │ Stopped          │ GET /transactions│ 503 (retry attempt) │
│ T4          │ Started          │ GET /transactions│ 200 (success!) │
│ T5          │ Started          │ GET /transactions│ 200 (cached)   │
└─────────────┴──────────────────┴─────────────────┴────────────────┘
```

### Scenario 2: Connection Drops Mid-Session

```
┌─────────────┬──────────────────┬─────────────────┬────────────────┐
│ Time        │ Database Status  │ Request         │ Response       │
├─────────────┼──────────────────┼─────────────────┼────────────────┤
│ T0          │ Running          │ GET /transactions│ 200 (success)  │
│ T1          │ Stopped          │ GET /transactions│ 503 (detected) │
│ T2          │ Stopped          │ GET /transactions│ 503 (rate limited) │
│ T3          │ Started          │ GET /transactions│ 200 (recovered!) │
└─────────────┴──────────────────┴─────────────────┴────────────────┘
```

### Scenario 3: Health Check States

```
┌──────────────────┬─────────────────────────────────────────────────┐
│ Database Status  │ Health Check Response                           │
├──────────────────┼─────────────────────────────────────────────────┤
│ Never connected  │ {"status": "healthy", "database": "not_configured"} │
│ Failed recently  │ {"status": "healthy", "database": "cooling_down"}   │
│ Ready to retry   │ {"status": "healthy", "database": "ready_to_retry"} │
│ Connected        │ {"status": "healthy", "database": "healthy"}        │
│ Connection lost  │ {"status": "unhealthy", "database": "disconnected"} │
└──────────────────┴─────────────────────────────────────────────────┘
```

## Key Features

### 1. **Fast Failure Detection** (2 seconds)
```go
connect_timeout=2 // Database connection timeout
PingContext(1s)   // Health check timeout
```

### 2. **Rate Limiting** (30 seconds cooldown)
```go
retryInterval = 30 * time.Second
```
Prevents overwhelming failed database with connection attempts.

### 3. **Connection Health Monitoring**
```go
// Before each request, check if connection is still alive
if err := sqlDB.PingContext(ctx); err == nil {
    return gormDB  // Use existing connection
}
// Otherwise, attempt reconnection
```

### 4. **Graceful Recovery**
- Closes broken connections properly
- Attempts fresh connection when database becomes available
- No Lambda restart required

## Manual Recovery

### Force Connection Reset
```bash
curl -X POST https://api-ahorro-transactions-savak.vkdev1.com/db-reset
```

Response:
```json
{
  "message": "Database connection reset - next request will attempt fresh connection",
  "status": "reset"
}
```

### Check Connection Status
```bash
curl https://api-ahorro-transactions-savak.vkdev1.com/health
```

## Database Management Commands

```bash
# Check status
make db-status

# Start database (takes ~30-60 seconds)
make db-start

# Stop database (immediate)
make db-stop

# Quick operations
make db-quick-start  # Start without waiting
make db-quick-stop   # Stop without confirmation
```

## Performance Characteristics

### Lambda Cold Start
- **Before**: 25+ seconds (hung waiting for database)
- **After**: ~2-3 seconds (lazy initialization)

### Request Response Times
- **Database Available**: ~200-500ms (normal)
- **Database Stopped**: ~2-3 seconds (fast failure)
- **During Rate Limit**: ~10ms (immediate 503)

### Recovery Time
- **Database Startup**: ~30-60 seconds (AWS RDS constraint)
- **First Success After Recovery**: ~2-3 seconds (connection establishment)
- **Subsequent Requests**: ~200-500ms (connection reuse)

## Cost Optimization

### Development Workflow
```bash
# Daily development cycle
make db-start        # Start when you begin work
# ... develop/test ...
make db-stop         # Stop when done (~$1.50/day saved)

# Quick testing
make db-quick-start  # Start for quick test
# ... test API ...
make db-quick-stop   # Stop immediately after
```

### Production Considerations
- Keep database running for production/staging
- Use db-stop/start cycle only for development/testing
- Connection recovery ensures zero-downtime when database starts

## Troubleshooting

### Common Issues

1. **Still getting 500 errors**
   - Check if maintenance middleware is active
   - Verify error handling in all endpoints

2. **Long response times**
   - Database might be starting up (wait 30-60s)
   - Check AWS RDS console for database status

3. **Rate limiting too aggressive**
   - Use `/db-reset` endpoint to force immediate retry
   - Adjust `retryInterval` if needed

4. **Connection not recovering**
   - Check database security groups
   - Verify database is in "available" state
   - Try manual reset: `curl -X POST .../db-reset`
