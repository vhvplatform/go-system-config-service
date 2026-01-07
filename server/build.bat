@echo off
REM Build script for Windows
REM Usage: build.bat [command]
REM Commands: build, test, clean, run, deps, fmt, vet, lint, help

setlocal enabledelayedexpansion

set SERVICE_NAME=system-config-service
set GO_VERSION=1.25.5

if "%1"=="" goto help
if /I "%1"=="build" goto build
if /I "%1"=="test" goto test
if /I "%1"=="test-coverage" goto test-coverage
if /I "%1"=="clean" goto clean
if /I "%1"=="run" goto run
if /I "%1"=="deps" goto deps
if /I "%1"=="fmt" goto fmt
if /I "%1"=="vet" goto vet
if /I "%1"=="lint" goto lint
if /I "%1"=="install-tools" goto install-tools
if /I "%1"=="help" goto help
goto help

:help
echo Available commands:
echo   build           - Build the service
echo   test            - Run tests
echo   test-coverage   - Run tests with coverage
echo   clean           - Clean build artifacts
echo   run             - Run the service locally
echo   deps            - Download dependencies
echo   fmt             - Format code
echo   vet             - Run go vet
echo   lint            - Run linters (requires golangci-lint)
echo   install-tools   - Install development tools
echo   help            - Display this help screen
goto end

:build
echo Building %SERVICE_NAME%...
if not exist bin mkdir bin
go build -o bin\%SERVICE_NAME%.exe .\cmd\main.go
if %ERRORLEVEL% EQU 0 (
    echo Build complete! Binary: bin\%SERVICE_NAME%.exe
) else (
    echo Build failed!
    exit /b 1
)
goto end

:test
echo Running tests...
go test -v -race .\...
goto end

:test-coverage
echo Running tests with coverage...
go test -v -race -coverprofile=coverage.txt -covermode=atomic .\...
if %ERRORLEVEL% EQU 0 (
    go tool cover -html=coverage.txt -o coverage.html
    echo Coverage report generated: coverage.html
)
goto end

:clean
echo Cleaning...
if exist bin rmdir /s /q bin
if exist dist rmdir /s /q dist
if exist coverage.txt del /f /q coverage.txt
if exist coverage.html del /f /q coverage.html
if exist *.out del /f /q *.out
go clean -testcache
echo Clean complete!
goto end

:run
echo Running %SERVICE_NAME%...
go run .\cmd\main.go
goto end

:deps
echo Downloading dependencies...
go mod download
go mod tidy
echo Dependencies updated!
goto end

:fmt
echo Formatting code...
go fmt .\...
gofmt -s -w .
echo Code formatted!
goto end

:vet
echo Running go vet...
go vet .\...
goto end

:lint
echo Running linters...
golangci-lint run .\...
goto end

:install-tools
echo Installing development tools...
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
if exist proto (
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
)
echo Tools installed!
goto end

:end
endlocal
