# TheraClosure Docker Infrastructure

A comprehensive Docker-based microservices infrastructure for the TheraClosure platform, providing production-ready containerization with health checks, networking, and service orchestration.

## 🏗️ Architecture Overview

### Infrastructure Services (`docker-compose.infrastructure.yml`)
Core infrastructure services that support the application microservices:

- **PostgreSQL 15**: Multi-database setup with PostGIS extension
- **Redis 7**: Caching and session management
- **MailHog**: Email testing and development
- **pgAdmin 4**: Database administration interface

### Application Services (`services/docker-compose.services.yml`)
Business logic microservices:

- **Auth Service** (Port 3001): JWT authentication and session management
- **Users Service** (Port 3002): User profiles and enrollment workflow
- **Geolocation Service** (Port 3003): Global geographic data with automated seeding
- **Payments Service** (Port 3004): Stripe integration and payment processing

## 🚀 Quick Start

### 1. Complete Infrastructure Setup
```bash
# Navigate to project root
cd thera-closure/web-app

# Start infrastructure services
docker compose -f infra/docker/docker-compose.infrastructure.yml up -d

# Start all application services
docker compose -f infra/docker/services/docker-compose.services.yml up -d

# Verify all services are running
docker ps
```

### 2. Individual Service Management
```bash
# Start only infrastructure
make infra

# Start specific services
docker compose -f infra/docker/services/docker-compose.services.yml up auth-service -d
docker compose -f infra/docker/services/docker-compose.services.yml up payments-service -d

# View service logs
docker logs -f theraclosure-payments-service
docker logs -f theraclosure-auth-service
```

## 📊 Service Configuration

### Network Architecture
All services communicate through the `theraclosure-network` Docker network:
```yaml
networks:
  theraclosure-network:
    external: true
    name: docker_theraclosure-network
```

### Port Mapping
| Service | Internal Port | External Port | Protocol |
|---------|---------------|---------------|----------|
| PostgreSQL | 5432 | 5432 | TCP |
| Redis | 6379 | 6379 | TCP |
| MailHog SMTP | 1025 | 1025 | TCP |
| MailHog Web | 8025 | 8025 | HTTP |
| pgAdmin | 80 | 5050 | HTTP |
| Auth Service | 3001 | 3001 | HTTP |
| Users Service | 3002 | 3002 | HTTP |
| Geolocation Service | 3003 | 3003 | HTTP |
| Payments Service | 3004 | 3004 | HTTP |

### Health Checks
All services include comprehensive health monitoring:

#### Infrastructure Services
```yaml
# PostgreSQL
healthcheck:
  test: ["CMD-SHELL", "pg_isready -U theraclosure -d theraclosure"]
  interval: 30s
  timeout: 10s
  retries: 5

# Redis
healthcheck:
  test: ["CMD", "redis-cli", "ping"]
  interval: 30s
  timeout: 3s
  retries: 5
```

#### Application Services
```yaml
# Example: Payments Service
healthcheck:
  test: ["CMD", "wget", "--no-verbose", "--tries=1", "--spider", "http://localhost:3004/health"]
  interval: 30s
  timeout: 10s
  retries: 5
  start_period: 40s
```

## 🗄️ Database Configuration

### Multi-Database Architecture
Each microservice has its own dedicated database following Domain-Driven Design:

```sql
-- Created by init-db.sql
CREATE DATABASE theraclosure_auth;        -- Authentication service
CREATE DATABASE theraclosure_users;       -- User management service  
CREATE DATABASE theraclosure_payments;    -- Payment processing service
CREATE DATABASE theraclosure_geolocation; -- Geographic data service
CREATE DATABASE theraclosure_core;        -- Core gateway service
```

### Database Migrations
Automated schema setup with comprehensive migrations:

1. **`init-db.sql`**: Creates databases and basic tables
2. **`003_payments_schema.sql`**: Advanced payments schema with constraints
3. **Auto-migration**: Services automatically create tables on startup

### pgAdmin Access
Pre-configured database administration:
- **URL**: http://localhost:5050
- **Email**: admin@theraclosure.com  
- **Password**: admin123
- **Auto-configured**: All databases pre-connected

