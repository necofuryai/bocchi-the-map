# 🎨 Bocchi The Map Web

> **Modern React frontend with Next.js 15** - Server Components, Turbopack, and cutting-edge UX patterns

[![Next.js](https://img.shields.io/badge/Next.js-15.3.2-000000?style=flat&logo=next.js)](https://nextjs.org/)
[![React](https://img.shields.io/badge/React-19.0.0-61DAFB?style=flat&logo=react)](https://reactjs.org/)
[![TypeScript](https://img.shields.io/badge/TypeScript-5.0+-3178C6?style=flat&logo=typescript)](https://www.typescriptlang.org/)
[![Tailwind CSS](https://img.shields.io/badge/Tailwind-4.0+-06B6D4?style=flat&logo=tailwindcss)](https://tailwindcss.com/)

A **performance-first PWA** for solo location discovery, built with React Server Components, edge-optimized rendering, and design system excellence.

## ⚡ Quick Start

```bash
# Prerequisites: Node.js 20+, pnpm
pnpm install                    # Install dependencies (auto-installs Playwright)

# 🔐 Setup authentication (required)
cp .env.local.example .env.local
# Add your OAuth credentials (see Authentication Setup below)

pnpm dev                        # Start with Turbopack 🚀

# App ready at http://localhost:3000
# Lightning-fast HMR with Turbopack
```

## 🔐 Authentication Status

### ✅ PRODUCTION READY - SUPABASE AUTH MIGRATION COMPLETE (2025-07-13)

- ✅ **Supabase Auth**: Modern email/password authentication with email verification
- ✅ **Authentication UI**: Complete login/signup pages with modern UX
- ✅ **Session Management**: Cookie-based sessions with automatic token refresh
- ✅ **SSR Support**: Full Next.js App Router compatibility
- ✅ **Protected Routes**: Authentication guards and middleware integration
- ✅ **User Profile**: Header with profile dropdown and logout
- ✅ **Backend Integration**: API authentication with Supabase tokens
- ✅ **Type Safety**: Complete TypeScript integration for auth flows

### 🎯 RECENT MAJOR MIGRATION (2025-07-13)

- **Complete Migration**: Migrated from Auth0 to Supabase Auth
- **Modernization**: SSR-compatible authentication with Next.js App Router
- **Simplified Setup**: Reduced configuration complexity and improved DX
- **Enhanced Security**: Cookie-based sessions with automatic token management

### ✅ READY FOR PRODUCTION

- Complete authentication flow with email/password and email verification
- All user management features functional with modern architecture
- SSR-compatible authentication for optimal performance
- Type-safe authentication context and components

### Supabase Auth Setup Required

**Supabase Project Setup:**
1. Go to [Supabase Dashboard](https://supabase.com/dashboard)
2. Create a new project
3. Navigate to Authentication > Settings
4. Enable Email provider and configure settings

**Environment Configuration:**

Add credentials to `web/.env.local`:
```bash
NEXT_PUBLIC_SUPABASE_URL=your_supabase_project_url
NEXT_PUBLIC_SUPABASE_ANON_KEY=your_supabase_anon_key
```

## 🚀 Modern React Patterns

### Next.js 15 App Router

```tsx
// app/layout.tsx - Root layout with providers
export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="ja" suppressHydrationWarning>
      <body>
        <AuthProvider>
          <ThemeProvider>
            {children}
          </ThemeProvider>
        </AuthProvider>
      </body>
    </html>
  )
}
```

### React Server Components

```tsx
// app/spots/page.tsx - Zero-JS by default
async function SpotsPage() {
  const spots = await getSpots() // Runs on server
  
  return (
    <div>
      <SpotsList spots={spots} />
      <InteractiveMap /> {/* Client component only when needed */}
    </div>
  )
}
```

### Streaming & Suspense

```tsx
// Progressive loading with granular suspense
<Suspense fallback={<MapSkeleton />}>
  <MapView />
</Suspense>
<Suspense fallback={<SpotListSkeleton />}>
  <SpotsList />
</Suspense>
```

## 🏗️ Architecture Features

### Performance-First Design

- **Turbopack** - 700x faster than Webpack for dev builds
- **React Server Components** - Zero-JS server-rendered components
- **Streaming SSR** - Progressive page hydration
- **Edge Runtime** - Deploy to 300+ edge locations globally
- **Automatic Code Splitting** - Ship only what users need

### Developer Experience

- **TypeScript** - End-to-end type safety with API contracts  
- **Tailwind CSS 4.0** - Modern CSS with design tokens
- **shadcn/ui** - Copy-paste component system
- **ESLint + Prettier** - Consistent code formatting
- **Hot Module Replacement** - Instant feedback loop

### Production-Ready

- **Progressive Web App** - Native-like mobile experience
- **Dark Mode** - System preference aware theming
- **Internationalization** - Japanese/English support
- **Accessibility** - WCAG 2.1 AA compliance
- **SEO Optimized** - Meta tags, structured data, sitemap

## 📁 Project Structure

```text
web/
├── 🎯 src/app/                  # 📱 APP ROUTER (Next.js 15)
│   ├── layout.tsx               # Root layout + providers
│   ├── page.tsx                 # Home page (RSC)
│   ├── spots/                   # Spot discovery pages
│   ├── profile/                 # User profile pages
│   └── api/                     # API routes (Edge Runtime)
├── 🧩 src/components/           # 🎨 REUSABLE COMPONENTS
│   ├── ui/                      # shadcn/ui primitives
│   ├── auth-provider.tsx        # Authentication context
│   ├── theme-provider.tsx       # Dark mode support
│   └── map/                     # Map-related components
├── 🎣 src/hooks/                # ⚡ CUSTOM REACT HOOKS
│   ├── useLocalStorage.ts       # Persistent client state
│   ├── useDebounce.ts           # Input debouncing
│   └── useMapControls.ts        # Map interaction logic
├── 🔧 src/lib/                  # 🛠️ UTILITIES & CONFIG
│   ├── auth.ts                  # Supabase Auth configuration
│   ├── utils.ts                 # Shared utility functions
│   └── validations.ts           # Zod schemas
├── 🎨 src/styles/               # 💄 GLOBAL STYLES
│   └── globals.css              # Tailwind base + custom CSS
└── 🏗️ public/                  # 📦 STATIC ASSETS
    ├── icons/                   # PWA icons + favicons
    └── images/                  # Optimized images
```

## ⚡ Performance Optimizations

### Bundle Analysis

```bash
# Analyze bundle size
pnpm build && pnpm analyze

# Performance profiling
pnpm dev -- --experimental-https
lighthouse https://localhost:3000
```

### Core Web Vitals

| Metric | Target | Strategy |
|--------|--------|----------|
| **LCP** | < 1.2s | Image optimization, edge caching |
| **FID** | < 100ms | Code splitting, minimal JS |
| **CLS** | < 0.1 | Reserved space, font loading |

### Optimization Techniques

- **Image Optimization** - Next.js Image component with WebP/AVIF
- **Font Loading** - `next/font` with display swap
- **Lazy Loading** - Components and routes on-demand
- **Service Worker** - Offline support and caching
- **CDN Caching** - Static assets via Cloudflare

## 🎨 Design System

### Tailwind CSS 4.0 Features

```css
/* Design tokens with CSS variables */
@theme {
  --color-primary-50: #fef2f2;
  --color-primary-500: #ef4444;
  --color-primary-900: #7f1d1d;
}

/* Container queries */
@container (min-width: 768px) {
  .card { @apply p-8; }
}
```

### Component System

```tsx
// Type-safe component variants
const Button = cva("px-4 py-2 rounded-md", {
  variants: {
    variant: {
      primary: "bg-blue-500 text-white",
      secondary: "bg-gray-200 text-gray-900"
    },
    size: {
      sm: "text-sm px-3 py-1",
      lg: "text-lg px-6 py-3"
    }
  }
})
```

### Dark Mode Implementation

```tsx
// System-aware dark mode
function ThemeProvider({ children }: { children: React.ReactNode }) {
  return (
    <NextThemesProvider
      attribute="class"
      defaultTheme="system"
      enableSystem
      disableTransitionOnChange
    >
      {children}
    </NextThemesProvider>
  )
}
```

## 🗺️ Map Integration

### MapLibre GL JS

```tsx
// High-performance vector maps
import maplibregl from 'maplibre-gl'

function MapView() {
  const mapRef = useRef<maplibregl.Map>()
  
  useEffect(() => {
    mapRef.current = new maplibregl.Map({
      container: 'map',
      style: '/api/maps/style.json', // Custom vector style
      center: [139.7671, 35.6812],   // Tokyo
      zoom: 12
    })
  }, [])
  
  return <div id="map" className="w-full h-96" />
}
```

### Custom Map Features

- **Vector Tiles** - Smooth zooming and styling
- **Clustering** - Efficient rendering of thousands of spots
- **Custom Markers** - React components as map overlays
- **Geolocation** - Find spots near user location
- **Offline Support** - Cached tiles for offline browsing

## 🔐 Authentication & Security

### Supabase Auth

```tsx
// Supabase Auth with SSR support
import { createClient } from '@/utils/supabase/client'

// Email/password authentication
export const signInWithPassword = async (email: string, password: string) => {
  const supabase = createClient()
  const { data, error } = await supabase.auth.signInWithPassword({
    email,
    password,
  })
  return { data, error }
}

export const signUpWithPassword = async (email: string, password: string) => {
  const supabase = createClient()
  const { data, error } = await supabase.auth.signUp({
    email,
    password,
    options: {
      emailRedirectTo: `${window.location.origin}/auth/callback`,
    },
  })
  return { data, error }
}

export const signOut = async () => {
  const supabase = createClient()
  await supabase.auth.signOut()
}
```

### Security Features

- **CSRF Protection** - Built-in token validation
- **Content Security Policy** - XSS attack prevention
- **Secure Headers** - HSTS, frame options, etc.
- **Input Validation** - Zod schemas for type safety
- **Rate Limiting** - API endpoint protection

## 🧪 Testing Strategy

### Testing Pyramid

```bash
# Unit tests with Vitest
pnpm test:unit

# Component tests with Testing Library
pnpm test:components

# E2E tests with Playwright
pnpm test:e2e

# Visual regression tests
pnpm test:visual
```

### Test Examples

```tsx
// Component testing with Testing Library
import { render, screen } from '@testing-library/react'
import { SpotCard } from './SpotCard'

test('displays spot information correctly', () => {
  const spot = { name: 'Shibuya Coffee', rating: 4.5 }
  render(<SpotCard spot={spot} />)
  
  expect(screen.getByText('Shibuya Coffee')).toBeInTheDocument()
  expect(screen.getByText('4.5')).toBeInTheDocument()
})

// E2E testing with Playwright
test('user can search for spots', async ({ page }) => {
  await page.goto('/')
  await page.fill('[data-testid=search-input]', 'coffee')
  await page.click('[data-testid=search-button]')
  
  await expect(page.locator('[data-testid=spot-list]')).toBeVisible()
})
```

## 🚀 Development Workflow

### Local Development

```bash
# Development server with Turbopack
pnpm dev

# Development with HTTPS (for PWA testing)
pnpm dev:https

# Type checking
pnpm type-check

# Linting and formatting
pnpm lint                   # ESLint + TypeScript checking
pnpm format                 # Prettier formatting
```

### Environment Configuration

```bash
# .env.local
NEXT_PUBLIC_API_URL=http://localhost:8080
NEXT_PUBLIC_MAP_STYLE_URL=/api/maps/style.json

# Supabase Authentication
NEXT_PUBLIC_SUPABASE_URL=your-supabase-project-url
NEXT_PUBLIC_SUPABASE_ANON_KEY=your-supabase-anon-key
SUPABASE_SERVICE_ROLE_KEY=your-supabase-service-role-key

# Analytics (optional)
NEXT_PUBLIC_GA_ID=your-google-analytics-id
```

## 📱 Progressive Web App

### PWA Features

- **Installable** - Add to home screen on mobile/desktop
- **Offline Support** - Service worker caching
- **Push Notifications** - New spot alerts (future)
- **Background Sync** - Queue actions when offline
- **Share Target** - Receive shared locations

### Manifest Configuration

```json
{
  "name": "Bocchi The Map",
  "short_name": "Bocchi Map",
  "display": "standalone",
  "start_url": "/",
  "theme_color": "#3b82f6",
  "background_color": "#ffffff",
  "icons": [
    {
      "src": "/icons/icon-192.png",
      "sizes": "192x192",
      "type": "image/png"
    }
  ]
}
```

## 🌐 Internationalization

### Multi-language Support

```tsx
// Using next-intl for i18n
import { useTranslations } from 'next-intl'

function SpotCard({ spot }: { spot: Spot }) {
  const t = useTranslations('spots')
  
  return (
    <div>
      <h3>{spot.name}</h3>
      <p>{t('rating')}: {spot.rating}</p>
    </div>
  )
}
```

### Supported Locales

- 🇯🇵 **Japanese** - Primary market
- 🇺🇸 **English** - International users
- 🔄 **Future**: Korean, Chinese (Traditional/Simplified)

## 📊 Analytics & Monitoring

### Performance Monitoring

```tsx
// Core Web Vitals tracking
import { getCLS, getFID, getFCP, getLCP, getTTFB } from 'web-vitals'

function sendToAnalytics(metric: any) {
  // Send to Google Analytics, Vercel Analytics, etc.
  gtag('event', metric.name, {
    value: Math.round(metric.name === 'CLS' ? metric.value * 1000 : metric.value),
    event_label: metric.id,
  })
}

getCLS(sendToAnalytics)
getFID(sendToAnalytics)
getFCP(sendToAnalytics)
getLCP(sendToAnalytics)
getTTFB(sendToAnalytics)
```

### Error Tracking

- **Sentry** - Production error monitoring
- **Vercel Analytics** - Performance insights
- **Console Logs** - Development debugging
- **User Feedback** - In-app error reporting

## 🚢 Deployment

### Vercel

```bash
# Build and deploy
pnpm build
pnpm export

# Environment variables
NEXT_PUBLIC_API_URL=https://api.bocchi-map.com
NEXT_PUBLIC_SUPABASE_URL=https://your-project.supabase.co
```

### Performance Optimization

- **Edge Functions** - API routes at the edge
- **Static Generation** - Pre-built pages for performance
- **Image Optimization** - Automatic WebP/AVIF conversion
- **CDN Caching** - Global content distribution

## 🤝 Contributing

### Code Standards

- **TypeScript** - Strict mode enabled
- **ESLint** - Airbnb + Next.js rules
- **Prettier** - Consistent formatting
- **Husky** - Pre-commit hooks
- **Conventional Commits** - Semantic versioning

### Component Guidelines

```tsx
// Component structure example
interface SpotCardProps {
  spot: Spot
  onClick?: (spot: Spot) => void
  className?: string
}

export function SpotCard({ spot, onClick, className }: SpotCardProps) {
  return (
    <div className={cn("p-4 border rounded-lg", className)}>
      {/* Component content */}
    </div>
  )
}
```

## 🎯 Roadmap

- [x] **v1.0** - Core spot discovery with server components
- [x] **v1.1** - Map integration with MapLibre GL JS
- [ ] **v1.2** - Real-time collaboration features
- [ ] **v1.3** - Advanced PWA capabilities
- [ ] **v2.0** - React Native mobile app with shared components

---

**🎨 Crafted with modern React patterns for exceptional UX**

[🌟 Live Demo](https://bocchi-map.com) • [📖 Storybook](https://storybook.bocchi-map.com) • [🎨 Design System](https://design.bocchi-map.com)
