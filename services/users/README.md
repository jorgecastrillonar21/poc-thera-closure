# TheraClosure Users Service

A comprehensive user management microservice built with Go, Gin, and PostgreSQL following hexagonal architecture principles.

## Features

✅ **User Profile Management** - Complete CRUD operations for therapist profiles  
✅ **Enrollment Process** - Multi-step enrollment with progress tracking  
✅ **Profile Validation** - Business rule validation and completion checking  
✅ **Search & Filtering** - Profile search and status-based filtering  
✅ **Plan Management** - Support for Essential, Professional, and Enterprise plans  
✅ **PostgreSQL Integration** - GORM ORM with automatic migrations  
✅ **Hexagonal Architecture** - Clean separation of concerns  
✅ **RESTful API** - Complete user management endpoints  
✅ **CORS Support** - Frontend integration ready  
✅ **Environment Configuration** - Viper-based config management

## Architecture

```
├── cmd/
│   └── main.go                 # Application entry point
├── internal/
│   ├── core/                   # Business logic (domain layer)
│   │   ├── domain/             # Domain entities
│   │   │   └── user.go         # UserProfile, EnrollmentData models
│   │   ├── ports/              # Interface definitions
│   │   │   └── ports.go        # Repository and service interfaces
│   │   └── services/           # Business logic implementation
│   │       ├── user_service.go       # User profile business logic
│   │       └── enrollment_service.go # Enrollment business logic
│   └── adapters/               # Infrastructure layer
│       ├── config/             # Configuration management
│       │   └── config.go       # Viper configuration
│       ├── http/               # HTTP transport layer
│       │   ├── server.go       # Gin HTTP server
│       │   └── handlers.go     # HTTP request handlers
│       └── persistence/        # Data persistence layer
│           ├── database.go     # Database connection
│           ├── user_repository.go      # User profile repository
│           └── enrollment_repository.go # Enrollment repository
├── Dockerfile                  # Docker configuration
├── .env                        # Environment variables
├── go.mod                      # Go module definition
└── README.md                   # This file
```

## Domain Models

### UserProfile
Complete therapist profile information including:
- **Personal Information**: Name, email, phone, address
- **Professional Information**: License details, specializations, experience
- **Practice Information**: Practice name, type, location, contact details
- **Emergency Contacts**: Emergency contact information
- **Status Tracking**: Profile completion, status management

### EnrollmentData
Multi-step enrollment process tracking:
- **Step Progress**: 5-step enrollment process (Personal Info, Licensure, Practice Info, Admin Setup, Schedule Config)
- **Plan Selection**: Essential ($297), Professional ($497), Enterprise ($897)
- **Status Management**: In-progress, completed, paused enrollment states
- **Payment Tracking**: Payment status and completion

## API Endpoints

### User Profiles
```
POST   /api/v1/users/profiles              # Create profile
GET    /api/v1/users/profiles/:userId      # Get profile by user ID
PUT    /api/v1/users/profiles/:userId      # Update profile
DELETE /api/v1/users/profiles/:userId      # Delete profile
GET    /api/v1/users/profiles              # List profiles (paginated)
GET    /api/v1/users/profiles/search       # Search profiles
```

### Enrollment Management
```
POST   /api/v1/enrollments/start                    # Start enrollment
GET    /api/v1/enrollments/:userId                  # Get enrollment data
PUT    /api/v1/enrollments/:userId                  # Update enrollment
POST   /api/v1/enrollments/:userId/steps/:step/complete # Complete step
GET    /api/v1/enrollments/:userId/progress         # Get progress
POST   /api/v1/enrollments/:userId/complete         # Complete enrollment
PUT    /api/v1/enrollments/:userId/plan            # Update selected plan
```

### Health Check
```
GET    /health                            # Service health status
```

## Configuration

### Environment Variables
```bash
# Server
PORT=3002                    # Server port
HOST=0.0.0.0                # Server host

# Database
DB_HOST=localhost            # PostgreSQL host
DB_PORT=5432                # PostgreSQL port  
DB_USER=theraclosure        # Database user
DB_PASSWORD=theraclosure123 # Database password
DB_NAME=theraclosure_users  # Database name
DB_SSLMODE=disable          # SSL mode

# Environment
ENVIRONMENT=development      # Environment (development/production)

# CORS
CORS_ALLOWED_ORIGINS=http://localhost:3000,http://localhost:3001
```

