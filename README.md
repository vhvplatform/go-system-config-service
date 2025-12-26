# System Config Service

The System Config Service is an enterprise-grade microservice responsible for managing all common system configurations, master data, and secrets for the entire SaaS platform with advanced features like hot reload, version control, and multi-environment support.

## Overview

This service provides centralized, secure, and dynamic management of:

- **Application Components**: Core application components and their configurations
- **SaaS Modules**: Available modules for the SaaS platform
- **Service Packages**: Subscription tiers and pricing packages
- **Admin Menus**: Dynamic menu configurations for admin interfaces
- **Permissions**: Fine-grained permissions for RBAC
- **Roles**: Role definitions and their associated permissions
- **Master Data**: Countries, ethnicities, provinces, districts, wards, currencies
- **Secrets Management**: Secure storage and retrieval of sensitive data
- **Configuration Versioning**: Track all changes with rollback capability

## Key Features

### Core Features
- **Multi-tenancy Support**: Tenant-specific configurations for customizable entities
- **Global Master Data**: Shared master data across all tenants
- **Redis Caching**: Performance optimization with intelligent caching strategies
- **MongoDB Storage**: Flexible document storage with optimized indexing
- **REST APIs**: Complete RESTful API endpoints for all entities
- **gRPC Support**: High-performance inter-service communication
- **Internationalization**: Multi-language support (en, vi) for master data

### Advanced Configuration Management
- **Hot Reload**: Real-time configuration updates without service restart
- **Configuration Versioning**: Complete version history with rollback capability
- **Multi-Environment Support**: Separate configs for dev, staging, and production
- **Configuration Validation**: Schema validation before applying changes
- **Change Approval Workflow**: Multi-step approval for critical configurations
- **Configuration Templates**: Reusable configuration templates
- **Import/Export**: Bulk configuration import and export

### Security & Compliance
- **Secret Encryption**: AES-256-GCM encryption for secrets at rest
- **Access Control**: Fine-grained RBAC for configuration access
- **Audit Logging**: Complete audit trail for all configuration changes
- **Secret Rotation**: Automatic and manual secret rotation policies
- **Sensitive Data Masking**: Automatic masking in logs and responses
- **Compliance Ready**: Supports PCI-DSS, GDPR, and SOC 2 requirements

### Reliability & Performance
- **Atomic Updates**: ACID-compliant configuration updates
- **Automatic Rollback**: Rollback on validation or deployment failures
- **Distributed Caching**: Redis-based caching with TTL and invalidation
- **Watch Mechanism**: Subscribe to configuration changes via webhooks
- **Batch Operations**: Efficient bulk updates and queries
- **High Availability**: Designed for multi-region deployments

## Tech Stack

- **Language**: Go 1.25+
- **Framework**: Gin
- **Database**: MongoDB
- **Cache**: Redis
- **Message Queue**: RabbitMQ
- **gRPC**: For inter-service communication

## Architecture

For detailed architecture diagrams, see [docs/diagrams/](docs/diagrams/).

```
┌─────────────────────────────────────────────────────────────┐
│                     API Layer                                │
│  HTTP Server (Gin) │ gRPC Server │ Health Checks             │
└─────────────────────────────────────────────────────────────┘
                          ↓
┌─────────────────────────────────────────────────────────────┐
│                   Handler Layer                              │
│  Config │ Secret │ Country │ Component │ Validation          │
└─────────────────────────────────────────────────────────────┘
                          ↓
┌─────────────────────────────────────────────────────────────┐
│                   Service Layer                              │
│  Business Logic │ Validation │ Encryption │ Audit            │
│  Watch Service │ Observer Pattern │ Strategy Pattern         │
└─────────────────────────────────────────────────────────────┘
                          ↓
┌─────────────────────────────────────────────────────────────┐
│                 Repository Layer                             │
│  Config Repo │ Version Repo │ Audit Repo │ Secret Repo       │
└─────────────────────────────────────────────────────────────┘
                          ↓
┌──────────────────┬──────────────────┬──────────────────────┐
│    MongoDB       │   Redis Cache    │   RabbitMQ Queue     │
│  Primary Store   │   Performance    │   Event Streaming    │
└──────────────────┴──────────────────┴──────────────────────┘
```

