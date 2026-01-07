# System Config Service

The System Config Service is an enterprise-grade microservice responsible for managing all common system configurations, master data, and secrets for the entire SaaS platform with advanced features like hot reload, version control, and multi-environment support.

## Repository Structure

This repository is organized into the following main directories:

```
.
├── server/          # Golang backend microservice
├── client/          # ReactJS frontend microservice
├── flutter/         # Flutter mobile application
└── docs/            # Project documentation
```

### Server (Backend)

The `server` directory contains the Golang backend microservice code. See [server/README.md](server/README.md) for detailed information about the backend service, including:
- Architecture and design patterns
- Setup and installation
- API documentation
- Development guidelines
- Build and deployment instructions

### Client (Frontend)

The `client` directory contains the ReactJS frontend microservice code for the web interface. See [client/README.md](client/README.md) for frontend-specific documentation.

### Flutter (Mobile)

The `flutter` directory contains the Flutter mobile application code for iOS and Android platforms. See [flutter/README.md](flutter/README.md) for mobile app documentation.

### Documentation

The `docs` directory contains shared project documentation, including:
- Architecture diagrams
- Development guides
- Windows setup instructions
- Performance optimization guides
- Dependencies documentation

## Quick Start

### Backend (Server)

```bash
cd server
make setup
make run
```

For detailed instructions, see [server/README.md](server/README.md).

### Frontend (Client)

Coming soon...

### Mobile App (Flutter)

Coming soon...

## Architecture

> **Architecture Alignment**: This service conforms to the architectural standards defined in [go-infrastructure](https://github.com/vhvplatform/go-infrastructure). See [server/ARCHITECTURE_ALIGNMENT.md](server/ARCHITECTURE_ALIGNMENT.md) for details.

For detailed architecture diagrams and documentation, see the [docs/](docs/) directory.

## Contributing

Please read [CONTRIBUTING.md](CONTRIBUTING.md) for details on our code of conduct and the process for submitting pull requests.

## License

This project is proprietary and confidential.
