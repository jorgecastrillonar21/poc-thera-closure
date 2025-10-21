# TheraClosure - Professional Practice Closure Platform

A comprehensive full-stack web application for helping therapists close their practices with dignity, care, and professional integrity. Built with modern microservices architecture using Go, React, and containerized deployment.

## 🏗️ Architecture Overview

### Monorepo Structure
```
theraclosure-web-app/
├── frontend/                 # React + Vite + TypeScript + Material-UI
├── services/
│   ├── auth/                # Authentication service (Go/Gin)
│   ├── users/               # User management service (Go/Gin)
│   ├── payments/            # Stripe integration service (Go/Gin)
│   └── core/                # API Gateway service (Go/Gin)
├── shared/                  # Shared TypeScript interfaces/DTOs
├── infra/                   # Infrastructure configurations
├── docker-compose.yml       # Multi-service orchestration
└── README.md
```

### Technology Stack

#### Frontend
- **React 18** with TypeScript
- **Vite** for build tooling
- **Material-UI (MUI v5)** for component library
- **React Router v6** for routing
- **React Query** for state management
- **Axios** for HTTP requests with JWT interceptors

#### Backend Services
- **Go 1.21** with Gin framework
- **Hexagonal Architecture** (Ports & Adapters)
- **PostgreSQL** for persistent data
- **Redis** for sessions and caching
- **JWT** for authentication
- **AWS Cognito** for SSO (OAuth2/OpenID Connect)
- **Stripe** for payment processing

#### Infrastructure
- **Docker Compose** for local development
- **PostgreSQL 15** with multiple databases
- **Redis 7** for caching and sessions
- **Mailhog** for email testing
- **NGINX** for reverse proxy (optional)

## 🚀 Quick Start

### Prerequisites
- Docker and Docker Compose
- Git
- Node.js 18+ (for frontend development)
- Go 1.21+ (for backend development)

### 1. Clone and Setup
```bash
git clone <repository-url>
cd theraclosure-web-app

# Copy environment files
cp frontend/.env.local.example frontend/.env.local
cp services/auth/.env.example services/auth/.env
# Edit the environment files with your configurations
```

### 2. Start All Services
```bash
# Build and start all services
docker-compose up --build

# Or start in detached mode
docker-compose up -d --build
```

### 3. Access the Application
- **Frontend**: http://localhost:3000
- **API Gateway**: http://localhost:8080
- **Mailhog UI**: http://localhost:8025
- **Auth Service**: http://localhost:3001
- **Users Service**: http://localhost:3002
- **Payments Service**: http://localhost:3003

### 4. Default Admin Login
```
Email: admin@theraclosure.com
Password: admin123
```

## 📋 Available Pages & Features

### Public Pages
1. **Home** - Hero section with CTAs
2. **About** - Origin story, team, partners
3. **Coverage & Costs** - Pricing information
4. **How We Work** - Process explanation
5. **Retirement** - Retirement-specific guidance
6. **Testimonials** - Client success stories
7. **FAQ** - Frequently asked questions
8. **Contact** - Contact form and information
9. **Login/Register** - Authentication with SSO

### Protected Dashboard
10. **Dashboard** - Overview and navigation
11. **Enrollment Checklist** - Multi-step wizard:
    - Personal Information
    - Licensure Details
    - Practice Information
    - Administrative Setup
    - Schedule Configuration
12. **Templates** - Closure letter templates
13. **Billing Management** - Stripe customer portal
14. **Support** - Help desk and tickets

## 🔐 Authentication & Authorization

### Features
- **Local Authentication** - Email/password with bcrypt
- **AWS Cognito SSO** - OAuth2/OpenID Connect integration
- **JWT Tokens** - Access tokens (1h) + Refresh tokens (7d)
- **Role-Based Access Control (RBAC)**:
  - `ADMIN` - Full system access
  - `THERAPIST` - Standard user access
  - `STAFF` - Limited support access
- **Session Management** - Redis-backed sessions
- **Secure Cookies** - HTTP-only JWT storage

### API Endpoints
```
POST /api/v1/auth/register      # User registration
POST /api/v1/auth/login         # Email/password login
POST /api/v1/auth/refresh       # Token refresh
POST /api/v1/auth/logout        # Logout
GET  /api/v1/auth/me           # Get current user
GET  /api/v1/auth/cognito/login # Initiate Cognito SSO
```

## 💳 Payment Integration

### Stripe Features
- **Subscription Management** - TheraClosure Plan
- **Checkout Sessions** - Hosted checkout flow
- **Customer Portal** - Self-service billing
- **Webhook Handling** - Automatic status updates
- **Test Mode** - Development-safe transactions

### Payment Flow
1. User initiates checkout from frontend
2. Core service creates Stripe checkout session
3. User completes payment on Stripe
4. Webhook updates subscription status
5. User gains access to premium features

## 🏥 Hexagonal Architecture

### Core Domain (Business Logic)
- **Domain Models** - User, Session, Subscription entities
- **Business Rules** - Authentication, authorization logic
- **Use Cases** - Application-specific operations

### Ports (Interfaces)
- **Primary Ports** - Service interfaces (AuthService, UserService)
- **Secondary Ports** - Repository interfaces (UserRepository, SessionRepository)

### Adapters (Infrastructure)
- **HTTP Adapters** - Gin web framework handlers
- **Persistence Adapters** - GORM database repositories  
- **External Adapters** - Stripe, AWS Cognito integrations

## 🐳 Docker Services

### Service Dependencies
```
Frontend → Core Service → Auth/Users/Payments Services → PostgreSQL/Redis
```

