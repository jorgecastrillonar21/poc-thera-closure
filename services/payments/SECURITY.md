# Security Configuration for Payments Service

## Environment Variables for Production

```bash
# JWT Authentication
export JWT_SECRET="your-very-secure-jwt-secret-key-change-in-production"

# Enable Authentication
export ENABLE_AUTHENTICATION="true"

# Rate Limiting
export RATE_LIMIT_WINDOW="1m"    # 1 minute window
export RATE_LIMIT_RPS="100"      # 100 requests per window

# Stripe Security
export STRIPE_WEBHOOK_SECRET="whsec_your_webhook_secret_from_stripe"

# Request Size Limits
export MAX_REQUEST_SIZE="10485760"  # 10MB
```

## API Key Authentication

API keys can be configured for service-to-service authentication:

```yaml
security:
  api_keys:
    "frontend-service": "api_key_for_frontend_12345"
    "admin-service": "api_key_for_admin_67890"
    "monitoring": "api_key_for_monitoring_abcdef"
```

## Basic Authentication

For simple authentication scenarios:

```yaml
security:
  basic_auth_users:
    "admin": "secure_password_123"
    "monitor": "monitoring_password_456"
```

## Rate Limiting Configuration

The service supports multiple rate limiting strategies:

### 1. IP-based Rate Limiting (Default)
- Limits requests per IP address
- Fallback to local storage if Redis unavailable

### 2. User-based Rate Limiting
- Limits requests per authenticated user
- Requires authentication middleware

### 3. Endpoint-based Rate Limiting
- Different limits per endpoint
- Prevents abuse of specific endpoints

## Security Headers

The service automatically adds security headers:

- `X-RateLimit-Limit`: Current rate limit
- `X-RateLimit-Remaining`: Remaining requests in window
- `X-RateLimit-Reset`: When the rate limit resets

## Webhook Security

Stripe webhook signatures are automatically verified:

1. **Signature Verification**: HMAC-SHA256 verification
2. **Timestamp Validation**: Prevents replay attacks (5-minute window)
3. **Request Body Integrity**: Ensures payload hasn't been tampered with

## Request Validation

Automatic validation includes:

- **Content-Type validation**: Only allows JSON for API endpoints
- **Request size limits**: Prevents large payload attacks
- **HTTP method validation**: Blocks unsupported methods
- **User-Agent filtering**: Blocks suspicious/malicious agents
- **Required headers**: Ensures necessary headers are present

## Security Logging

All security events are logged:

- Failed authentication attempts
- Rate limit violations
- Invalid webhook signatures
- Blocked requests
- Suspicious user agents

## Production Security Checklist

- [ ] Change default JWT secret
- [ ] Enable authentication (`ENABLE_AUTHENTICATION=true`)
- [ ] Configure appropriate rate limits
- [ ] Set up proper Stripe webhook secrets
- [ ] Configure API keys for service-to-service communication
- [ ] Review and adjust request size limits
- [ ] Monitor security logs for suspicious activity
- [ ] Use HTTPS in production
- [ ] Configure proper CORS settings
- [ ] Implement proper secret management (e.g., HashiCorp Vault, AWS Secrets Manager)

## Security Middleware Order

The middleware is applied in this order for optimal security:

1. Request ID generation
2. Recovery middleware
3. Logging middleware
4. Error handling middleware
5. Security logging middleware
6. Request size validation
7. Rate limiting
8. Request validation
9. Webhook signature verification
10. Authentication (JWT/API Key/Basic Auth)
11. CORS handling

This order ensures proper security validation while maintaining performance and debuggability.