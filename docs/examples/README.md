# Configuration Examples

This directory contains comprehensive examples for using the System Config Service.

## Directory Structure

```
examples/
├── configs/           # Configuration file examples
│   ├── base.yaml             # Base configuration (shared)
│   ├── development.yaml      # Development environment
│   ├── staging.yaml          # Staging environment
│   ├── production.yaml       # Production environment
│   └── feature-flags.yaml    # Feature flags example
├── scripts/           # Example scripts
│   ├── create-config.sh      # Create and activate configuration
│   ├── manage-secrets.sh     # Secret management operations
│   └── subscribe-watch.sh    # Subscribe to config changes
└── README.md          # This file
```

## Configuration Files

### Base Configuration (`configs/base.yaml`)
Shared configuration across all environments. Contains:
- Service settings
- Feature flags
- Database configuration
- Caching strategy
- API configuration
- Security settings
- Monitoring configuration

### Environment-Specific Configurations

**Development** (`configs/development.yaml`)
- Local development settings
- Debug logging enabled
- Relaxed security for testing
- Local database and cache

**Staging** (`configs/staging.yaml`)
- Pre-production environment
- Similar to production configuration
- Test data and services
- Moderate logging

**Production** (`configs/production.yaml`)
- Production-ready configuration
- High-performance settings
- Strict security
- Error-only logging
- Backup and monitoring enabled

### Feature Flags (`configs/feature-flags.yaml`)
Examples of dynamic feature management:
- Gradual rollouts (percentage-based)
- User/tenant allowlists
- Beta access controls
- Feature dependencies

## Usage Examples

### 1. Configuration Hierarchy

Configurations are merged in priority order:
```
Production values = Base + Production overrides + Runtime env vars
```

Example:
```yaml
# base.yaml
database:
  max_connections: 50
  timeout: "30s"

# production.yaml
database:
  max_connections: 200  # Overrides base

# Result in production:
database:
  max_connections: 200  # From production.yaml
  timeout: "30s"        # From base.yaml
```

### 2. Environment Variables

Use environment variables for sensitive data:
```yaml
database:
  uri: "${MONGODB_URI}"
  password: "${DB_PASSWORD}"

security:
  jwt:
    secret: "${JWT_SECRET}"
```

Set in environment:
```bash
export MONGODB_URI="mongodb://prod-server:27017"
export DB_PASSWORD="secure-password"
export JWT_SECRET="your-jwt-secret"
```

### 3. Loading Configurations

```go
// Load base + environment-specific config
cfg, err := config.LoadWithEnvironment("configs/base.yaml", "production")

// Or load single file
cfg, err := config.Load("configs/production.yaml")

// With environment variable substitution
cfg, err := config.LoadWithEnvSubstitution("configs/production.yaml")
```

## Example Scripts

### Create Configuration (`scripts/create-config.sh`)

Creates a new configuration with versioning:

```bash
export API_URL="http://localhost:8085"
export AUTH_TOKEN="your-auth-token"

./scripts/create-config.sh
```

This script:
1. Creates a new configuration
2. Activates the configuration version
3. Displays version history

### Manage Secrets (`scripts/manage-secrets.sh`)

Demonstrates secret management operations:

```bash
export API_URL="http://localhost:8085"
export AUTH_TOKEN="your-auth-token"

./scripts/manage-secrets.sh
```

This script:
1. Creates an encrypted secret
2. Retrieves the secret
3. Rotates the secret
4. Views audit log

### Subscribe to Changes (`scripts/subscribe-watch.sh`)

Sets up a webhook for configuration changes:

```bash
export API_URL="http://localhost:8085"
export AUTH_TOKEN="your-auth-token"
export CALLBACK_URL="https://your-service.com/webhook"

./scripts/subscribe-watch.sh
```

This script:
1. Subscribes to config patterns
2. Configures retry policy
3. Lists active subscriptions

## API Examples

### Create Configuration

```bash
curl -X POST http://localhost:8085/api/v1/configs \
  -H "Authorization: Bearer $AUTH_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "key": "api.timeout",
    "value": "30s",
    "environment": "production",
    "tenant_id": "tenant-123"
  }'
```

### Get Configuration

```bash
curl -X GET http://localhost:8085/api/v1/configs/api.timeout \
  -H "Authorization: Bearer $AUTH_TOKEN" \
  -H "X-Environment: production" \
  -H "X-Tenant-ID: tenant-123"
```

### Update Configuration

```bash
curl -X PUT http://localhost:8085/api/v1/configs/api.timeout \
  -H "Authorization: Bearer $AUTH_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "value": "60s",
    "reason": "Increase timeout for slow clients"
  }'
```

### Rollback Configuration

```bash
curl -X POST http://localhost:8085/api/v1/configs/rollback \
  -H "Authorization: Bearer $AUTH_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "config_key": "api.timeout",
    "target_version": "v2",
    "reason": "Previous version was unstable"
  }'
```

### Export Configurations

```bash
# Export all configs
curl -X GET http://localhost:8085/api/v1/configs/export \
  -H "Authorization: Bearer $AUTH_TOKEN" \
  -o backup.json

# Export specific environment
curl -X GET "http://localhost:8085/api/v1/configs/export?env=production" \
  -H "Authorization: Bearer $AUTH_TOKEN" \
  -o prod-backup.json
```

### Import Configurations

```bash
curl -X POST http://localhost:8085/api/v1/configs/import \
  -H "Authorization: Bearer $AUTH_TOKEN" \
  -H "Content-Type: application/json" \
  -d @backup.json
```

## Best Practices

### 1. Configuration Organization
- Keep base configuration DRY
- Use environment-specific files for overrides only
- Document all configuration options
- Use meaningful configuration keys

### 2. Secret Management
- Never commit secrets to version control
- Use environment variables for sensitive data
- Rotate secrets regularly (e.g., every 90 days)
- Enable audit logging for secret access

### 3. Versioning
- Tag versions with meaningful descriptions
- Test in staging before promoting to production
- Keep a rollback plan ready
- Document breaking changes

### 4. Hot Reload
- Test configuration changes in staging first
- Use gradual rollouts for critical changes
- Monitor metrics after configuration updates
- Have rollback automation ready

### 5. Multi-Environment
- Maintain parity between environments
- Use infrastructure as code
- Automate configuration deployment
- Test configuration migrations

## Troubleshooting

### Configuration Not Loading

Check file paths and permissions:
```bash
ls -la configs/
cat configs/base.yaml
```

### Environment Variable Not Substituted

Verify environment variables are set:
```bash
printenv | grep MONGODB_URI
echo $JWT_SECRET
```

### Configuration Not Updating

Clear cache and reload:
```bash
curl -X DELETE http://localhost:8085/api/v1/cache/clear \
  -H "Authorization: Bearer $AUTH_TOKEN"
```

Check hot reload status:
```bash
curl -X GET http://localhost:8085/api/v1/watch/status \
  -H "Authorization: Bearer $AUTH_TOKEN"
```

## Additional Resources

- [Main Documentation](../../README.md)
- [Architecture Diagrams](../diagrams/)
- [API Reference](../API.md)
- [Contributing Guide](../../CONTRIBUTING.md)

## Support

For questions or issues:
- GitHub Issues: https://github.com/vhvplatform/go-system-config-service/issues
- Email: support@vhvplatform.com
