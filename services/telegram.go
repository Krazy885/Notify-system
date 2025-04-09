package services

import (
	"fmt"
	"strconv"
	"notify_bot/database"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramBot struct {
	Bot              *tg.BotAPI
	DB               *database.DB
	PendingDiscordID map[int64]string // Хранит Discord ID до ввода начального интервала
	PendingTimeout   map[int64]bool   // Флаг, что пользователь в процессе ввода таймаута
}

func NewTelegramBot(token string, db *database.DB) (*TelegramBot, error) {
	bot, err := tg.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	tgBot := &TelegramBot{
		Bot:              bot,
		DB:               db,
		PendingDiscordID: make(map[int64]string),
		PendingTimeout:   make(map[int64]bool),
	}
	go tgBot.handleUpdates()
	return tgBot, nil
}

func (b *TelegramBot) SendMessage(telegramID int64, text string) error {
	msg := tg.NewMessage(telegramID, text)
	_, err := b.Bot.Send(msg)
	return err
}

// Создаем клавиатуру с кнопкой "Изменить интервал"
func (b *TelegramBot) getKeyboard() tg.ReplyKeyboardMarkup {
	return tg.NewReplyKeyboard(
		tg.NewKeyboardButtonRow(
			tg.NewKeyboardButton("Изменить интервал"),
		),
	)
}

func (b *TelegramBot) handleUpdates() {
	u := tg.NewUpdate(0)
	u.Timeout = 60
	updates := b.Bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		userID := update.Message.From.ID
		if !b.DB.IsUserRegistered(userID) {
			if update.Message.Text == "/start" {
				msg := tg.NewMessage(userID, "Привет! Чтобы зарегистрироваться, отправь мне свой Discord ID.")
				b.Bot.Send(msg)
			} else if _, exists := b.PendingDiscordID[userID]; !exists {
				// Сохраняем Discord ID и запрашиваем начальный интервал
				b.PendingDiscordID[userID] = update.Message.Text
				msg := tg.NewMessage(userID, "Отлично! Теперь укажи начальный интервал в минутах (например, 15), через сколько тебя уведомлять, если ты не ответишь.")
				b.Bot.Send(msg)
			} else {
				// Обрабатываем ввод начального интервала
				timeout, err := strconv.Atoi(update.Message.Text)
				if err != nil || timeout <= 0 {
					msg := tg.NewMessage(userID, "Пожалуйста, введи корректное число минут (например, 15).")
					b.Bot.Send(msg)
					continue
				}
				discordID := b.PendingDiscordID[userID]
				err = b.DB.RegisterUser(discordID, userID, timeout)
				if err != nil {
					msg := tg.NewMessage(userID, "Ошибка при регистрации. Попробуй еще раз с /start.")
					b.Bot.Send(msg)
					continue
				}
				msg := tg.NewMessage(userID, fmt.Sprintf("Ты успешно зарегистрирован!\nDiscord ID: %s\nИнтервал: %d минут\nНажми кнопку ниже, чтобы изменить интервал.", discordID, timeout))
				msg.ReplyMarkup = b.getKeyboard() // Добавляем клавиатуру
				b.Bot.Send(msg)
				delete(b.PendingDiscordID, userID)
			}
		} else {
			// Если пользователь ожидает ввода таймаута
			if pending, exists := b.PendingTimeout[userID]; exists && pending {
				timeout, err := strconv.Atoi(update.Message.Text)
				if err != nil || timeout <= 0 {
					msg := tg.NewMessage(userID, "Пожалуйста, введи корректное число минут (например, 15).")
					msg.ReplyMarkup = b.getKeyboard()
					b.Bot.Send(msg)
					continue
				}
				err = b.DB.RegisterUser("", userID, timeout) // Обновляем только таймаут
				if err != nil {
					msg := tg.NewMessage(userID, "Ошибка при обновлении интервала. Попробуй еще раз.")
					msg.ReplyMarkup = b.getKeyboard()
					b.Bot.Send(msg)
					continue
				}
				msg := tg.NewMessage(userID, fmt.Sprintf("Интервал обновлен: %d минут", timeout))
				msg.ReplyMarkup = b.getKeyboard()
				b.Bot.Send(msg)
				delete(b.PendingTimeout, userID) // Сбрасываем флаг
			} else if update.Message.Text == "Изменить интервал" {
				// Пользователь нажал кнопку
				b.PendingTimeout[userID] = true
				msg := tg.NewMessage(userID, "Введи новый интервал в минутах (например, 15):")
				msg.ReplyMarkup = b.getKeyboard()
				b.Bot.Send(msg)
			} else {
				// Если ввели что-то другое
				msg := tg.NewMessage(userID, "Нажми кнопку ниже, чтобы изменить интервал уведомлений.")
				msg.ReplyMarkup = b.getKeyboard()
				b.Bot.Send(msg)
			}
		}
	}
}