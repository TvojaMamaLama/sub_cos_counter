# SSH Setup для GitHub Actions

## Быстрая настройка

### 1. Создание SSH ключа
```bash
# Создать SSH ключ для GitHub Actions
ssh-keygen -t rsa -b 4096 -f ~/.ssh/github_actions_key -N ""

# Скопировать публичный ключ на сервер
ssh-copy-id -i ~/.ssh/github_actions_key.pub root@212.23.211.72
```

### 2. Добавление секретов в GitHub

1. Откройте https://github.com/[ваш-username]/[название-репо]/settings/secrets/actions

2. Добавьте следующие секреты:

**SSH_PRIVATE_KEY**:
```bash
cat ~/.ssh/github_actions_key
```
Скопируйте весь вывод включая строки `-----BEGIN OPENSSH PRIVATE KEY-----` и `-----END OPENSSH PRIVATE KEY-----`

**SERVER_HOST**: `212.23.211.72`

**SERVER_USER**: `root`

**TELEGRAM_BOT_TOKEN**: `7860783058:AAF8j8NdOPSeHuHLK2pSLs54iG-G52vfmJE`

**DATABASE_PASSWORD**: `your_secure_password_here` (создайте надежный пароль)

**TELEGRAM_ALLOWED_USER**: `your_telegram_user_id` (опционально)

### 3. Проверка SSH подключения
```bash
# Проверить подключение
ssh -i ~/.ssh/github_actions_key root@212.23.211.72
```

### 4. Подготовка сервера (одноразово)
```bash
# Подключиться к серверу
ssh root@212.23.211.72

# Создать необходимые директории
mkdir -p /opt/subscription-bot
mkdir -p /opt/backups
mkdir -p /var/lib/subscription-bot

# Обновить систему
apt update && apt upgrade -y

# Установить базовые утилиты
apt install -y curl wget git htop
```

## После настройки

1. Сделайте коммит в ветку `main`
2. GitHub Actions автоматически развернет приложение на сервере
3. Проверьте работу бота в Telegram

## Безопасность

- Приватный ключ SSH хранится только в GitHub Secrets
- Все пароли генерируются случайным образом
- Публичный ключ добавлен только на нужный сервер
- .env файлы исключены из git через .gitignore