### Design Patterns
- **Observer Pattern**: Configuration change notifications
- **Strategy Pattern**: Multiple configuration sources (file, DB, remote)
- **Factory Pattern**: Configuration loaders and parsers
- **Repository Pattern**: Data access abstraction

## API Endpoints

### Configuration Management
- `GET    /api/v1/configs` - List all configurations
- `GET    /api/v1/configs/:key` - Get specific configuration
- `POST   /api/v1/configs` - Create new configuration
- `PUT    /api/v1/configs/:key` - Update configuration
- `DELETE /api/v1/configs/:key` - Delete configuration
- `POST   /api/v1/configs/:key/activate` - Activate configuration version
- `POST   /api/v1/configs/:key/rollback` - Rollback to previous version
- `GET    /api/v1/configs/:key/history` - Get version history
- `GET    /api/v1/configs/compare?v1=:v1&v2=:v2` - Compare versions

### Secret Management
- `GET    /api/v1/secrets` - List secrets (masked values)
- `GET    /api/v1/secrets/:key` - Get secret value
- `POST   /api/v1/secrets` - Create new secret
- `PUT    /api/v1/secrets/:key` - Update secret
- `DELETE /api/v1/secrets/:key` - Delete secret
- `POST   /api/v1/secrets/:key/rotate` - Rotate secret
- `GET    /api/v1/secrets/:key/audit` - Get secret access audit log

### Watch Subscriptions
- `POST   /api/v1/watch/subscribe` - Subscribe to config changes
- `DELETE /api/v1/watch/unsubscribe/:id` - Unsubscribe from notifications
- `GET    /api/v1/watch/subscriptions` - List active subscriptions

### Application Components
- `GET    /api/v1/system-config/app-components`
- `GET    /api/v1/system-config/app-components/:id`
- `POST   /api/v1/system-config/app-components`
- `PUT    /api/v1/system-config/app-components/:id`
- `DELETE /api/v1/system-config/app-components/:id`

### Countries
- `GET    /api/v1/system-config/countries`
- `GET    /api/v1/system-config/countries/:code`
- `POST   /api/v1/system-config/countries`
- `PUT    /api/v1/system-config/countries/:code`
- `DELETE /api/v1/system-config/countries/:code`

### Health Checks
- `GET /health` - Service health check
- `GET /ready` - Readiness probe
- `GET /metrics` - Prometheus metrics

## Environment Variables

```bash
# Service Configuration
SYSTEM_CONFIG_SERVICE_PORT=50055       # gRPC port
SYSTEM_CONFIG_SERVICE_HTTP_PORT=8085   # HTTP port
ENVIRONMENT=development                 # Environment: development|staging|production

# MongoDB
MONGODB_URI=mongodb://localhost:27017
MONGODB_DATABASE=saas_framework
MONGODB_MAX_POOL_SIZE=100
MONGODB_MIN_POOL_SIZE=10

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0
REDIS_MAX_RETRIES=3
REDIS_POOL_SIZE=10

# RabbitMQ (Optional - for event streaming)
RABBITMQ_URL=amqp://guest:guest@localhost:5672/
RABBITMQ_EXCHANGE=config-events

# Secret Management
ENCRYPTION_KEY_PATH=/path/to/encryption/key
VAULT_ADDR=https://vault.example.com
VAULT_TOKEN=your-vault-token
SECRET_ROTATION_DAYS=90

# Hot Reload
WATCH_ENABLED=true
WATCH_CONFIG_DIR=/etc/config
WATCH_DEBOUNCE_MS=1000

# Caching
CACHE_TTL_DEFAULT=3600              # 1 hour in seconds
CACHE_TTL_MASTER_DATA=86400         # 24 hours
CACHE_TTL_SECRETS=300               # 5 minutes

# Logging
LOG_LEVEL=info                      # debug|info|warn|error
LOG_FORMAT=json                     # json|text

# Security
ENABLE_AUTH=true
JWT_SECRET=your-jwt-secret
CORS_ALLOWED_ORIGINS=http://localhost:3000,https://app.example.com

# Performance
MAX_CONCURRENT_REQUESTS=1000
REQUEST_TIMEOUT_SECONDS=30
```

## Running Locally

### Prerequisites
- Go 1.25+
- MongoDB
- Redis

### Build and Run

```bash
# Install dependencies
cd services/system-config-service
go mod download

# Build
go build -o bin/system-config-service ./cmd/main.go

# Run
./bin/system-config-service
```

