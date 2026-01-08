# Windows Development Guide

This guide provides detailed instructions for setting up and working with the System Config Service on Windows.

## Table of Contents

- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Building the Service](#building-the-service)
- [Running Tests](#running-tests)
- [Development Workflow](#development-workflow)
- [Docker on Windows](#docker-on-windows)
- [Troubleshooting](#troubleshooting)
- [IDE Setup](#ide-setup)

## Prerequisites

### Required Software

1. **Go 1.25.5 or higher**
   - Download from: https://golang.org/dl/
   - Add to PATH during installation
   - Verify: `go version`

2. **Git for Windows**
   - Download from: https://git-scm.com/download/win
   - Use Git Bash or PowerShell
   - Verify: `git --version`

3. **MongoDB**
   - Option A: MongoDB Community Server
     - Download from: https://www.mongodb.com/try/download/community
     - Install as Windows Service
   - Option B: MongoDB Atlas (cloud)
     - Sign up at: https://www.mongodb.com/cloud/atlas
   - Option C: Docker Desktop (see below)

4. **Redis**
   - Option A: Redis for Windows (via WSL2 or Docker)
   - Option B: Docker Desktop (recommended)
   - Option C: Memurai (Redis-compatible for Windows)
     - Download from: https://www.memurai.com/

### Optional but Recommended

5. **Docker Desktop for Windows**
   - Download from: https://www.docker.com/products/docker-desktop
   - Requires Windows 10/11 Pro, Enterprise, or Education (64-bit)
   - Enables WSL2 backend for better performance

6. **Windows Terminal** (for better command-line experience)
   - Install from Microsoft Store
   - Or download from: https://github.com/microsoft/terminal

7. **Visual Studio Code**
   - Download from: https://code.visualstudio.com/
   - Recommended extensions:
     - Go (by Go Team at Google)
     - Docker
     - GitLens
     - REST Client

## Installation

### 1. Clone the Repository

```powershell
# Using PowerShell
cd C:\workspace
git clone https://github.com/vhvplatform/go-system-config-service.git
cd go-system-config-service
```

```cmd
:: Using Command Prompt
cd C:\workspace
git clone https://github.com/vhvplatform/go-system-config-service.git
cd go-system-config-service
```

### 2. Install Go Dependencies

```powershell
# Using PowerShell
.\build.ps1 deps
```

```cmd
:: Using Command Prompt
build.bat deps
```

Or manually:
```powershell
go mod download
go mod tidy
```

### 3. Install Development Tools

```powershell
# Using PowerShell
.\build.ps1 install-tools
```

```cmd
:: Using Command Prompt
build.bat install-tools
```

This installs:
- `golangci-lint` - Code linter
- `protoc-gen-go` - Protocol buffer compiler (if needed)
- `protoc-gen-go-grpc` - gRPC code generator (if needed)

**Important**: Add Go binary path to your system PATH:
```powershell
# Check your Go binary path
go env GOPATH

# Add to PATH (typically C:\Users\YourName\go\bin)
# System Properties > Environment Variables > Path > Add
```

### 4. Set Up Environment Variables

Create a `.env` file in the project root:

```env
# Service Configuration
SYSTEM_CONFIG_SERVICE_PORT=50055
SYSTEM_CONFIG_SERVICE_HTTP_PORT=8085
ENVIRONMENT=development

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

# RabbitMQ (Optional)
RABBITMQ_URL=amqp://guest:guest@localhost:5672/
RABBITMQ_EXCHANGE=config-events

# Logging
LOG_LEVEL=info
LOG_FORMAT=json
```

### 5. Start Dependencies with Docker (Recommended)

If using Docker Desktop:

```powershell
# Start MongoDB, Redis, and RabbitMQ
docker run -d --name mongodb -p 27017:27017 mongo:7.0
docker run -d --name redis -p 6379:6379 redis:7-alpine
docker run -d --name rabbitmq -p 5672:5672 -p 15672:15672 rabbitmq:3.12-management
```

Or use Docker Compose (create `docker-compose.yml`):

```yaml
version: '3.8'
services:
  mongodb:
    image: mongo:7.0
    ports:
      - "27017:27017"
    volumes:
      - mongodb_data:/data/db

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"

  rabbitmq:
    image: rabbitmq:3.12-management
    ports:
      - "5672:5672"
      - "15672:15672"
    environment:
      RABBITMQ_DEFAULT_USER: guest
      RABBITMQ_DEFAULT_PASS: guest

volumes:
  mongodb_data:
```

Start with:
```powershell
docker-compose up -d
```

## Building the Service

### Using PowerShell (Recommended)

```powershell
# Build the service
.\build.ps1 build

# The binary will be created at: bin\system-config-service.exe
```

### Using Batch Script

```cmd
:: Build the service
build.bat build

:: The binary will be created at: bin\system-config-service.exe
```

### Using Make (with Git Bash or WSL)

If you have Make installed (via Git Bash, WSL, or Chocolatey):

```bash
# Build the service
make build
```

### Manual Build

```powershell
# Create bin directory if it doesn't exist
if (!(Test-Path "bin")) { New-Item -ItemType Directory -Path "bin" }

# Build
go build -o bin\system-config-service.exe .\cmd\main.go
```

## Running Tests

### Using PowerShell

```powershell
# Run all tests
.\build.ps1 test

# Run tests with coverage
.\build.ps1 test-coverage

# This will generate coverage.html and open it in your browser
```

### Using Batch Script

```cmd
:: Run all tests
build.bat test

:: Run tests with coverage
build.bat test-coverage
```

### Manual Test Commands

```powershell
# Run all tests
go test -v -race .\...

# Run tests with coverage
go test -v -race -coverprofile=coverage.txt -covermode=atomic .\...
go tool cover -html=coverage.txt -o coverage.html

# Run specific package tests
go test -v .\internal\service\...

# Run specific test
go test -v -run TestConfigService .\internal\service\
```

## Development Workflow

### Running the Service Locally

```powershell
# Using PowerShell script
.\build.ps1 run

# Using batch script
build.bat run

# Or directly
go run .\cmd\main.go
```

### Code Formatting

```powershell
# Format all Go files
.\build.ps1 fmt

# Or manually
go fmt .\...
gofmt -s -w .
```

### Linting

```powershell
# Run linters
.\build.ps1 lint

# Or manually
golangci-lint run .\...
```

### Code Vetting

```powershell
# Run go vet
.\build.ps1 vet

# Or manually
go vet .\...
```

### Cleaning Build Artifacts

```powershell
# Clean up
.\build.ps1 clean
```

## Docker on Windows

### Prerequisites

- Docker Desktop for Windows installed and running
- WSL2 backend enabled (recommended)

### Build Docker Image

```powershell
# Using PowerShell script
.\build.ps1 docker-build

# Or manually
docker build -t system-config-service:latest .
```

### Run Docker Container

```powershell
# Run the service in Docker
docker run --rm -p 8085:8085 -p 50055:50055 `
  -e MONGODB_URI=mongodb://host.docker.internal:27017 `
  -e REDIS_HOST=host.docker.internal `
  --name system-config-service `
  system-config-service:latest
```

**Note**: Use `host.docker.internal` to connect to services running on your Windows host from within Docker containers.

## Troubleshooting

### Common Issues

#### 1. "go: command not found"

**Solution**: 
- Ensure Go is installed and added to PATH
- Restart your terminal/PowerShell after installation
- Verify with: `go version`

#### 2. "golangci-lint: command not found"

**Solution**:
```powershell
# Install tools
.\build.ps1 install-tools

# Verify GOPATH\bin is in PATH
$env:Path -split ';' | Select-String "go\bin"

# If not in PATH, add it:
$goPath = go env GOPATH
[Environment]::SetEnvironmentVariable("Path", $env:Path + ";$goPath\bin", "User")
```

#### 3. MongoDB Connection Errors

**Solution**:
- Verify MongoDB is running: `mongosh --eval "db.version()"`
- Check connection string in `.env`
- If using Docker: `docker ps | Select-String mongodb`
- Try connecting: `mongosh mongodb://localhost:27017`

#### 4. Redis Connection Errors

**Solution**:
- Verify Redis is running
- If using Docker: `docker ps | Select-String redis`
- Test connection: `redis-cli ping` (should return PONG)
- Or use: `docker exec -it redis redis-cli ping`

#### 5. Line Ending Issues (CRLF vs LF)

**Solution**:
```powershell
# Configure Git to handle line endings
git config --global core.autocrlf true

# For existing files
git add --renormalize .
```

#### 6. PowerShell Execution Policy

If you get "cannot be loaded because running scripts is disabled":

```powershell
# Check current policy
Get-ExecutionPolicy

# Set to RemoteSigned (run as Administrator)
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser

# Or run scripts with bypass
powershell -ExecutionPolicy Bypass -File .\build.ps1 build
```

#### 7. Port Already in Use

**Solution**:
```powershell
# Find process using port 8085
netstat -ano | Select-String ":8085"

# Kill process (replace PID)
Stop-Process -Id <PID> -Force

# Or use different ports
$env:SYSTEM_CONFIG_SERVICE_HTTP_PORT=8086
.\build.ps1 run
```

#### 8. Path Too Long Errors

**Solution**:
```powershell
# Enable long paths in Windows (requires admin)
New-ItemProperty -Path "HKLM:\SYSTEM\CurrentControlSet\Control\FileSystem" `
  -Name "LongPathsEnabled" -Value 1 -PropertyType DWORD -Force

# Or clone to shorter path like C:\dev\
```

### Firewall Issues

If you have firewall issues:

1. Allow Go through Windows Firewall
2. Add rules for ports 8085 (HTTP) and 50055 (gRPC)

```powershell
# Run as Administrator
New-NetFirewallRule -DisplayName "Go System Config Service HTTP" `
  -Direction Inbound -Protocol TCP -LocalPort 8085 -Action Allow

New-NetFirewallRule -DisplayName "Go System Config Service gRPC" `
  -Direction Inbound -Protocol TCP -LocalPort 50055 -Action Allow
```

## IDE Setup

### Visual Studio Code

1. **Install Extensions**:
   - Go (ms-vscode.go)
   - Docker (ms-azuretools.vscode-docker)
   - GitLens (eamodio.gitlens)
   - REST Client (humao.rest-client)

2. **Configure Settings** (`.vscode/settings.json`):
```json
{
  "go.toolsManagement.autoUpdate": true,
  "go.useLanguageServer": true,
  "go.lintTool": "golangci-lint",
  "go.lintOnSave": "workspace",
  "go.formatTool": "gofmt",
  "go.formatOnSave": true,
  "editor.formatOnSave": true,
  "files.eol": "\n",
  "go.testFlags": ["-v", "-race"],
  "go.coverOnSave": true
}
```

3. **Add Tasks** (`.vscode/tasks.json`):
```json
{
  "version": "2.0.0",
  "tasks": [
    {
      "label": "build",
      "type": "shell",
      "command": "powershell",
      "args": ["-File", ".\\build.ps1", "build"],
      "group": {
        "kind": "build",
        "isDefault": true
      }
    },
    {
      "label": "test",
      "type": "shell",
      "command": "powershell",
      "args": ["-File", ".\\build.ps1", "test"],
      "group": {
        "kind": "test",
        "isDefault": true
      }
    }
  ]
}
```

4. **Add Launch Configuration** (`.vscode/launch.json`):
```json
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Launch Service",
      "type": "go",
      "request": "launch",
      "mode": "auto",
      "program": "${workspaceFolder}\\cmd\\main.go",
      "env": {
        "MONGODB_URI": "mongodb://localhost:27017",
        "REDIS_HOST": "localhost",
        "LOG_LEVEL": "debug"
      },
      "console": "integratedTerminal"
    }
  ]
}
```

### GoLand / IntelliJ IDEA

1. Import project as Go Module
2. Set Go SDK to 1.25.5 or higher
3. Enable Go Modules integration
4. Configure Run Configuration:
   - Working directory: project root
   - Environment variables from `.env`

## Performance Tips

### WSL2 vs Native Windows

For better performance on Windows:

1. **Use WSL2** for development (recommended)
   - Install WSL2: `wsl --install`
   - Install Ubuntu: `wsl --install -d Ubuntu`
   - Clone repo in WSL filesystem: `/home/username/projects/`
   - Use VS Code Remote-WSL extension

2. **Exclude from Windows Defender**
   - Add project folder to exclusions
   - Add Go installation folder
   - Significantly improves build times

```powershell
# Add exclusions (run as Administrator)
Add-MpPreference -ExclusionPath "C:\workspace\go-system-config-service"
Add-MpPreference -ExclusionPath "C:\Go"
Add-MpPreference -ExclusionPath "$env:USERPROFILE\go"
```

### Build Performance

```powershell
# Enable Go module caching
$env:GOCACHE="$env:LOCALAPPDATA\go-build"

# Use parallel builds
go build -p 8 -o bin\system-config-service.exe .\cmd\main.go
```

## Additional Resources

- [Go on Windows](https://golang.org/doc/install/windows)
- [Docker Desktop for Windows](https://docs.docker.com/desktop/windows/)
- [WSL2 Setup](https://docs.microsoft.com/en-us/windows/wsl/install)
- [MongoDB on Windows](https://docs.mongodb.com/manual/tutorial/install-mongodb-on-windows/)
- [Visual Studio Code Go Extension](https://code.visualstudio.com/docs/languages/go)

## Getting Help

If you encounter issues not covered here:

1. Check the main [README.md](../README.md)
2. Search [GitHub Issues](https://github.com/vhvplatform/go-system-config-service/issues)
3. Create a new issue with:
   - Windows version
   - Go version (`go version`)
   - Error messages
   - Steps to reproduce

---

**Last Updated**: 2025-12-30  
**Tested on**: Windows 10/11 (64-bit)
