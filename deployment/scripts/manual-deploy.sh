#!/bin/bash

# Manual deployment script for Subscription Bot
# Run this script on the production server to deploy manually

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}🚀 Manual Deployment of Subscription Bot${NC}"
echo ""

# Configuration
APP_DIR="/opt/subscription-bot"
REPO_URL="https://github.com/[YOUR_USERNAME]/sub_cos_counter.git"  # Replace with actual repo URL
DOCKER_COMPOSE="docker-compose"

# Check if running as root
if [[ $EUID -ne 0 ]]; then
   echo -e "${RED}This script must be run as root${NC}"
   exit 1
fi

# Step 1: Install Docker and Docker Compose via apt (simple!)
echo -e "${BLUE}1. Installing Docker and Docker Compose...${NC}"
apt-get update
apt-get install -y docker.io docker-compose
systemctl start docker
systemctl enable docker
echo -e "${GREEN}✅ Docker and Docker Compose installed via apt${NC}"

# Step 2: Create app directory
echo -e "\n${BLUE}2. Setting up application directory...${NC}"
mkdir -p "$APP_DIR"
cd "$APP_DIR"

# Step 3: Get the code (if not already present)
if [ ! -d ".git" ]; then
    echo -e "${YELLOW}Cloning repository... (you may need to provide credentials)${NC}"
    echo "If repository is private, you can manually copy files instead"
    # git clone "$REPO_URL" . || echo "Failed to clone. Please copy files manually."
    echo "Please copy your project files to $APP_DIR manually"
    echo "Or clone with: git clone [YOUR_REPO_URL] ."
    read -p "Press Enter when files are ready..."
fi

# Step 4: Create .env file
echo -e "\n${BLUE}3. Setting up environment configuration...${NC}"
if [ ! -f ".env" ]; then
    echo "Creating .env file..."
    cat > .env << 'EOF'
# Production Environment Variables
TELEGRAM_BOT_TOKEN=7860783058:AAF8j8NdOPSeHuHLK2pSLs54iG-G52vfmJE
TELEGRAM_ALLOWED_USER=0

# Database settings
DATABASE_NAME=sub_cos_counter
DATABASE_USER=postgres
DATABASE_PASSWORD=postgres
DATABASE_HOST=postgres
DATABASE_PORT=5432

# Application settings
APP_NAME=Subscription Tracker Bot
APP_VERSION=1.0.0
APP_ENVIRONMENT=production
APP_DEBUG=false

# Logging settings
LOGGING_LEVEL=info
LOGGING_FORMAT=json
LOGGING_OUTPUT=stdout

# Connection pool settings
DATABASE_MAX_CONNECTIONS=25
DATABASE_MAX_IDLE_TIME=15
DATABASE_CONN_MAX_LIFETIME=60
EOF
    echo -e "${GREEN}✅ .env file created${NC}"
    echo -e "${YELLOW}⚠️  Please review and update .env file if needed${NC}"
else
    echo -e "${GREEN}✅ .env file already exists${NC}"
fi

# Step 5: Stop any existing containers
echo -e "\n${BLUE}4. Stopping existing containers...${NC}"
$DOCKER_COMPOSE -f deployment/docker/docker-compose.yml down || echo "No containers to stop"

# Step 6: Build and start containers
echo -e "\n${BLUE}5. Building and starting containers...${NC}"
$DOCKER_COMPOSE -f deployment/docker/docker-compose.yml build --no-cache
$DOCKER_COMPOSE -f deployment/docker/docker-compose.yml up -d

# Step 7: Wait and check status
echo -e "\n${BLUE}6. Waiting for services to start...${NC}"
sleep 15

echo -e "\n${BLUE}7. Checking deployment status...${NC}"
echo "Container status:"
$DOCKER_COMPOSE -f deployment/docker/docker-compose.yml ps

echo -e "\nRecent logs:"
$DOCKER_COMPOSE -f deployment/docker/docker-compose.yml logs --tail=10

# Step 8: Final verification
echo -e "\n${BLUE}8. Final verification...${NC}"
if $DOCKER_COMPOSE -f deployment/docker/docker-compose.yml ps | grep -q "Up"; then
    echo -e "${GREEN}🎉 Deployment successful!${NC}"
    echo ""
    echo "Your bot should now be running. Test it in Telegram:"
    echo "• Find your bot with token: 7860783058:AAF8j8NdOPSeHuHLK2pSLs54iG-G52vfmJE"
    echo "• Send /start command"
    echo ""
    echo "Useful commands:"
    echo "• View logs: $DOCKER_COMPOSE -f deployment/docker/docker-compose.yml logs -f bot"
    echo "• Restart: $DOCKER_COMPOSE -f deployment/docker/docker-compose.yml restart"
    echo "• Stop: $DOCKER_COMPOSE -f deployment/docker/docker-compose.yml down"
else
    echo -e "${RED}❌ Deployment may have issues. Check logs above.${NC}"
    exit 1
fi