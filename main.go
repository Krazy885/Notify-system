package main

import (
	"notify_bot/config"
	"notify_bot/database"
	"notify_bot/services"
	"log"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Ошибка загрузки конфигурации:", err)
	}

	db, err := database.NewDB(cfg.DBUser, cfg.DBPass, cfg.DBHost, cfg.DBName)
	if err != nil {
		log.Fatal("Ошибка подключения к базе данных:", err)
	}

	tgBot, err := services.NewTelegramBot(cfg.TelegramToken, db)
	if err != nil {
		log.Fatal("Ошибка создания Telegram-бота:", err)
	}

	discordBot, err := services.NewDiscordBot(cfg.DiscordToken, db, tgBot)
	if err != nil {
		log.Fatal("Ошибка создания Discord-бота:", err)
	}

	if err := discordBot.Start(); err != nil {
		log.Fatal("Ошибка запуска Discord-бота:", err)
	}

	select {}
}


//sudo supervisorctl restart notify
//sudo tail -f /var/log/notify.err.log