## Running with Docker

```bash
# Build Docker image
docker build -f services/system-config-service/Dockerfile -t system-config-service:latest .

# Run container
docker run -p 8085:8085 -p 50055:50055 \
  -e MONGODB_URI=mongodb://mongo:27017 \
  -e REDIS_HOST=redis \
  system-config-service:latest
```

## Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html

# Run specific test
go test -v -run TestConfigService ./internal/service/

# Run integration tests
go test -tags=integration ./...

# Run benchmarks
go test -bench=. ./...
```

## Configuration Management

### Hot Reload Mechanism

The service supports hot reload of configurations without requiring a restart:

1. **File-based watching**: Monitors config files using `fsnotify`
2. **Database change streams**: Listens to MongoDB change streams
3. **Validation**: Validates new configs before applying
4. **Atomic updates**: Ensures consistency during updates
5. **Event notification**: Notifies subscribers via RabbitMQ

**Example: Subscribe to config changes**
```bash
curl -X POST http://localhost:8085/api/v1/watch/subscribe \
  -H "Content-Type: application/json" \
  -d '{
    "patterns": ["db.*", "api.*.timeout"],
    "callback_url": "https://your-service.com/webhook/config-change",
    "environments": ["production"]
  }'
```

### Configuration Versioning

Every configuration change creates a new version:

```bash
# Create new config version
curl -X POST http://localhost:8085/api/v1/configs \
  -H "Content-Type: application/json" \
  -d '{
    "key": "db.timeout",
    "value": "30s",
    "environment": "production",
    "tenant_id": "tenant-123"
  }'

# Activate a version
curl -X POST http://localhost:8085/api/v1/configs/db.timeout/activate

# View version history
curl -X GET http://localhost:8085/api/v1/configs/db.timeout/history

# Rollback to previous version
curl -X POST http://localhost:8085/api/v1/configs/rollback \
  -H "Content-Type: application/json" \
  -d '{"target_version": "v2"}'

# Compare versions
curl -X GET "http://localhost:8085/api/v1/configs/compare?v1=v2&v2=v3"
```

### Environment-Based Configurations

Configurations can be environment-specific:

```yaml
# base-config.yaml (shared across all environments)
app:
  name: "System Config Service"
  features:
    hot_reload: true
    versioning: true

# production.yaml (production overrides)
database:
  pool_size: 100
  timeout: "30s"
  
logging:
  level: "error"
  
cache:
  ttl: 3600

# development.yaml (development overrides)
database:
  pool_size: 10
  timeout: "5s"
  
logging:
  level: "debug"
  
cache:
  ttl: 60
```

**Configuration Priority** (highest to lowest):
1. Runtime overrides
2. Environment-specific secrets
3. Environment-specific configs
4. Global configuration
5. Base/default values

### Configuration Validation

All configurations are validated before applying:

```go
// Define validation rules
type ConfigValidationRule struct {
    Key         string
    Type        string  // string, int, bool, duration, url, email
    Required    bool
    MinValue    interface{}
    MaxValue    interface{}
    Pattern     string  // regex pattern
    AllowedValues []string
}

// Example validation rules
rules := []ConfigValidationRule{
    {
        Key:      "db.timeout",
        Type:     "duration",
        Required: true,
        MinValue: "1s",
        MaxValue: "60s",
    },
    {
        Key:      "api.rate_limit",
        Type:     "int",
        Required: true,
        MinValue: 10,
        MaxValue: 10000,
    },
}
```

### Secret Management

Secure storage and management of sensitive data:

```bash
# Store a secret
curl -X POST http://localhost:8085/api/v1/secrets \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{
    "key": "db_password",
    "value": "super-secret-password",
    "environment": "production",
    "tenant_id": "tenant-123",
    "metadata": {
      "description": "Production database password",
      "rotation_days": 90
    }
  }'

# Retrieve a secret
curl -X GET http://localhost:8085/api/v1/secrets/db_password \
  -H "Authorization: Bearer <token>"

# Rotate a secret
curl -X POST http://localhost:8085/api/v1/secrets/db_password/rotate \
  -H "Authorization: Bearer <token>"

# View secret audit log
curl -X GET http://localhost:8085/api/v1/secrets/db_password/audit \
  -H "Authorization: Bearer <token>"
