# TheraClosure - Professional Practice Closure Platform

A comprehensive full-stack web application for helping therapists close their practices with dignity, care, and professional integrity. Built with modern microservices architecture using Go, React, and enterprise-grade containerized infrastructure.

## 🏗️ Architecture Overview

### Monorepo Structure
```
theraclosure-web-app/
├── frontend/                 # React + Vite + TypeScript + Material-UI
├── services/
│   ├── auth/                # Authentication service (Go/Gin) ✅ DOCKERIZED
│   ├── users/               # User management service (Go/Gin) ✅ DOCKERIZED
│   ├── geolocation/         # World geographic data service (Go/Gin) ✅ DOCKERIZED + DATA SEEDED
│   ├── payments/            # Stripe integration service (Go/Gin) ✅ PRODUCTION READY + DOCKERIZED
│   └── core/                # API Gateway service (Go/Gin + ActiveMQ) 📋 PLANNED
├── shared/                  # Shared TypeScript interfaces/DTOs
├── infra/                   # Infrastructure configurations
│   ├── docker/              # Docker Compose files
│   │   ├── docker-compose.infrastructure.yml
│   │   └── services/docker-compose.services.yml
│   ├── migrations/          # Database migrations
│   └── pgadmin/             # pgAdmin configuration
├── Makefile                 # Enterprise infrastructure management
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
- **Go 1.23** with Gin framework
- **Hexagonal Architecture** (Ports & Adapters)
- **PostgreSQL** for persistent data
- **Redis** for sessions and caching
- **ActiveMQ** for pub/sub messaging between services
- **JWT** for authentication and authorization
- **Comprehensive Geolocation Service** with complete world data (250+ countries, 3000+ subdivisions, major cities)
- **Stripe** for payment processing and subscriptions

#### Infrastructure
- **Docker Compose** for orchestrated services
- **PostgreSQL 15** with specialized databases per service
- **Redis 7** for caching and session management
- **ActiveMQ** for microservices communication
- **pgAdmin 4** for database management
- **Mailhog** for email testing
- **Enterprise Makefile** for infrastructure management

### 🎨 Design System & Branding

#### TheraClosure Brand Colors
- **Primary Teal**: `#2C5F5D` - Professional healthcare tone
- **Secondary Cream**: `#F5F1E8` - Warm, approachable backgrounds
- **Accent Gold**: `#D4A574` - Premium service highlights
- **Text Dark**: `#2C2C2C` - High contrast readability

#### Material-UI Theme
```typescript
// Custom theme with healthcare-focused design
const theme = createTheme({
  palette: {
    primary: { main: '#2C5F5D' },
    secondary: { main: '#F5F1E8' },
    text: { primary: '#2C2C2C' }
  },
  typography: {
    fontFamily: '"Inter", "Roboto", sans-serif',
    h1: { fontWeight: 600, color: '#2C2C2C' },
    // Optimized for accessibility and trust-building
  }
})
```

#### Design Principles
- **Professional Healthcare Aesthetic**: Clean, trustworthy, accessible
- **WCAG AA Compliance**: High contrast ratios for accessibility
- **Mobile-First Responsive**: Optimized for all device types
- **Trust-Building Elements**: Professional testimonials, credentials, guarantees

## �️ Database Architecture

### Multi-Database Design
Each microservice has its own dedicated database following Domain-Driven Design principles:

#### **`theraclosure_auth`** 🔐 Authentication Service
- `users` - User accounts, authentication credentials, roles (ADMIN, THERAPIST, STAFF)
- `sessions` - JWT session management, refresh tokens with SHA-256 hashing

#### **`theraclosure_users`** 👥 User Management Service  
- `user_profiles` - Extended profile data, licenses, specializations, practice info
- `enrollment_data` - Multi-step onboarding progress tracking

