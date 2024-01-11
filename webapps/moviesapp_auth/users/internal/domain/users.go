package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID               int       `json:"id"`
	Login            string    `json:"login"`
	Password         string    `json:"password"`
	PasswordExpires  time.Time `json:"passwordExpires"`
	NotificationSent bool      `json:"notificationSent"`
}

type Session struct {
	ID        int       `json:"id"`
	Key       uuid.UUID `json:"key"`
	UserID    int       `json:"userId"`
	StartedAt time.Time `json:"startedAt"`
}

var ErrNotFound = errors.New("not found")
var ErrConflict = errors.New("already exists")
