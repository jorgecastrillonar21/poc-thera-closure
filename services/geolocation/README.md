# TheraClosure Geolocation Service

A comprehensive geolocation management microservice with **complete world geographic data** built using Go, Gin, and PostgreSQL following hexagonal architecture principles.

## 🌍 World Data Coverage

### **Complete Global Database**
- **250+ Countries**: All UN recognized countries using ISO 3166-1 standards
- **3000+ Subdivisions**: States, provinces, regions, departments using ISO 3166-2
- **Major World Cities**: Metropolitan areas across all continents with population data
- **Comprehensive Seeding**: Fully automated data population from authoritative APIs

### **Automated Data Sources**
1. **pycountry Library** - Official ISO 3166 country and subdivision standards
2. **REST Countries API** - Enhanced country metadata (currencies, regions, capitals)
3. **GeoNames API** - Major world cities with population and coordinate data
4. **geonamescache Library** - Offline geographic database for reliability
5. **Smart Fallbacks** - Hardcoded data only used when APIs are unavailable

## 🌟 Features

✅ **Complete World Coverage** - Pre-populated with comprehensive global geographic data  
✅ **Country Management** - ISO standard countries with enhanced metadata  
✅ **Subdivision Management** - All world states, provinces, regions with hierarchical relationships  
✅ **City Management** - Major cities with coordinates, population, and postal codes  
✅ **Automated Data Seeding** - One-command world data population from multiple APIs  
✅ **Hierarchical API** - Get complete country/state/city hierarchies  
✅ **Advanced Search & Filtering** - Search across all geographic entities with pagination  
✅ **Bulk Operations** - High-performance batch operations for data import  
✅ **PostgreSQL Integration** - GORM ORM with automatic migrations and indexing  
✅ **Hexagonal Architecture** - Clean separation of concerns and testable design  
✅ **RESTful API** - Complete geolocation management with OpenAPI documentation  
✅ **CORS Support** - Frontend integration ready with proper headers  
✅ **Environment Configuration** - Flexible configuration management with Viper

## 🚀 Quick Start with World Data

### 1. Start the Service
```bash
cd services/geolocation
go run cmd/main.go
```

### 2. Populate World Data (One Command)
```bash
cd data/seeds
./run_seeder.sh --full
```

This will automatically populate:
- **250+ countries** from ISO 3166 standards + REST Countries API
- **3000+ subdivisions** from pycountry ISO 3166-2 database  
- **Major cities** from GeoNames API + comprehensive fallbacks

### 3. Verify Data
```bash
curl http://localhost:3003/api/v1/countries | jq '.countries | length'     # Should show 250+
curl http://localhost:3003/api/v1/countries | jq '.countries[0]'           # Sample country
```

### Custom Seeding Options
```bash
./run_seeder.sh --countries-only          # Countries only
./run_seeder.sh --subdivisions-only       # States/provinces only
./run_seeder.sh --cities-only             # Major cities only  
./run_seeder.sh --reset --full            # Clear and reseed all data
```

## Architecture

```
├── cmd/
│   └── main.go                 # Application entry point
├── internal/
│   ├── core/                   # Business logic (domain layer)
│   │   ├── domain/             # Domain entities
│   │   │   └── geolocation.go  # Country, State, City models
│   │   ├── ports/              # Interface definitions
│   │   │   └── ports.go        # Repository and service interfaces
│   │   └── services/           # Business logic implementation
│   │       └── geolocation_service.go # Geolocation business logic
│   └── adapters/               # Infrastructure layer
│       ├── config/             # Configuration management
│       │   └── config.go       # Viper-based configuration
│       ├── http/               # HTTP adapters
│       │   ├── server.go       # Gin server setup and routing
│       │   └── handlers.go     # HTTP request handlers
│       └── persistence/        # Data persistence
│           ├── database.go     # Database connection and setup
│           └── geolocation_repository.go # GORM repository implementation
├── tests/                      # Test files
├── Dockerfile                  # Docker configuration
├── go.mod                      # Go module definition
└── README.md                   # This file
```

## API Endpoints

### Countries
- `POST /api/v1/countries` - Create a new country
- `GET /api/v1/countries` - List all countries (paginated)
- `GET /api/v1/countries/search?q={query}` - Search countries
- `GET /api/v1/countries/{id}` - Get country by ID
- `PUT /api/v1/countries/{id}` - Update country
- `DELETE /api/v1/countries/{id}` - Delete country
- `GET /api/v1/countries/{id}/states` - Get states in country
- `GET /api/v1/countries/{id}/hierarchy` - Get country with states

