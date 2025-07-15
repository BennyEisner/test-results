# System Patterns

## Architecture

The system follows **Hexagonal Architecture** (Ports and Adapters pattern) with three main layers:

1. **Domain Layer** - Core business logic and entities
2. **Application Layer** - Use cases and application services  
3. **Infrastructure Layer** - External concerns (databases, HTTP, etc.)

### Directory Structure Pattern

```
api/internal/{module}/
├── domain/
│   ├── models/          # Domain entities
│   ├── ports/           # Port interfaces
│   └── errors/          # Domain-specific errors
├── application/
│   └── service.go       # Application services (use cases)
└── infrastructure/
    ├── database/        # Repository implementations
    ├── http/           # HTTP handlers
    └── middleware/     # HTTP middleware
```

## Authentication Patterns

### OAuth2 Flow
1. User clicks "Continue with GitHub"
2. Redirect to GitHub for authorization
3. GitHub redirects to `/auth/github/callback`
4. Session created and stored in database
5. `session_id` cookie set in browser

### API Key Authentication
- CLI tools and CI/CD systems use Bearer tokens
- Format: `Authorization: Bearer {api_key}`
- Keys are managed through web interface

## Dependency Injection

Uses a container pattern to wire up dependencies:
- Repositories (Secondary/Driven Adapters)
- Application Services (Use Cases)
- HTTP Handlers (Primary/Driving Adapters)

## Frontend Patterns

- **React Context** for authentication state management
- **Protected Routes** component for route protection
- **TypeScript** for type safety
- **Vite** for build tooling

## Database Patterns

- PostgreSQL with migration scripts
- Separate tables for users, sessions, and API keys
- Repository pattern for data access