## 🔧 Environment Configuration

### Infrastructure Services
```yaml
# PostgreSQL
POSTGRES_DB: theraclosure
POSTGRES_USER: theraclosure
POSTGRES_PASSWORD: password123
POSTGRES_MULTIPLE_DATABASES: theraclosure_auth,theraclosure_users,theraclosure_payments,theraclosure_core,theraclosure_geolocation

# Redis  
REDIS_PASSWORD: (none - development setup)

# pgAdmin
PGADMIN_DEFAULT_EMAIL: admin@theraclosure.com
PGADMIN_DEFAULT_PASSWORD: admin123
```

### Application Services Environment
All microservices use consistent environment patterns:

```yaml
# Database Configuration (consistent across all services)
DB_HOST: postgres
DB_PORT: 5432
DB_USER: theraclosure
DB_PASSWORD: password123
DB_SSL_MODE: disable

# Service-specific database names
DB_NAME: theraclosure_{service}  # e.g., theraclosure_payments

# Server Configuration
SERVER_HOST: 0.0.0.0
SERVER_PORT: 300{X}  # X = service number (1-4)
SERVER_MODE: release

# Application Configuration
APP_NAME: theraclosure-{service}-service
APP_VERSION: 1.0.0
APP_LOG_LEVEL: info
```

## 🔍 Monitoring and Observability

### Health Endpoints
Every service exposes standardized health endpoints:
```bash
# Check individual service health
curl http://localhost:3001/health  # Auth
curl http://localhost:3002/health  # Users  
curl http://localhost:3003/health  # Geolocation
curl http://localhost:3004/health  # Payments

# Check infrastructure health
curl http://localhost:5050         # pgAdmin
curl http://localhost:8025         # MailHog
```

### Prometheus Metrics
All application services expose metrics at `/metrics`:
```bash
# View service metrics
curl http://localhost:3004/metrics  # Payments metrics
curl http://localhost:3001/metrics  # Auth metrics
```

### Swagger Documentation
Interactive API documentation for all services:
- **Auth**: http://localhost:3001/swagger/index.html
- **Users**: http://localhost:3002/swagger/index.html
- **Geolocation**: http://localhost:3003/swagger/index.html
- **Payments**: http://localhost:3004/swagger/index.html

## 🛠️ Development Workflow

### Service Development
```bash
# Develop with hot reload
cd services/payments
docker compose -f ../../infra/docker/docker-compose.infrastructure.yml up -d
go run cmd/main.go

# Build and test in container
docker build -t theraclosure-payments-service .
docker run --rm --network docker_theraclosure-network \
  -p 3004:3004 \
  -e DB_HOST=postgres \
  -e STRIPE_SECRET_KEY=sk_test_... \
  theraclosure-payments-service
```

### Database Management
```bash
# Connect to any database
docker exec -it theraclosure-postgres psql -U theraclosure -d theraclosure_payments

# Run migrations
docker exec -i theraclosure-postgres psql -U theraclosure -d theraclosure_payments < infra/migrations/003_payments_schema.sql

# Backup database
docker exec theraclosure-postgres pg_dump -U theraclosure theraclosure_payments > backup.sql

# Restore database  
docker exec -i theraclosure-postgres psql -U theraclosure -d theraclosure_payments < backup.sql
```

### Log Management
```bash
# View logs for all services
docker compose -f infra/docker/services/docker-compose.services.yml logs -f

# View logs for specific service
docker logs -f theraclosure-payments-service

# View infrastructure logs
docker compose -f infra/docker/docker-compose.infrastructure.yml logs -f postgres
```

## 🚀 Production Considerations

### Resource Requirements
| Service | CPU | Memory | Storage | Notes |
|---------|-----|--------|---------|--------|
| PostgreSQL | 2 vCPU | 2GB RAM | 50GB SSD | Can scale with read replicas |
| Redis | 1 vCPU | 512MB RAM | 10GB SSD | Primarily for caching |
| Auth Service | 1 vCPU | 256MB RAM | Minimal | Stateless, horizontally scalable |
| Users Service | 1 vCPU | 256MB RAM | Minimal | Stateless, horizontally scalable |
| Geolocation Service | 1 vCPU | 512MB RAM | Minimal | Large dataset in memory |
| Payments Service | 1 vCPU | 256MB RAM | Minimal | Stateless, horizontally scalable |

