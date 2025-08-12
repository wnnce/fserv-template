# fserv-template

A modern, production-ready Go microservice template based on [Fiber V2](https://github.com/gofiber/fiber), supporting RESTful APIs, WebSocket, PostgreSQL, MongoDB, Redis, Kafka, and more.

This template is designed for rapid development of scalable backend services with robust configuration, logging, validation, and CI/CD support.

---

## Table of Contents
- [Features](#features)
- [Architecture](#architecture)
- [Getting Started](#getting-started)
- [Configuration](#configuration)
- [Usage](#usage)
- [Testing](#testing)
- [Deployment](#deployment)
- [Recommended Enhancements](#recommended-enhancements)
- [Contributing](#contributing)
- [License](#license)

---

## Features
- ⚡️ Fast HTTP server powered by [Fiber v2](https://github.com/gofiber/fiber)
- 🧩 Modular structure for business logic, middleware, and configuration
- 🗄️ Built-in support for PostgreSQL, MongoDB, Redis, and Kafka
- 🔒 Centralized error handling and request validation (with Chinese translation)
- 🧪 Comprehensive unit tests and CI/CD pipeline (GitHub Actions)
- 🐳 Docker & docker-compose ready for local and production deployment
- 🌐 RESTful API and WebSocket support
- 📦 Modern Go modules and dependency management
- 📝 Configurable logging with file rotation

## Architecture
```
.
├── main.go                # Application entrypoint
├── config/                # Configuration loading and management
├── biz/                   # Business logic (routes, handlers, DAL)
│   ├── dal/               # Data access layer (Postgres, Mongo, etc.)
│   ├── mw/                # Middleware (Kafka, Redis, etc.)
│   └── route/             # API and WebSocket route registration
├── internal/              # Internal utilities and middleware
├── logging/               # Logging setup (slog, lumberjack)
├── pkg/                   # Reusable packages (tool functions, SQL builder)
├── configs/               # Configuration files (YAML)
├── Dockerfile             # Docker build instructions
├── docker-compose.yaml    # Multi-service orchestration
└── ...
```

## Getting Started
### Prerequisites
- [Go 1.24+](https://golang.org/dl/)
- [Docker](https://www.docker.com/) (optional, for containerized deployment)
- PostgreSQL, MongoDB, Redis, Kafka (optional, for full feature set)

### Quick Start
1. Clone the repository
   ```bash
   git clone https://github.com/wnnce/fserv-template.git
   cd fserv-template
   ```
2. Copy and edit configuration
   ```bash
   cp config.yaml.example configs/config.yaml
   # Edit configs/config.yaml as needed
   ```
3. Run with Go
   ```bash
   go mod tidy
   go run main.go
   ```
   The server will start on the port specified in your config (default: `7000`).
4. Or run with Docker
   ```bash
   docker build -t fserv-template:latest .
   docker run -p 7000:7000 -v $(pwd)/logs:/app/logs -v $(pwd)/configs:/app/configs fserv-template:latest
   ```
5. Or use docker-compose
   ```bash
   docker-compose up -d
   ```

## Configuration
All configuration is managed via `configs/config.yaml`.
See `config.yaml.example` for all available options, including:
- Server host, port, environment
- Logger settings
- Database (PostgreSQL), MongoDB, Redis, Kafka connection info

Example:
```yaml
server:
  name: fserv-template
  environment: dev
  version: 1.0.0
  host: 0.0.0.0
  port: 7000
# ...
```

## Usage
- **REST API**: Register your routes in `biz/route/`.
- **WebSocket**: Example endpoint at `/ws/echo`.
- **Validation**: Uses [go-playground/validator](https://github.com/go-playground/validator) with Chinese translation.
- **Logging**: Configurable via YAML, supports file rotation.
- **Database**: PostgreSQL and MongoDB clients are initialized via config.
- **Error Handling**: Centralized error handler chain, customizable.
- **Testing**: Place your tests in `*_test.go` files. Run `go test ./... -v` for all tests.

## Testing
- Run all unit tests:
  ```bash
  go test ./... -v
  ```
- CI will also run tests and static analysis on every push (see `.github/workflows/main.yaml`).

## Deployment
- **Docker**: See [Dockerfile](Dockerfile)
- **docker-compose**: See [docker-compose.yaml](docker-compose.yaml)
- **Kubernetes**: You can adapt the Docker image for K8s deployment.

## Recommended Enhancements
To make this template even more powerful and production-ready, consider integrating:
- **API Documentation**: Integrate [swaggo/swag](https://github.com/swaggo/swag) or [go-fiber/swagger](https://github.com/gofiber/swagger) for automatic Swagger/OpenAPI docs.
- **Database Migration**: Add [golang-migrate/migrate](https://github.com/golang-migrate/migrate) or [pressly/goose](https://github.com/pressly/goose) for DB schema management.
- **Health Checks**: Add `/health` endpoint for readiness/liveness probes.
- **Authentication/Authorization**: Integrate JWT/OAuth2 for secure APIs.
- **Metrics/Monitoring**: Add Prometheus metrics endpoint for observability.
- **API Versioning**: Use route groups like `/api/v1/` for versioned APIs.

## Contributing
Contributions are welcome! Please open issues or submit pull requests.
For major changes, please open an issue first to discuss what you would like to change.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
