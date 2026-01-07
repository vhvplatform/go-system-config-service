# System Config Service - Monorepo

Đây là kho lưu trữ tổng hợp (monorepo) cho hệ thống System Config Service, bao gồm tất cả các thành phần của ứng dụng.

This is the monorepo for the System Config Service system, containing all application components.

## Cấu trúc dự án / Project Structure

```
.
├── server/          # Backend microservice (Golang)
├── client/          # Frontend application (ReactJS)
├── flutter/         # Mobile application (Flutter)
├── docs/            # Tài liệu dự án / Project documentation
├── .github/         # GitHub workflows and configurations
└── README.md        # Tài liệu này / This file
```

## Các thành phần / Components

### 1. Server (Backend - Golang)

Backend microservice được xây dựng bằng Golang, cung cấp REST API và gRPC endpoints cho việc quản lý cấu hình hệ thống.

Backend microservice built with Golang, providing REST API and gRPC endpoints for system configuration management.

📁 **Thư mục**: [`server/`](./server/)
📖 **Chi tiết**: Xem [server/README.md](./server/README.md)

**Công nghệ chính / Key Technologies:**
- Go 1.25+
- Gin Framework
- MongoDB
- Redis
- gRPC

### 2. Client (Frontend - ReactJS)

Ứng dụng web frontend được xây dựng bằng ReactJS, cung cấp giao diện người dùng để quản lý cấu hình.

Frontend web application built with ReactJS, providing user interface for configuration management.

📁 **Thư mục**: [`client/`](./client/)
📖 **Chi tiết**: Xem [client/README.md](./client/README.md)

**Công nghệ chính / Key Technologies:**
- React 18+
- Redux / Context API
- Material-UI / Ant Design
- Axios

**Trạng thái**: 🚧 *Đang phát triển / In Development*

### 3. Flutter (Mobile App)

Ứng dụng di động đa nền tảng được xây dựng bằng Flutter cho Android và iOS.

Cross-platform mobile application built with Flutter for Android and iOS.

📁 **Thư mục**: [`flutter/`](./flutter/)
📖 **Chi tiết**: Xem [flutter/README.md](./flutter/README.md)

**Công nghệ chính / Key Technologies:**
- Flutter 3.x+
- Dart 3.x+
- Provider / Riverpod

**Trạng thái**: 🚧 *Đang phát triển / In Development*

### 4. Documentation

Tài liệu dự án, bao gồm tài liệu kỹ thuật, hướng dẫn sử dụng, và các sơ đồ kiến trúc.

Project documentation, including technical docs, user guides, and architecture diagrams.

📁 **Thư mục**: [`docs/`](./docs/)

## Bắt đầu / Getting Started

### Yêu cầu hệ thống / Prerequisites

- **Server**: Go 1.25+, MongoDB, Redis
- **Client**: Node.js 18+, npm/yarn
- **Flutter**: Flutter SDK 3.x+, Dart 3.x+

### Cài đặt / Installation

Mỗi thành phần có hướng dẫn cài đặt riêng trong thư mục tương ứng:

Each component has its own installation instructions in the respective directory:

1. **Server**: Xem [server/README.md](./server/README.md)
2. **Client**: Xem [client/README.md](./client/README.md)
3. **Flutter**: Xem [flutter/README.md](./flutter/README.md)

## Phát triển / Development

### Cấu trúc Monorepo

Dự án này sử dụng cấu trúc monorepo để quản lý tất cả các thành phần trong một kho lưu trữ duy nhất. Mỗi thành phần có thể phát triển, build và deploy độc lập.

This project uses a monorepo structure to manage all components in a single repository. Each component can be developed, built, and deployed independently.

### Quy trình làm việc / Workflow

1. Clone repository
2. Chọn thành phần bạn muốn làm việc / Choose the component you want to work on
3. Làm theo hướng dẫn trong README của từng thành phần / Follow the instructions in each component's README
4. Commit và push code

### Git Branches

- `main` hoặc `master`: Nhánh chính, code production
- `develop`: Nhánh phát triển
- `feature/*`: Nhánh tính năng mới
- `bugfix/*`: Nhánh sửa lỗi
- `hotfix/*`: Nhánh sửa lỗi khẩn cấp

## Kiểm tra / Testing

Mỗi thành phần có quy trình kiểm tra riêng. Xem README của từng thành phần để biết chi tiết.

Each component has its own testing process. See each component's README for details.

## Triển khai / Deployment

Thông tin về quy trình triển khai sẽ được cập nhật trong tài liệu riêng.

Deployment process information will be updated in separate documentation.

## Đóng góp / Contributing

Vui lòng đọc [CONTRIBUTING.md](./server/CONTRIBUTING.md) để biết chi tiết về quy trình đóng góp code.

Please read [CONTRIBUTING.md](./server/CONTRIBUTING.md) for details on our code contribution process.

## Changelog

Xem [CHANGELOG.md](./server/CHANGELOG.md) để biết lịch sử thay đổi của dự án.

See [CHANGELOG.md](./server/CHANGELOG.md) for project change history.

## Hỗ trợ / Support

- **Issues**: [GitHub Issues](https://github.com/vhvplatform/go-system-config-service/issues)
- **Discussions**: [GitHub Discussions](https://github.com/vhvplatform/go-system-config-service/discussions)
- **Email**: support@vhvplatform.com

## Giấy phép / License

MIT License - xem file [LICENSE](LICENSE) để biết chi tiết.

MIT License - see [LICENSE](LICENSE) file for details.

---

**Maintained by**: VHV Corp Development Team  
**Last Updated**: 2026-01-07
