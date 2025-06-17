# 🤝 Development Handoff Guide

> **For the next Claude agent** - Everything you need to continue OAuth authentication implementation

## 📊 Current Status: 80% Complete

### ✅ COMPLETED (12/16 tasks)

**🏗️ Infrastructure**
- [x] Colima + Docker development environment  
- [x] MySQL container with automated setup
- [x] golang-migrate for database management
- [x] Environment variable configuration
- [x] Makefile automation (`make dev-setup`)

**🔧 Backend (Go + Huma)**
- [x] Complete Onion Architecture implementation
- [x] sqlc integration for type-safe SQL
- [x] User authentication API (`POST /api/users`)
- [x] Database schema (users, spots, reviews)
- [x] gRPC service layer with database integration

**🎨 Frontend (Next.js + Auth.js)**
- [x] Auth.js v5 configuration (Google/X OAuth)
- [x] Authentication UI (`/auth/signin`, `/auth/error`)
- [x] Header with authentication state management
- [x] Session management with useSession

### 🔄 REMAINING TASKS (4/16 tasks)

1. **Frontend-Backend Integration Testing** (Priority: HIGH)
2. **Live OAuth Credentials Setup** (Priority: MEDIUM)  
3. **E2E Test Updates** (Priority: MEDIUM)
4. **Full Integration Testing** (Priority: LOW)

---

## 🚀 5-Minute Quick Start

### 1. Start Development Environment

```bash
# Backend (Terminal 1)
cd api
make dev-setup          # Starts MySQL + migrations + API server
# ✅ API running at http://localhost:8080

# Frontend (Terminal 2)  
cd web
cp .env.local.example .env.local
# ⚠️ Add OAuth credentials (see step 2)
pnpm dev               # Starts Next.js with Turbopack
# ✅ Web app at http://localhost:3000
```

### 2. OAuth Credentials (Required for Testing)

