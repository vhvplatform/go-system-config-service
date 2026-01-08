# Hướng dẫn Checkout Repository / Repository Checkout Guide

## Thông tin về cấu trúc mới / Information about new structure

Repository đã được tái cấu trúc thành monorepo với 3 thư mục chính:
- `server/` - Backend code bằng Golang
- `client/` - Frontend code bằng ReactJS (đang phát triển)
- `flutter/` - Mobile app code bằng Flutter (đang phát triển)
- `docs/` - Tài liệu dự án

The repository has been restructured into a monorepo with 3 main directories:
- `server/` - Backend code in Golang
- `client/` - Frontend code in ReactJS (under development)
- `flutter/` - Mobile app code in Flutter (under development)
- `docs/` - Project documentation

---

## 1. Checkout Repository Hiện Tại (Existing Repository)

### Nếu bạn đã có repository trên máy / If you already have the repository cloned:

```bash
# Di chuyển vào thư mục repository
cd go-system-config-service

# Fetch tất cả các thay đổi từ remote
git fetch origin

# Checkout nhánh mới với cấu trúc đã được tái cấu trúc
git checkout copilot/refactor-repo-structure

# Pull các thay đổi mới nhất
git pull origin copilot/refactor-repo-structure
```

### Sau khi checkout:

```bash
# Kiểm tra cấu trúc thư mục mới
ls -la

# Bạn sẽ thấy:
# - client/   (thư mục ReactJS frontend)
# - server/   (thư mục Golang backend)
# - flutter/  (thư mục Flutter mobile app)
# - docs/     (thư mục tài liệu)
# - README.md (tài liệu monorepo)
```

---

## 2. Clone Repository Mới (New Clone)

### Clone toàn bộ repository với cấu trúc mới:

```bash
# Clone repository
git clone https://github.com/vhvplatform/go-system-config-service.git

# Di chuyển vào thư mục
cd go-system-config-service

# Checkout nhánh với cấu trúc mới
git checkout copilot/refactor-repo-structure
```

### Hoặc clone trực tiếp nhánh mới:

```bash
# Clone repository với nhánh cụ thể
git clone -b copilot/refactor-repo-structure https://github.com/vhvplatform/go-system-config-service.git

# Di chuyển vào thư mục
cd go-system-config-service
```

---

## 3. Làm việc với từng Component / Working with Each Component

### A. Server (Golang Backend)

```bash
# Di chuyển vào thư mục server
cd server

# Cài đặt dependencies
go mod download

# Build
make build
# hoặc trên Windows:
# build.bat build

# Chạy
make run
# hoặc trên Windows:
# build.bat run

# Test
make test
# hoặc:
# go test ./...
```

### B. Client (ReactJS Frontend) - Sẽ được phát triển

```bash
# Di chuyển vào thư mục client
cd client

# Sẽ cài đặt dependencies khi có code
# npm install
# hoặc
# yarn install

# Chạy development server
# npm start
# hoặc
# yarn start
```

### C. Flutter (Mobile App) - Sẽ được phát triển

```bash
# Di chuyển vào thư mục flutter
cd flutter

# Sẽ cài đặt dependencies khi có code
# flutter pub get

# Chạy trên emulator/device
# flutter run
```

---

## 4. Cấu trúc Thư mục / Directory Structure

```
go-system-config-service/
├── README.md                    # Tài liệu monorepo chính
├── .gitignore                   # Git ignore cho tất cả components
├── .github/                     # GitHub workflows
├── docs/                        # Tài liệu dự án
│   ├── diagrams/               # Sơ đồ kiến trúc
│   ├── examples/               # Ví dụ
│   └── *.md                    # Tài liệu khác
├── server/                      # Golang Backend
│   ├── cmd/                    # Entry points
│   ├── internal/               # Internal packages
│   │   ├── domain/            # Domain models
│   │   ├── handler/           # HTTP handlers
│   │   ├── repository/        # Data access layer
│   │   ├── router/            # Route definitions
│   │   └── service/           # Business logic
│   ├── migrations/            # Database migrations
│   ├── go.mod                 # Go dependencies
│   ├── Makefile               # Build commands
│   ├── Dockerfile             # Docker config
│   └── README.md              # Server documentation
├── client/                      # ReactJS Frontend (Đang phát triển)
│   └── README.md               # Client documentation
└── flutter/                     # Flutter Mobile App (Đang phát triển)
    └── README.md               # Flutter documentation
```

---

## 5. Thông tin Commit / Commit Information

**Branch mới**: `copilot/refactor-repo-structure`

**Commit mới nhất**: 
- Message: "Restructure repository into monorepo with server, client, and flutter directories"
- SHA: f57c19c

**Thay đổi chính**:
- ✅ Di chuyển tất cả Go code vào `server/`
- ✅ Tạo thư mục `client/` với README placeholder
- ✅ Tạo thư mục `flutter/` với README placeholder
- ✅ Giữ nguyên thư mục `docs/` ở root
- ✅ Tạo README.md mới ở root giải thích cấu trúc monorepo
- ✅ Cập nhật .gitignore cho cả 3 platforms

---

## 6. Migration Notes / Ghi chú về Migration

### ⚠️ Quan trọng / Important:

1. **Paths đã thay đổi**: 
   - Tất cả Go code giờ nằm trong `server/`
   - Import paths không thay đổi (vẫn là `github.com/vhvplatform/go-system-config-service`)

2. **Build commands**: 
   - Phải chạy từ thư mục `server/`: `cd server && make build`

3. **Docker**: 
   - Dockerfile đã di chuyển vào `server/`
   - Build command: `docker build -f server/Dockerfile -t system-config-service .`

4. **CI/CD**: 
   - Có thể cần cập nhật workflows trong `.github/` để trỏ đúng đường dẫn

---

## 7. Hỗ trợ / Support

Nếu gặp vấn đề, vui lòng:
- Xem tài liệu chi tiết trong mỗi thư mục component
- Tạo issue trên GitHub
- Liên hệ team qua email: support@vhvplatform.com

If you encounter issues, please:
- Check detailed documentation in each component directory
- Create an issue on GitHub
- Contact the team via email: support@vhvplatform.com

---

**Last Updated**: 2026-01-07  
**Maintained by**: VHV Corp Development Team
