# Easy Attend Service API Documentation

## Overview
REST API service สำหรับระบบ Easy Attend ที่ใช้ในการจัดการข้อมูลครูและนักเรียน

## Project Structure
```
easy-attend-serviceV2/
├── cmd/                    # Command line interface
├── configs/                # Database configuration
├── controller/             # HTTP handlers
├── middlewares/            # Authentication middleware
├── models/                 # Database models
├── requests/               # Request validation structs
├── response/               # Response helper functions
├── services/               # Business logic
└── utils/                  # Utility functions
    └── jwt/                # JWT token management
```

## Technologies Used
- **Go** - Programming language
- **Gin** - HTTP web framework
- **GORM** - ORM library
- **PostgreSQL** - Database
- **JWT** - Authentication
- **bcrypt** - Password hashing

## Environment Variables
Create a `.env` file in the root directory:
```env
DB_HOST=localhost
DB_PORT=5432
DB_DATABASE=easy-attend-serviceV2
DB_USER=postgres
DB_PASSWORD=1234

PORT=8080
GIN_MODE=debug

JWT_SECRET=your_jwt_secret_key_here
JWT_EXPIRE_HOURS=24

APP_NAME=Easy Attend Service
APP_VERSION=1.0.0
```

## Running the Application

### 1. Database Migration
```bash
./easy-attend-service.exe migrate
```

### 2. Start the Server
```bash
./easy-attend-service.exe serve
```

The server will start on `http://localhost:8080`

## API Endpoints

### Authentication Endpoints

#### POST /api/v1/auth/login
Login with email and password
```json
{
  "email": "teacher@example.com",
  "password": "password123"
}
```

#### POST /api/v1/auth/register
Register a new teacher
```json
{
  "email": "teacher@example.com",
  "password": "password123",
  "first_name": "John",
  "last_name": "Doe",
  "phone": "0812345678"
}
```

#### GET /api/v1/auth/profile
Get authenticated user profile (requires authentication)

### Teacher Endpoints (Protected)
**All teacher endpoints require authentication header:**
```
Authorization: Bearer <JWT_TOKEN>
```

#### GET /api/v1/teachers
Get all teachers with pagination
- Query Parameters:
  - `page` (optional): Page number (default: 1)
  - `limit` (optional): Items per page (default: 10, max: 100)

#### POST /api/v1/teachers
Create a new teacher
```json
{
  "email": "teacher@example.com",
  "password": "password123",
  "first_name": "John",
  "last_name": "Doe",
  "phone": "0812345678"
}
```

#### GET /api/v1/teachers/:id
Get teacher by ID

#### PUT /api/v1/teachers/:id
Update teacher information
```json
{
  "email": "newemail@example.com",
  "password": "newpassword123",
  "first_name": "John",
  "last_name": "Smith",
  "phone": "0812345679"
}
```

#### DELETE /api/v1/teachers/:id
Delete teacher by ID

### Student Endpoints (Protected)
**All student endpoints require authentication header:**
```
Authorization: Bearer <JWT_TOKEN>
```

#### GET /api/v1/students
Get all students with pagination
- Query Parameters:
  - `page` (optional): Page number (default: 1)
  - `limit` (optional): Items per page (default: 10, max: 100)

#### POST /api/v1/students
Create a new student
```json
{
  "student_no": "STD001",
  "first_name": "Jane",
  "last_name": "Doe",
  "email": "student@example.com",
  "phone": "0812345678"
}
```

#### GET /api/v1/students/:id
Get student by ID

#### PUT /api/v1/students/:id
Update student information
```json
{
  "student_no": "STD001",
  "first_name": "Jane",
  "last_name": "Smith",
  "email": "newstudent@example.com",
  "phone": "0812345679"
}
```

#### DELETE /api/v1/students/:id
Delete student by ID

### Health Check

#### GET /health
Check API health status

## Response Format

### Success Response
```json
{
  "status": {
    "code": 200,
    "message": "Success message"
  },
  "data": {
    // Response data
  }
}
```

### Paginated Response
```json
{
  "status": {
    "code": 200,
    "message": "Success message"
  },
  "data": {
    "teachers": [...],
    "pagination": {
      "page": 1,
      "limit": 10,
      "total": 100,
      "total_pages": 10
    }
  }
}
```

