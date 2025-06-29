# MySQL Readiness Wait Loop - Debug Enhancement

## Problem Fixed
The MySQL readiness wait loop in `setup-go-test-env/action.yml` lacked detailed logging when failures occurred, making it difficult to diagnose issues when MySQL failed to start within the timeout period.

❌ **Previous Issues:**
- No command tracing during the wait process
- Minimal error information on timeout
- No visibility into MySQL container status or logs
- Difficult to debug connection failures

## Enhancements Made

### 1. Command Tracing (Line 76)
```bash
set -o xtrace  # Enable command tracing for debugging
```
**Benefits:**
- ✅ Shows every command being executed
- ✅ Helps trace script execution flow
- ✅ Reveals exact command syntax and parameters

### 2. Enhanced Progress Logging (Lines 78-85)
```bash
echo "Starting MySQL readiness check..."
echo "Connection details: host=127.0.0.1, port=3306, user=root"
echo "Waiting for MySQL... (attempt $((counter + 1))/$max_wait)"
```
**Benefits:**
- ✅ Clear start indication with connection details
- ✅ Progress counter showing current attempt vs. maximum
- ✅ Better visibility into wait process

### 3. Periodic Detailed Checks (Lines 87-93)
```bash
# Show more detailed connection attempt every 10 seconds
if [ $((counter % 10)) -eq 0 ]; then
  echo "Attempting detailed MySQL connection test..."
  mysqladmin ping -h"127.0.0.1" -P"3306" -u"root" -p"password" || true
  echo "Checking if MySQL port is accessible..."
  nc -z 127.0.0.1 3306 && echo "Port 3306 is open" || echo "Port 3306 is not accessible"
fi
```
**Benefits:**
- ✅ Periodic detailed diagnostics during wait
- ✅ Port accessibility verification
- ✅ MySQL-specific connection test results

### 4. Comprehensive Failure Diagnostics (Lines 98-127)
When timeout occurs, the script now provides:

#### **Connection Diagnostics:**
```bash
echo "MySQL connection attempt result:"
mysqladmin ping -h"127.0.0.1" -P"3306" -u"root" -p"password" || true
echo "Port connectivity check:"
nc -z 127.0.0.1 3306 && echo "Port 3306 is accessible" || echo "Port 3306 is NOT accessible"
```

#### **Network Information:**
```bash
echo "Network connections:"
netstat -tlpn | grep :3306 || echo "No process listening on port 3306"
```

#### **Docker Container Status:**
```bash
echo "Available Docker containers:"
docker ps -a || echo "Could not list Docker containers"
```

#### **MySQL Container Logs:**
```bash
MYSQL_CONTAINER=$(docker ps --filter "ancestor=mysql:8.0" --format "{{.ID}}" | head -1)
if [ -n "$MYSQL_CONTAINER" ]; then
  echo "Found MySQL container: $MYSQL_CONTAINER"
  docker logs --tail 50 "$MYSQL_CONTAINER" || echo "Could not retrieve MySQL container logs"
else
  echo "MySQL container not found, trying alternative approach..."
  docker logs $(docker ps -q --filter "ancestor=mysql:8.0") 2>/dev/null || echo "No MySQL container logs available"
fi
```

### 5. Success Confirmation (Lines 131-132)
```bash
echo "✅ MySQL is ready and accepting connections!"
echo "Connection successful after $counter seconds"
```

## Debug Output Examples

### Success Case:
```
Starting MySQL readiness check...
Connection details: host=127.0.0.1, port=3306, user=root
+ mysqladmin ping -h127.0.0.1 -P3306 -uroot -p*** --silent
Waiting for MySQL... (attempt 1/60)
Waiting for MySQL... (attempt 2/60)
...
✅ MySQL is ready and accepting connections!
Connection successful after 15 seconds
```

### Failure Case:
```
Starting MySQL readiness check...
Connection details: host=127.0.0.1, port=3306, user=root
Waiting for MySQL... (attempt 1/60)
...
Attempting detailed MySQL connection test...
mysqladmin: connect to server at '127.0.0.1' failed
error: 'Can't connect to MySQL server on '127.0.0.1:3306' (111)'
Checking if MySQL port is accessible...
Port 3306 is not accessible
...
===============================================
ERROR: MySQL failed to start within 60 seconds
===============================================
Diagnostic information:
Current time: Thu Jan 1 12:00:00 UTC 2024
MySQL connection attempt result:
mysqladmin: connect to server at '127.0.0.1' failed
Port connectivity check:
Port 3306 is NOT accessible
Network connections:
No process listening on port 3306
Available Docker containers:
CONTAINER ID   IMAGE       STATUS
abc123def456   mysql:8.0   Exiting (1)
MySQL container logs (if available):
Found MySQL container: abc123def456
2024-01-01 12:00:00 [ERROR] Failed to initialize DD Storage Engine
2024-01-01 12:00:00 [ERROR] Plugin 'InnoDB' init function returned error
===============================================
```

## Benefits

✅ **Command Tracing** - See exactly what commands are executed  
✅ **Progress Visibility** - Clear indication of wait progress  
✅ **Periodic Diagnostics** - Regular health checks during wait  
✅ **Comprehensive Failure Info** - Detailed error context on timeout  
✅ **Container Log Access** - MySQL container logs for root cause analysis  
✅ **Network Diagnostics** - Port and process information  
✅ **Success Confirmation** - Clear indication when MySQL is ready  

## Debugging Workflow

1. **Monitor Progress** - Track wait attempts in real-time
2. **Check Periodic Diagnostics** - Review detailed checks every 10 seconds
3. **Analyze Failure Info** - If timeout occurs, review comprehensive diagnostics
4. **Examine Container Logs** - Check MySQL container logs for specific errors
5. **Verify Network Status** - Confirm port accessibility and process status

This enhanced debugging capability significantly improves the ability to diagnose and resolve MySQL startup issues in GitHub Actions workflows.