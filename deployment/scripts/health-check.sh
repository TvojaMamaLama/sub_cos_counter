#!/bin/bash

# Health Check Script for Subscription Bot
# Run this script on the production server to verify everything is working

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
APP_DIR="/opt/subscription-bot"
COMPOSE_FILE="deployment/docker/docker-compose.yml"
DOCKER_COMPOSE="docker-compose"

# Check if running from correct directory
if [ ! -d "$APP_DIR" ]; then
    echo -e "${RED}❌ App directory not found: $APP_DIR${NC}"
    exit 1
fi

cd "$APP_DIR"

echo -e "${BLUE}🔍 Starting health check for Subscription Bot...${NC}\n"

# Function to print status
print_status() {
    if [ $1 -eq 0 ]; then
        echo -e "${GREEN}✅ $2${NC}"
    else
        echo -e "${RED}❌ $2${NC}"
        return 1
    fi
}

# 1. Check Docker is running
echo -e "${BLUE}1. Checking Docker...${NC}"
if systemctl is-active --quiet docker; then
    print_status 0 "Docker service is running"
    docker --version
else
    print_status 1 "Docker service is not running"
    exit 1
fi

# 2. Check Docker Compose
echo -e "\n${BLUE}2. Checking Docker Compose...${NC}"
if command -v docker-compose &> /dev/null; then
    print_status 0 "Docker Compose found"
    docker-compose --version
else
    print_status 1 "Docker Compose not found in PATH"
    exit 1
fi

# 3. Check containers status
echo -e "\n${BLUE}3. Checking containers...${NC}"
CONTAINERS_STATUS=$($DOCKER_COMPOSE -f $COMPOSE_FILE ps --format "table" 2>/dev/null || echo "ERROR")

if [[ "$CONTAINERS_STATUS" == "ERROR" ]]; then
    print_status 1 "Failed to get container status"
else
    echo "$CONTAINERS_STATUS"
    
    # Check if containers are running
    RUNNING_CONTAINERS=$($DOCKER_COMPOSE -f $COMPOSE_FILE ps --services --filter "status=running" 2>/dev/null | wc -l)
    EXPECTED_CONTAINERS=2  # bot + postgres
    
    if [ "$RUNNING_CONTAINERS" -eq "$EXPECTED_CONTAINERS" ]; then
        print_status 0 "All containers are running ($RUNNING_CONTAINERS/$EXPECTED_CONTAINERS)"
    else
        print_status 1 "Not all containers are running ($RUNNING_CONTAINERS/$EXPECTED_CONTAINERS)"
    fi
fi

# 4. Check database connection
echo -e "\n${BLUE}4. Checking database connection...${NC}"
DB_CHECK=$($DOCKER_COMPOSE -f $COMPOSE_FILE exec -T postgres pg_isready -U postgres -d sub_cos_counter 2>/dev/null && echo "OK" || echo "FAIL")

if [[ "$DB_CHECK" == *"OK"* ]]; then
    print_status 0 "Database is accessible"
    
    # Check tables exist
    TABLES_COUNT=$($DOCKER_COMPOSE -f $COMPOSE_FILE exec -T postgres psql -U postgres -d sub_cos_counter -t -c "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'public';" 2>/dev/null | xargs)
    
    if [ "$TABLES_COUNT" -gt 0 ]; then
        print_status 0 "Database tables found ($TABLES_COUNT tables)"
    else
        print_status 1 "No tables found in database"
    fi
else
    print_status 1 "Database connection failed"
fi

# 5. Check bot process
echo -e "\n${BLUE}5. Checking bot process...${NC}"
BOT_PROCESS=$($DOCKER_COMPOSE -f $COMPOSE_FILE exec -T bot pgrep bot 2>/dev/null && echo "RUNNING" || echo "NOT_RUNNING")

if [[ "$BOT_PROCESS" == "RUNNING" ]]; then
    print_status 0 "Bot process is running"
else
    print_status 1 "Bot process is not running"
fi

# 6. Check recent logs for errors
echo -e "\n${BLUE}6. Checking recent logs...${NC}"
echo "Recent bot logs (last 10 lines):"
$DOCKER_COMPOSE -f $COMPOSE_FILE logs --tail=10 bot 2>/dev/null || echo "Failed to get logs"

