# InSavein Frontend

Frontend application for the InSavein financial discipline platform built with TanStack Start, React, TypeScript, and Tailwind CSS.

## Features

- **Authentication**: User registration, login, logout with JWT token management
- **Protected Routes**: Automatic authentication checks and redirects
- **API Client**: Type-safe API client with error handling and retry logic
- **Responsive Design**: Mobile-first design with Tailwind CSS
- **Server-Side Rendering**: Fast initial page loads with TanStack Start

## Project Structure

```
frontend/
├── src/
│   ├── components/        # Reusable React components
│   │   └── ProtectedRoute.tsx
│   ├── lib/              # Core libraries and utilities
│   │   ├── api/          # API client modules
│   │   │   ├── client.ts
│   │   │   ├── auth.ts
│   │   │   ├── user.ts
│   │   │   ├── savings.ts
│   │   │   ├── budget.ts
│   │   │   ├── goals.ts
│   │   │   ├── education.ts
│   │   │   ├── analytics.ts
│   │   │   ├── notifications.ts
│   │   │   └── index.ts
│   │   ├── auth/         # Authentication context and storage
│   │   │   ├── context.tsx
│   │   │   └── storage.ts
│   │   └── types/        # TypeScript type definitions
│   │       └── api.ts
│   ├── routes/           # TanStack Router routes
│   │   ├── __root.tsx
│   │   ├── index.tsx
│   │   ├── login.tsx
│   │   ├── register.tsx
│   │   └── dashboard.tsx
│   ├── router.tsx        # Router configuration
│   └── styles.css        # Global styles and Tailwind config
├── .env                  # Environment variables
├── .env.example          # Environment variables template
├── package.json
├── tsconfig.json
└── vite.config.ts
```

## Getting Started

### Prerequisites

- Node.js 18+ and npm
- Backend services running (see root README)

### Installation

1. Install dependencies:
```bash
npm install
```

2. Copy environment variables:
```bash
cp .env.example .env
```

3. Update `.env` with your backend service URLs (defaults to localhost:808X)

### Development

Start the development server:
```bash
npm run dev
```

The application will be available at http://localhost:3000

### Building for Production

Build the application:
```bash
npm run build
```

Preview the production build:
```bash
npm run preview
```

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `VITE_AUTH_SERVICE_URL` | Auth service URL | http://localhost:8081 |
| `VITE_USER_SERVICE_URL` | User service URL | http://localhost:8082 |
| `VITE_SAVINGS_SERVICE_URL` | Savings service URL | http://localhost:8083 |
| `VITE_BUDGET_SERVICE_URL` | Budget service URL | http://localhost:8084 |
| `VITE_GOAL_SERVICE_URL` | Goal service URL | http://localhost:8085 |
| `VITE_EDUCATION_SERVICE_URL` | Education service URL | http://localhost:8086 |
| `VITE_NOTIFICATION_SERVICE_URL` | Notification service URL | http://localhost:8087 |
| `VITE_ANALYTICS_SERVICE_URL` | Analytics service URL | http://localhost:8088 |

## API Client

The API client is organized into service modules:

- **AuthService**: Registration, login, logout, token refresh
- **UserService**: Profile management, preferences
- **SavingsService**: Savings transactions, streaks, history
- **BudgetService**: Budget management, spending tracking, alerts
- **GoalService**: Financial goals, milestones, progress
- **EducationService**: Lessons, progress tracking
- **AnalyticsService**: Spending analysis, recommendations, financial health
- **NotificationService**: Notifications, read status

All services include:
- Type-safe request/response interfaces
- Automatic JWT token injection
- Error handling with retry logic
- Unauthorized (401) handling with automatic logout

## Authentication

The application uses JWT-based authentication with:

- **Access tokens**: 15-minute expiry
- **Refresh tokens**: 7-day expiry
- **Automatic token refresh**: Checks every minute and refreshes when needed
- **Secure storage**: Tokens stored in localStorage
- **Protected routes**: Automatic redirect to login for unauthenticated users

### Using Authentication

```tsx
import { useAuth } from '../lib/auth/context'

function MyComponent() {
  const { user, isAuthenticated, login, logout, api } = useAuth()
  
  // Access user info
  console.log(user?.email)
  
  // Make authenticated API calls
  const savings = await api.savings.getSummary()
  
  return <div>...</div>
}
```

## Testing

Run tests:
```bash
npm test
```

## Tech Stack

- **TanStack Start**: Full-stack React framework with SSR
- **React 19**: UI library
- **TypeScript**: Type safety
- **Tailwind CSS**: Utility-first styling
- **TanStack Router**: Type-safe routing
- **Vite**: Build tool and dev server

## License

MIT
