# Production Deployment Guide

## Overview

Этот проект настроен для автоматического развертывания через GitHub Actions при пуше в ветку `main`/`master`.

## 🔧 Настройка GitHub Secrets

Для работы CI/CD необходимо добавить следующие секреты в настройки GitHub репозитория:

### 1. Переход к настройкам секретов
1. Откройте репозиторий на GitHub
2. Перейдите в `Settings` → `Secrets and variables` → `Actions`
3. Нажмите `New repository secret`

### 2. Обязательные секреты

#### SSH_PRIVATE_KEY
**Назначение**: Приватный SSH ключ для подключения к серверу

**Как получить**:
```bash
# На локальной машине создать SSH ключ
ssh-keygen -t rsa -b 4096 -C "github-actions@subscription-bot"

# Скопировать публичный ключ на сервер
ssh-copy-id -i ~/.ssh/id_rsa.pub root@212.23.211.72

# Содержимое приватного ключа добавить в GitHub Secret
cat ~/.ssh/id_rsa
```

**Значение**: Вставить весь содержимый приватного ключа, включая `-----BEGIN OPENSSH PRIVATE KEY-----` и `-----END OPENSSH PRIVATE KEY-----`

#### SERVER_HOST
**Значение**: `212.23.211.72`

#### SERVER_USER  
**Значение**: `root`

#### TELEGRAM_BOT_TOKEN
**Значение**: `7860783058:AAF8j8NdOPSeHuHLK2pSLs54iG-G52vfmJE`

#### DATABASE_PASSWORD
**Значение**: `postgres_secure_password_123!` (рекомендуется сгенерировать сложный пароль)

#### TELEGRAM_ALLOWED_USER (опционально)
**Значение**: Ваш Telegram User ID (если хотите ограничить доступ к боту)

## 🚀 Процесс деплоя

### Автоматический деплой
1. Закоммитьте изменения в ветку `main`
2. GitHub Actions автоматически:
   - Запустит тесты
   - Соберет Docker образы
   - Развернет на продакшн сервере
   - Выполнит проверку работоспособности

### Ручной деплой (если нужен)
```bash
# На сервере
cd /opt/subscription-bot
sudo ./deployment/scripts/deploy.sh
```

## 📋 Что делает деплой

1. **Тестирование**: Запуск Go тестов
2. **Сборка**: Создание Docker образов
3. **Резервное копирование**: Создание бэкапа текущей версии
4. **Развертывание**: 
   - Копирование файлов на сервер
   - Создание production .env файла
   - Установка Docker/Docker Compose (если нужно)
   - Запуск контейнеров
5. **Проверка**: Верификация работоспособности сервисов

## 🔍 Мониторинг и логи

### Проверка статуса на сервере
```bash
# Статус контейнеров
cd /opt/subscription-bot
sudo docker-compose -f deployment/docker/docker-compose.yml ps

# Логи бота
sudo docker-compose -f deployment/docker/docker-compose.yml logs -f bot

# Логи PostgreSQL
sudo docker-compose -f deployment/docker/docker-compose.yml logs -f postgres

# Здоровье сервисов
sudo ./deployment/scripts/deploy.sh health
```

### Полезные команды
```bash
# Перезапуск сервисов
sudo systemctl restart subscription-bot

# Просмотр логов деплоя
sudo tail -f /var/log/subscription-bot-deploy.log

# Бэкап базы данных
sudo ./deployment/scripts/deploy.sh backup

# Подключение к БД
sudo docker-compose -f deployment/docker/docker-compose.yml exec postgres psql -U postgres -d sub_cos_counter
```

## 🛠 Структура на продакшн сервере

```
/opt/subscription-bot/          # Основная директория приложения
├── deployment/                 # Docker конфигурация
├── internal/                   # Исходный код
├── .env                       # Production переменные окружения
└── logs/                      # Логи приложения

/opt/backups/subscription-bot/  # Бэкапы
└── backup_YYYYMMDD_HHMMSS.sql

/var/lib/subscription-bot/      # Persistent данные
└── postgres/                  # PostgreSQL данные

/var/log/                      # Системные логи
├── subscription-bot-deploy.log
└── subscription-bot.log
```

## 🔧 Настройка сервера (одноразово)

При первом развертывании GitHub Actions автоматически:

1. **Установит Docker и Docker Compose**
2. **Создаст необходимые директории**:
   ```bash
   sudo mkdir -p /opt/subscription-bot
   sudo mkdir -p /opt/backups/subscription-bot
   sudo mkdir -p /var/lib/subscription-bot/postgres
   ```

3. **Настроит systemd сервис** для автоматического запуска
4. **Настроит ротацию логов** через logrotate

## 🚨 Troubleshooting

### Если деплой упал
1. Проверьте логи GitHub Actions
2. Подключитесь к серверу и проверьте:
   ```bash
   sudo journalctl -u subscription-bot
   cd /opt/subscription-bot
   sudo docker-compose -f deployment/docker/docker-compose.yml logs
   ```

### Если бот не отвечает
1. Проверьте статус контейнеров:
   ```bash
   sudo docker-compose -f deployment/docker/docker-compose.yml ps
   ```

2. Проверьте логи бота:
   ```bash
   sudo docker-compose -f deployment/docker/docker-compose.yml logs bot
   ```

3. Перезапустите сервисы:
   ```bash
   sudo systemctl restart subscription-bot
   ```

### Восстановление из бэкапа
```bash
# Список бэкапов
ls -la /opt/backups/subscription-bot/

# Восстановление БД из бэкапа
cd /opt/subscription-bot
sudo docker-compose -f deployment/docker/docker-compose.yml exec -T postgres psql -U postgres -d sub_cos_counter < /opt/backups/subscription-bot/backup_YYYYMMDD_HHMMSS.sql
```

## 📊 Мониторинг Production

### Автоматические проверки
- Health checks каждые 30 секунд
- Автоматический перезапуск при падении
- Ротация логов каждый день

### Ручная проверка
```bash
# Проверка ресурсов
sudo docker stats

# Проверка дискового пространства
df -h

# Проверка памяти
free -h

# Проверка работы бота через Telegram
# Отправьте /start боту
```

## 🔄 Rollback процедура

Если нужно откатиться к предыдущей версии:

```bash
cd /opt/subscription-bot

# Остановить текущие сервисы
sudo docker-compose -f deployment/docker/docker-compose.yml down

# Восстановить из последнего бэкапа
sudo cp -r /opt/backups/subscription-bot/subscription-bot-backup-LATEST/* .

# Запустить сервисы
sudo docker-compose -f deployment/docker/docker-compose.yml -f deployment/docker/docker-compose.prod.yml up -d
```

---

**🔗 Полезные ссылки:**
- GitHub Actions: `https://github.com/[username]/[repo]/actions`
- Сервер: `212.23.211.72`
- Telegram Bot: Найти по токену в BotFather