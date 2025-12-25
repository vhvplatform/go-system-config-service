# System Config Service Dependencies

## Shared Packages (from saas-shared-go)

```go
require (
    github.com/longvhv/saas-shared-go/config
    github.com/longvhv/saas-shared-go/logger
    github.com/longvhv/saas-shared-go/mongodb
    github.com/longvhv/saas-shared-go/redis
    github.com/longvhv/saas-shared-go/errors
    github.com/longvhv/saas-shared-go/middleware
    github.com/longvhv/saas-shared-go/response
)
```

## External Dependencies

### Infrastructure
- **MongoDB**: System configurations, feature flags
  - Collections: `system_configs`, `feature_flags`, `settings`
- **Redis**: Configuration cache
  - Keys: `config:*`, `feature:*`

### Third-party Libraries
```go
require (
    github.com/gin-gonic/gin v1.10.0
    google.golang.org/grpc v1.69.2
    go.mongodb.org/mongo-driver v1.17.3
)
```

## Inter-service Communication

### Services Called by System Config Service
- None (leaf service)

### Services Calling System Config Service
- **All Services**: Configuration and feature flag queries

## Environment Variables

```bash
# Server
SYSTEM_CONFIG_SERVICE_PORT=50055
SYSTEM_CONFIG_SERVICE_HTTP_PORT=8085

# Database
MONGODB_URI=mongodb://localhost:27017
MONGODB_DATABASE=saas_framework

# Redis
REDIS_URL=redis://localhost:6379/0

# Logging
LOG_LEVEL=info
```

## Database Schema

### Collections

#### system_configs
```json
{
  "_id": "ObjectId",
  "key": "string (unique, indexed)",
  "value": "mixed",
  "description": "string",
  "updated_at": "timestamp"
}
```

#### feature_flags
```json
{
  "_id": "ObjectId",
  "name": "string (unique, indexed)",
  "enabled": "boolean",
  "description": "string",
  "updated_at": "timestamp"
}
```

## Resource Requirements

### Production
- CPU: 0.5 cores
- Memory: 512MB
- Replicas: 2
