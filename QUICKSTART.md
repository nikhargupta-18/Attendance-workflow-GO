# Quick Start Guide

## Prerequisites
- Go 1.16 or higher
- Docker and Docker Compose
- Make (optional)

## Get Started

### Step 1: Set Up Environment
```bash
# Copy example environment file (if not exists)
cp .env.example .env

# Install dependencies
go mod download
go mod tidy
```

### Step 2: Start the Services
```bash
# Start only database and redis (recommended for development)
docker-compose up -d postgres redis

# Or start all services including API (alternative)
docker-compose up -d
```

### Step 3: Run the Application
Choose ONE of these methods:

#### A. Run Locally (Recommended for Development)
```bash
# First, ensure the API container is not running
docker-compose stop api

# Then run the application locally
go run cmd/server/main.go
```

#### B. Run in Docker (Alternative)
```bash
# Run everything in Docker
docker-compose up -d
```

### Step 4: Verify Services
The following endpoints should be available:

- API: http://localhost:8080
- Health Check: http://localhost:8080/health
- Swagger UI: http://localhost:8080/swagger/index.html

### Step 5: Test the API

#### Login as Admin
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@university.edu",
    "password": "admin123"
  }'
```

#### Register a Student
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@student.edu",
    "password": "password123",
    "role": "student",
    "dept": "Computer Science"
  }'
```

## Next Steps

1. Apply for Leave - Use the token from login
2. Mark Attendance - As faculty/warden
3. View Analytics - As admin
4. Check Notifications - For leave updates

## Full Documentation

See [README.md](README.md) for complete API documentation.

## Troubleshooting

### Common Issues

**Port 8080 already in use?**
1. Check if the API container is running:
```bash
docker-compose ps
```
2. Stop the API container if running:
```bash
docker-compose stop api
```
3. Or change the port in .env:
```bash
PORT=8081
```

**Database connection error?**
1. Verify services are running:
```bash
docker-compose ps
```
2. Check logs:
```bash
docker-compose logs postgres
docker-compose logs redis
```

**Module not found?**
```bash
go mod tidy
go mod download
```

**Redis connection issues?**
1. Verify Redis is running:
```bash
docker-compose ps redis
```
2. Check Redis logs:
```bash
docker-compose logs redis
```

