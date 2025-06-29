    # CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Development Commands

### Quick Start (Docker - Recommended)
```bash
make setup      # Create .env from template
# Edit .env with TELEGRAM_BOT_TOKEN
make up         # Start PostgreSQL + Bot
make migrate    # Apply database migrations
make logs       # Check application logs
```

### Local Development
```bash
go mod download                    # Install dependencies
atlas migrate apply --env local   # Apply migrations locally
go run cmd/bot/main.go            # Run bot directly
go test ./...                     # Run all tests
```

### Build and Test
```bash
make build      # Build Docker images
make test       # Run Go tests
make docker-test # Run tests inside container
go build -o bin/bot cmd/bot/main.go  # Build binary locally
```

### Development Workflow
```bash
make dev-setup  # Setup development environment (.env with debug settings)
make up-dev     # Start in development mode (debug logs, source mounting)
make logs-bot   # Monitor bot logs only
make shell-bot  # Open shell in bot container
make shell-db   # Connect to PostgreSQL console
```

### Database Operations
```bash
make migrate           # Apply migrations via Docker
make migrate-status    # Check migration status
make backup-db         # Create database backup
atlas migrate diff --env local  # Generate new migration from schema changes
```

## Architecture Overview

### Clean Architecture Layers
The application follows clean architecture principles with clear separation of concerns:

- **cmd/bot/main.go**: Application entry point with dependency injection and graceful shutdown
- **internal/models/**: Domain entities (Subscription, Payment) with business rules and validation
- **internal/repository/**: Data access layer using pgx connection pooling for PostgreSQL
- **internal/services/**: Business logic layer that orchestrates domain operations
- **internal/bot/**: Telegram bot interface layer using Telebot v3 with state management

### Configuration System (Viper)
Multi-source configuration with priority: Environment Variables > YAML files > Defaults
- **internal/config/config.go**: Structured configuration with validation and type safety
- **configs/examples/**: Template configurations for different environments
- Supports both full DATABASE_URL and individual database connection parameters

### Database Architecture
- **PostgreSQL** with pgx v5 driver for performance and PostgreSQL-specific features
- **Atlas** for schema-as-code migrations instead of traditional SQL migration files
- **Connection pooling** with configurable parameters (MaxConnections, IdleTime, Lifetime)
- **migrations/schema.sql**: Single source of truth for database schema

### Telegram Bot Implementation
- **Button-driven interface**: No text commands, only inline keyboards for UX
- **State management**: Multi-step dialogs for subscription creation stored in memory
- **Personal bot mode**: Optional restriction to single Telegram user ID
- **Callback handlers**: Separate handlers for different button interactions and menu navigation

### Currency and Localization
- **Dual currency support**: USD and RUB with separate tracking and analytics
- **Category system**: Entertainment, Work, Education, Home, Other with emoji representations
- **Subscription periods**: Days, weeks, months, years with flexible custom periods

## Key Components Integration

### Subscription Lifecycle
1. **Creation**: Multi-step dialog (category → currency → period → auto-renewal → name → cost)
2. **Payment tracking**: Manual marking of payments with automatic next payment date calculation
3. **Analytics**: Monthly expenses and category breakdowns with currency separation
4. **Management**: View, pay, delete subscriptions through inline keyboards

### Docker Compose Architecture
- **postgres**: PostgreSQL 16 with health checks, persistent volumes, and initialization scripts
- **bot**: Go application with auto-restart, resource limits, and log management
- **atlas**: On-demand migration service using profiles
- **Networks**: Isolated network for service communication
- **Volumes**: Persistent PostgreSQL data and application logs

### Environment Management
- **Development**: `docker-compose.override.yml` with debug settings and source mounting
- **Production**: `docker-compose.prod.yml` with resource limits and optimized logging
- **Local**: Direct Go execution with local PostgreSQL and Atlas CLI

## File Organization

### Deployment Separation
All deployment-related files are organized in `deployment/` to keep the project root clean:
- **deployment/docker/**: All Docker and compose files
- **deployment/scripts/**: Production deployment automation with backup and health checks
- **configs/examples/**: Configuration templates and examples

### Critical Configuration Files
- **.env**: Runtime environment variables (created from templates, not committed)
- **atlas.hcl**: Local Atlas configuration for development
- **deployment/docker/atlas.docker.hcl**: Docker-specific Atlas configuration
- **Makefile**: Unified command interface that abstracts Docker Compose complexity

## Development Patterns

### Error Handling
- Repository layer: Wrap database errors with context using `fmt.Errorf("context: %w", err)`
- Service layer: Validate inputs and provide business-context error messages
- Bot handlers: Convert service errors to user-friendly Telegram messages

### Testing Strategy
- Unit tests for business logic in services layer
- Repository tests can use Docker PostgreSQL for integration testing
- Bot handlers are tested through mock services

### State Management
- User dialog states stored in memory (map[int64]*UserState)
- States include current step and collected data for multi-step flows
- State cleanup on completion or cancellation of dialogs

## Dependencies and Versions

- **Go 1.24+**: Required for modern features and performance
- **Telebot v3**: Modern Telegram bot framework with inline keyboard support
- **pgx v5**: PostgreSQL driver chosen over lib/pq for performance and feature set
- **Viper**: Configuration management with multiple source support
- **Atlas**: Modern database migration tool chosen over traditional tools like Goose