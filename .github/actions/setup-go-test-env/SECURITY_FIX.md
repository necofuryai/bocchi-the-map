# Security Fix: MySQL Root Password

## Problem Fixed
The `setup-go-test-env` action previously had a hardcoded default password of `'password'` which posed a security risk.

## Changes Made

### 1. Action Configuration (`.github/actions/setup-go-test-env/action.yml`)
- ❌ **Removed**: `default: 'password'` from mysql-root-password input
- ✅ **Added**: `required: true` to enforce explicit password passing
- 🔒 **Security**: No more plaintext default passwords in code

### 2. Workflow Configuration (`.github/workflows/bdd-tests.yml`)
- ❌ **Replaced**: `MYSQL_ROOT_PASSWORD: password` 
- ✅ **With**: `MYSQL_ROOT_PASSWORD: ${{ secrets.MYSQL_ROOT_PASSWORD }}`
- ❌ **Replaced**: `root:password@tcp(localhost:3306)/bocchi_test...`
- ✅ **With**: `root:${{ secrets.MYSQL_ROOT_PASSWORD }}@tcp(localhost:3306)/bocchi_test...`
- ✅ **Added**: `mysql-root-password: ${{ secrets.MYSQL_ROOT_PASSWORD }}` to action calls

## Required GitHub Secret

**⚠️ IMPORTANT**: You must create a repository secret to complete this security fix:

1. Go to: **Settings** → **Secrets and variables** → **Actions**
2. Click **New repository secret**
3. Name: `MYSQL_ROOT_PASSWORD`
4. Value: A secure password (e.g., generated with `openssl rand -base64 32`)
5. Save the secret

## Security Benefits

✅ **No hardcoded passwords** in source code  
✅ **Explicit password requirement** prevents accidental defaults  
✅ **Secret-based authentication** for CI/CD  
✅ **Improved compliance** with security best practices  

## Testing

After creating the `MYSQL_ROOT_PASSWORD` secret, the BDD tests workflow will:
- Use the secure password for MySQL service
- Pass the password securely to the setup action
- Connect to the database with the secret password

## Rollback (if needed)

If you need to temporarily rollback this change:
```yaml
# In action.yml - TEMPORARY ONLY
mysql-root-password:
  description: 'MySQL root password for database connection'
  required: false
  default: 'temp_password_123'
```

**Note**: Always remove temporary passwords before committing to production!