### Security Configuration
```bash
# Production environment variables
POSTGRES_PASSWORD: ${STRONG_PASSWORD}
REDIS_PASSWORD: ${REDIS_PASSWORD}
STRIPE_SECRET_KEY: ${PRODUCTION_STRIPE_KEY}
JWT_SECRET: ${RANDOM_JWT_SECRET}

# Network isolation
- Use Docker secrets for sensitive data
- Implement proper firewall rules
- Use TLS for all external communication
- Regular security updates for base images
```

### Scaling Guidelines
```yaml
# Horizontal scaling example
auth-service:
  deploy:
    replicas: 3
    restart_policy:
      condition: on-failure
      delay: 5s
      max_attempts: 3
    resources:
      limits:
        cpus: '1.0'
        memory: 512M
```

### Backup Strategy
```bash
# Automated database backups
docker exec theraclosure-postgres pg_dumpall -U theraclosure > daily_backup.sql

# Volume backups
docker run --rm -v docker_postgres_data:/data -v $(pwd):/backup alpine tar czf /backup/postgres_backup.tar.gz /data

# Redis persistence
docker exec theraclosure-redis redis-cli BGSAVE
```

## 🧪 Testing Infrastructure

### Integration Testing
```bash
# Start test environment
docker compose -f infra/docker/docker-compose.infrastructure.yml up -d
docker compose -f infra/docker/services/docker-compose.services.yml up -d

# Run integration tests
cd services/payments
go test -tags=integration ./tests/...

# Clean up test environment
docker compose -f infra/docker/services/docker-compose.services.yml down
docker compose -f infra/docker/docker-compose.infrastructure.yml down -v
```

### Load Testing
```bash
# Test service endpoints
curl -X POST http://localhost:3004/api/v1/customers \
  -H "Content-Type: application/json" \
  -d '{"user_id":"test","email":"test@example.com","name":"Test"}'

# Performance testing with ab
ab -n 1000 -c 10 http://localhost:3004/health
```

## 🔧 Troubleshooting

### Common Issues

#### Services Not Starting
```bash
# Check container status
docker ps -a

# Check logs for errors  
docker logs theraclosure-payments-service

# Verify network connectivity
docker network ls
docker network inspect docker_theraclosure-network
```

#### Database Connection Issues
```bash
# Test database connectivity
docker exec -it theraclosure-postgres psql -U theraclosure -c "SELECT version();"

# Check database exists
docker exec -it theraclosure-postgres psql -U theraclosure -c "\l"

# Verify user permissions
docker exec -it theraclosure-postgres psql -U theraclosure -c "\du"
```

#### Port Conflicts
```bash
# Check port usage
netstat -tulpn | grep :3004
lsof -i :3004

# Use different ports if needed
docker run -p 3005:3004 theraclosure-payments-service
```

#### Performance Issues
```bash
# Check resource usage
docker stats

# Monitor database performance
docker exec -it theraclosure-postgres psql -U theraclosure -c "SELECT * FROM pg_stat_activity;"

# Check Redis performance
docker exec -it theraclosure-redis redis-cli INFO stats
```

## 🚀 Deployment Automation

### Production Deployment
```bash
# Build all images
docker compose -f infra/docker/services/docker-compose.services.yml build

# Tag for registry
docker tag theraclosure-payments-service registry.example.com/theraclosure/payments:latest

# Deploy to production
docker stack deploy -c docker-compose.prod.yml theraclosure
```

### CI/CD Integration
```yaml
# Example GitHub Actions
- name: Build and Deploy
  run: |
    docker compose -f infra/docker/docker-compose.infrastructure.yml up -d
    docker compose -f infra/docker/services/docker-compose.services.yml up --build -d
    docker compose -f infra/docker/services/docker-compose.services.yml exec -T payments-service go test ./...
```

---

**TheraClosure Docker Infrastructure** - Production-ready microservices orchestration for helping therapists transition with dignity and care.