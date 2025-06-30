# Auth0 Environment Setup Guide

This guide provides step-by-step instructions for setting up Auth0 environment variables for testing the Bocchi The Map application.

## Prerequisites

- Auth0 account (free tier available at [auth0.com](https://auth0.com/))
- Node.js and npm/pnpm installed
- Go installed (for backend development)
- OpenSSL (for generating secrets)

## Environment Variables Overview

### Backend API Variables (`api/.env`)

| Variable | Required | Description | Example |
|----------|----------|-------------|---------|
| `JWT_SECRET` | ✅ Yes | JWT signing secret (min 32 chars) | Generated with `openssl rand -base64 32` |
| `AUTH0_DOMAIN` | ✅ Yes | Auth0 tenant domain | `your-tenant.auth0.com` |
| `AUTH0_AUDIENCE` | ✅ Yes | API identifier from Auth0 | `bocchi-the-map-api` |
| `AUTH0_CLIENT_ID` | ✅ Yes | Auth0 application client ID | `abc123...` |
| `AUTH0_CLIENT_SECRET` | ⚠️ Prod Only | Auth0 application client secret | `xyz789...` |

### Frontend Variables (`web/.env.local`)

| Variable | Required | Description | Example |
|----------|----------|-------------|---------|
| `AUTH0_SECRET` | ✅ Yes | NextAuth.js secret (min 32 chars) | Generated with `openssl rand -base64 32` |
| `APP_BASE_URL` | ✅ Yes | Application base URL | `http://localhost:3000` |
| `AUTH0_DOMAIN` | ✅ Yes | Auth0 tenant domain | `your-tenant.auth0.com` |
| `AUTH0_CLIENT_ID` | ✅ Yes | Auth0 application client ID | `abc123...` |
| `AUTH0_CLIENT_SECRET` | ✅ Yes | Auth0 application client secret | `xyz789...` |
| `AUTH0_AUDIENCE` | 🔧 Optional | API identifier (for API access) | `bocchi-the-map-api` |
| `AUTH0_SCOPE` | 🔧 Optional | OAuth scopes | `openid profile email` |
| `API_URL` | ✅ Yes | Backend API endpoint (server-side) | `http://localhost:8080` |
| `NEXT_PUBLIC_API_URL` | ✅ Yes | Backend API endpoint (client-side) | `http://localhost:8080` |

## Step-by-Step Auth0 Setup

### 1. Create Auth0 Account and Tenant

1. Visit [Auth0 Dashboard](https://manage.auth0.com/)
2. Sign up for a free account or log in
3. Create a new tenant (or use existing one)
4. Note your tenant domain: `your-tenant.auth0.com`

### 2. Create Auth0 Application

1. In Auth0 Dashboard, go to **Applications** → **Applications**
2. Click **+ Create Application**
3. Choose **Regular Web Application**
4. Name: `Bocchi The Map`
5. Click **Create**

### 3. Configure Application Settings

In your application's **Settings** tab:

#### Basic Information
- Copy **Domain**, **Client ID**, and **Client Secret**

#### Application URIs
```
Allowed Callback URLs:
http://localhost:3000/auth/callback
https://your-production-domain.com/auth/callback

Allowed Logout URLs:
http://localhost:3000
https://your-production-domain.com

Allowed Web Origins:
http://localhost:3000
https://your-production-domain.com
```

#### Advanced Settings
- **Grant Types**: Ensure `Authorization Code` is enabled
- **Token Expiration**: Set appropriate values for your needs

### 4. Create Auth0 API (Optional but Recommended)

1. Go to **Applications** → **APIs**
2. Click **+ Create API**
3. Name: `Bocchi The Map API`
4. Identifier: `bocchi-the-map-api` (this becomes your `AUTH0_AUDIENCE`)
5. Signing Algorithm: `RS256`
6. Click **Create**

### 5. Generate Secrets

Generate secure secrets for JWT signing:

```bash
# For AUTH0_SECRET (frontend)
openssl rand -base64 32

# For JWT_SECRET (backend)
openssl rand -base64 32
```

## Environment File Setup

### Backend Configuration

Copy `api/.env.example` to `api/.env` and update:

```env
# Server Configuration
PORT=8080
HOST=0.0.0.0
ENV=development

# Auth0 Configuration
JWT_SECRET=your-generated-jwt-secret-here
AUTH0_DOMAIN=your-tenant.auth0.com
AUTH0_AUDIENCE=bocchi-the-map-api
AUTH0_CLIENT_ID=your-auth0-client-id
AUTH0_CLIENT_SECRET=your-auth0-client-secret

# Database Configuration
TIDB_HOST=localhost
TIDB_PORT=3306
TIDB_USER=bocchi_user
TIDB_PASSWORD=your-db-password
TIDB_DATABASE=bocchi_the_map
```

### Frontend Configuration

Copy `web/.env.example` to `web/.env.local` and update:

```env
# Authentication (Auth0)
AUTH0_SECRET=your-generated-auth0-secret-here
APP_BASE_URL=http://localhost:3000
AUTH0_DOMAIN=your-tenant.auth0.com
AUTH0_CLIENT_ID=your-auth0-client-id
AUTH0_CLIENT_SECRET=your-auth0-client-secret
AUTH0_AUDIENCE=bocchi-the-map-api
AUTH0_SCOPE=openid profile email

# Backend API Configuration
API_URL=http://localhost:8080
NEXT_PUBLIC_API_URL=http://localhost:8080

# Map Configuration
NEXT_PUBLIC_MAP_STYLE_URL=https://pub-f7098f2a137d4fcc854d717d48a53615.r2.dev/worldmap.pmtiles
```

## Testing the Setup

### 1. Start the Backend

```bash
cd api
go run cmd/api/main.go
```

Verify Auth0 configuration loads correctly by checking the logs.

### 2. Start the Frontend

```bash
cd web
npm run dev
# or
pnpm dev
```

### 3. Test Authentication Flow

1. Visit `http://localhost:3000`
2. Click login button
3. Should redirect to Auth0 login page
4. After successful login, should return to app with user session

### 4. Test API Authentication

```bash
# Get access token from frontend developer tools
# Look for Auth0 token in localStorage or network requests

# Test authenticated API endpoint
curl -H "Authorization: Bearer YOUR_TOKEN" \
     http://localhost:8080/api/v1/protected-endpoint
```

## Common Issues and Solutions

### JWT_SECRET Validation Errors

**Error**: `JWT_SECRET must be at least 32 characters long`

**Solution**: Generate a proper secret:
```bash
openssl rand -base64 32
```

### Auth0 Domain Format Errors

**Error**: `AUTH0_DOMAIN should not include protocol`

**Solution**: Use domain only, not full URL:
- ✅ Correct: `your-tenant.auth0.com`
- ❌ Wrong: `https://your-tenant.auth0.com`

### Callback URL Mismatch

**Error**: `Invalid callback URL`

**Solution**: Ensure callback URLs in Auth0 dashboard exactly match your app URLs.

### CORS Issues

**Error**: Cross-origin requests blocked

**Solution**: Add your frontend domain to Auth0 **Allowed Web Origins**.

## Security Best Practices

### Development Environment
- Use different secrets for development and production
- Never commit `.env.local` or `.env` files to version control
- Rotate secrets regularly

### Production Environment
- Use environment-specific secrets
- Enable Auth0 security features (MFA, anomaly detection)
- Monitor authentication logs
- Use HTTPS for all callback URLs

### Secret Management
- Store production secrets in secure environment variable systems
- Use different Auth0 tenants for development and production
- Implement proper secret rotation policies

## Production Deployment

When deploying to production:

1. **Create Production Auth0 Tenant/Application**
2. **Update Environment Variables** with production values
3. **Configure Production Callback URLs** in Auth0 dashboard
4. **Use HTTPS** for all URLs
5. **Enable Auth0 Security Features**
6. **Monitor Authentication Metrics**

## Integration Architecture

The Auth0 integration follows this architecture pattern:

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   Next.js Web   │    │   Auth0 Service  │    │   Go API        │
│                 │    │                  │    │                 │
│ • Auth Button   │◄──►│ • Login/Logout   │    │ • JWT Validator │
│ • User Profile  │    │ • Token Issuing  │◄──►│ • Auth Handler  │
│ • Auth Guard    │    │ • User Info      │    │ • Middleware    │
│ • Middleware    │    │                  │    │ • Rate Limiter  │
└─────────────────┘    └──────────────────┘    └─────────────────┘
         │                                               │
         └───────────────────────────────────────────────┘
                    Direct API Communication
```

## Authentication Flow Details

### Login Flow
1. User clicks login button → Next.js Auth0 middleware
2. Redirect to Auth0 login page
3. Auth0 authentication → JWT token issued
4. Token stored in session cookies
5. User redirected back to application

### API Request Flow
1. Frontend makes API request with session
2. Next.js middleware extracts JWT token
3. Token sent to Go API in Authorization header
4. Go JWT middleware validates token against Auth0 JWKS
5. User context added to request
6. Protected resource accessed

### Logout Flow
1. User clicks logout → Next.js logout endpoint
2. Session cleared from cookies
3. Optional: API call to blacklist token (future enhancement)
4. User redirected to login page

## Dependencies Analysis

### Frontend Dependencies
- `@auth0/nextjs-auth0`: Core Auth0 integration
- `next`: Next.js framework
- `react`: React framework
- `typescript`: Type safety

### Backend Dependencies
- `github.com/golang-jwt/jwt`: JWT token handling
- `github.com/danielgtaylor/huma/v2`: API framework
- Database drivers and ORM

## Implementation Status

### Completed Components
- **Frontend**: Auth0 client, middleware, React components, API routes, authentication pages
- **Backend**: Authentication handler, JWT middleware, auth service, JWT validator
- **Database**: User migrations, SQL queries, generated Go code
- **Configuration**: Environment setup, CORS, type definitions

### Test Results
- **Status**: 33/34 tests passing (97% success rate)
- **Coverage**: Component structure, build system, dependencies, implementation, configuration

## Additional Resources

- [Auth0 Next.js SDK Documentation](https://auth0.com/docs/quickstart/webapp/nextjs)
- [Auth0 Go SDK Documentation](https://auth0.com/docs/quickstart/backend/golang)
- [Auth0 Dashboard](https://manage.auth0.com/)
- [JWT.io](https://jwt.io/) for token debugging