/**
 * Auth0 configuration for Bocchi The Map
 * 
 * In v4 of @auth0/nextjs-auth0, we create an Auth0Client instance
 * and use it with middleware for automatic route handling.
 */

import { Auth0Client } from "@auth0/nextjs-auth0/server";

// Create Auth0 client instance
export const auth0 = new Auth0Client({
  // Configuration is handled via environment variables
  // AUTH0_DOMAIN, AUTH0_CLIENT_ID, AUTH0_CLIENT_SECRET, AUTH0_SECRET, APP_BASE_URL
  authorizationParameters: {
    scope: process.env.AUTH0_SCOPE || 'openid profile email',
    audience: process.env.AUTH0_AUDIENCE,
  }
});