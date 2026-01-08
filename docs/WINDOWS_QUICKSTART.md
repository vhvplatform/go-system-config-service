# Quick Start - Windows Development

A concise guide to get started with the System Config Service on Windows quickly.

## Prerequisites (Quick Install)

1. **Install Go**: https://golang.org/dl/ (1.25.5 or higher)
2. **Install Git**: https://git-scm.com/download/win
3. **Install Docker Desktop** (optional): https://www.docker.com/products/docker-desktop

## Quick Setup (5 minutes)

### 1. Clone and Setup
```powershell
# Clone repository
git clone https://github.com/vhvplatform/go-system-config-service.git
cd go-system-config-service

# Install dependencies
.\build.ps1 deps

# Install dev tools
.\build.ps1 install-tools
```

### 2. Start Dependencies (Docker)
```powershell
# Start MongoDB
docker run -d --name mongodb -p 27017:27017 mongo:7.0

# Start Redis
docker run -d --name redis -p 6379:6379 redis:7-alpine
```

### 3. Build and Run
```powershell
# Build
.\build.ps1 build

# Run
.\build.ps1 run
```

## Common Commands

### PowerShell
```powershell
.\build.ps1 build          # Build service
.\build.ps1 test           # Run tests
.\build.ps1 test-coverage  # Run tests with coverage
.\build.ps1 run            # Run service
.\build.ps1 clean          # Clean artifacts
.\build.ps1 fmt            # Format code
.\build.ps1 lint           # Run linters
```

### Command Prompt
```cmd
build.bat build          :: Build service
build.bat test           :: Run tests
build.bat run            :: Run service
build.bat clean          :: Clean artifacts
```

## Default Ports
- HTTP: `http://localhost:8085`
- gRPC: `localhost:50055`
- MongoDB: `localhost:27017`
- Redis: `localhost:6379`

## Quick Test
```powershell
# Health check
curl http://localhost:8085/health

# List countries
curl http://localhost:8085/api/v1/system-config/countries

# List app components
curl http://localhost:8085/api/v1/system-config/app-components
```

## Troubleshooting

### "go: command not found"
- Restart terminal after installing Go
- Verify: `go version`

### "execution policy" error
```powershell
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser
```

### Port already in use
```powershell
# Find and kill process on port 8085
netstat -ano | Select-String ":8085"
Stop-Process -Id <PID>
```

## Next Steps

- Read full guide: [docs/WINDOWS_SETUP.md](WINDOWS_SETUP.md)
- API documentation: [README.md](../README.md)
- Architecture: [ARCHITECTURE_ALIGNMENT.md](../ARCHITECTURE_ALIGNMENT.md)

---

**Need help?** See [docs/WINDOWS_SETUP.md](WINDOWS_SETUP.md) for detailed instructions.
