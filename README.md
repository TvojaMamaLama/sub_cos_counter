# Subscription Cost Counter Bot

Телеграм-бот для трекинга подписок с кнопочным интерфейсом.

## Возможности

- ✅ Добавление подписок с указанием стоимости, периодичности и категории
- 💰 Подсчет месячных расходов по валютам  
- 📊 Аналитика по категориям (развлечения, работа, обучение, дом)
- 📅 Отслеживание дат платежей и отметка об оплате
- 🔄 Поддержка автопродления подписок
- 💵 Поддержка USD и RUB валют
- 📜 История всех платежей

## Технологии

- **Backend**: Go 1.24+
- **База данных**: PostgreSQL
- **Миграции**: Atlas  
- **Telegram Bot**: Telebot v3
- **БД драйвер**: pgx v5
- **Конфигурация**: Viper (YAML + env vars)

## Установка и запуск

### 🐳 Быстрый старт с Docker (рекомендуется)

#### 1. Подготовка
```bash
git clone <repository>
cd sub_cos_counter
make setup  # Создаст .env файл
```

#### 2. Настройка
Отредактируйте `.env` файл:
```bash
# Обязательно укажите токен бота
TELEGRAM_BOT_TOKEN=your_bot_token_here

# Опционально: ваш Telegram user ID для персонального бота
TELEGRAM_ALLOWED_USER=123456789
```

#### 3. Запуск
```bash
# Запустить всё (PostgreSQL + Bot)
make up

# Применить миграции БД
make migrate

# Проверить логи
make logs
```

#### 4. Управление
```bash
make help          # Показать все команды
make logs-bot       # Логи только бота
make logs-db        # Логи только БД
make restart        # Перезапустить сервисы
make down          # Остановить всё
```

---

### 🔧 Ручная установка (для разработки)

#### 1. Установка зависимостей

```bash
go mod download
```

### 2. Настройка PostgreSQL

Создайте базу данных:
```sql
CREATE DATABASE sub_cos_counter;
```

### 3. Установка Atlas (для миграций)

```bash
# macOS
brew install ariga/tap/atlas

# Linux
curl -sSf https://atlasgo.sh | sh

# Windows
# Скачайте бинарик с GitHub releases
```

### 4. Выполнение миграций

```bash
# Создание миграции из схемы
atlas migrate diff --env local

# Применение миграций
atlas migrate apply --env local
```

### 5. Настройка конфигурации

#### Вариант 1: Переменные окружения (рекомендуется)

Создайте файл `.env`:
```env
# Скопируйте из .env.example и заполните свои значения
TELEGRAM_BOT_TOKEN=your_bot_token_here
DATABASE_URL=postgres://username:password@localhost:5432/sub_cos_counter?sslmode=disable
APP_ENVIRONMENT=development
```

#### Вариант 2: YAML конфигурация

Создайте файл `config.yaml`:
```yaml
telegram:
  bot_token: "your_bot_token_here"
  
database:
  url: "postgres://username:password@localhost:5432/sub_cos_counter?sslmode=disable"
  
app:
  environment: "development"
  debug: true
```

#### Вариант 3: Смешанный подход

Используйте `config/config.yaml` для базовых настроек и переменные окружения для секретов:
```yaml
# config/config.yaml
app:
  name: "My Subscription Bot"
  debug: false

database:
  host: "localhost"
  port: 5432
  name: "sub_cos_counter"
  # password задается через DATABASE_PASSWORD env var
```

### 6. Запуск бота

```bash
go run cmd/bot/main.go
```

## Создание Telegram бота

1. Найдите @BotFather в Telegram
2. Отправьте команду `/newbot`
3. Следуйте инструкциям для создания бота
4. Получите токен и добавьте его в `.env`

## Структура проекта

```
├── cmd/bot/main.go              # Точка входа
├── internal/                    # Исходный код приложения
│   ├── config/config.go         # Конфигурация (Viper)
│   ├── models/                  # Модели данных
│   ├── repository/              # Слой данных (pgx)
│   ├── services/                # Бизнес-логика
│   └── bot/                     # Telegram bot (Telebot)
├── deployment/                  # Деплой и контейнеризация
│   ├── docker/                  # Docker конфигурация
│   │   ├── Dockerfile
│   │   ├── docker-compose.yml
│   │   ├── docker-compose.override.yml  # Development
│   │   └── docker-compose.prod.yml      # Production
│   └── scripts/                 # Скрипты деплоя
│       └── deploy.sh
├── configs/                     # Конфигурационные файлы
│   └── examples/                # Примеры конфигурации
│       ├── config.yaml.example
│       └── env.docker.example
├── migrations/                  # Миграции БД (Atlas)
│   ├── atlas.hcl
│   └── schema.sql
├── Makefile                     # Команды управления
├── QUICKSTART.md               # Быстрый старт
└── README.md
```

