#!/bin/bash

# Subscription Tracker Bot - Production Deployment Script
set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
PROJECT_NAME="subscription-tracker-bot"
APP_DIR="/opt/subscription-bot"
DEPLOYMENT_DIR="deployment/docker"
COMPOSE_FILE="$DEPLOYMENT_DIR/docker-compose.yml"
COMPOSE_PROD="$DEPLOYMENT_DIR/docker-compose.prod.yml"
BACKUP_DIR="/opt/backups/subscription-bot"
LOG_FILE="/var/log/subscription-bot-deploy.log"

# Functions
log() {
    echo -e "${BLUE}[$(date +'%Y-%m-%d %H:%M:%S')]${NC} $1" | tee -a $LOG_FILE
}

success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1" | tee -a $LOG_FILE
}

warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1" | tee -a $LOG_FILE
}

error() {
    echo -e "${RED}[ERROR]${NC} $1" | tee -a $LOG_FILE
    exit 1
}

# Pre-deployment checks
check_requirements() {
    log "Checking requirements..."
    
    command -v docker >/dev/null 2>&1 || error "Docker is required but not installed"
    command -v docker-compose >/dev/null 2>&1 || error "Docker Compose is required but not installed"
    
    if [ ! -f .env ]; then
        error ".env file not found. Run 'make setup' first."
    fi
    
    if ! grep -q "TELEGRAM_BOT_TOKEN=" .env; then
        error "TELEGRAM_BOT_TOKEN not set in .env file"
    fi
    
    success "Requirements check passed"
}

# Backup database
backup_database() {
    log "Creating database backup..."
    
    mkdir -p $BACKUP_DIR
    
    if docker-compose -f $COMPOSE_FILE ps postgres | grep -q "Up"; then
        BACKUP_FILE="$BACKUP_DIR/backup_$(date +%Y%m%d_%H%M%S).sql"
        docker-compose -f $COMPOSE_FILE exec -T postgres pg_dump -U postgres sub_cos_counter > $BACKUP_FILE
        success "Database backup created: $BACKUP_FILE"
    else
        warning "Database not running, skipping backup"
    fi
}

# Deploy application
deploy() {
    log "Starting deployment..."
    
    # Pull latest code (if in git repo)
    if [ -d .git ]; then
        log "Pulling latest code..."
        git pull origin main || warning "Failed to pull latest code"
    fi
    
    # Build new images
    log "Building Docker images..."
    docker-compose -f $COMPOSE_FILE -f $COMPOSE_PROD build --no-cache
    
    # Stop existing services
    log "Stopping existing services..."
    docker-compose -f $COMPOSE_FILE -f $COMPOSE_PROD down
    
    # Start services
    log "Starting services..."
    docker-compose -f $COMPOSE_FILE -f $COMPOSE_PROD up -d
    
    # Run migrations
    log "Running database migrations..."
    sleep 10  # Wait for database to be ready
    docker-compose -f $COMPOSE_FILE --profile migration run --rm atlas || warning "Migration failed or no changes"
    
    success "Deployment completed successfully"
}

# Health check
health_check() {
    log "Performing health check..."
    
    sleep 15  # Wait for services to start
    
    # Check if containers are running
    if ! docker-compose -f $COMPOSE_FILE ps | grep -q "Up"; then
        error "Some services are not running"
    fi
    
    # Check database connection
    if ! docker-compose -f $COMPOSE_FILE exec -T postgres pg_isready -U postgres -d sub_cos_counter >/dev/null 2>&1; then
        error "Database health check failed"
    fi
    
    # Check bot process
    if ! docker-compose -f $COMPOSE_FILE exec -T bot pgrep bot >/dev/null 2>&1; then
        error "Bot process not running"
    fi
    
    success "Health check passed"
}

# Cleanup old images
cleanup() {
    log "Cleaning up old Docker images..."
    docker image prune -f
    docker system prune -f
    success "Cleanup completed"
}

# Main deployment flow
main() {
    log "Starting deployment of $PROJECT_NAME"
    
    check_requirements
    backup_database
    deploy
    health_check
    cleanup
    
    success "🎉 Deployment completed successfully!"
    log "Check logs with: make logs"
    log "Monitor with: make health"
}

# Handle script arguments
case "${1:-deploy}" in
    "check")
        check_requirements
        ;;
    "backup")
        backup_database
        ;;
    "deploy")
        main
        ;;
    "health")
        health_check
        ;;
    *)
        echo "Usage: $0 {check|backup|deploy|health}"
        echo "  check  - Check deployment requirements"
        echo "  backup - Backup database only"
        echo "  deploy - Full deployment (default)"
        echo "  health - Health check only"
        exit 1
        ;;
esac