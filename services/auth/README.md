# TheraClosure Auth Service

A production-ready authentication microservice built with Go, Gin, and PostgreSQL following hexagonal architecture principles.

## Features

✅ **JWT Authentication** - Access and refresh token management  
✅ **Session Management** - PostgreSQL-based session tracking with token validation  
✅ **Password Security** - Bcrypt hashing with secure storage  
✅ **Role-Based Access** - ADMIN, THERAPIST, STAFF role management  
✅ **PostgreSQL Integration** - GORM ORM with automatic migrations  
✅ **Hexagonal Architecture** - Clean separation of concerns  
✅ **RESTful API** - Complete authentication endpoints  
✅ **CORS Support** - Frontend integration ready  
✅ **Environment Configuration** - Viper-based config management  
✅ **Production Ready** - Comprehensive error handling and logging

## 🚀 Recent Updates (2025)
- ✅ **Redis Session Repository**: Fixed session storage and retrieval bugs
- ✅ **Frontend Integration**: Enhanced JWT token handling for React frontend
- ✅ **CORS Configuration**: Optimized for TheraClosure frontend at localhost:3000
- ✅ **Docker Integration**: Fully containerized service with health checks
- ✅ **Production Testing**: Comprehensive authentication flow validation

## Architecture

```
├── cmd/
│   └── main.go                 # Application entry point
├── internal/
│   ├── core/                   # Business logic (domain layer)
│   │   ├── domain/             # Domain entities
│   │   │   └── user.go         # User, Session, Token models
│   │   ├── ports/              # Interface definitions
│   │   │   └── ports.go        # Repository and service interfaces
│   │   └── services/           # Business logic implementation
│   │       ├── auth_service.go # Authentication business logic
│   │       ├── jwt_service.go  # JWT token management
│   │       └── user_service.go # User management
│   └── adapters/               # External interfaces (adapters layer)
│       ├── config/             # Configuration management
│       │   └── config.go       # Viper configuration
│       ├── http/               # HTTP handlers
│       │   └── server.go       # Gin HTTP server & endpoints
│       └── persistence/        # Data persistence
│           ├── database.go     # PostgreSQL connection
│           ├── user_repository.go      # User data access
│           └── session_repository.go   # Session management
└── tests/
    ├── complete-test.sh        # Full authentication flow test
    ├── debug-refresh.sh        # Refresh token testing
    └── test-auth.sh           # Basic auth endpoint test
```

## API Endpoints

### Health Check
- `GET /api/v1/health` - Service health status

### Authentication
- `POST /api/v1/auth/register` - User registration
- `POST /api/v1/auth/login` - User login
- `POST /api/v1/auth/refresh` - Refresh access token
- `POST /api/v1/auth/logout` - User logout
- `GET /api/v1/auth/me` - Get current user profile

## Database Schema

### Users Table
```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL DEFAULT 'THERAPIST',
    subscription_status VARCHAR(50) DEFAULT 'trialing',
    stripe_customer_id VARCHAR(255),
    cognito_id VARCHAR(255),
    is_active BOOLEAN DEFAULT true,
    email_verified BOOLEAN DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
```

### Sessions Table
```sql
CREATE TABLE sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    refresh_token_hash VARCHAR(255) NOT NULL UNIQUE,
    access_token_jti VARCHAR(255) NOT NULL UNIQUE,
    user_agent TEXT,
    ip_address INET,
    is_active BOOLEAN NOT NULL DEFAULT true,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
```

## Session Management

The service implements a robust session management system:

### Token Security
- **Refresh tokens** are hashed using SHA-256 before storage
- **Access tokens** include unique JTI (JWT ID) for tracking
- **Session validation** checks token hash, expiration, and active status
- **Automatic cleanup** of expired and inactive sessions

### Session Lifecycle
1. **Registration/Login**: Creates new session with hashed refresh token
2. **Token Refresh**: Validates existing session and generates new tokens
3. **Logout**: Marks session as inactive (soft delete)
4. **Expiration**: Automatic cleanup of expired sessions via database function

### JWT Claims Structure
```json
{
  "user_id": "uuid",
  "email": "user@example.com", 
  "role": "THERAPIST|ADMIN|STAFF",
  "session_id": "uuid",
  "type": "access|refresh",
  "exp": "timestamp",
  "iat": "timestamp",
  "jti": "jwt-id"
}
```

## Configuration