## Использование бота

### Главное меню
- 📝 **Добавить подписку** - Создание новой подписки
- 📋 **Мои подписки** - Просмотр и управление подписками  
- 💰 **Месячные расходы** - Сумма трат за месяц
- 📊 **Аналитика** - Разбивка по категориям
- 📜 **История платежей** - Последние операции
- ⚙️ **Настройки** - Информация о боте

### Добавление подписки

Пошаговый процесс с кнопками:
1. Выбор категории (🎮 Развлечения, 💼 Работа, 📚 Обучение, 🏠 Дом)
2. Выбор валюты (💵 USD, 🔹 RUB) 
3. Выбор периода (🗓️ Неделя, 📅 Месяц, 📆 Год, ⚡ Другое)
4. Настройка автопродления (✅ Да, ❌ Нет)
5. Ввод названия подписки
6. Ввод стоимости

### Управление подписками

Для каждой подписки доступны действия:
- ✅ **Оплатить** - Отметить платеж как выполненный
- ❌ **Удалить** - Деактивировать подписку

## Разработка

### Добавление новых миграций

```bash
# Внесите изменения в migrations/schema.sql
# Создайте миграцию
atlas migrate diff --env local
```

### Тестирование

```bash
go test ./...
```

## Особенности архитектуры

- **Чистая архитектура** с разделением на слои
- **Repository pattern** для работы с БД
- **Service layer** для бизнес-логики  
- **Connection pooling** через pgx с настраиваемыми параметрами
- **Graceful shutdown** для корректного завершения
- **Inline keyboards** для удобного UI
- **State management** для многошаговых диалогов
- **Гибкая конфигурация** через Viper (YAML + env vars + defaults)

## Конфигурация

### Приоритет настроек (от высшего к низшему):
1. **Переменные окружения** (например: `TELEGRAM_BOT_TOKEN`)
2. **Конфигурационный файл** (`config.yaml`, `config/config.yaml`)
3. **Значения по умолчанию**

### Расширенные настройки:

#### Персональный бот
```yaml
telegram:
  allowed_user: 123456789  # Ваш Telegram user ID
```

#### Настройки пула соединений БД
```yaml
database:
  max_connections: 25
  max_idle_time: 15        # минуты
  conn_max_lifetime: 60    # минуты
```

#### Логирование
```yaml
logging:
  level: "debug"           # debug, info, warn, error
  format: "json"           # json, text
  output: "file"           # stdout, stderr, file
  filename: "bot.log"      # если output = file
```

### Переменные окружения

Все настройки можно переопределить через env vars с префиксом разделённым подчёркиванием:
- `APP_NAME` → `app.name`
- `TELEGRAM_BOT_TOKEN` → `telegram.bot_token`
- `DATABASE_MAX_CONNECTIONS` → `database.max_connections`

## 🐳 Docker команды

### Основные команды
```bash
make setup         # Первоначальная настройка
make up           # Запустить все сервисы
make down         # Остановить все сервисы
make logs         # Показать логи всех сервисов
make migrate      # Выполнить миграции БД
make clean        # Очистить контейнеры
```

### Разработка
```bash
make dev-setup    # Настройка для разработки
make up-dev       # Запуск в режиме разработки
make rebuild      # Пересборка с нуля
make shell-bot    # Войти в контейнер бота
make shell-db     # Подключиться к PostgreSQL
```

### Мониторинг
```bash
make health       # Проверить здоровье сервисов
make stats        # Показать статистику контейнеров
make logs-bot     # Логи только бота
make logs-db      # Логи только БД
```

### Production
```bash
make deploy       # Полный деплой (setup + build + migrate + up)
make backup-db    # Создать бэкап БД
```

## Docker Compose архитектура

### Сервисы:
- **postgres** - PostgreSQL 16 с health check и persistent storage
- **bot** - Telegram бот с автоматическим restart
- **atlas** - Сервис для миграций (запускается по требованию)

### Сети:
- **subscription-bot-network** - изолированная сеть для всех сервисов

### Volumes:
- **postgres_data** - постоянное хранилище данных PostgreSQL
- **logs** - логи приложения (монтируется в ./logs)

### Переменные окружения:
Все настройки настраиваются через `.env` файл или переменные окружения.

## Развитие проекта

Планируемые улучшения:
- 🔔 Уведомления о предстоящих платежах
- 📈 Графики и расширенная аналитика  
- 💱 Автоматическая конвертация валют
- 📱 Web интерфейс для настроек
- 🏷️ Теги и расширенная категоризация
- 🐳 Kubernetes манифесты для production
- 📊 Мониторинг и метрики (Prometheus/Grafana)