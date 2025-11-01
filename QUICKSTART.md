# Quick Start Guide

## Get Started

### Step 1: Start the Database
```bash
docker-compose up -d postgres
```

### Step 2: Run the Application
```bash
go run cmd/server/main.go
```

### Step 3: Test the API

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

**Port already in use?**
```bash
# Change PORT in .env file
PORT=8081
```

**Database connection error?**
```bash
docker-compose ps
```

**Module not found?**
```bash
go mod tidy
```

