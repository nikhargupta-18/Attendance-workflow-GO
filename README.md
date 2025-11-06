# Attendance Workflow API

RESTful API for managing student attendance and leave requests. Built with Go, PostgreSQL, and Docker.

## Features

- JWT Authentication with role-based access (Admin, Faculty, Warden, Student)
- Attendance tracking and marking
- Leave request management with approval workflow
- Notifications system
- Analytics dashboard
- Swagger API documentation

## Quick Start

### 1. Start Database

```bash
docker-compose up -d postgres
```

### 2. Run Server

```bash
go mod tidy
go run cmd/server/main.go
```

Server runs on `http://localhost:8080`

**Default Admin:**
- Email: `admin@university.edu`
- Password: `admin123`

## Getting Bearer Tokens

### ⭐ Recommended: Use Swagger UI (Easiest)

1. **Open Swagger:** http://localhost:8080/swagger/index.html
2. **Find `POST /api/v1/auth/login`** and click "Try it out"
3. **Enter credentials:**
   ```json
   {
     "email": "admin@university.edu",
     "password": "admin123"
   }
   ```
4. **Click "Execute"** - Copy the `token` from response
5. **Click "Authorize" button** (🔒 top right), paste token, click "Authorize"
6. **Now you can test all endpoints** directly in Swagger!

### Alternative: PowerShell (Windows)

**Login:**
```powershell
$body = '{"email":"admin@university.edu","password":"admin123"}'
$response = Invoke-RestMethod -Uri "http://localhost:8080/api/v1/auth/login" `
  -Method POST -ContentType "application/json" -Body $body
$token = $response.token
Write-Host "Token: $token"
```

**Use Token:**
```powershell
Invoke-RestMethod -Uri "http://localhost:8080/api/v1/users" `
  -Headers @{"Authorization"="Bearer $token"}
```

### Alternative: curl (Linux/Mac)

```bash
# Login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@university.edu","password":"admin123"}'

# Use token
curl -X GET http://localhost:8080/api/v1/users \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

**Note:** Don't paste URLs in browser - they only do GET requests. Use Swagger UI or the commands above.


## API Endpoints

**Base URL:** `/api/v1`

### Authentication

- `POST /auth/register` - Register new user
- `POST /auth/login` - Login and get token

### Users (Protected)

- `GET /users` - List users (Admin only, paginated)
- `GET /users/:id` - Get user by ID
- `PUT /users/:id` - Update user
- `DELETE /users/:id` - Delete user (Admin only)

**Query params:** `page`, `limit`, `role`, `dept`

### Leaves (Protected)

- `POST /leaves/apply` - Apply for leave (Student)
- `GET /leaves/my` - Get my leaves
- `GET /leaves/pending` - Get pending leaves (Faculty/Warden/Admin)
- `PUT /leaves/:id/approve` - Approve/reject leave (Faculty/Warden/Admin)
- `GET /leaves` - Get all leaves (Admin only)
- `DELETE /leaves/:id` - Delete leave (Admin only)

**Leave types:** `Medical`, `Personal`, `Emergency`, `Other`

**Date format:** Use `YYYY-MM-DD` (e.g., `2025-11-05`)

### Attendance (Protected)

- `POST /attendance/mark` - Mark attendance
- `GET /attendance/student/:id` - Get student attendance
- `GET /attendance/my` - Get my attendance (Student)
- `GET /attendance/daily` - Get daily attendance

**Query params:** `start_date`, `end_date`, `date` (all in `YYYY-MM-DD` format)

**Date format for JSON:** Use `YYYY-MM-DD` (e.g., `2025-11-05`)

### Notifications (Protected)

- `GET /notifications/my` - Get my notifications
- `PUT /notifications/:id/read` - Mark as read
- `GET /notifications/unread-count` - Get unread count

### Analytics (Admin Only)

- `GET /analytics/dashboard` - Dashboard stats
- `GET /analytics/leave-breakdown` - Leave breakdown
- `GET /analytics/department?dept=CS` - Department stats
- `GET /analytics/absentees` - Frequent absentees

### Health Check

- `GET /health` - Server health status

## Configuration

Create `.env` file (optional - defaults provided):

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=attendance_user
DB_PASSWORD=attendance_pass
DB_NAME=attendance_db
PORT=8080
JWT_SECRET=your-secret-key-change-in-production
JWT_EXPIRY=24h
```

## Troubleshooting

**Port 8080 in use:**
```bash
PORT=8081 go run cmd/server/main.go
```

**Database connection error:**
```bash
docker-compose up -d postgres
docker-compose logs postgres
```

**Swagger shows "No operations":**
```bash
go install github.com/swaggo/swag/cmd/swag@latest
swag init -g cmd/server/main.go -o docs
```

**Getting 404 on login endpoint:**
- Don't paste the URL in browser (browsers only do GET requests)
- Use Swagger UI: http://localhost:8080/swagger/index.html
- Or use PowerShell/curl commands shown above

**Module errors:**
```bash
go mod tidy
go mod download
```

## Project Structure

```
attendance-workflow/
├── cmd/server/       # Entry point
├── internal/
│   ├── api/         # Routes
│   ├── auth/        # Authentication
│   ├── users/       # User handlers
│   ├── leaves/      # Leave handlers
│   ├── attendance/  # Attendance handlers
│   ├── notifications/
│   └── analytics/
├── pkg/
│   ├── config/      # Configuration
│   └── db/          # Database models
└── docs/            # Swagger docs
```

## Example Workflow

**Using Swagger UI (Recommended):**
1. Open http://localhost:8080/swagger/index.html
2. Login via `POST /api/v1/auth/login`
3. Click "Authorize" and paste token
4. Test endpoints directly in Swagger

**Using PowerShell:**
```powershell
# Login and get token
$body = '{"email":"admin@university.edu","password":"admin123"}'
$response = Invoke-RestMethod -Uri "http://localhost:8080/api/v1/auth/login" `
  -Method POST -ContentType "application/json" -Body $body
$token = $response.token

# Use token
Invoke-RestMethod -Uri "http://localhost:8080/api/v1/users" `
  -Headers @{"Authorization"="Bearer $token"}
```

## Links

- **API:** http://localhost:8080
- **Swagger:** http://localhost:8080/swagger/index.html
- **Health:** http://localhost:8080/health
