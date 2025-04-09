package database

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

type DB struct {
	*sql.DB
}

func NewDB(user, pass, host, name string) (*DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s", user, pass, host, name)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	return &DB{db}, nil
}

func (db *DB) GetTelegramID(discordID string) (int64, error) {
	var telegramID int64
	err := db.QueryRow("SELECT telegram_id FROM users WHERE discord_id = ?", discordID).Scan(&telegramID)
	if err != nil {
		return 0, err
	}
	return telegramID, nil
}

func (db *DB) RegisterUser(discordID string, telegramID int64, timeoutMinutes int) error {
	_, err := db.Exec("INSERT INTO users (discord_id, telegram_id, timeout_minutes) VALUES (?, ?, ?) ON DUPLICATE KEY UPDATE telegram_id = ?, timeout_minutes = ?", 
		discordID, telegramID, timeoutMinutes, telegramID, timeoutMinutes)
	return err
}

func (db *DB) IsUserRegistered(telegramID int64) bool {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM users WHERE telegram_id = ?", telegramID).Scan(&count)
	return err == nil && count > 0
}

func (db *DB) GetUserTimeout(discordID string) (int, error) {
	var timeoutMinutes int
	err := db.QueryRow("SELECT timeout_minutes FROM users WHERE discord_id = ?", discordID).Scan(&timeoutMinutes)
	if err != nil {
		return 0, err
	}
	return timeoutMinutes, nil
}

func (db *DB) GetAllUsers() (map[string]int64, error) {
	rows, err := db.Query("SELECT discord_id, telegram_id FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make(map[string]int64)
	for rows.Next() {
		var discordID string
		var telegramID int64
		if err := rows.Scan(&discordID, &telegramID); err != nil {
			return nil, err
		}
		users[discordID] = telegramID
	}
	return users, nil
}