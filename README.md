# Campus Leave & Attendance Management System

A modern, scalable REST API for managing student attendance and leave requests in educational institutions. Built with Go, PostgreSQL, and Redis.

## Key Features

### 👥 User Management
- Role-based access (Admin, Faculty, Warden, Student)
- Secure JWT authentication
- Profile management

### 📊 Attendance System
- Real-time attendance marking
- Individual & department tracking
- Historical records
- Absence notifications

### 📝 Leave Management
- Multiple leave types (Medical, Personal, Emergency)
- Multi-level approval workflow
- Status tracking & notifications
- Leave history

### 📈 Analytics Dashboard
- Department-wise attendance stats
- Leave analysis and trends
- Absentee monitoring
- Monthly summary reports

### 🔔 Smart Notifications
- Real-time updates via WebSocket
- Async email notifications
- Scheduled reminders
- Notification center

## Tech Stack

- **Backend**: Go 1.21 (Gin Framework)
- **Database**: PostgreSQL
- **Cache/Queue**: Redis
- **Documentation**: Swagger
- **Authentication**: JWT
- **Container**: Docker

## Demo Accounts and Tokens

For testing and demonstration purposes, you can use the following accounts. These accounts are created automatically when you run the database reset script (`go run scripts/reset_db.go`).

### Demo Credentials

1. **Admin User**
   - Name: admin
   - Email: admin@college.edu
   - Password: admin123

2. **Warden**
   - Name: warden
   - Email: warden@college.edu
   - Password: warden123

3. **Faculty**
   - Name: faculty
   - Email: faculty@college.edu
   - Password: faculty123

4. **Student 1**
   - Name: stud1
   - Email: stud1@college.edu
   - Password: student123

5. **Student 2**
   - Name: stud2
   - Email: stud2@college.edu
   - Password: student123

Note: Bearer tokens for all users are automatically generated when you run the database reset script. Each token is valid for 24 hours.

### Usage

1. Reset the database and get fresh tokens:
```bash
go run scripts/reset_db.go
```

2. Use the generated tokens in API requests:
```bash
curl -H "Authorization: Bearer <token>" http://localhost:8080/api/v1/users/profile
```

The reset script will:
- Reset the database with fresh demo data
- Create all demo users with their roles
- Generate and display new bearer tokens for each user
- Tokens are valid for 24 hours from generation time

## Quick Start

### Using Docker (Recommended)

```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f api
```

### Local Development

```bash
# Start database
docker-compose up -d postgres

# Install dependencies
go mod tidy

# Run server
go run cmd/server/main.go
```

**Default Admin Credentials:**
- Email: admin@university.edu
- Password: admin123

## API Documentation (Swagger)

Once the server is running, access the interactive Swagger UI:

```
http://localhost:8080/swagger/index.html
```

**Features:**
- View all API endpoints organized by tags
- Try out endpoints directly from the browser
- See request/response schemas
- Authenticate using Bearer token (click "Authorize" button)

**To regenerate docs after adding annotations:**
```bash
go install github.com/swaggo/swag/cmd/swag@latest
swag init -g cmd/server/main.go -o docs
```

## API Endpoints

### Authentication

**Register**
```http
POST /api/v1/auth/register
Content-Type: application/json

{
  "name": "John Doe",
  "email": "john@student.edu",
  "password": "password123",
  "role": "student",
  "dept": "Computer Science"
}
```

**Login**
```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "email": "admin@university.edu",
  "password": "admin123"
}
```

### User Management (Admin Only)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/users` | List all users (with pagination, role, dept filters) |
| GET | `/api/v1/users/:id` | Get user by ID |
| PUT | `/api/v1/users/:id` | Update user |
| DELETE | `/api/v1/users/:id` | Delete user |

### Leave Management

**Apply for Leave (Student)**
```http
POST /api/v1/leaves/apply
Authorization: Bearer <token>
Content-Type: application/json

{
  "leave_type": "Medical",
  "reason": "Doctor appointment",
  "start_date": "2025-10-25T00:00:00Z",
  "end_date": "2025-10-25T00:00:00Z"
}
```

**Leave Types:** Medical, Personal, Emergency, Other

**Other Endpoints:**
| Method | Endpoint | Description | Role |
|--------|----------|-------------|------|
| GET | `/api/v1/leaves/my` | Get my leave requests | All |
| GET | `/api/v1/leaves/pending` | Get pending leaves | Faculty/Warden/Admin |
| PUT | `/api/v1/leaves/:id/approve` | Approve/reject leave | Faculty/Warden/Admin |
| GET | `/api/v1/leaves` | Get all leaves | Admin |
| DELETE | `/api/v1/leaves/:id` | Delete leave | Admin |

