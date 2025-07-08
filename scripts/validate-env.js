#!/usr/bin/env node

/**
 * Environment Variables Validation Script
 * 
 * This script validates that all required Auth0 environment variables
 * are properly configured for both frontend and backend.
 * 
 * Usage:
 *   node scripts/validate-env.js
 */

const fs = require('fs');
const path = require('path');

// Color codes for terminal output
const colors = {
  red: '\x1b[31m',
  green: '\x1b[32m',
  yellow: '\x1b[33m',
  blue: '\x1b[34m',
  cyan: '\x1b[36m',
  reset: '\x1b[0m',
  bold: '\x1b[1m'
};

function log(message, color = 'reset') {
  console.log(`${colors[color]}${message}${colors.reset}`);
}

function loadEnvFile(filePath) {
  if (!fs.existsSync(filePath)) {
    return null;
  }
  
  const content = fs.readFileSync(filePath, 'utf8');
  const env = {};
  
  content.split('\n').forEach(line => {
    const trimmed = line.trim();
    if (trimmed && !trimmed.startsWith('#')) {
      const [key, ...valueParts] = trimmed.split('=');
      if (key && valueParts.length > 0) {
        env[key.trim()] = valueParts.join('=').trim();
      }
    }
  });
  
  return env;
}

function validateBackendEnv() {
  log('\n📡 Validating Backend Environment (.env)', 'cyan');
  log('=' .repeat(50), 'cyan');
  
  const envPath = path.join(__dirname, '..', 'api', '.env');
  const env = loadEnvFile(envPath);
  
  if (!env) {
    log('❌ Backend .env file not found!', 'red');
    log('   Create api/.env from api/.env.example', 'yellow');
    return false;
  }
  
  const requiredVars = {
    'JWT_SECRET': {
      required: true,
      minLength: 32,
      description: 'JWT signing secret'
    },
    'AUTH0_DOMAIN': {
      required: true,
      pattern: /^[a-zA-Z0-9-]+\.auth0\.com$/,
      description: 'Auth0 tenant domain'
    },
    'AUTH0_AUDIENCE': {
      required: true,
      description: 'Auth0 API audience'
    },
    'AUTH0_CLIENT_ID': {
      required: true,
      description: 'Auth0 client ID'
    },
    'AUTH0_CLIENT_SECRET': {
      required: false, // Only required in production
      description: 'Auth0 client secret (required in production)'
    }
  };
  
  let allValid = true;
  
  for (const [varName, config] of Object.entries(requiredVars)) {
    const value = env[varName];
    
    if (config.required && (!value || value.startsWith('your-'))) {
      log(`❌ ${varName}: Missing or placeholder value`, 'red');
      log(`   ${config.description}`, 'yellow');
      allValid = false;
    } else if (value && config.minLength && value.length < config.minLength) {
      log(`❌ ${varName}: Too short (min ${config.minLength} chars)`, 'red');
      allValid = false;
    } else if (value && config.pattern && !config.pattern.test(value)) {
      log(`❌ ${varName}: Invalid format`, 'red');
      if (varName === 'AUTH0_DOMAIN') {
        log('   Should be: your-tenant.auth0.com (no https://)', 'yellow');
      }
      allValid = false;
    } else if (value) {
      log(`✅ ${varName}: Configured`, 'green');
    } else if (!config.required) {
      log(`⚠️  ${varName}: Optional (${config.description})`, 'yellow');
    }
  }
  
  return allValid;
}