**Google OAuth Setup:**
1. [Google Cloud Console](https://console.cloud.google.com/) → APIs & Services → Credentials
2. Create OAuth 2.0 Client ID
3. Add to authorized origins: `http://localhost:3000`
4. Add to redirect URIs: `http://localhost:3000/api/auth/callback/google`

**X (Twitter) OAuth Setup:**
1. [Twitter Developer Portal](https://developer.twitter.com/) → Create App
2. Add callback URL: `http://localhost:3000/api/auth/callback/twitter`

**Add to `web/.env.local`:**
```bash
GOOGLE_CLIENT_ID=your_google_client_id
GOOGLE_CLIENT_SECRET=your_google_client_secret
TWITTER_CLIENT_ID=your_twitter_client_id
TWITTER_CLIENT_SECRET=your_twitter_client_secret
NEXTAUTH_SECRET=your_random_secret_key
API_URL=http://localhost:8080
```

### 3. Test Authentication Flow

```bash
# 1. Visit http://localhost:3000
# 2. Click "ログイン" button
# 3. Click "Googleでログイン" or "Xでログイン"
# 4. Complete OAuth flow
# 5. Verify user appears in database:

docker exec bocchi-the-map-mysql mysql -u bocchi_user -pbocchi_password bocchi_the_map -e "SELECT * FROM users;"
```

---

## 🔧 Technical Architecture

### Authentication Flow
```
1. User clicks login → /auth/signin
2. User selects provider → Auth.js OAuth flow  
3. OAuth callback → Auth.js processes
4. Auth.js calls → POST /api/users (creates user)
5. User stored in MySQL → Session established
6. User redirected → / (authenticated)
```

### Database Schema
```sql
-- Users table (fully implemented)
CREATE TABLE users (
    id VARCHAR(36) PRIMARY KEY,
    anonymous_id VARCHAR(36),
    email VARCHAR(255) UNIQUE NOT NULL,
    display_name VARCHAR(100) NOT NULL,
    avatar_url TEXT,
    auth_provider ENUM('google', 'x') NOT NULL,
    auth_provider_id VARCHAR(255) NOT NULL,
    preferences JSON,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY unique_provider_user (auth_provider, auth_provider_id)
);
```

### Key Files Modified

**Backend:**
- `api/interfaces/http/handlers/user_handler.go` - User creation endpoint
- `api/infrastructure/grpc/user_service.go` - Database integration
- `api/infrastructure/database/` - sqlc generated code
- `api/migrations_new/000001_initial_schema.up.sql` - Database schema

**Frontend:**  
- `web/src/lib/auth.ts` - Auth.js configuration
- `web/src/components/header.tsx` - Authentication state management
- `web/src/app/auth/signin/page.tsx` - Sign-in page
- `web/src/app/auth/error/page.tsx` - Error handling

---

## 🐛 Known Issues & Solutions

### Docker/Database Issues
**Problem:** `docker-compose` not found  
**Solution:** Use Colima - `brew install colima && colima start`

**Problem:** MySQL connection failed  
**Solution:** Ensure Docker context - `docker context use colima`

**Problem:** Migration errors  
**Solution:** Check DATABASE_URL format in `.env`

### Authentication Issues  
**Problem:** OAuth redirect mismatch  
**Solution:** Verify callback URLs in provider console match `http://localhost:3000/api/auth/callback/{provider}`

**Problem:** NEXTAUTH_SECRET missing  
**Solution:** Generate secret - `openssl rand -base64 32`

### Frontend-Backend Connection
**Problem:** CORS errors  
**Solution:** API_URL in frontend `.env.local` should match backend port

---

## 📋 Next Steps Checklist

### Immediate Tasks (Next 30 minutes)

- [ ] **Test OAuth Flow End-to-End**
  - Set up Google OAuth credentials
  - Test complete login → user creation → session
  - Verify user data in MySQL database

- [ ] **Fix Frontend-Backend Integration**
  - Ensure Auth.js correctly calls `/api/users`
  - Debug any CORS or API connection issues
  - Test user creation via API endpoint

### Medium-term Tasks (Next 2 hours)

- [ ] **Update E2E Tests**
  - Modify Playwright tests for real authentication
  - Test login/logout functionality
  - Update test data and expectations

- [ ] **Error Handling Improvements**
  - Test various OAuth error scenarios
  - Improve error messages and UX
  - Add better loading states

### Future Enhancements

- [ ] **User Profile Management**
  - Add user settings page
  - Implement preference updates
  - Add avatar upload functionality

- [ ] **Advanced Authentication**
  - Add email verification
  - Implement account linking
  - Add two-factor authentication

---

## 📞 Quick Reference

### Useful Commands
```bash
# Backend
make dev-setup          # Full environment setup
make docker-up          # Start MySQL only
make migrate-up         # Run migrations
make run               # Start API server

# Frontend  
pnpm dev               # Start Next.js dev server
pnpm test:e2e          # Run Playwright tests
pnpm build             # Production build

# Database
docker exec bocchi-the-map-mysql mysql -u bocchi_user -pbocchi_password bocchi_the_map
```

### Important URLs
- **Frontend:** http://localhost:3000
- **Backend API:** http://localhost:8080  
- **API Docs:** http://localhost:8080/docs
- **Sign-in Page:** http://localhost:3000/auth/signin

### Environment Files
- `api/.env` - Backend configuration (TiDB credentials included)
- `api/.env.local` - Local MySQL override  
- `web/.env.local` - Frontend OAuth credentials (needs setup)

---

## 🎯 Success Metrics

**You'll know the implementation is complete when:**

1. ✅ User can click "ログイン" → select provider → complete OAuth
2. ✅ User data appears in MySQL `users` table
3. ✅ Header shows authenticated user name/avatar
4. ✅ User can logout and login again successfully  
5. ✅ E2E tests pass with real authentication flow

**Expected completion time:** 2-4 hours for remaining tasks

---

**🚀 Ready to continue! The foundation is solid, just need to connect the final pieces.**