#### **`theraclosure_geolocation`** 🌍 Geolocation Service
- `countries` - Complete world countries (250+) with ISO codes, regions, currencies
- `states` - Global subdivisions (3000+) including states, provinces, regions, departments
- `cities` - Major world cities with coordinates, population data, postal codes
- **Pre-populated** via comprehensive automated data seeding system

#### **`theraclosure_payments`** 💳 Payment Processing Service ✅ PRODUCTION READY
- `customers` - Stripe customer management with user mapping
- `subscriptions` - Complete subscription lifecycle management, billing cycles, trials
- `payments` - Payment transactions, intent tracking, refunds, metadata
- **Stripe Integration**: Full webhook support, real-time event processing
- **Security**: JWT authentication, request validation, comprehensive error handling
- **Monitoring**: Prometheus metrics, health checks, performance tracking

#### **`theraclosure_core`** 🎯 Core Gateway Service
- `templates` - Email/document templates for client notifications
- `support_tickets` - Customer support and help desk system

#### **`theraclosure`** 📊 Main Database (Reserved)
*Currently empty - reserved for future cross-service data, global configuration, or system-wide audit logs*

### Database Management
- **pgAdmin**: http://localhost:5050 
  - **Email**: `admin@theraclosure.com`
  - **Password**: `admin123`  
  - **Pre-configured** with all 5 databases
  - **Auto-authentication** via pgpass file

## �🚀 Quick Start

### Prerequisites
- Docker and Docker Compose v2+
- Git
- Node.js 18+ (for frontend development)
- Go 1.23+ (for backend development)

### 1. Clone and Setup
```bash
git clone <repository-url>
cd theraclosure-web-app
```

### 2. Enterprise Infrastructure Management
```bash
# Complete infrastructure setup (recommended)
make setup

# Or step by step:
make infra      # Start core infrastructure (PostgreSQL, Redis, etc.)
make migrate    # Run database migrations  
make services   # Start application services (auth, users, geolocation, payments, core)
make users      # Start only users service
make auth       # Start only auth service
```

### 3. Development Workflow
```bash
# Infrastructure management
make help       # Show all available commands
make status     # Check service status  
make health     # Health check all services
make logs       # View logs from all services
make clean      # Clean shutdown

# Service-specific operations
make auth       # Start only auth service
make build-auth # Build auth service image
make logs-auth  # View auth service logs
```

### 4. Access the Application

#### **Frontend & APIs**
- **Frontend App**: http://localhost:3000 ✅ 
- **Auth Service**: http://localhost:3001 ✅ 
- **Users Service**: http://localhost:3002
- **Payments Service**: http://localhost:3003  
- **Core Gateway**: http://localhost:8080

#### **Infrastructure Services**
- **pgAdmin (Database Management)**: http://localhost:5050 ✅
- **MailHog (Email Testing)**: http://localhost:8025 ✅
- **PostgreSQL**: localhost:5432 ✅
- **Redis**: localhost:6379 ✅

### 4. Default Admin Login
```
Email: admin@theraclosure.com
Password: admin123
```

## 📋 Available Pages & Features ✨ RECENTLY UPDATED

### Public Pages (Professional TheraClosure Branding)
1. **Home** - Professional hero section with real TheraClosure content, service cards, and benefit highlights
2. **About** - Complete origin story with mission, values, and professional credentials
3. **Coverage & Costs** - Three-tier pricing structure ($297 Essential, $497 Professional, $897 Enterprise)
4. **How We Work** - Comprehensive service portfolio with detailed process explanation
5. **Retirement** - Full retirement planning services with timeline and consultation options
6. **Testimonials** - Professional client testimonials with authentic healthcare industry voices
7. **FAQ** - Comprehensive questions covering service overview, pricing, legal compliance
8. **Contact** - Professional contact form with practice type selection and business information
9. **Login/Register** - Secure authentication with professional healthcare styling

