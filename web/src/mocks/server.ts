import { setupServer } from 'msw/node'
import { spotHandlers } from './handlers/spot-handlers'
import { authHandlers } from './handlers/auth-handlers'
import { reviewHandlers } from './handlers/review-handlers'

// Setup the server for Node.js environments (testing)
export const server = setupServer(
  ...spotHandlers,
  ...authHandlers,
  ...reviewHandlers
)

// Enable request mocking in Node.js
if (typeof process !== 'undefined' && process.env.NODE_ENV === 'test') {
  server.listen({
    onUnhandledRequest: 'warn',
  })
}