### States
- `POST /api/v1/states` - Create a new state
- `GET /api/v1/states/search?q={query}&country_id={id}` - Search states
- `GET /api/v1/states/{id}` - Get state by ID
- `PUT /api/v1/states/{id}` - Update state
- `DELETE /api/v1/states/{id}` - Delete state
- `GET /api/v1/states/{id}/cities` - Get cities in state
- `GET /api/v1/states/{id}/hierarchy` - Get state with cities

### Cities
- `POST /api/v1/cities` - Create a new city
- `GET /api/v1/cities/search?q={query}&state_id={id}` - Search cities
- `GET /api/v1/cities/{id}` - Get city by ID
- `PUT /api/v1/cities/{id}` - Update city
- `DELETE /api/v1/cities/{id}` - Delete city

### Hierarchy
- `GET /api/v1/hierarchy` - Get complete country/state hierarchy

### Bulk Operations
- `POST /api/v1/bulk/countries` - Bulk create countries
- `POST /api/v1/bulk/states` - Bulk create states
- `POST /api/v1/bulk/cities` - Bulk create cities

### Health Check
- `GET /health` - Service health status

## Configuration

### Environment Variables
```bash
# Server Configuration
SERVER_HOST=0.0.0.0
SERVER_PORT=3003

# Database Configuration  
DB_HOST=localhost
DB_PORT=5432
DB_USER=theraclosure
DB_PASSWORD=password123
DB_NAME=theraclosure_geolocation
DB_SSL_MODE=disable

# Application Configuration
APP_NAME=theraclosure-geolocation-service
APP_VERSION=1.0.0
APP_LOG_LEVEL=info
```

## Development

### Prerequisites
- Go 1.23 or later
- PostgreSQL 15+
- Docker (optional)

### Running Locally
```bash
# Clone and navigate to the service directory
cd services/geolocation

# Install dependencies
go mod tidy

# Set environment variables
export DB_PASSWORD=password123

# Run the service
go run cmd/main.go
```

### Database Setup
The service automatically creates the required database tables and indexes on startup.

Required database: `theraclosure_geolocation`

### API Examples

#### Create a Country
```bash
curl -X POST http://localhost:3003/api/v1/countries \
  -H "Content-Type: application/json" \
  -d '{
    "name": "United States",
    "code": "USA",
    "code2": "US",
    "region": "North America",
    "currency": "USD"
  }'
```

#### Create a State
```bash
curl -X POST http://localhost:3003/api/v1/states \
  -H "Content-Type: application/json" \
  -d '{
    "country_id": "country-uuid",
    "name": "California",
    "code": "CA"
  }'
```

#### Create a City
```bash
curl -X POST http://localhost:3003/api/v1/cities \
  -H "Content-Type: application/json" \
  -d '{
    "state_id": "state-uuid",
    "name": "Los Angeles",
    "zip_code": "90210",
    "latitude": 34.0522,
    "longitude": -118.2437
  }'
```

#### Search Countries
```bash
curl "http://localhost:3003/api/v1/countries/search?q=united&limit=10&offset=0"
```

## Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific test
go test ./internal/core/services/
```

## 🐳 Docker Deployment

### Build Docker Image
```bash
docker build -t theraclosure-geolocation-service .
```

### Run with Docker
```bash
docker run -p 3003:3003 \
  -e DB_HOST=host.docker.internal \
  -e DB_PASSWORD=password123 \
  theraclosure-geolocation-service
```

### Microservices Integration
- **Users Service**: Provides geographic data for user profiles ✅
- **Auth Service**: No direct integration needed ✅
- **Payments Service**: Country/region data for billing 📋
- **Core Gateway**: API orchestration and caching 📋

## Production Considerations

### Performance
- Database indexes on frequently queried fields
- Pagination for large result sets
- Efficient search using database LIKE operations

### Data Sources
The service can be populated with data from:
- GeoNames.org free dataset
- Countries API
- Custom CSV imports via bulk operations

### Scalability
- Stateless service design for horizontal scaling
- Database connection pooling
- Configurable pagination limits

## Contributing

1. Follow hexagonal architecture principles
2. Maintain comprehensive test coverage
3. Use conventional commit messages
4. Update API documentation for new endpoints

---

**Part of the TheraClosure Microservices Architecture**