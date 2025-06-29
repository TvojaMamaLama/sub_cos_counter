#!/bin/bash

# Quick Docker and Docker Compose installation script
# Run this on the server: curl -sSL https://raw.githubusercontent.com/[user]/[repo]/main/deployment/scripts/install-docker.sh | bash

set -e

echo "🐳 Installing Docker and Docker Compose..."

# Update system
echo "📦 Updating system packages..."
apt-get update

# Install Docker
if ! command -v docker &> /dev/null; then
    echo "🐳 Installing Docker..."
    
    # Install prerequisites
    apt-get install -y ca-certificates curl
    
    # Add Docker's official GPG key
    install -m 0755 -d /etc/apt/keyrings
    curl -fsSL https://download.docker.com/linux/ubuntu/gpg -o /etc/apt/keyrings/docker.asc
    chmod a+r /etc/apt/keyrings/docker.asc
    
    # Add the repository to Apt sources
    echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/ubuntu $(. /etc/os-release && echo "$VERSION_CODENAME") stable" | tee /etc/apt/sources.list.d/docker.list > /dev/null
    
    # Update package index
    apt-get update
    
    # Install Docker
    apt-get install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin
    
    # Start and enable Docker
    systemctl start docker
    systemctl enable docker
    
    echo "✅ Docker installed successfully"
else
    echo "✅ Docker is already installed"
fi

# Install Docker Compose (standalone)
if [ ! -f /usr/local/bin/docker-compose ]; then
    echo "🔧 Installing Docker Compose..."
    
    # Download latest Docker Compose
    DOCKER_COMPOSE_VERSION=$(curl -s https://api.github.com/repos/docker/compose/releases/latest | grep 'tag_name' | cut -d\" -f4)
    curl -L "https://github.com/docker/compose/releases/download/${DOCKER_COMPOSE_VERSION}/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
    
    # Make it executable
    chmod +x /usr/local/bin/docker-compose
    
    # Create symlink
    ln -sf /usr/local/bin/docker-compose /usr/bin/docker-compose
    
    echo "✅ Docker Compose installed successfully"
else
    echo "✅ Docker Compose is already installed"
fi

# Verify installations
echo ""
echo "🔍 Verifying installations..."
echo "Docker version:"
docker --version

echo "Docker Compose version:"
/usr/local/bin/docker-compose --version

echo "Docker service status:"
systemctl is-active docker && echo "✅ Docker service is running" || echo "❌ Docker service is not running"

echo ""
echo "🎉 Installation completed!"
echo ""
echo "Next steps:"
echo "1. Clone your repository to /opt/subscription-bot"
echo "2. Create .env file with your configuration"
echo "3. Run: docker-compose -f deployment/docker/docker-compose.yml up -d"