#!/bin/bash

# Authentication System - Local Development Setup Script
# This script helps set up the authentication system for local development

set -e

echo "🔐 Setting up Authentication System for Local Development"
echo "========================================================"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if we're in the right directory
if [ ! -f "docker-compose.yml" ]; then
    print_error "Please run this script from the project root directory"
    exit 1
fi

print_status "Starting PostgreSQL database..."
docker-compose up -d postgres

# Wait for PostgreSQL to be ready
print_status "Waiting for PostgreSQL to be ready..."
sleep 5

# Check if database is accessible
if ! docker-compose exec -T postgres pg_isready -U postgres > /dev/null 2>&1; then
    print_error "PostgreSQL is not ready. Please check the logs:"
    docker-compose logs postgres
    exit 1
fi

print_success "PostgreSQL is ready!"

# Create database if it doesn't exist
print_status "Creating database if it doesn't exist..."
docker-compose exec -T postgres psql -U postgres -c "CREATE DATABASE test_results;" 2>/dev/null || true

# Run authentication migration
print_status "Running authentication migration..."
if [ -f "db/add_auth_tables.sql" ]; then
    docker-compose exec -T postgres psql -U postgres -d test_results -f /docker-entrypoint-initdb.d/add_auth_tables.sql
    print_success "Authentication tables created!"
else
    print_error "Authentication migration file not found: db/add_auth_tables.sql"
    exit 1
fi

# Create .env file if it doesn't exist
if [ ! -f "api/.env" ]; then
    print_status "Creating .env file..."
    cat > api/.env << EOF
# Database
DATABASE_URL=postgres://postgres:password@localhost:5432/test_results?sslmode=disable

# Session Management
SESSION_SECRET=$(openssl rand -hex 32)

# OAuth2 Providers
# For local development, we'll use GitHub
GITHUB_CLIENT_ID=your_github_client_id
GITHUB_CLIENT_SECRET=your_github_client_secret

# Application Settings
BASE_URL=http://localhost:8080
ENVIRONMENT=development
EOF
    print_success ".env file created!"
    print_warning "Please update GITHUB_CLIENT_ID and GITHUB_CLIENT_SECRET in api/.env"
else
    print_status ".env file already exists"
fi

# Check if Go dependencies are installed
print_status "Checking Go dependencies..."
cd api
if ! go mod tidy; then
    print_error "Failed to tidy Go modules"
    exit 1
fi

# Check if frontend dependencies are installed
print_status "Checking frontend dependencies..."
cd ../frontend
if [ ! -d "node_modules" ]; then
    print_status "Installing frontend dependencies..."
    npm install
fi

print_success "Setup complete!"
echo ""
echo "🎉 Next steps:"
echo "1. Set up GitHub OAuth app:"
echo "   - Go to https://github.com/settings/developers"
echo "   - Create new OAuth app with callback URL: http://localhost:8080/auth/github/callback"
echo "   - Update GITHUB_CLIENT_ID and GITHUB_CLIENT_SECRET in api/.env"
echo ""
echo "2. Start the application:"
echo "   # Terminal 1 - API"
echo "   cd api && go run cmd/server/main.go"
echo ""
echo "   # Terminal 2 - Frontend"
echo "   cd frontend && npm run dev"
echo ""
echo "3. Access the application:"
echo "   - Frontend: http://localhost:3000"
echo "   - API: http://localhost:8080"
echo "   - API Docs: http://localhost:8080/swagger/"
echo ""
echo "📚 For more information, see docs/auth-development.md" 