```

**Security Features**:
- AES-256-GCM encryption at rest
- Unique encryption key per secret or tenant
- HMAC for integrity verification
- Automatic rotation policies
- Complete audit trail
- Role-based access control

### Configuration Backup & Restore

```bash
# Export all configurations
curl -X GET http://localhost:8085/api/v1/configs/export \
  -H "Authorization: Bearer <token>" \
  -o configs-backup.json

# Export specific environment
curl -X GET "http://localhost:8085/api/v1/configs/export?env=production" \
  -H "Authorization: Bearer <token>" \
  -o prod-configs.json

# Import configurations
curl -X POST http://localhost:8085/api/v1/configs/import \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d @configs-backup.json

# Restore to specific point in time
curl -X POST http://localhost:8085/api/v1/configs/restore \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "timestamp": "2024-01-01T00:00:00Z",
    "environment": "production"
  }'
```

## Caching Strategy

The service implements a sophisticated multi-level caching strategy:

### Cache Layers
1. **Master Data** (Countries, Currencies, Ethnicities): TTL 24 hours
2. **Configuration Data** (App Components, Modules): TTL 1 hour  
3. **Secrets** (Decrypted values): TTL 5 minutes
4. **API Responses**: TTL configurable per endpoint

### Cache Keys Pattern
```
system-config:{type}:{tenant_id}:{environment}:{key}
```

### Cache Invalidation
Automatic invalidation on:
- Create operations
- Update operations
- Delete operations
- Version activation
- Secret rotation
- Manual cache clear

### Cache Optimization
```go
// Batch cache operations
keys := []string{"db.timeout", "db.pool_size", "db.max_conn"}
configs := cacheService.GetMulti(keys)

// Cache warming on startup
cacheService.WarmUpCache([]string{
    "master_data:countries",
    "master_data:currencies",
    "app_components",
})

// Distributed locking for cache updates
lock := cacheService.AcquireLock("config:db.timeout", 10*time.Second)
defer lock.Release()
```

## Monitoring & Observability

### Metrics
The service exposes Prometheus metrics:

```prometheus
# Configuration metrics
config_changes_total{environment="production",tenant="xyz"} 156
config_versions_total{environment="production"} 1024
config_validation_failures_total 23
config_rollbacks_total 5

# Secret metrics
secret_access_total{secret_key="db_password",result="success"} 45678
secret_rotation_total{secret_key="db_password"} 12
secret_encryption_duration_seconds 0.023

# Cache metrics
cache_hits_total{cache_type="config"} 98765
cache_misses_total{cache_type="config"} 1234
cache_evictions_total 456
cache_size_bytes{cache_type="config"} 2048576

# Watch metrics
watch_subscriptions_active 45
watch_notifications_sent_total 12345
watch_notification_failures_total 23

# Performance metrics
http_request_duration_seconds{endpoint="/api/v1/configs",method="GET"} 0.045
mongodb_query_duration_seconds{operation="find"} 0.012
```

### Logging
Structured logging with contextual information:

```json
{
  "timestamp": "2024-01-15T10:30:00Z",
  "level": "info",
  "message": "Configuration updated",
  "config_key": "db.timeout",
  "old_value": "20s",
  "new_value": "30s",
  "version": "v5",
  "user_id": "user-123",
  "tenant_id": "tenant-xyz",
  "environment": "production",
  "trace_id": "abc-123-def"
}
```

### Tracing
Distributed tracing with OpenTelemetry:
- Trace config change propagation
- Track hot reload flow
- Monitor secret access patterns
- Identify performance bottlenecks

### Alerting Rules

```yaml
# Critical alerts
- alert: ConfigValidationFailureRate
  expr: rate(config_validation_failures_total[5m]) > 0.1
  severity: critical
  
- alert: SecretAccessFailureRate
  expr: rate(secret_access_total{result="failure"}[5m]) > 0.05
  severity: critical

- alert: WatchNotificationFailureRate  
  expr: rate(watch_notification_failures_total[5m]) > 0.1
  severity: warning

# Performance alerts
- alert: HighConfigQueryLatency
  expr: histogram_quantile(0.95, http_request_duration_seconds) > 1
  severity: warning

- alert: CacheHitRateLow
  expr: rate(cache_hits_total[5m]) / rate(cache_requests_total[5m]) < 0.8
  severity: info
