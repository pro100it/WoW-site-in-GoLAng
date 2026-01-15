#!/bin/bash

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}========================================${NC}"
echo -e "${GREEN}WoW Registration Server Setup${NC}"
echo -e "${BLUE}========================================${NC}"

# Проверка зависимостей
check_dependency() {
    if ! command -v $1 &> /dev/null; then
        echo -e "${RED}Error: $1 is not installed${NC}"
        echo "Please install $1 before continuing."
        exit 1
    fi
}

echo -e "${YELLOW}Checking dependencies...${NC}"
check_dependency docker
check_dependency docker-compose
check_dependency git
check_dependency curl

# Создание необходимых директорий
echo -e "${YELLOW}Creating directories...${NC}"
mkdir -p logs
mkdir -p nginx/ssl
mkdir -p prometheus
mkdir -p grafana/provisioning
mkdir -p frontend/static
mkdir -p scripts/backup

# Копирование конфигурационных файлов
echo -e "${YELLOW}Copying configuration files...${NC}"
if [ ! -f .env ]; then
    cp .env.example .env
    echo -e "${GREEN}Created .env file from example${NC}"
    echo -e "${YELLOW}Please edit .env file with your settings${NC}"
fi

if [ ! -f nginx/nginx.conf ]; then
    cp nginx/nginx.conf.example nginx/nginx.conf
fi

if [ ! -f prometheus/prometheus.yml ]; then
    cp prometheus/prometheus.yml.example prometheus/prometheus.yml
fi

# Генерация SSL сертификатов (для разработки)
echo -e "${YELLOW}Generating SSL certificates...${NC}"
if [ ! -f nginx/ssl/server.crt ]; then
    openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
        -keyout nginx/ssl/server.key \
        -out nginx/ssl/server.crt \
        -subj "/C=US/ST=State/L=City/O=Company/CN=localhost"
    echo -e "${GREEN}Generated self-signed SSL certificate${NC}"
fi

# Установка прав
echo -e "${YELLOW}Setting permissions...${NC}"
chmod +x scripts/*.sh
chmod 644 .env
chmod 600 nginx/ssl/*

# Загрузка и сборка
echo -e "${YELLOW}Pulling Docker images...${NC}"
docker-compose pull

echo -e "${YELLOW}Building application...${NC}"
docker-compose build

echo -e "${YELLOW}Starting services...${NC}"
docker-compose up -d

# Ожидание запуска сервисов
echo -e "${YELLOW}Waiting for services to start...${NC}"
sleep 10

# Проверка статуса
echo -e "${YELLOW}Checking service status...${NC}"
if curl -s http://localhost:${PORT:-8080}/health > /dev/null; then
    echo -e "${GREEN}✓ Application is running${NC}"
else
    echo -e "${RED}✗ Application failed to start${NC}"
    docker-compose logs app
    exit 1
fi

echo -e "${BLUE}========================================${NC}"
echo -e "${GREEN}Setup completed successfully!${NC}"
echo ""
echo -e "${YELLOW}Services:${NC}"
echo -e "  Application:  ${GREEN}http://localhost:${PORT:-8080}${NC}"
echo -e "  Nginx:        ${GREEN}http://localhost:80${NC}"
echo -e "  MySQL:        ${GREEN}localhost:3306${NC}"
echo -e "  Redis:        ${GREEN}localhost:6379${NC}"
echo -e "  MailHog:      ${GREEN}http://localhost:8025${NC}"
echo -e "  Prometheus:   ${GREEN}http://localhost:9090${NC}"
echo -e "  Grafana:      ${GREEN}http://localhost:3000${NC}"
echo ""
echo -e "${YELLOW}Default credentials:${NC}"
echo -e "  Grafana: admin / admin"
echo -e "  MySQL:   ${DB_USER:-root} / ${DB_PASSWORD}"
echo ""
echo -e "${YELLOW}Next steps:${NC}"
echo "  1. Edit .env file with your production settings"
echo "  2. Run: docker-compose logs -f app"
echo "  3. Visit http://localhost to verify installation"
echo -e "${BLUE}========================================${NC}"