# Check for errors in logs
ERROR_COUNT=$($DOCKER_COMPOSE -f $COMPOSE_FILE logs bot 2>/dev/null | grep -i "error\|fail\|exception" | wc -l)
if [ "$ERROR_COUNT" -eq 0 ]; then
    print_status 0 "No errors found in bot logs"
else
    print_status 1 "Found $ERROR_COUNT error(s) in bot logs"
    echo "Recent errors:"
    $DOCKER_COMPOSE -f $COMPOSE_FILE logs bot 2>/dev/null | grep -i "error\|fail\|exception" | tail -5
fi

# 7. Check system resources
echo -e "\n${BLUE}7. Checking system resources...${NC}"

# Memory
MEMORY_USAGE=$(free | grep '^Mem:' | awk '{printf "%.1f", $3/$2 * 100.0}')
echo "Memory usage: ${MEMORY_USAGE}%"

# Disk space
DISK_USAGE=$(df /opt | tail -1 | awk '{print $5}' | sed 's/%//')
echo "Disk usage: ${DISK_USAGE}%"

if [ "$DISK_USAGE" -lt 90 ]; then
    print_status 0 "Disk usage is normal (${DISK_USAGE}%)"
else
    print_status 1 "Disk usage is high (${DISK_USAGE}%)"
fi

# 8. Check environment file
echo -e "\n${BLUE}8. Checking configuration...${NC}"
if [ -f ".env" ]; then
    print_status 0 ".env file exists"
    
    # Check if bot token is set
    if grep -q "TELEGRAM_BOT_TOKEN=" .env && [ -n "$(grep "TELEGRAM_BOT_TOKEN=" .env | cut -d'=' -f2)" ]; then
        print_status 0 "Telegram bot token is configured"
    else
        print_status 1 "Telegram bot token is not configured"
    fi
else
    print_status 1 ".env file not found"
fi

# 9. Test Telegram bot (basic connectivity)
echo -e "\n${BLUE}9. Testing Telegram bot connectivity...${NC}"
BOT_TOKEN=$(grep "TELEGRAM_BOT_TOKEN=" .env 2>/dev/null | cut -d'=' -f2 | tr -d ' "' || echo "")

if [ -n "$BOT_TOKEN" ]; then
    # Test bot API
    API_RESPONSE=$(curl -s "https://api.telegram.org/bot${BOT_TOKEN}/getMe" || echo "FAILED")
    
    if [[ "$API_RESPONSE" == *'"ok":true'* ]]; then
        BOT_USERNAME=$(echo "$API_RESPONSE" | grep -o '"username":"[^"]*"' | cut -d'"' -f4)
        print_status 0 "Telegram bot API is accessible (@${BOT_USERNAME})"
    else
        print_status 1 "Telegram bot API is not accessible"
    fi
else
    print_status 1 "Cannot test Telegram API - bot token not found"
fi

# Summary
echo -e "\n${BLUE}========================${NC}"
echo -e "${BLUE}  HEALTH CHECK SUMMARY  ${NC}"
echo -e "${BLUE}========================${NC}"

# Count passed/failed checks (approximate)
echo -e "\n📊 ${YELLOW}Quick commands for monitoring:${NC}"
echo "• View logs: $DOCKER_COMPOSE -f $COMPOSE_FILE logs -f bot"
echo "• Container status: $DOCKER_COMPOSE -f $COMPOSE_FILE ps"
echo "• Restart bot: $DOCKER_COMPOSE -f $COMPOSE_FILE restart bot"
echo "• Database console: $DOCKER_COMPOSE -f $COMPOSE_FILE exec postgres psql -U postgres -d sub_cos_counter"

echo -e "\n🤖 ${YELLOW}Test your bot:${NC}"
if [ -n "$BOT_USERNAME" ]; then
    echo "• Find your bot: @${BOT_USERNAME}"
else
    echo "• Find your bot using token: ${BOT_TOKEN:0:10}..."
fi
echo "• Send /start command"
echo "• Try adding a subscription"

echo -e "\n${GREEN}✅ Health check completed!${NC}"