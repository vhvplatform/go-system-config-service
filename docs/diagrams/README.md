# System Configuration Service - Architecture Diagrams

This directory contains PlantUML diagrams documenting the architecture, flows, and mechanisms of the System Configuration Service.

## Diagrams Overview

### 1. Configuration Architecture (`01-configuration-architecture.puml`)
**Purpose**: Shows the overall system architecture including layers, patterns, and external systems.

**Key Components**:
- API Layer (HTTP/gRPC servers)
- Handler Layer (Request processing)
- Service Layer (Business logic)
- Repository Layer (Data access)
- External Systems (MongoDB, Redis, RabbitMQ)
- Design Patterns (Observer, Strategy, Factory)

**Use Cases**: Understanding the overall system structure, component relationships, and data flow.

---

### 2. Hot Reload Flow (`02-hot-reload-flow.puml`)
**Purpose**: Demonstrates how configuration changes are detected and propagated to subscribers in real-time.

**Key Steps**:
1. Configuration file modification detection
2. Validation of new configuration
3. Version creation and storage
4. Cache invalidation and update
5. Subscriber notification via events
6. Rollback mechanism on failure

**Use Cases**: Implementing hot reload, debugging configuration propagation issues, understanding the event-driven architecture.

---

### 3. Config Versioning (`03-config-versioning.puml`)
**Purpose**: Illustrates the version control system for configurations, including creation, activation, and rollback.

**Key Features**:
- Sequential version management (v1, v2, v3...)
- Version activation workflow
- Complete version history
- Rollback to previous versions
- Version comparison

**Use Cases**: Managing configuration versions, performing safe rollbacks, auditing configuration changes.

---

### 4. Multi-Environment Diagram (`04-multi-environment.puml`)
**Purpose**: Shows how configurations are managed across different environments (dev, staging, production).

**Key Concepts**:
- Configuration inheritance hierarchy
- Environment-specific overrides
- Merge strategies
- Shared global configurations
- Environment isolation

**Use Cases**: Deploying to multiple environments, managing environment-specific settings, understanding configuration precedence.

---

### 5. Secret Management (`05-secret-management.puml`)
**Purpose**: Details the secure storage, retrieval, and rotation of secrets.

**Security Features**:
- AES-256-GCM encryption
- Key management with external Vault/KMS
- Access control and permissions
- Audit logging
- Automatic secret rotation
- HMAC integrity verification

**Use Cases**: Storing sensitive data securely, implementing secret rotation policies, compliance requirements (PCI-DSS, GDPR).

---

### 6. Database Schema (`06-database-schema.puml`)
**Purpose**: Provides the complete MongoDB database schema with all collections, fields, and relationships.

**Collections**:
- `configs`: Main configuration storage
- `config_versions`: Version history
- `config_audit_log`: Audit trail (TTL: 2 years)
- `secrets`: Encrypted secret storage
- `secret_access_log`: Secret access tracking (TTL: 1 year)
- `app_components`: Application components
- `countries`: Country master data
- `watch_subscriptions`: Configuration watchers
- `change_approvals`: Approval workflow

**Use Cases**: Database design, index optimization, query planning, data modeling.

---

### 7. Watch Mechanism (`07-watch-mechanism.puml`)
**Purpose**: Explains how services subscribe to configuration changes and receive notifications.

**Features**:
- File-based watching (fsnotify)
- Database change streams (MongoDB)
- Pattern-based subscriptions (glob patterns)
- Batch notifications
- Retry with exponential backoff
- Health monitoring of subscribers

**Use Cases**: Implementing configuration watchers, setting up webhooks, debugging notification issues.

---

## How to View Diagrams

### Online Viewers
1. **PlantUML Web Server**: http://www.plantuml.com/plantuml/uml/
   - Copy and paste the diagram code
2. **PlantText**: https://www.planttext.com/
   - Paste the code and view instantly

### VS Code Extension
Install the "PlantUML" extension:
```bash
code --install-extension jebbs.plantuml
```

### Command Line
Install PlantUML and generate images:
```bash
# Install PlantUML (requires Java)
brew install plantuml  # macOS
apt-get install plantuml  # Ubuntu/Debian

# Generate PNG images
plantuml docs/diagrams/*.puml

# Generate SVG images (scalable)
plantuml -tsvg docs/diagrams/*.puml
```

### Docker
Use PlantUML Docker image:
```bash
docker run -v $(pwd)/docs/diagrams:/data plantuml/plantuml:latest -tsvg /data/*.puml
```

---

## Diagram Conventions

### Color Coding
- **Blue**: API/Interface components
- **Green**: Service/Business logic
- **Yellow**: Data storage
- **Red**: External systems
- **Gray**: Utility/Helper components

### Notation
- **Solid Lines**: Synchronous calls
- **Dashed Lines**: Asynchronous/Event-driven
- **Thick Lines**: Primary data flow
- **Thin Lines**: Secondary/Optional flow

### Status Indicators
- ✓ : Success/Completed
- ✗ : Failed/Error
- ⚠ : Warning/Attention needed

---

## Updating Diagrams

When updating these diagrams:

1. Keep them synchronized with actual implementation
2. Update the version/date in comments if making significant changes
3. Test rendering before committing
4. Update this README if adding new diagrams
5. Follow PlantUML best practices for readability

---

## Related Documentation

- [Main README](../../README.md) - Service overview and setup
- [API Documentation](../API.md) - REST API reference (if exists)
- [Contributing Guide](../../CONTRIBUTING.md) - How to contribute
- [Dependencies](../DEPENDENCIES.md) - Dependency information

---

## Questions or Issues?

If you have questions about these diagrams or notice inconsistencies:
1. Check the main documentation
2. Review the actual implementation
3. Create an issue in the repository
4. Contact the development team

---

**Last Updated**: 2025-12-25  
**Maintained By**: VHV Corp Development Team
