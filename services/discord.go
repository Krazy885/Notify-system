package services

import (
	"log"
	"notify_bot/database"
	"notify_bot/models"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

type DiscordBot struct {
	Session     *discordgo.Session
	DB          *database.DB
	Users       map[string]*models.UserState
	TelegramBot TelegramSender
}

type TelegramSender interface {
	SendMessage(telegramID int64, text string) error
}

func NewDiscordBot(token string, db *database.DB, tgBot TelegramSender) (*DiscordBot, error) {
	s, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}
	s.Identify.Intents = discordgo.IntentsGuildMessages
	bot := &DiscordBot{
		Session:     s,
		DB:          db,
		Users:       make(map[string]*models.UserState),
		TelegramBot: tgBot,
	}
	s.AddHandler(bot.messageCreate)
	go bot.checkTimeouts()
	log.Println("Дискорд-бот инициализирован для всех зарегистрированных пользователей")
	return bot, nil
}

func (b *DiscordBot) Start() error {
	err := b.Session.Open()
	if err != nil {
		log.Printf("Не удалось запустить Дискорд-бот: %v", err)
		return err
	}
	log.Println("Дискорд-бот успешно запущен")
	go b.updateStatusLoop()
	return nil
}

func (b *DiscordBot) messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	log.Printf("Получено сообщение от %s в канале %s: %s", m.Author.ID, m.ChannelID, m.Content)

	// Игнорируем сообщения от нашего бота
	if m.Author.ID == s.State.User.ID {
		log.Println("Игнорируем сообщение от самого бота")
		return
	}

	// Игнорируем сообщения от другого бота с ID 1246124994325643305
	if m.Author.ID == "1246124994325643305" {
		log.Println("Игнорируем сообщение от бота с ID 1246124994325643305")
		return
	}

	// Обновляем время последнего сообщения для автора, если он зарегистрирован
	users, err := b.DB.GetAllUsers()
	if err != nil {
		log.Printf("Ошибка получения списка пользователей: %v", err)
		return
	}

	if _, exists := users[m.Author.ID]; exists {
		userState, exists := b.Users[m.Author.ID]
		if !exists {
			userState = &models.UserState{}
			b.Users[m.Author.ID] = userState
		}
		userState.LastMessage = time.Now()
		log.Printf("Обновлено время последнего сообщения для пользователя %s на %v", m.Author.ID, userState.LastMessage)
	}

	// Проверяем явные теги в содержимом сообщения
	log.Println("Проверяем явные теги в содержимом сообщения...")
	for userID := range users {
		tag := "<@" + userID + ">"
		if strings.Contains(m.Content, tag) {
			userState, exists := b.Users[userID]
			if !exists {
				userState = &models.UserState{}
				b.Users[userID] = userState
			}
			userState.LastMention = time.Now()
			log.Printf("Зарегистрированный пользователь %s упомянут явным тегом в %v", userID, userState.LastMention)
		}
	}
}

func (b *DiscordBot) checkTimeouts() {
	for {
		time.Sleep(1 * time.Minute)
		log.Println("Проверка таймаутов...")
		users, err := b.DB.GetAllUsers()
		if err != nil {
			log.Printf("Ошибка получения списка пользователей для проверки таймаутов: %v", err)
			continue
		}

		for userID, state := range b.Users {
			log.Printf("Пользователь %s: Последнее упоминание=%v, Последнее сообщение=%v", userID, state.LastMention, state.LastMessage)
			if !state.LastMention.IsZero() && state.LastMessage.Before(state.LastMention) {
				timeoutMinutes, err := b.DB.GetUserTimeout(userID)
				if err != nil {
					log.Printf("Ошибка получения интервала для %s: %v", userID, err)
					continue
				}
				log.Printf("Интервал для %s составляет %d минут", userID, timeoutMinutes)
				if time.Since(state.LastMention) > time.Duration(timeoutMinutes)*time.Minute {
					telegramID, ok := users[userID]
					if !ok {
						log.Printf("Пользователь %s не найден в базе для отправки уведомления", userID)
						continue
					}
					log.Printf("Отправка уведомления на Telegram ID %d", telegramID)
					err = b.TelegramBot.SendMessage(telegramID, "Тебя упомянули на сервере ShinDau, скорее заходи и ответь!")
					if err != nil {
						log.Printf("Ошибка отправки сообщения в Telegram на %d: %v", telegramID, err)
						continue
					}
					log.Printf("Уведомление отправлено пользователю %s (Telegram ID: %d)", userID, telegramID)
					state.LastMention = time.Time{}
				}
			}
		}
	}
}


func (b *DiscordBot) updateStatusLoop() {
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Статус 1: "Модерации сервера"
			err := b.Session.UpdateStatusComplex(discordgo.UpdateStatusData{
				Status: "online",
				Activities: []*discordgo.Activity{
					{
						Name: "регистрации в тг @ShinDau_notify_bot",
						Type: discordgo.ActivityTypeWatching,
					},
				},
			})
			if err != nil {
				log.Printf("Ошибка установки статуса 1: %v", err)
			}

			time.Sleep(15 * time.Second)

			// Статус 2: "за {members} участниками"
			err = b.Session.UpdateStatusComplex(discordgo.UpdateStatusData{
				Status: "online",
				Activities: []*discordgo.Activity{
					{
						Name: "за чатом",
						Type: discordgo.ActivityTypeWatching,
					},
				},
			})
			if err != nil {
				log.Printf("Ошибка установки статуса 2: %v", err)
			}

			time.Sleep(15 * time.Second)

			
		}
	}
}