### Environment Variables
```bash
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=theraclosure
DB_PASSWORD=password123
DB_NAME=theraclosure_auth
DB_SSL_MODE=disable

# JWT Configuration  
JWT_SECRET=your-super-secret-key
JWT_ACCESS_TOKEN_DURATION=1h
JWT_REFRESH_TOKEN_DURATION=168h

# Server Configuration
SERVER_HOST=0.0.0.0
SERVER_PORT=3001
SERVER_MODE=debug

# Application Configuration
APP_NAME=theraclosure-auth-service
APP_VERSION=1.0.0
APP_LOG_LEVEL=info
```

## Usage Examples

### Registration
```bash
curl -X POST http://localhost:3001/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "therapist@example.com",
    "password": "securepassword",
    "firstName": "John",
    "lastName": "Doe"
  }'
```

### Login
```bash
curl -X POST http://localhost:3001/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "therapist@example.com", 
    "password": "securepassword"
  }'
```

### Token Refresh
```bash
curl -X POST http://localhost:3001/api/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{
    "refreshToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }'
```

### Get Current User
```bash
curl -X GET http://localhost:3001/api/v1/auth/me \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

## Development

### Prerequisites
- Go 1.23+
- PostgreSQL 15+
- Docker & Docker Compose (for development)

### Setup
1. **Clone and navigate**:
   ```bash
   cd services/auth
   ```

2. **Install dependencies**:
   ```bash
   go mod tidy
   ```

3. **Start infrastructure**:
   ```bash
   docker compose up -d postgres
   ```

4. **Run migrations**:
   ```bash
   docker exec -i theraclosure-postgres psql -U theraclosure -d theraclosure_auth < ../../infra/migrations/001_users.sql
   docker exec -i theraclosure-postgres psql -U theraclosure -d theraclosure_auth < ../../infra/migrations/002_sessions.sql
   ```

5. **Start service**:
   ```bash
   go run cmd/main.go
   ```

### Testing
Run comprehensive authentication flow tests:
```bash
./tests/complete-test.sh     # Full authentication flow
./tests/debug-refresh.sh     # Refresh token testing  
./tests/test-auth.sh        # Basic endpoint testing
```

### Building
```bash
go build -o auth-service cmd/main.go
```

## Dependencies

### Core Dependencies
- **Gin** - HTTP web framework
- **GORM** - ORM for PostgreSQL 
- **Viper** - Configuration management
- **JWT-Go** - JWT token handling
- **Bcrypt** - Password hashing
- **UUID** - Unique identifier generation

### Database
- **PostgreSQL Driver** - `gorm.io/driver/postgres`
- **GORM** - `gorm.io/gorm`

### Security  
- **JWT** - `github.com/golang-jwt/jwt/v5`
- **Crypto** - `golang.org/x/crypto/bcrypt`

## Security Features

### Password Security
- **Bcrypt hashing** with cost factor 10
- **Password validation** with minimum length requirements
- **Secure storage** - passwords never stored in plaintext

### Token Security
- **JWT signing** with HS256 algorithm
- **Unique JTI** per access token for tracking
- **Refresh token hashing** before database storage
- **Automatic expiration** handling

### Session Security
- **Database-backed sessions** for immediate revocation
- **IP address tracking** (optional)
- **User agent tracking** (optional) 
- **Inactive session cleanup** via database functions

## Production Considerations

### Performance
- **Connection pooling** configured (10 idle, 100 max open)
- **Database indexes** on frequently queried fields
- **JWT validation** optimized with caching opportunities

### Monitoring
- **Structured logging** throughout the application
- **Health check endpoint** for load balancer integration  
- **Database connection monitoring** via GORM callbacks

### Scalability
- **Stateless design** - sessions stored in database, not memory
- **Horizontal scaling** ready - no shared state between instances
- **Database-backed** session management for multi-instance deployments

## Future Enhancements

### Planned Features
- **AWS Cognito integration** for enterprise SSO
- **Rate limiting** for authentication endpoints  
- **Email verification** workflow
- **Password reset** functionality
- **Multi-factor authentication** (MFA)
- **Session device tracking** and management
- **OAuth2/OIDC provider** implementation

### Configuration Improvements
- **Vault integration** for secret management
- **Environment-specific configs** (dev/staging/prod)
- **Dynamic configuration** reload without restart

## Contributing

1. Follow **hexagonal architecture** principles
2. Maintain **test coverage** for new features
3. Use **conventional commits** for clear history
4. Update **API documentation** for endpoint changes
5. Ensure **backward compatibility** for existing clients

## License

This project is part of the TheraClosure platform and follows the main project's licensing terms.