### Error Response
```json
{
  "code": 400,
  "message": "Error message: Error details"
}
```

## Authentication
This API uses JWT (JSON Web Token) for authentication. After successful login, you'll receive a token that must be included in the Authorization header for protected endpoints.

Token format: `Bearer <JWT_TOKEN>`

### School Endpoints (Protected)

#### GET /api/v1/schools
Get all schools with pagination

#### POST /api/v1/schools
Create a new school
```json
{
  "name": "ABC School"
}
```

#### GET /api/v1/schools/:id
Get school by ID

#### PUT /api/v1/schools/:id
Update school information
```json
{
  "name": "New School Name"
}
```

#### DELETE /api/v1/schools/:id
Delete school by ID

## Database Models

### School
```go
type School struct {
    ID        uuid.UUID `json:"id"`
    Name      string    `json:"name"`
    CreatedAt int64     `json:"created_at"`
    UpdatedAt int64     `json:"updated_at"`
    DeletedAt *int64    `json:"deleted_at,omitempty"`
}
```

### Teacher
```go
type Teacher struct {
    ID        uuid.UUID `json:"id"`
    SchoolID  uuid.UUID `json:"school_id"`
    Email     string    `json:"email"`
    FirstName string    `json:"first_name"`
    LastName  string    `json:"last_name"`
    Phone     string    `json:"phone"`
    CreatedAt int64     `json:"created_at"`
    UpdatedAt int64     `json:"updated_at"`
    DeletedAt *int64    `json:"deleted_at,omitempty"`
}
```

### Student
```go
type Student struct {
    ID        uuid.UUID `json:"id"`
    StudentNo string    `json:"student_no"`
    FirstName string    `json:"first_name"`
    LastName  string    `json:"last_name"`
    CreatedAt int64     `json:"created_at"`
    UpdatedAt int64     `json:"updated_at"`
    DeletedAt *int64    `json:"deleted_at,omitempty"`
}
```

### Classroom
```go
type Classroom struct {
    ID        uuid.UUID `json:"id"`
    SchoolID  uuid.UUID `json:"school_id"`
    TeacherID uuid.UUID `json:"teacher_id"`
    StudentID uuid.UUID `json:"student_id"`
    Name      string    `json:"name"`
    CreatedAt int64     `json:"created_at"`
    UpdatedAt int64     `json:"updated_at"`
    DeletedAt *int64    `json:"deleted_at,omitempty"`
}
```

### ClassroomMember
```go
type ClassroomMember struct {
    TeacherID   *uuid.UUID `json:"teacher_id"`
    StudentID   *uuid.UUID `json:"student_id"`
    ClassroomID uuid.UUID  `json:"classroom_id"`
}
```

### Attendance
```go
type Attendance struct {
    ID          uuid.UUID        `json:"id"`
    ClassroomID uuid.UUID        `json:"classroom_id"`
    TeacherID   uuid.UUID        `json:"teacher_id"`
    StudentID   uuid.UUID        `json:"student_id"`
    SessionDate string           `json:"session_date"`
    Status      AttendanceStatus `json:"status"` // "มาเรียน", "ขาด", "มาสาย", "ลา"
    CheckedAt   int64            `json:"checked_at"`
    Remark      string           `json:"remark"`
}
```

### Log
```go
type Log struct {
    ID        uuid.UUID  `json:"id"`
    TeacherID uuid.UUID  `json:"teacher_id"`
    Action    LogAction  `json:"action"`
    Detail    string     `json:"detail"`
    CreatedAt int64      `json:"created_at"`
    SchoolID  *uuid.UUID `json:"school_id,omitempty"`
}
```

## Development

### Building the Application
```bash
go build -o easy-attend-service.exe
```

### Running Tests
```bash
go test ./...
```

### Code Structure
- **Controllers**: Handle HTTP requests and responses
- **Services**: Contain business logic
- **Models**: Database entity definitions
- **Middlewares**: Authentication and other middleware functions
- **Utils**: Helper functions (JWT, password hashing, etc.)

## Security Features
- JWT authentication
- Password hashing with bcrypt
- Input validation
- CORS support (can be configured)
- SQL injection protection (via GORM)

## Error Handling
The API includes comprehensive error handling with appropriate HTTP status codes and descriptive error messages.

## License
This project is licensed under the MIT License.