```

## MongoDB Indexes

The service creates optimized indexes for query performance:

```javascript
// Configurations
db.configs.createIndex({ "tenant_id": 1, "config_key": 1, "environment": 1 }, { unique: true });
db.configs.createIndex({ "tenant_id": 1, "environment": 1, "status": 1 });
db.configs.createIndex({ "updated_at": 1 });
db.configs.createIndex({ "tags": 1 });

// Configuration Versions
db.config_versions.createIndex({ "config_id": 1, "version_number": 1 });
db.config_versions.createIndex({ "config_id": 1, "status": 1 });
db.config_versions.createIndex({ "tenant_id": 1, "created_at": -1 });

// Audit Logs (with TTL)
db.config_audit_log.createIndex({ "config_id": 1, "timestamp": -1 });
db.config_audit_log.createIndex({ "tenant_id": 1, "timestamp": -1 });
db.config_audit_log.createIndex({ "user_id": 1, "timestamp": -1 });
db.config_audit_log.createIndex({ "timestamp": 1 }, { expireAfterSeconds: 63072000 }); // 2 years

// Secrets
db.secrets.createIndex({ "tenant_id": 1, "secret_key": 1, "environment": 1 }, { unique: true });
db.secrets.createIndex({ "tenant_id": 1, "environment": 1 });
db.secrets.createIndex({ "expires_at": 1 });
db.secrets.createIndex({ "last_rotated_at": 1 });

// App Components
db.app_components.createIndex({ "tenant_id": 1, "code": 1 }, { unique: true });
db.app_components.createIndex({ "tenant_id": 1, "status": 1 });

// Countries
db.countries.createIndex({ "code": 1 }, { unique: true });
db.countries.createIndex({ "status": 1 });

// Watch Subscriptions
db.watch_subscriptions.createIndex({ "subscriber_id": 1, "tenant_id": 1 });
db.watch_subscriptions.createIndex({ "service_name": 1, "status": 1 });
```

## Best Practices

### Configuration Management
1. **Use semantic versioning** for config versions
2. **Always validate** before activating configurations
3. **Test in staging** before promoting to production
4. **Use approval workflows** for critical configs
5. **Document changes** in version commit messages
6. **Keep configs DRY** using inheritance and templates
7. **Separate secrets** from regular configurations

### Security
1. **Encrypt all secrets** at rest using AES-256-GCM
2. **Rotate secrets regularly** (every 90 days recommended)
3. **Use RBAC** for configuration access control
4. **Enable audit logging** for compliance
5. **Mask sensitive data** in logs and API responses
6. **Use strong authentication** (JWT tokens)
7. **Implement rate limiting** to prevent abuse

### Performance
1. **Use caching** aggressively for read-heavy configs
2. **Batch operations** when possible
3. **Index frequently queried fields**
4. **Monitor cache hit rates** (target >80%)
5. **Use connection pooling** for databases
6. **Implement circuit breakers** for external dependencies
7. **Profile and optimize** slow queries

### Reliability
1. **Validate before apply** to prevent bad configs
2. **Enable automatic rollback** on failures
3. **Use health checks** for service monitoring
4. **Implement retry logic** with exponential backoff
5. **Test disaster recovery** procedures
6. **Monitor error rates** and set up alerts
7. **Maintain backup schedules**

### Development
1. **Use feature flags** for gradual rollouts
2. **Write comprehensive tests** (unit + integration)
3. **Document API changes** immediately
4. **Version your APIs** (use /v1/, /v2/ paths)
5. **Use consistent naming** conventions
6. **Review code** before merging
7. **Automate CI/CD** pipelines

## Migration & Rollback Strategies

### Configuration Migration

#### 1. Blue-Green Deployment
```bash
# Deploy new version alongside old
kubectl apply -f deployment-v2.yaml

# Switch traffic gradually
kubectl patch service config-service -p '{"spec":{"selector":{"version":"v2"}}}'

# Rollback if issues detected
kubectl patch service config-service -p '{"spec":{"selector":{"version":"v1"}}}'
```

#### 2. Canary Deployment
```bash
# Deploy canary with 10% traffic
kubectl apply -f deployment-canary.yaml
kubectl set image deployment/config-service-canary app=config-service:v2

