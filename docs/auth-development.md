# Authentication System - Local Development Guide

This guide explains how to set up and use the authentication system for local development.

## Overview

The authentication system supports two authentication methods:
- **OAuth2 Authentication** - For web users (GitHub for development, Okta for production)
- **API Key Authentication** - For CLI tools and CI/CD systems (Jenkins, etc.)

## Prerequisites

- Go 1.21+
- PostgreSQL database
- GitHub account (for OAuth2 development)
- Docker and Docker Compose (optional, for easy setup)

## Local Development Setup

### 1. Database Setup

First, ensure your database is running and has the authentication tables:

```bash
# If using Docker Compose
docker-compose up -d postgres

# Run the authentication migration
psql -h localhost -U postgres -d test_results -f db/add_auth_tables.sql
```

### 2. Environment Variables

Create a `.env` file in the `api` directory with the following variables:

```bash
# Database
DATABASE_URL=postgres://postgres:password@localhost:5432/test_results?sslmode=disable

# Session Management
SESSION_SECRET=your-super-secret-session-key-change-this-in-production

# OAuth2 Providers
# For local development, we'll use GitHub
GITHUB_CLIENT_ID=your_github_client_id
GITHUB_CLIENT_SECRET=your_github_client_secret

# Application Settings
BASE_URL=http://localhost:8080
ENVIRONMENT=development
```

### 3. GitHub OAuth2 Setup (for local development)

1. Go to [GitHub Developer Settings](https://github.com/settings/developers)
2. Click "New OAuth App"
3. Fill in the details:
   - **Application name**: `Test Results Local Dev`
   - **Homepage URL**: `http://localhost:3000`
   - **Authorization callback URL**: `http://localhost:8080/auth/github/callback`
4. Copy the Client ID and Client Secret to your `.env` file

### 4. Running the Application

```bash
# Start the API server
cd api
go run cmd/server/main.go

# In another terminal, start the frontend
cd frontend
npm install
npm run dev
```

The application will be available at:
- **Frontend**: http://localhost:3000
- **API**: http://localhost:8080
- **API Documentation**: http://localhost:8080/swagger/

## Authentication Flow

### OAuth2 Authentication (Web Users)

1. **Login**: User clicks "Continue with GitHub" on the login page
2. **Redirect**: User is redirected to GitHub for authorization
3. **Callback**: GitHub redirects back to `/auth/github/callback`
4. **Session Creation**: A session is created and stored in the database
5. **Cookie Set**: A `session_id` cookie is set in the browser
6. **Access**: User can now access protected routes

### API Key Authentication (CLI/Jenkins)

1. **Create API Key**: User creates an API key through the web interface
2. **Use API Key**: CLI tools or Jenkins use the API key in the Authorization header:
   ```
   Authorization: Bearer your_api_key_here
   ```

## API Endpoints

### Authentication Endpoints

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/auth/github` | Start GitHub OAuth2 flow | No |
| GET | `/auth/github/callback` | GitHub OAuth2 callback | No |
| POST | `/auth/logout` | Logout and clear session | Yes |
| GET | `/auth/me` | Get current user info | Yes |
| GET | `/auth/api-keys` | List user's API keys | Yes |
| POST | `/auth/api-keys` | Create new API key | Yes |
| DELETE | `/auth/api-keys/{id}` | Delete API key | Yes |

### Protected Endpoints

All endpoints under `/api/` require authentication:
- `/api/projects/*`
- `/api/builds/*`
- `/api/test-suites/*`
- `/api/test-cases/*`
- etc.

## Development Workflow

### 1. Testing OAuth2 Authentication

1. Start the application
2. Navigate to http://localhost:3000
3. Click "Continue with GitHub"
4. Authorize the application on GitHub
5. You should be redirected back and logged in

### 2. Testing API Key Authentication

1. Log in through the web interface
2. Go to your profile page (http://localhost:3000/profile)
3. Create a new API key
4. Copy the plain text key (you won't see it again!)
5. Use it in your CLI tools or API requests:

```bash
# Example API request with API key
curl -H "Authorization: Bearer your_api_key_here" \
     http://localhost:8080/api/projects
```

### 3. Testing Protected Routes

```bash
# This should fail without authentication
curl http://localhost:8080/api/projects

# This should work with API key
curl -H "Authorization: Bearer your_api_key_here" \
     http://localhost:8080/api/projects

# This should work with session cookie (after web login)
curl -H "Cookie: session_id=your_session_id" \
     http://localhost:8080/api/projects
```

## Frontend Development

### Authentication Context

The frontend uses React Context for authentication state management:

```typescript
import { useAuth } from '../context/AuthContext';

const MyComponent = () => {
  const { user, isAuthenticated, login, logout } = useAuth();
  
  if (!isAuthenticated) {
    return <LoginPage />;
  }
  
  return <div>Welcome, {user?.name}!</div>;
};
```

### Protected Routes

Use the `ProtectedRoute` component to protect pages:

```typescript
import ProtectedRoute from '../components/auth/ProtectedRoute';

<Route path="/dashboard" element={
  <ProtectedRoute>
    <DashboardPage />
  </ProtectedRoute>
} />
```

## Troubleshooting

### Common Issues

1. **"OAuth2 provider not found"**
   - Check that your GitHub OAuth app is properly configured
   - Verify the callback URL matches exactly

2. **"Session not found"**
   - Clear browser cookies and try again
   - Check that the database is running and accessible

3. **"API key invalid"**
   - Make sure you're using the correct API key format
   - Check that the API key hasn't expired
   - Verify the Authorization header format: `Bearer your_key`

4. **Database connection issues**
   - Ensure PostgreSQL is running
   - Check the DATABASE_URL in your .env file
   - Verify the authentication tables exist

### Debug Mode

Enable debug logging by setting the log level:

```bash
export LOG_LEVEL=debug
go run cmd/server/main.go
```

### Database Inspection

Check authentication data in the database:

```sql
-- Check users
SELECT * FROM users;

-- Check sessions
SELECT * FROM auth_sessions;

-- Check API keys
SELECT * FROM api_keys;
```

## Security Considerations

### For Local Development

1. **Session Secret**: Use a strong, unique session secret
2. **OAuth App**: Create a separate OAuth app for development
3. **Database**: Use a local database, not production data
4. **HTTPS**: Local development typically uses HTTP, but be aware of security implications

### For Production

1. **Environment Variables**: Use proper secret management
2. **HTTPS**: Always use HTTPS in production
3. **Session Security**: Configure secure cookie settings
4. **API Key Rotation**: Implement API key rotation policies
5. **Rate Limiting**: Add rate limiting to authentication endpoints

## Next Steps

1. **Add More OAuth Providers**: Implement additional OAuth2 providers (Google, Microsoft, etc.)
2. **Role-Based Access Control**: Add user roles and permissions
3. **Audit Logging**: Log authentication events for security monitoring
4. **Password Reset**: Implement password reset functionality for OAuth users
5. **Multi-Factor Authentication**: Add MFA support for enhanced security

## Support

If you encounter issues:
1. Check the troubleshooting section above
2. Review the application logs
3. Verify your environment configuration
4. Check the API documentation at http://localhost:8080/swagger/ 