### ✨ Recent Frontend Improvements (2025)
- **Complete UI/UX Overhaul**: Professional healthcare design with TheraClosure brand colors
- **Real Content Integration**: Authentic service descriptions and pricing from theraclosure.com
- **Enhanced Accessibility**: Improved color contrast and WCAG AA compliance
- **Responsive Design**: Mobile-first approach with tablet and desktop optimization
- **Material-UI Theme**: Custom teal (#2C5F5D), cream (#F5F1E8), and gold (#D4A574) branding
- **Navigation Improvements**: React Router Link integration for seamless page transitions
- **Trust Building Elements**: Professional testimonials, credentials, and service guarantees

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

## 🌍 Geolocation Service - Comprehensive World Data

### World Geographic Database
The Geolocation Service provides complete global geographic data with a sophisticated automated seeding system:

#### **Data Coverage**
- **250+ Countries**: All UN recognized countries with ISO 3166-1 standards
  - Official names, 2-letter (US) and 3-letter (USA) codes
  - Regions, subregions, currencies, capitals
  - Population and area data
  
- **3000+ Subdivisions**: States, provinces, regions using ISO 3166-2
  - US states (50 states + DC)
  - Canadian provinces and territories (13)
  - German states, French regions, etc.
  - Administrative divisions for all countries
  
- **Major World Cities**: Comprehensive metropolitan areas
  - Population data and coordinates
  - Postal/ZIP codes
  - Global distribution across all continents

#### **Automated Data Seeding System**
```bash
# Navigate to geolocation service
cd services/geolocation/data/seeds

# Run comprehensive world data seeding
./run_seeder.sh --full

# Or run specific data types
./run_seeder.sh --countries-only          # 250+ countries
./run_seeder.sh --subdivisions-only       # 3000+ subdivisions  
./run_seeder.sh --cities-only             # Major cities
```

#### **Data Sources (Fully Automated)**
1. **pycountry Library**: Official ISO 3166 standards for countries and subdivisions
2. **REST Countries API**: Enhanced metadata (currencies, regions, capitals, population)
3. **GeoNames API**: Major world cities with population thresholds
4. **geonamescache Library**: Offline city database for reliability
5. **Hardcoded Fallbacks**: Only used if all APIs fail

#### **API Endpoints**
```
GET  /api/v1/countries                    # List all countries
GET  /api/v1/countries/{id}              # Get country details
GET  /api/v1/countries/{id}/states       # Get states/provinces
GET  /api/v1/states/{id}/cities          # Get cities in state
POST /api/v1/bulk/countries              # Bulk insert countries
POST /api/v1/bulk/states                 # Bulk insert subdivisions
POST /api/v1/bulk/cities                 # Bulk insert cities
```

#### **Integration Benefits**
- **Practice Registration**: Accurate location selection for therapists
- **License Validation**: State-specific professional licensing requirements  
- **Service Areas**: Geographic coverage for therapy services
- **Billing Addresses**: International address validation
- **Analytics**: Geographic distribution of user base

## 🐳 Docker Services

### Service Dependencies
```
Frontend → Core Gateway → Auth/Users/Geolocation/Payments Services → PostgreSQL/Redis/ActiveMQ
```

### Current Service Status
- ✅ **Auth Service**: JWT authentication, session management (Port 3001, Dockerized)
- ✅ **Users Service**: User profiles, enrollment workflow (Port 3002, Dockerized) 
- ✅ **Geolocation Service**: Complete world geographic data with comprehensive seeding system (Port 3003, Dockerized)
  - 250+ countries using ISO 3166 standards
  - 3000+ subdivisions (states, provinces, regions) 
  - Major world cities with coordinates and population data
  - Fully automated data seeding from multiple APIs (GeoNames, REST Countries, pycountry)
- ✅ **Payments Service**: Stripe integration, subscriptions (Port 3004, Production Ready + Dockerized)
  - Full Stripe Payment Intents API support
  - Comprehensive webhook event handling
  - Customer, subscription, and payment lifecycle management
  - Production-grade monitoring, security, and error handling
  - Real-time metrics and health checks
- 📋 **Core Gateway**: API Gateway with ActiveMQ messaging (Planned)

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
- **Auth Service**: http://localhost:3001/swagger/index.html
- **Users Service**: http://localhost:3002/swagger/index.html  
- **Geolocation Service**: http://localhost:3003/swagger/index.html
- **Payments Service**: http://localhost:3004/swagger/index.html ✅ PRODUCTION READY
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

## � Recent Updates & Improvements (2025)

### Frontend Revamp & Professional Branding
- ✅ **Complete UI/UX Overhaul**: Redesigned all pages with professional healthcare aesthetic
- ✅ **Real TheraClosure Content**: Integrated authentic service descriptions and pricing
- ✅ **Brand Color Implementation**: Applied teal (#2C5F5D) primary theme throughout
- ✅ **Enhanced Typography**: Improved readability with proper contrast ratios
- ✅ **Responsive Design**: Mobile-first approach with tablet/desktop optimization

### Bug Fixes & Technical Improvements
- ✅ **Color Contrast Issues**: Fixed green-on-green text visibility problems
- ✅ **Navigation Enhancement**: Implemented React Router Link components for seamless routing
- ✅ **Import Resolution**: Resolved TypeScript module import issues
- ✅ **Pricing Card UX**: Fixed "Most Popular" flag z-index and positioning
- ✅ **Theme Optimization**: Updated Material-UI theme for better accessibility

### Content & Feature Updates
- ✅ **Retirement Services**: Expanded retirement planning page with comprehensive content
- ✅ **Pricing Structure**: Professional three-tier pricing ($297/$497/$897 annually)
- ✅ **Contact Forms**: Enhanced contact page with practice type selection
- ✅ **Testimonials**: Added authentic healthcare industry client testimonials
- ✅ **FAQ Expansion**: Comprehensive questions covering all service aspects

### Development Documentation
- ✅ **v0.ly.ai Integration Guide**: Created comprehensive prompt for ConvexDB migration
- ✅ **Updated READMEs**: Documented all recent changes and improvements
- ✅ **Design System**: Established brand guidelines and color specifications

### Next Steps
- 🔄 Users Service Implementation
- 🔄 Payments Service with Stripe Integration  
- 🔄 Core Gateway Development
- 🔄 Complete Authentication Flow

## �📚 Additional Resources

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

## � Project Status

### Completed ✅
- **Frontend Foundation**: React app with Material-UI, routing, and authentication
- **Auth Service**: JWT authentication, session management, user registration/login
- **Users Service**: User profiles, 5-step enrollment workflow, CRUD operations
- **Geolocation Service**: Complete world geographic data with 250+ countries, 3000+ subdivisions, major cities
- **Payments Service**: Production-ready Stripe integration with comprehensive payment processing
- **Infrastructure**: PostgreSQL, Redis, Docker orchestration with health checks
- **Database Design**: Specialized databases per service with proper migrations
- **Containerization**: All core services dockerized and orchestrated

### In Progress 🚧
- **Frontend Integration**: React components for payments and geolocation services
- **Service Integration**: Inter-service communication patterns

### Planned 📋
- **Core Gateway**: API Gateway with ActiveMQ pub/sub messaging
- **Complete Frontend**: All 14 pages with full functionality
- **Production Deployment**: AWS infrastructure and CI/CD pipelines

### Architecture Progress
```
✅ Infrastructure Layer    - PostgreSQL, Redis, Docker, Makefile
✅ Auth Microservice      - JWT, sessions, user authentication  
✅ Users Microservice     - Profiles, enrollment, workflow management
✅ Geolocation Service   - Complete world geographic data with automated seeding
✅ Payments Service      - Production-ready Stripe integration with full lifecycle management
📋 Core Gateway         - API orchestration with ActiveMQ messaging
� Frontend Integration  - Complete React app with all services
```

## �📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 💬 Support

For support and questions:
- Create an issue in the repository
- Contact: support@theraclosure.com
- Documentation: [Wiki](wiki-url)

---

**TheraClosure Team** - Helping therapists transition with dignity and care.