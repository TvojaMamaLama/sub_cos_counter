# CI/CD Setup Instructions

Это руководство поможет настроить автоматический деплой через GitHub Actions.

## Что делает CI/CD pipeline

1. **Test** - запускает тесты Go при каждом push/PR
2. **Build** - собирает Docker image и пушит в GitHub Container Registry  
3. **Deploy** - автоматически разворачивает на продакшн при push в main (если настроены секреты)
4. **Notify** - сообщает о статусе деплоя

**Важно**: Pipeline работает даже без настроенных секретов деплоя - будет выполняться тестирование и сборка образов.

## Требования для сервера

- Ubuntu 20.04+ / Debian 11+
- Docker и Docker Compose
- SSH доступ
- Git установлен

## Настройка GitHub Secrets

В настройках репозитория (Settings > Secrets and variables > Actions) добавьте:

### Обязательные секреты:

```bash
# Сервер для деплоя
SERVER_HOST=your-server-ip-or-domain
SERVER_USER=root  # или другой пользователь с sudo
SSH_PRIVATE_KEY=-----BEGIN OPENSSH PRIVATE KEY-----
...
-----END OPENSSH PRIVATE KEY-----

# Путь к приложению на сервере (опционально)
DEPLOY_PATH=/opt/subscription-bot  # по умолчанию

# SSH порт (опционально)
DEPLOY_PORT=22  # по умолчанию

# Telegram Bot Token
TELEGRAM_BOT_TOKEN=your-bot-token-here

# База данных
DATABASE_PASSWORD=your-secure-password

# Опционально - ограничение доступа к боту
TELEGRAM_ALLOWED_USER=123456789  # ваш Telegram ID
```

**Текущие настроенные секреты**: SERVER_HOST, SERVER_USER, SSH_PRIVATE_KEY, TELEGRAM_BOT_TOKEN, DATABASE_PASSWORD ✅

## Настройка SSH ключей

### 1. Создайте SSH ключи (на локальной машине):

```bash
ssh-keygen -t ed25519 -f ~/.ssh/deploy_key -N ""
```

### 2. Скопируйте публичный ключ на сервер:

```bash
ssh-copy-id -i ~/.ssh/deploy_key.pub user@your-server
```

### 3. Добавьте приватный ключ в GitHub Secrets:

```bash
cat ~/.ssh/deploy_key
# Скопируйте весь вывод в SSH_PRIVATE_KEY
```

## Настройка сервера

### 1. Клонируйте репозиторий на сервер:

```bash
sudo mkdir -p /opt/subscription-bot
sudo chown $USER:$USER /opt/subscription-bot
cd /opt/subscription-bot
git clone https://github.com/YOUR_USERNAME/sub_cos_counter.git .
```

### 2. Создайте .env файл:

```bash
cp configs/examples/env.docker.example .env
# Отредактируйте .env с вашими настройками
```

### 3. Установите Docker (если не установлен):

```bash
curl -fsSL https://get.docker.com | sh
sudo usermod -aG docker $USER
sudo systemctl enable docker
sudo systemctl start docker
```

## Workflow процесс

### Автоматический деплой:

1. Пушите изменения в `main` ветку
2. GitHub Actions автоматически:
   - Запускает тесты
   - Собирает Docker image
   - Пушит image в GitHub Container Registry
   - Подключается к серверу по SSH
   - Останавливает старые контейнеры
   - Запускает новые с последним image
   - Проверяет статус деплоя

### Ручной деплой:

Если нужно задеплоить вручную:

```bash
# На сервере
cd /opt/subscription-bot
git pull origin main
docker-compose -f deployment/docker/docker-compose.yml down
docker-compose -f deployment/docker/docker-compose.yml -f deployment/docker/docker-compose.prod.yml up -d --build
```

## Мониторинг деплоя

### Проверка статуса:

```bash
# На сервере
cd /opt/subscription-bot
docker-compose -f deployment/docker/docker-compose.yml ps
docker-compose -f deployment/docker/docker-compose.yml logs bot
```

### Полезные команды:

```bash
# Перезапуск бота
docker-compose -f deployment/docker/docker-compose.yml restart bot

# Просмотр логов в реальном времени
docker-compose -f deployment/docker/docker-compose.yml logs -f bot

# Проверка использования ресурсов
docker stats
```

## Troubleshooting

### Проблемы с деплоем:

1. **SSH подключение не работает**:
   - Проверьте `DEPLOY_HOST`, `DEPLOY_USER`, `DEPLOY_SSH_KEY`
   - Убедитесь что SSH ключ добавлен на сервер

2. **Docker build fails**:
   - Проверьте тесты локально: `go test ./...`
   - Проверьте Dockerfile

3. **Контейнеры не запускаются**:
   - Проверьте .env файл на сервере
   - Проверьте логи: `docker-compose logs`

4. **База данных не подключается**:
   - Проверьте `DATABASE_PASSWORD` в секретах
   - Убедитесь что PostgreSQL контейнер запущен

### Откат изменений:

```bash
# На сервере - откат к предыдущему image
docker-compose -f deployment/docker/docker-compose.yml down
docker pull ghcr.io/YOUR_USERNAME/sub_cos_counter:previous-tag
# Обновите docker-compose.yml с нужным тегом
docker-compose -f deployment/docker/docker-compose.yml up -d
```

## Security Best Practices

1. **Используйте отдельного пользователя** для деплоя (не root)
2. **Ограничьте SSH доступ** только к нужному IP
3. **Регулярно обновляйте** зависимости и базовые образы
4. **Используйте secrets** для всех чувствительных данных
5. **Мониторьте логи** на предмет подозрительной активности

## Расширенная настройка

### Добавление уведомлений в Telegram:

Можно добавить в workflow отправку уведомлений о деплое в Telegram чат.

### Настройка staging окружения:

Создайте отдельную ветку `staging` для тестирования перед продакшном.

### Backup автоматизация:

Добавьте в workflow создание backup'ов базы данных перед деплоем.