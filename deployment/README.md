# Deployment

This directory contains all deployment-related files for the Subscription Tracker Bot.

## Structure

```
deployment/
├── docker/                 # Docker configuration files
│   ├── Dockerfile          # Main application Dockerfile
│   ├── .dockerignore       # Docker build context ignore rules
│   ├── docker-compose.yml  # Base Docker Compose configuration
│   ├── docker-compose.override.yml  # Development overrides
│   ├── docker-compose.prod.yml      # Production overrides
│   └── atlas.docker.hcl    # Atlas configuration for Docker
└── scripts/                # Deployment scripts
    └── deploy.sh           # Production deployment script
```

## Usage

All deployment commands should be run from the **project root directory**, not from this deployment directory.

### Quick Start
```bash
# From project root
make setup
make up
make migrate
```

### Development
```bash
make up-dev     # Start in development mode
make logs-bot   # Monitor bot logs
```

### Production
```bash
make deploy-prod                    # Full production deployment
# or
./deployment/scripts/deploy.sh     # Manual deployment script
```

## Docker Compose Files

### Base Configuration (`docker-compose.yml`)
- PostgreSQL database with health checks
- Bot application with auto-restart
- Atlas migration service (on-demand)
- Shared network and volumes

### Development Overrides (`docker-compose.override.yml`)
- Debug logging enabled
- Source code mounting for development
- Different database name to avoid conflicts
- Development-friendly settings

### Production Overrides (`docker-compose.prod.yml`)
- Resource limits and reservations
- Enhanced health checks
- Production logging configuration
- Optimized restart policies

## Environment Variables

All configuration is managed through the `.env` file in the project root. 

Example setup:
```bash
# Copy template
cp configs/examples/env.docker.example .env

# Edit with your settings
vim .env

# Key variables to set:
TELEGRAM_BOT_TOKEN=your_bot_token
DATABASE_PASSWORD=secure_password
TELEGRAM_ALLOWED_USER=your_telegram_id  # Optional
```

## Scripts

### `deploy.sh`
Production deployment script with:
- Requirements checking
- Automatic database backup
- Zero-downtime deployment
- Health checks
- Rollback capabilities

Usage:
```bash
./deployment/scripts/deploy.sh [check|backup|deploy|health]
```

## Monitoring

### Health Checks
```bash
make health     # Check service health
make logs       # View all logs
make ps         # Show container status
```

### Database Management
```bash
make backup-db  # Create database backup
make shell-db   # Connect to PostgreSQL
```

## Security Notes

- Never commit `.env` files
- Use strong passwords for production
- Consider using Docker secrets for sensitive data
- Regularly update base images
- Monitor container resource usage

## Troubleshooting

### Container won't start
```bash
make logs-bot   # Check bot logs
make logs-db    # Check database logs
make health     # Overall health check
```

### Database issues
```bash
make shell-db   # Connect to database
make migrate    # Run migrations
```

### Complete reset
```bash
make clean-all  # Remove everything including volumes
make setup      # Reinitialize
```