### Health Checks
All services include health check endpoints:
- `GET /api/v1/health` - Service health status
- Automatic service discovery and recovery
- Dependency waiting with health check conditions

### Volumes
- **postgres_data** - Database persistence
- **redis_data** - Cache persistence
- **Frontend hot reload** - Development volume mounting

## 🔧 Development

### Backend Development
```bash
# Enter a service directory
cd services/auth

# Install dependencies
go mod download

# Run locally (requires local DB/Redis)
go run cmd/main.go

# Build
go build -o bin/auth cmd/main.go

# Test
go test ./...
```

### Frontend Development
```bash
cd frontend

# Install dependencies
npm install

# Start development server
npm run dev

# Build for production
npm run build

# Preview production build
npm run preview
```

### Database Operations
```bash
# Connect to PostgreSQL
docker exec -it theraclosure-postgres psql -U theraclosure -d theraclosure_auth

# View logs
docker-compose logs -f postgres

# Reset database
docker-compose down -v
docker-compose up -d postgres
```

## 🔐 Environment Configuration

### Frontend (.env.local)
```bash
VITE_API_URL=http://localhost:8080/api
VITE_STRIPE_PUBLISHABLE_KEY=pk_test_...
VITE_AWS_COGNITO_USER_POOL_ID=us-east-1_...
VITE_AWS_COGNITO_CLIENT_ID=...
```

### Auth Service (.env)
```bash
# Database
DB_HOST=postgres
DB_USER=theraclosure
DB_PASSWORD=password123
DB_NAME=theraclosure_auth

# JWT
JWT_SECRET=your-super-secret-jwt-key
JWT_ACCESS_TOKEN_DURATION=1h
JWT_REFRESH_TOKEN_DURATION=168h

# AWS Cognito
AWS_REGION=us-east-1
AWS_COGNITO_USER_POOL_ID=us-east-1_...
AWS_COGNITO_CLIENT_ID=...
AWS_COGNITO_CLIENT_SECRET=...
```

## 📊 API Documentation

### Swagger/OpenAPI
- **Auth Service**: http://localhost:3001/api/docs
- **Users Service**: http://localhost:3002/api/docs  
- **Payments Service**: http://localhost:3003/api/docs
- **Core Gateway**: http://localhost:8080/api/docs

### Key API Patterns
- **RESTful Design** - Standard HTTP methods and status codes
- **JSON API** - Consistent request/response formats
- **Error Handling** - Structured error responses
- **Pagination** - Limit/offset pagination for lists
- **Filtering** - Query parameter filtering

## 🧪 Testing

### Stripe Test Data
```bash
# Test Credit Card Numbers
4242424242424242  # Visa (succeeds)
4000000000000002  # Visa (declined)

# Test Keys (use in .env files)
STRIPE_PUBLISHABLE_KEY=pk_test_51...
STRIPE_SECRET_KEY=sk_test_51...
STRIPE_WEBHOOK_SECRET=whsec_...
```

### Test Users
```bash
# Admin User
Email: admin@theraclosure.com
Password: admin123
Role: ADMIN

# Test Therapist (create via registration)
Role: THERAPIST (default)
```

## 🚀 AWS Deployment

### Prerequisites
- AWS Account with appropriate IAM permissions
- ECR repositories for each service
- RDS PostgreSQL instance
- ElastiCache Redis cluster
- ALB (Application Load Balancer)

### Deployment Steps
1. **Build and push Docker images to ECR**
2. **Deploy infrastructure using AWS CDK/CloudFormation**
3. **Configure environment variables in ECS/EKS**
4. **Set up RDS and ElastiCache**
5. **Configure ALB with SSL certificates**
6. **Set up Route53 for custom domains**

## 🔧 Adding New Features

### Adding a New Page
1. Create React component in `frontend/src/pages/`
2. Add route in `frontend/src/App.tsx`
3. Update navigation in `frontend/src/components/Layout/Navbar.tsx`

### Adding a New API Endpoint
1. Define domain models in `internal/core/domain/`
2. Create service interfaces in `internal/core/ports/`
3. Implement business logic in `internal/core/services/`
4. Add HTTP handlers in `internal/adapters/http/`
5. Update Swagger documentation

### Adding RBAC Roles
1. Update `UserRole` enum in shared types
2. Add role validation in auth middleware
3. Update frontend route protection
4. Test with different user roles

## 🐛 Troubleshooting

### Common Issues

#### Database Connection Issues
```bash
# Check if PostgreSQL is running
docker-compose ps postgres

# View PostgreSQL logs
docker-compose logs postgres

# Connect to database manually
docker exec -it theraclosure-postgres psql -U theraclosure
```

#### Service Not Starting
```bash
# Check service logs
docker-compose logs auth-service

# Rebuild specific service
docker-compose up --build auth-service

# Check health status
curl http://localhost:3001/api/v1/health
```

#### Frontend Build Issues
```bash
# Clear node modules and reinstall
cd frontend
rm -rf node_modules package-lock.json
npm install

# Check for TypeScript errors
npm run build
```

## 📚 Additional Resources

- [Go Gin Documentation](https://gin-gonic.com/docs/)
- [React Documentation](https://react.dev/)
- [Material-UI Documentation](https://mui.com/)
- [Stripe API Documentation](https://stripe.com/docs/api)
- [AWS Cognito Documentation](https://docs.aws.amazon.com/cognito/)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [Redis Documentation](https://redis.io/documentation)

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 💬 Support

For support and questions:
- Create an issue in the repository
- Contact: support@theraclosure.com
- Documentation: [Wiki](wiki-url)

---

**TheraClosure Team** - Helping therapists transition with dignity and care.