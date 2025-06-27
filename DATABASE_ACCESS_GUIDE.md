# Temporary Database Public Access Guide

## ⚠️ **SECURITY WARNING**
This guide shows how to temporarily make your Aurora database accessible from the internet for development/debugging purposes. **NEVER use this in production!**

## 🔧 Quick Setup Instructions

### 1. Get Your Public IP Address
```bash
make get-my-ip
```

### 2. Enable Public Database Access
```bash
make enable-db-public-access
```
This will:
- Prompt you to enter your IP address in CIDR format (e.g., `203.0.113.1/32`)
- Update the Aurora security group to allow your IP
- Set the Aurora instances to `publicly_accessible = true`
- Apply the Terraform changes

### 3. Connect to the Database
```bash
make connect-db
```
This will automatically connect you to the PostgreSQL database using psql.

### 4. **IMPORTANT: Disable Public Access When Done**
```bash
make disable-db-public-access
```

## 🛠️ Manual Database Connection

If you prefer to connect manually, use these details:

```bash
# Get connection details
make show-db-config

# Connect manually
PGPASSWORD="your_password" psql \
  --host=your-aurora-endpoint \
  --port=5432 \
  --username=your_username \
  --dbname=your_database_name
```

## 📊 Database Schema Exploration

Once connected, you can explore the PostgreSQL schema:

```sql
-- List all tables
\dt

-- View table structures
\d category
\d merchant
\d transaction
\d transaction_entry

-- Sample queries
SELECT * FROM category LIMIT 10;
SELECT * FROM transaction ORDER BY created_at DESC LIMIT 5;
```

## 🔍 Alternative Tools

### pgAdmin
1. Enable public access: `make enable-db-public-access`
2. Open pgAdmin and create a new server connection:
   - Host: Use `make get-db-endpoint`
   - Port: 5432
   - Database: Use `make get-db-name`
   - Username: Check Secrets Manager
   - Password: Check Secrets Manager

### DBeaver
1. Enable public access: `make enable-db-public-access`
2. Create new PostgreSQL connection with Aurora endpoint details

### TablePlus (macOS)
1. Enable public access: `make enable-db-public-access`
2. Create new PostgreSQL connection

## 🚨 Security Checklist

**Before enabling public access:**
- [ ] Confirm you're working in a development environment
- [ ] Have your specific IP address ready
- [ ] Plan to disable access when finished

**After enabling public access:**
- [ ] Verify only your IP can access the database
- [ ] Don't share the connection details
- [ ] Monitor AWS CloudTrail for unexpected access

**When finished:**
- [ ] Run `make disable-db-public-access`
- [ ] Verify the database is no longer publicly accessible
- [ ] Check security group rules are reverted

## 🔧 Troubleshooting

### Connection Issues
```bash
# Test if the endpoint is reachable
telnet $(make get-db-endpoint -s) $(make get-db-port)

# Check security group rules
aws ec2 describe-security-groups --group-ids sg-xxxxxxxxx
```

### IP Address Changes
If your IP changes, simply run:
```bash
make enable-db-public-access
```
And enter your new IP address.

### Access Denied
- Verify your IP is correct: `make get-my-ip`
- Check if public access is enabled in Aurora instance
- Verify security group allows your IP

## 📝 What the Configuration Does

### Security Group Changes
- Adds a dynamic ingress rule for port 5432
- Only allows your specific IP address
- Maintains existing Lambda access rules

### Aurora Instance Changes
- Sets `publicly_accessible = true` temporarily
- Allows connections from outside the VPC
- Keeps all other security configurations

### Terraform Variables
- `enable_db_public_access`: Boolean to control public access
- `my_ip_cidr`: Your IP address in CIDR format
- Both default to secure values (false and empty)

## 🔄 Automation Scripts

The Makefile includes these helpful targets:

| Target | Description |
|--------|-------------|
| `make get-my-ip` | Shows your current public IP |
| `make enable-db-public-access` | Enables public access with IP prompt |
| `make disable-db-public-access` | Disables public access |
| `make connect-db` | Direct psql connection |
| `make show-db-config` | Shows all connection details |

---
**Remember: Always disable public access when finished!**
