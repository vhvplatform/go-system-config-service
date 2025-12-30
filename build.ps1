<#
.SYNOPSIS
    Build script for Windows (PowerShell)
.DESCRIPTION
    Provides build, test, and development commands for the System Config Service on Windows.
.PARAMETER Command
    The command to execute: build, test, test-coverage, clean, run, deps, fmt, vet, lint, install-tools, docker-build, help
.EXAMPLE
    .\build.ps1 build
    .\build.ps1 test
    .\build.ps1 clean
#>

param(
    [Parameter(Position=0)]
    [ValidateSet('build', 'test', 'test-coverage', 'clean', 'run', 'deps', 'fmt', 'vet', 'lint', 'install-tools', 'docker-build', 'help')]
    [string]$Command = 'help'
)

$ErrorActionPreference = 'Stop'
$ServiceName = 'system-config-service'
$GoVersion = '1.25.5'
$DockerRegistry = $env:DOCKER_REGISTRY
if (-not $DockerRegistry) {
    $DockerRegistry = 'ghcr.io/vhvplatform'
}

function Write-ColorOutput($ForegroundColor) {
    $fc = $host.UI.RawUI.ForegroundColor
    $host.UI.RawUI.ForegroundColor = $ForegroundColor
    if ($args) {
        Write-Output $args
    }
    $host.UI.RawUI.ForegroundColor = $fc
}

function Show-Help {
    Write-ColorOutput Green "`nSystem Config Service - Windows Build Script"
    Write-Output "`nAvailable commands:"
    Write-Output "  build           - Build the service"
    Write-Output "  test            - Run tests"
    Write-Output "  test-coverage   - Run tests with coverage"
    Write-Output "  clean           - Clean build artifacts"
    Write-Output "  run             - Run the service locally"
    Write-Output "  deps            - Download dependencies"
    Write-Output "  fmt             - Format code"
    Write-Output "  vet             - Run go vet"
    Write-Output "  lint            - Run linters (requires golangci-lint)"
    Write-Output "  install-tools   - Install development tools"
    Write-Output "  docker-build    - Build Docker image (requires Docker Desktop)"
    Write-Output "  help            - Display this help screen"
    Write-Output ""
}

function Build-Service {
    Write-ColorOutput Cyan "Building $ServiceName..."
    
    if (-not (Test-Path "bin")) {
        New-Item -ItemType Directory -Path "bin" | Out-Null
    }
    
    try {
        go build -o "bin\$ServiceName.exe" .\cmd\main.go
        Write-ColorOutput Green "✓ Build complete! Binary: bin\$ServiceName.exe"
    }
    catch {
        Write-ColorOutput Red "✗ Build failed!"
        Write-Error $_
        exit 1
    }
}

function Test-Service {
    Write-ColorOutput Cyan "Running tests..."
    go test -v -race .\...
}

function Test-Coverage {
    Write-ColorOutput Cyan "Running tests with coverage..."
    
    try {
        go test -v -race -coverprofile=coverage.txt -covermode=atomic .\...
        go tool cover -html=coverage.txt -o coverage.html
        Write-ColorOutput Green "✓ Coverage report generated: coverage.html"
        
        # Open coverage report in default browser
        if (Test-Path "coverage.html") {
            Write-Output "Opening coverage report in browser..."
            Start-Process "coverage.html"
        }
    }
    catch {
        Write-ColorOutput Red "✗ Tests failed!"
        Write-Error $_
        exit 1
    }
}

function Clean-Artifacts {
    Write-ColorOutput Cyan "Cleaning..."
    
    @('bin', 'dist') | ForEach-Object {
        if (Test-Path $_) {
            Remove-Item -Recurse -Force $_
            Write-Output "Removed: $_"
        }
    }
    
    @('coverage.txt', 'coverage.html', '*.out') | ForEach-Object {
        Get-ChildItem -Filter $_ | Remove-Item -Force
    }
    
    go clean -testcache
    Write-ColorOutput Green "✓ Clean complete!"
}

function Run-Service {
    Write-ColorOutput Cyan "Running $ServiceName..."
    Write-Output "Press Ctrl+C to stop the service"
    Write-Output ""
    go run .\cmd\main.go
}

function Update-Dependencies {
    Write-ColorOutput Cyan "Downloading dependencies..."
    go mod download
    go mod tidy
    Write-ColorOutput Green "✓ Dependencies updated!"
}

function Format-Code {
    Write-ColorOutput Cyan "Formatting code..."
    go fmt .\...
    gofmt -s -w .
    Write-ColorOutput Green "✓ Code formatted!"
}

function Run-Vet {
    Write-ColorOutput Cyan "Running go vet..."
    go vet .\...
}

function Run-Lint {
    Write-ColorOutput Cyan "Running linters..."
    
    # Check if golangci-lint is installed
    $golangciLint = Get-Command golangci-lint -ErrorAction SilentlyContinue
    if (-not $golangciLint) {
        Write-ColorOutput Yellow "⚠ golangci-lint not found. Installing..."
        Install-Tools
    }
    
    golangci-lint run .\...
}

function Install-Tools {
    Write-ColorOutput Cyan "Installing development tools..."
    
    Write-Output "Installing golangci-lint..."
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
    
    if (Test-Path "proto") {
        Write-Output "Installing protobuf tools..."
        go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
        go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
    }
    
    Write-ColorOutput Green "✓ Tools installed!"
    Write-Output ""
    Write-Output "Make sure the following directory is in your PATH:"
    Write-Output "  $(go env GOPATH)\bin"
}

function Build-Docker {
    Write-ColorOutput Cyan "Building Docker image..."
    
    # Check if Docker is available
    $docker = Get-Command docker -ErrorAction SilentlyContinue
    if (-not $docker) {
        Write-ColorOutput Red "✗ Docker not found. Please install Docker Desktop for Windows."
        exit 1
    }
    
    try {
        # Get version from git
        $version = git describe --tags --always --dirty 2>$null
        if (-not $version) {
            $version = "latest"
        }
        
        docker build -t "$DockerRegistry/${ServiceName}:$version" .
        docker tag "$DockerRegistry/${ServiceName}:$version" "$DockerRegistry/${ServiceName}:latest"
        
        Write-ColorOutput Green "✓ Docker image built: $DockerRegistry/${ServiceName}:$version"
    }
    catch {
        Write-ColorOutput Red "✗ Docker build failed!"
        Write-Error $_
        exit 1
    }
}

# Main execution
switch ($Command) {
    'build' { Build-Service }
    'test' { Test-Service }
    'test-coverage' { Test-Coverage }
    'clean' { Clean-Artifacts }
    'run' { Run-Service }
    'deps' { Update-Dependencies }
    'fmt' { Format-Code }
    'vet' { Run-Vet }
    'lint' { Run-Lint }
    'install-tools' { Install-Tools }
    'docker-build' { Build-Docker }
    'help' { Show-Help }
    default { Show-Help }
}
