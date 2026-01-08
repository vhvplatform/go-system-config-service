# Stage 1: Build
FROM golang:1.25.5-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk --no-cache add ca-certificates git

# 1. Copy file go.mod/sum của cả hai để cache layer (Tăng tốc build)
COPY go-shared/go.mod go-shared/go.sum ./go-shared/
COPY go-system-config-service/go.mod go-system-config-service/go.sum ./go-system-config-service/

# 2. Download dependencies (Go sẽ tự xử lý mối quan hệ giữa các module)
RUN cd go-system-config-service && go mod download

# 3. Copy toàn bộ mã nguồn cần thiết
COPY go-shared/ ./go-shared/
COPY go-system-config-service/ ./go-system-config-service/

# 4. Build ứng dụng
WORKDIR /app/go-system-config-service
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-w -s" -o /app/bin/system-config-service ./cmd/main.go

# Stage 2: Runtime
FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
RUN addgroup -S appgroup && adduser -S appuser -u 1000 -G appgroup
WORKDIR /app
RUN chown 1000:1000 /app
COPY --from=builder /app/bin/system-config-service .
# Expose ports
USER 1000
EXPOSE 8085 50055
CMD ["./system-config-service"]
