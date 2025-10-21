# Payments Service

A comprehensive payments microservice for TheraFlex Closure with Stripe integration, built using Go and following hexagonal architecture principles.

## Features

- **Customer Management**: Create, update, and manage customer profiles
- **Subscription Management**: Handle recurring subscriptions with Stripe
- **Payment Processing**: Process one-time and recurring payments
- **Payment Intents**: Support for modern Stripe Payment Intents API
- **Webhook Handling**: Real-time updates from Stripe webhooks
- **Database Integration**: PostgreSQL with GORM ORM
- **RESTful API**: Comprehensive REST endpoints with proper error handling
- **Health Checks**: Service health monitoring endpoints

## Architecture

The service follows hexagonal (ports and adapters) architecture:

- **Domain Layer** (`internal/core/domain`): Business entities and logic
- **Ports Layer** (`internal/core/ports`): Interfaces for external dependencies
- **Services Layer** (`internal/core/services`): Application business logic
- **Adapters Layer** (`internal/adapters`): External integrations
  - `config`: Configuration management
  - `http`: REST API server
  - `persistence`: Database repositories
  - `stripe`: Stripe API client

## Domain Models

### Customer
- Represents a payment customer linked to a user
- Contains Stripe customer ID for external integration
- Manages payment methods and profile information

### Subscription
- Handles recurring subscription billing
- Supports trial periods and cancellation
- Tracks subscription lifecycle and status

### Payment
- Individual payment transactions
- Supports both one-time and subscription payments
- Tracks payment status and metadata

## API Endpoints

### Health Check
```
GET /health
GET /api/v1/health
```

### Customers
```
POST   /api/v1/customers           # Create customer
GET    /api/v1/customers           # List customers
GET    /api/v1/customers/:id       # Get customer by ID
PUT    /api/v1/customers/:id       # Update customer
DELETE /api/v1/customers/:id       # Delete customer
GET    /api/v1/customers/user/:userID # Get customer by user ID
```

### Subscriptions
```
POST   /api/v1/subscriptions       # Create subscription
GET    /api/v1/subscriptions       # List subscriptions
GET    /api/v1/subscriptions/:id   # Get subscription by ID
PUT    /api/v1/subscriptions/:id   # Update subscription
DELETE /api/v1/subscriptions/:id   # Cancel subscription
```

### Payments
```
POST   /api/v1/payments            # Create payment record
GET    /api/v1/payments            # List payments
GET    /api/v1/payments/:id        # Get payment by ID
POST   /api/v1/payments/:id/refund # Refund payment
```

### Payment Intents
```
POST   /api/v1/payment-intents              # Create payment intent
POST   /api/v1/payment-intents/:id/confirm  # Confirm payment intent
```

### Webhooks
```
POST   /api/v1/webhooks/stripe     # Stripe webhook endpoint
```

## Configuration

The service uses environment variables for configuration:

### Database Configuration
- `DB_HOST`: Database host (default: localhost)
- `DB_PORT`: Database port (default: 5432)
- `DB_USER`: Database user (default: theraclosure)
- `DB_PASSWORD`: Database password (default: password123)
- `DB_NAME`: Database name (default: theraclosure_payments)
- `DB_SSL_MODE`: SSL mode (default: disable)

### Server Configuration
- `SERVER_HOST`: Server host (default: 0.0.0.0)
- `SERVER_PORT`: Server port (default: 3004)

### Stripe Configuration
- `STRIPE_PUBLIC_KEY`: Stripe publishable key
- `STRIPE_SECRET_KEY`: Stripe secret key
- `STRIPE_WEBHOOK_SECRET`: Stripe webhook endpoint secret

### Application Configuration
- `APP_NAME`: Application name
- `APP_VERSION`: Application version
- `APP_LOG_LEVEL`: Log level (debug, info, warn, error)

## Development

### Prerequisites
- Go 1.23 or later
- PostgreSQL 15 or later
- Stripe account for testing

### Setup
1. Clone the repository
2. Install dependencies:
   ```bash
   go mod download
   ```

3. Set up environment variables:
   ```bash
   export DB_HOST=localhost
   export DB_NAME=theraclosure_payments
   export STRIPE_SECRET_KEY=sk_test_...
   export STRIPE_WEBHOOK_SECRET=whsec_...
   ```

4. Run database migrations (see `../../infra/migrations/003_payments_schema.sql`)

5. Start the service:
   ```bash
   go run cmd/main.go
   ```

### Testing

Run health checks:
```bash
./tests/health-check.sh
```

Run unit tests:
```bash
go test ./...
```

### Building

Build the service:
```bash
go build -o payments-service cmd/main.go
```

Build Docker image:
```bash
docker build -t theraclosure/payments-service .
```

## Database Schema

The service uses three main tables:

### `customers`
- Customer profile information
- Links to user service via `user_id`
- Contains Stripe customer ID for payment processing

### `subscriptions`
- Subscription records with Stripe integration
- Tracks subscription lifecycle and billing periods
- Links to customers table

### `payments`
- Individual payment transaction records
- Can be linked to subscriptions for recurring payments
- Tracks payment status and metadata

## Stripe Integration

### Supported Operations
- Customer creation and management
- Subscription lifecycle management
- Payment Intent creation and confirmation
- Payment processing and refunds
- Webhook event handling

### Webhook Events
The service handles these Stripe webhook events:
- `payment_intent.succeeded`
- `payment_intent.payment_failed`
- `invoice.payment_succeeded`
- `invoice.payment_failed`
- `customer.subscription.created`
- `customer.subscription.updated`
- `customer.subscription.deleted`

## Docker Deployment

The service includes Docker configuration and can be deployed as part of the microservices stack:

```yaml
# In docker-compose.yml
payments-service:
  build: ./services/payments
  environment:
    DB_HOST: postgres
    STRIPE_SECRET_KEY: sk_test_...
  ports:
    - "3004:3004"
```

## Security Considerations

- Stripe webhook signature verification
- Input validation and sanitization
- SQL injection protection via GORM
- Environment-based secret management
- CORS configuration for API access

## Monitoring

- Health check endpoints for service monitoring
- Structured logging with configurable levels
- Database connection health checks
- Stripe API integration status monitoring

## Error Handling

- Comprehensive error responses with proper HTTP status codes
- Stripe API error propagation and handling
- Database transaction rollback on failures
- Graceful handling of webhook delivery failures

## Future Enhancements

- Unit and integration test coverage
- Metrics and observability improvements
- Rate limiting for API endpoints
- Advanced Stripe features (Connect, Marketplace)
- Multi-currency support
- Advanced webhook retry logic