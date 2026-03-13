# go-starter-utils

A CLI toolkit for scaffolding Go microservices with hexagonal / clean architecture.

## Installation

### From source

```bash
go install github.com/rahulbalajee/go-starter-utils@latest
```

### Build locally

```bash
git clone https://github.com/rahulbalajee/go-starter-utils.git
cd go-starter-utils
make install
```

## Usage

### Create a service

```bash
go-starter-utils create service order
```

This generates a fully structured service at `services/order-service/`:

```
services/order-service/
├── cmd/
│   └── main.go              # Application entry point
├── internal/
│   ├── domain/
│   │   └── order.go         # Domain models and interfaces
│   ├── service/
│   │   └── service.go       # Business logic implementation
│   └── infrastructure/
│       ├── events/           # Event handling
│       ├── grpc/             # gRPC server handlers
│       ├── http/             # HTTP server handlers
│       └── repository/       # Data persistence
├── pkg/
│   └── types/
│       └── types.go          # Shared types and models
├── go.mod
└── README.md
```

### Custom output directory

```bash
go-starter-utils create service order --output ./my-monorepo
```

### Version

```bash
go-starter-utils --version
```

## Development

```bash
# Run tests
make test

# Run linter
make lint

# Build binary
make build

# Build with version info
make build VERSION=v1.0.0

# Format code
make fmt

# Generate test coverage report
make coverage
```

## Architecture

The scaffolded services follow Clean Architecture / Hexagonal Architecture:

- **Domain Layer** — pure business logic, interfaces, and models with zero external dependencies
- **Service Layer** — business logic implementation using domain interfaces
- **Infrastructure Layer** — concrete implementations for persistence, messaging, and transport
- **Public Types** — shared types that other services can import

## License

See [LICENSE](LICENSE) for details.