function validateFrontendEnv() {
  log('\n🌐 Validating Frontend Environment (.env.local)', 'cyan');
  log('=' .repeat(50), 'cyan');
  
  const envPath = path.join(__dirname, '..', 'web', '.env.local');
  const env = loadEnvFile(envPath);
  
  if (!env) {
    log('❌ Frontend .env.local file not found!', 'red');
    log('   Create web/.env.local from web/.env.example', 'yellow');
    return false;
  }
  
  const requiredVars = {
    'AUTH0_SECRET': {
      required: true,
      minLength: 32,
      description: 'NextAuth.js secret for JWT signing'
    },
    'APP_BASE_URL': {
      required: true,
      pattern: /^https?:\/\/.+/,
      description: 'Application base URL'
    },
    'AUTH0_DOMAIN': {
      required: true,
      pattern: /^[a-zA-Z0-9-]+\.auth0\.com$/,
      description: 'Auth0 tenant domain'
    },
    'AUTH0_CLIENT_ID': {
      required: true,
      description: 'Auth0 client ID'
    },
    'AUTH0_CLIENT_SECRET': {
      required: true,
      description: 'Auth0 client secret'
    },
    'API_URL': {
      required: true,
      pattern: /^https?:\/\/.+/,
      description: 'Backend API URL (server-side)'
    },
    'NEXT_PUBLIC_API_URL': {
      required: true,
      pattern: /^https?:\/\/.+/,
      description: 'Backend API URL (client-side)'
    },
    'AUTH0_AUDIENCE': {
      required: false,
      description: 'Auth0 API audience (optional)'
    },
    'AUTH0_SCOPE': {
      required: false,
      description: 'OAuth scopes (optional)'
    }
  };
  
  let allValid = true;
  
  for (const [varName, config] of Object.entries(requiredVars)) {
    const value = env[varName];
    
    if (config.required && (!value || value.startsWith('your-'))) {
      log(`❌ ${varName}: Missing or placeholder value`, 'red');
      log(`   ${config.description}`, 'yellow');
      allValid = false;
    } else if (value && config.minLength && value.length < config.minLength) {
      log(`❌ ${varName}: Too short (min ${config.minLength} chars)`, 'red');
      allValid = false;
    } else if (value && config.pattern && !config.pattern.test(value)) {
      log(`❌ ${varName}: Invalid format`, 'red');
      if (varName === 'AUTH0_DOMAIN') {
        log('   Should be: your-tenant.auth0.com (no https://)', 'yellow');
      } else if (varName.includes('URL')) {
        log('   Should include protocol: http:// or https://', 'yellow');
      }
      allValid = false;
    } else if (value) {
      log(`✅ ${varName}: Configured`, 'green');
    } else if (!config.required) {
      log(`⚠️  ${varName}: Optional (${config.description})`, 'yellow');
    }
  }
  
  return allValid;
}

function showNextSteps() {
  log('\n🚀 Next Steps:', 'blue');
  log('=' .repeat(50), 'blue');
  log('1. Follow AUTH0_SETUP_GUIDE.md for detailed Auth0 configuration', 'blue');
  log('2. Generate secrets: openssl rand -base64 32', 'blue');
  log('3. Start backend: cd api && go run cmd/api/main.go', 'blue');
  log('4. Start frontend: cd web && npm run dev', 'blue');
  log('5. Test login at http://localhost:3000', 'blue');
}

function main() {
  log('🔧 Bocchi The Map - Environment Validation', 'bold');
  log('Checking Auth0 configuration...', 'cyan');
  
  const backendValid = validateBackendEnv();
  const frontendValid = validateFrontendEnv();
  
  log('\n📊 Summary:', 'bold');
  log('=' .repeat(50), 'cyan');
  
  if (backendValid && frontendValid) {
    log('✅ All environment variables are properly configured!', 'green');
    log('🎉 Ready to start development!', 'green');
    showNextSteps();
  } else {
    log('❌ Some environment variables need attention', 'red');
    log('\n💡 Quick fixes:', 'yellow');
    if (!backendValid) {
      log('   - Copy api/.env.example to api/.env', 'yellow');
      log('   - Update placeholder values with real Auth0 credentials', 'yellow');
    }
    if (!frontendValid) {
      log('   - Copy web/.env.example to web/.env.local', 'yellow');
      log('   - Update placeholder values with real Auth0 credentials', 'yellow');
    }
    log('\n📖 See AUTH0_SETUP_GUIDE.md for detailed instructions', 'blue');
    process.exit(1);
  }
}

if (require.main === module) {
  main();
}

module.exports = { validateBackendEnv, validateFrontendEnv };