**Approve/Reject Leave**
```http
PUT /api/v1/leaves/:id/approve
Authorization: Bearer <token>
Content-Type: application/json

{
  "approved": true,
  "remarks": "Approved"
}
```

### Attendance Tracking

**Mark Attendance**
```http
POST /api/v1/attendance/mark
Authorization: Bearer <token>
Content-Type: application/json

{
  "student_id": 42,
  "date": "2025-10-22T00:00:00Z",
  "present": true
}
```

**Get Student Attendance**
```http
GET /api/v1/attendance/student/:id?start_date=2025-01-01&end_date=2025-12-31
Authorization: Bearer <token>
```

**Response includes:**
- Attendance records
- Statistics (present_days, total_days, attendance_percentage)

**Other Endpoints:**
| Method | Endpoint | Description | Role |
|--------|----------|-------------|------|
| GET | `/api/v1/attendance/my` | Get my attendance | Student |
| GET | `/api/v1/attendance/daily?date=2025-10-22` | Get daily attendance | All |

### Notifications

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/notifications/my` | Get my notifications |
| PUT | `/api/v1/notifications/:id/read` | Mark as read |
| GET | `/api/v1/notifications/unread-count` | Get unread count |

### Analytics (Admin Only)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/analytics/dashboard` | Dashboard statistics |
| GET | `/api/v1/analytics/leave-breakdown` | Leave breakdown by type |
| GET | `/api/v1/analytics/department?dept=CS` | Department statistics |
| GET | `/api/v1/analytics/absentees` | Frequent absentees |

### Health Check

```http
GET /health
```

## Complete Example Workflow

**1. Register a Student**
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

**2. Login**
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"john@student.edu","password":"password123"}'
```

**3. Apply for Leave**
```bash
curl -X POST http://localhost:8080/api/v1/leaves/apply \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{
    "leave_type": "Medical",
    "reason": "Doctor appointment",
    "start_date": "2025-10-25T00:00:00Z",
    "end_date": "2025-10-25T00:00:00Z"
  }'
```

**4. Faculty Approves Leave**
```bash
curl -X PUT http://localhost:8080/api/v1/leaves/1/approve \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <faculty-token>" \
  -d '{"approved": true, "remarks": "Approved"}'
```

**5. Check Notifications**
```bash
curl -X GET http://localhost:8080/api/v1/notifications/my \
  -H "Authorization: Bearer <token>"
```

## Configuration

Create `.env` file:

```env
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=attendance_user
DB_PASSWORD=attendance_pass
DB_NAME=attendance_db
DB_SSL_MODE=disable

# Server Configuration
PORT=8080
GIN_MODE=debug

# Authentication
JWT_SECRET=your-secret-key-change-in-production
JWT_EXPIRY=24h

# Redis Configuration
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=
REDIS_DB=0

# Email Configuration (Optional)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=your-email@gmail.com
SMTP_PASSWORD=your-app-specific-password
SMTP_FROM=noreply@university.edu
```

## Project Structure

```
attendance-workflow/
├── cmd/
│   └── server/          # Application entrypoint
├── internal/
│   ├── analytics/       # Analytics & reporting
│   ├── api/            # Route definitions
│   ├── attendance/     # Attendance management
│   ├── auth/           # Authentication & authorization
│   ├── dto/            # Data transfer objects
│   ├── leaves/         # Leave management
│   ├── notifications/  # Notification system
│   └── users/          # User management
├── pkg/
│   ├── config/         # Configuration management
│   └── db/            # Database models & connection
├── scripts/           # Database scripts
├── docs/             # Swagger documentation
├── docker-compose.yml # Container orchestration
├── Dockerfile        # API container build
├── go.mod           # Dependencies
└── README.md        # Documentation
```

## Quick Commands

```bash
# Start all services
docker-compose up -d

# Reset database with demo data
go run scripts/reset_db.go

# View API logs
docker-compose logs -f api

# Access Swagger UI
http://localhost:8080/swagger/index.html

# Health check
curl http://localhost:8080/health
```

## Tech Stack

- Language: Go 1.21
- Framework: Gin
- ORM: GORM
- Database: PostgreSQL
- Authentication: JWT
- Containerization: Docker

## Commands

```bash
# Using Makefile
make build          # Build application
make run            # Run application
make docker-up      # Start Docker services
make docker-down    # Stop Docker services
make docker-logs    # View Docker logs
make clean          # Clean build files
```

## Troubleshooting

**Port already in use:**
```bash
PORT=8081 go run cmd/server/main.go
```

**Database connection error:**
```bash
docker-compose ps
docker-compose logs postgres
```

**Module errors:**
```bash
go mod tidy
go mod download
```