# Monitor metrics for 30 minutes
# If successful, roll out to 100%
kubectl scale deployment/config-service-canary --replicas=10
kubectl scale deployment/config-service --replicas=0
```

#### 3. Database Migration
```go
// migration_v2.go
func MigrateV1ToV2(ctx context.Context) error {
    // Add new fields with default values
    _, err := db.Collection("configs").UpdateMany(
        ctx,
        bson.M{"schema_version": "v1"},
        bson.M{
            "$set": bson.M{
                "schema_version": "v2",
                "new_field": "default_value",
            },
        },
    )
    return err
}

// Rollback function
func RollbackV2ToV1(ctx context.Context) error {
    _, err := db.Collection("configs").UpdateMany(
        ctx,
        bson.M{"schema_version": "v2"},
        bson.M{
            "$set": bson.M{"schema_version": "v1"},
            "$unset": bson.M{"new_field": ""},
        },
    )
    return err
}
```

### Version Rollback Strategies

#### 1. Immediate Rollback
For critical issues requiring instant rollback:
```bash
# Rollback to previous version
curl -X POST http://localhost:8085/api/v1/configs/rollback \
  -H "Content-Type: application/json" \
  -d '{"target_version": "v_previous", "reason": "Critical bug"}'
```

#### 2. Gradual Rollback
For non-critical issues with staged rollback:
```bash
# Rollback 10% of tenants first
curl -X POST http://localhost:8085/api/v1/configs/rollback \
  -H "Content-Type: application/json" \
  -d '{
    "target_version": "v2",
    "strategy": "gradual",
    "percentage": 10,
    "tenant_filter": ["tenant-1", "tenant-2"]
  }'

# Monitor for issues
# If stable, continue rollback
curl -X POST http://localhost:8085/api/v1/configs/rollback/continue
```

#### 3. Automated Rollback
Configure automatic rollback on validation failures:
```yaml
rollback_policy:
  enabled: true
  triggers:
    - validation_failure
    - error_rate_threshold: 0.05
    - response_time_p95: 1000ms
  cooldown_period: 5m
  max_attempts: 3
```

### Breaking Change Management

When introducing breaking changes:

1. **Version the API**: Use `/api/v2/` for new endpoints
2. **Maintain backward compatibility**: Support both old and new formats
3. **Deprecation period**: Give users 90 days notice
4. **Migration tools**: Provide automated migration scripts
5. **Clear communication**: Document changes in CHANGELOG.md

Example:
```go
// Support both old and new formats
func (s *ConfigService) GetConfig(key string) (*Config, error) {
    config, err := s.repo.Get(key)
    if err != nil {
        return nil, err
    }
    
    // Migrate old format to new if needed
    if config.SchemaVersion == "v1" {
        config = s.migrateV1ToV2(config)
    }
    
    return config, nil
}
```

## Configuration Examples

### Example 1: Database Configuration with Multi-Environment

```yaml
# configs/base.yaml
database:
  driver: "mongodb"
  max_connections: 50
  timeout: "30s"
  retry_attempts: 3
  
# configs/development.yaml  
database:
  max_connections: 10
  log_queries: true
  
# configs/production.yaml
database:
  max_connections: 200
  enable_ssl: true
  replica_set: "rs0"
```

### Example 2: Feature Flags

```json
{
  "features": {
    "new_dashboard": {
      "enabled": true,
      "rollout_percentage": 10,
      "allowed_tenants": ["tenant-alpha", "tenant-beta"]
    },
    "advanced_analytics": {
      "enabled": false,
      "beta_access": true
    }
  }
}
```

### Example 3: API Configuration

```yaml
api:
  rate_limiting:
    enabled: true
    requests_per_second: 100
    burst_size: 200
  
  timeouts:
    read_timeout: "10s"
    write_timeout: "10s"
    idle_timeout: "120s"
    
  cors:
    allowed_origins:
      - "https://app.example.com"
      - "https://admin.example.com"
    allowed_methods: ["GET", "POST", "PUT", "DELETE"]
    allow_credentials: true
```

### Example 4: Secret Configuration

```bash
# Store database credentials
{
  "key": "db_credentials",
  "value": {
    "username": "admin",
    "password": "super-secret-password",
    "connection_string": "mongodb://admin:password@host:27017"
  },
  "environment": "production",
  "metadata": {
    "rotation_days": 90,
    "description": "Production database credentials"
  }
}
```

## Troubleshooting

### Common Issues

#### 1. Configuration Not Updating
```bash
# Check if configuration is cached
redis-cli GET "system-config:configs:tenant-123:production:db.timeout"

