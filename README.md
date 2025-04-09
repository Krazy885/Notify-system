<h1 align="center">ShinDau Notify Bot</h1>

<p align="center">
  Бот для уведомлений между Discord и Telegram. Следит за упоминаниями в Discord и напоминает в Telegram, если вы не ответили вовремя! 🚀
</p>


## 📝 Описание

**ShinDau Notify Bot** — это удобный инструмент для тех, кто хочет оставаться в курсе упоминаний на сервере Discord, даже если не находится онлайн. Бот интегрируется с Telegram, позволяя регистрировать пользователей, настраивать интервал уведомлений и получать напоминания, если вы не ответили на тег в Discord.

---

## ✨ Основные возможности

- **Регистрация через Telegram** 📱: Укажи свой Discord ID и интервал уведомлений.
- **Отслеживание тегов** 🎯: Реагирует только на явные упоминания (`<@userID>`), игнорируя ответы через reply.
- **Гибкий интервал** ⏱️: Меняй время ожидания ответа через удобную кнопку в Telegram.
- **База данных** 🗄️: Использует MySQL для хранения данных пользователей.

---

## ⚙️ Требования

- **Go**: 1.16 или выше
- **MySQL**: 5.7+ 
- **Токены**: Discord Bot Token и Telegram Bot Token
- **Зависимости**:
  - `github.com/joho/godotenv`
  - `github.com/go-sql-driver/mysql`
  - `github.com/bwmarrin/discordgo`
  - `github.com/go-telegram-bot-api/telegram-bot-api/v5`

---

## 🛠️ Установка

### 1. Клонирование репозитория
```bash
git clone https://github.com/Krazy885/Notify-system.git
cd shindau-notify-bot
```

### 2. Настройка окружения
Создай файл `.env` в корне проекта:
```bash
touch .env
nano .env
```
Вставь:
```
DISCORD_TOKEN=your_discord_token
TELEGRAM_TOKEN=your_telegram_token
DB_NAME=notify_db
DB_HOST=localhost:3306
DB_USER=your_db_user
DB_PASS=your_db_password
```

### 3. Установка зависимостей
Инициализируй модуль и подтяни зависимости:
```bash
go mod init shindau-notify-bot
go get github.com/joho/godotenv
go get github.com/go-sql-driver/mysql
go get github.com/bwmarrin/discordgo
go get github.com/go-telegram-bot-api/telegram-bot-api/v5
```

### 4. Настройка базы данных
Подключись к MySQL и создай базу с таблицей:
```sql
CREATE DATABASE notify_db;
USE notify_db;
CREATE TABLE users (
    discord_id VARCHAR(255) PRIMARY KEY,
    telegram_id BIGINT NOT NULL UNIQUE,
    timeout_minutes INT NOT NULL
);
```
Настрой пользователя:
```sql
CREATE USER 'your_db_user'@'localhost' IDENTIFIED BY 'your_db_password';
GRANT ALL PRIVILEGES ON notify_db.* TO 'your_db_user'@'localhost';
FLUSH PRIVILEGES;
```

### 5. Сборка и запуск
Скомпилируй и запусти бота:
```bash
go build -o notify_bot main.go
./notify_bot
```

---

## 🚀 Использование

### Регистрация
1. Напиши боту в Telegram: `/start`.
2. Укажи свой Discord ID.
3. Введи интервал уведомлений в минутах (например, `15`).
4. Получи сообщение:  
   ```
   Ты успешно зарегистрирован!
   Discord ID: 1213123123123123
   Интервал: 15 минут
   Нажми кнопку ниже, чтобы изменить интервал.
   ```

### Изменение интервала
1. Нажми кнопку "Изменить интервал".
2. Введи новое значение (например, `30`).
3. Получи подтверждение: `Интервал обновлен: 30 минут`.

### Уведомления
- Если тебя тегнут в Discord (`<@your_discord_id>`), и ты не ответишь в течение заданного времени, бот отправит в Telegram:  
  ```
  Тебя упомянули на сервере ShinDau, скорее заходи и ответь!
  ```
---

## 📂 Структура проекта

```
shindau-notify-bot/
├── config/         # Загрузка конфигурации из .env
├── database/       # Интерфейс для работы с MySQL
├── services/       # Логика Discord и Telegram ботов
├── models/         # Модели данных (UserState)
├── main.go         # Точка входа
├── .env            # Файл окружения (игнорируется Git)
├── .gitignore      # Игнорируемые файлы
└── README.md       # Эта документация
```

---

## 🤝 Контрибьютинг

Хочешь улучшить бота? Форкни репозиторий, внеси изменения и создай Pull Request! Все идеи приветствуются.

---

## 📜 Лицензия

Проект распространяется под [MIT License](LICENSE). Используй его свободно!

---

## 📬 Контакты

- **Issues**: [Сообщить о проблеме](https://t.me/Krazy885_support_bot)

<p align="center">Сделано с ❤️ для сообщества ShinDau!</p>