## Database Schema

### user_profiles Table
- `id` (UUID) - Primary key
- `user_id` (UUID) - Foreign key to auth service
- `first_name`, `last_name` (String) - Personal information
- `email` (String) - Contact email (unique)
- `phone` (String) - Phone number
- `license_number`, `license_state` (String) - Professional licensing
- `practice_name`, `practice_type` (String) - Practice information
- `profile_complete` (Boolean) - Completion status
- `created_at`, `updated_at` (Timestamp) - Audit fields

### enrollment_data Table
- `id` (UUID) - Primary key
- `user_id` (UUID) - Foreign key to auth service
- `personal_info_complete` through `schedule_config_complete` (Boolean) - Step completion
- `enrollment_status` (String) - Overall enrollment status
- `current_step` (Integer) - Current enrollment step
- `selected_plan` (String) - Chosen plan (essential/professional/enterprise)
- `payment_status` (String) - Payment completion status
- `completion_date` (Timestamp) - Enrollment completion date

## Development

### Local Setup
```bash
# Clone the repository
git clone <repository-url>
cd services/users

# Install dependencies
go mod tidy

# Set up environment
cp .env.example .env
# Edit .env with your database credentials

# Run the service
go run cmd/main.go
```

### Docker Setup
```bash
# Build image
docker build -t theraclosure-users .

# Run container
docker run -p 3002:3002 --env-file .env theraclosure-users
```

### Testing the API

#### Create a Profile
```bash
curl -X POST http://localhost:3002/api/v1/users/profiles \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "123e4567-e89b-12d3-a456-426614174000",
    "first_name": "Dr. Sarah",
    "last_name": "Johnson",
    "email": "sarah.johnson@example.com",
    "license_number": "LIC123456",
    "license_state": "CA",
    "professional_title": "Licensed Clinical Psychologist"
  }'
```

#### Start Enrollment
```bash
curl -X POST http://localhost:3002/api/v1/enrollments/start \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "123e4567-e89b-12d3-a456-426614174000",
    "selected_plan": "professional"
  }'
```

#### Complete Enrollment Step
```bash
curl -X POST http://localhost:3002/api/v1/enrollments/123e4567-e89b-12d3-a456-426614174000/steps/1/complete
```

#### Get Enrollment Progress
```bash
curl http://localhost:3002/api/v1/enrollments/123e4567-e89b-12d3-a456-426614174000/progress
```

## Integration Points

### Auth Service Integration
- Relies on `user_id` from auth service for user identification
- Profile creation should be triggered after successful user registration

### Frontend Integration
- Provides complete user management API for React frontend
- Supports multi-step enrollment wizard UI
- Real-time progress tracking for enrollment completion

### Microservices Integration
- **Auth Service**: JWT authentication and session management ✅
- **Geolocation Service**: Countries/States/Cities data for address management 🚧 
- **Payments Service**: Payment status updates for enrollment completion 📋
- **Core Gateway**: API orchestration and template generation 📋

## 🐳 Docker Deployment ✅

### Current Status
The Users Service is **fully dockerized** and integrated with the main infrastructure.

### Quick Start
```bash
# Build and start the service
make build-users  # Build Docker image
make users        # Start containerized service

# Or start all services
make services     # Includes users service
```

### Docker Configuration
```yaml
users-service:
  build: ./services/users
  ports:
    - "3002:3002"
  environment:
    - DB_PASSWORD=password123
    - DB_HOST=postgres
    - DB_NAME=theraclosure_users
  depends_on:
    postgres:
      condition: service_healthy
  healthcheck:
    test: ["CMD", "curl", "-f", "http://localhost:3002/health"]
```

### Database Setup
Requires a dedicated PostgreSQL database `theraclosure_users` with proper migration support.

## Contributing

1. Follow hexagonal architecture principles
2. Maintain comprehensive test coverage
3. Use conventional commit messages
4. Update API documentation for new endpoints

---

**TheraClosure Users Service** - Comprehensive user profile and enrollment management for healthcare professionals.