# Clear cache manually
redis-cli DEL "system-config:configs:tenant-123:production:db.timeout"

# Verify version is activated
curl http://localhost:8085/api/v1/configs/db.timeout/history
```

#### 2. Hot Reload Not Working
```bash
# Check watch subscriptions
curl http://localhost:8085/api/v1/watch/subscriptions

# Verify RabbitMQ connection
rabbitmqctl list_queues

# Check file watcher status
curl http://localhost:8085/api/v1/watch/status
```

#### 3. Secret Decryption Failure
```bash
# Verify encryption key is accessible
ls -la /path/to/encryption/key

# Check Vault connection
curl https://vault.example.com/v1/sys/health

# View secret audit log
curl http://localhost:8085/api/v1/secrets/db_password/audit
```

#### 4. High Latency
```bash
# Check MongoDB indexes
db.configs.getIndexes()

# Monitor cache hit rate
curl http://localhost:8085/metrics | grep cache_hits

# Check connection pool usage
curl http://localhost:8085/metrics | grep mongodb_pool
```

### Debug Mode

Enable debug logging:
```bash
# Set log level to debug
export LOG_LEVEL=debug

# Or update via API
curl -X PUT http://localhost:8085/api/v1/admin/log-level \
  -H "Content-Type: application/json" \
  -d '{"level": "debug"}'
```

## Performance Benchmarks

Expected performance metrics:

| Operation | Latency (p95) | Throughput |
|-----------|--------------|------------|
| Get Config (cached) | 5ms | 10,000 req/s |
| Get Config (uncached) | 20ms | 2,000 req/s |
| Update Config | 50ms | 500 req/s |
| List Configs | 30ms | 1,000 req/s |
| Get Secret | 15ms | 3,000 req/s |
| Hot Reload Notification | 100ms | N/A |

## Contributing

Please read [CONTRIBUTING.md](CONTRIBUTING.md) for details on our code of conduct and the process for submitting pull requests.

### Development Setup

1. Clone the repository
2. Install dependencies: `go mod download`
3. Set up MongoDB and Redis locally
4. Copy `.env.example` to `.env` and configure
5. Run tests: `go test ./...`
6. Start the service: `make run`

### Code Style

- Follow Go best practices and idioms
- Use `gofmt` and `golangci-lint`
- Write tests for new features
- Document exported functions
- Keep functions small and focused

## Roadmap

### Completed ✓
- [x] Multi-tenancy support
- [x] Redis caching
- [x] MongoDB storage with indexes
- [x] REST and gRPC APIs
- [x] Health checks and readiness probes
- [x] Configuration versioning system
- [x] Hot reload mechanism
- [x] Secret management with encryption
- [x] Multi-environment support
- [x] Audit logging
- [x] Watch subscriptions

### In Progress 🚧
- [ ] Complete unit and integration test suite (target >80% coverage)
- [ ] Comprehensive API documentation with Swagger/OpenAPI
- [ ] Performance optimization and benchmarking
- [ ] SonarQube integration and code quality improvements

### Planned 📋
- [ ] GraphQL API support
- [ ] Config diff and merge tools
- [ ] Advanced RBAC with custom roles
- [ ] Config templates and inheritance
- [ ] Disaster recovery automation
- [ ] Multi-region replication
- [ ] Config approval workflow UI
- [ ] Real-time config validation IDE plugin
- [ ] Config dependency graph visualization
- [ ] AI-powered config recommendations
- [ ] Automated config optimization
- [ ] Integration with external secret managers (AWS Secrets Manager, Azure Key Vault)

## Documentation

- [Architecture Diagrams](docs/diagrams/) - PlantUML diagrams
- [API Reference](docs/API.md) - Detailed API documentation
- [Dependencies](docs/DEPENDENCIES.md) - Dependency information
- [Contributing Guide](CONTRIBUTING.md) - How to contribute
- [Changelog](CHANGELOG.md) - Version history

## Support

- **Issues**: [GitHub Issues](https://github.com/vhvplatform/go-system-config-service/issues)
- **Discussions**: [GitHub Discussions](https://github.com/vhvplatform/go-system-config-service/discussions)
- **Email**: support@vhvplatform.com

## License

MIT License - see [LICENSE](LICENSE) file for details

---

**Maintained by**: VHV Corp Development Team  
**Last Updated**: 2025-12-25
