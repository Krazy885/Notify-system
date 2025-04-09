package models

import "time"

type UserState struct {
	LastMessage time.Time
	LastMention time.Time
}