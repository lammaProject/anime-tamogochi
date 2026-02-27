package model

import "time"

type User struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
}

type LikedCatgirl struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	CatgirlID string    `json:"catgirl_id"`
	Data      string    `json:"data"` // JSON blob NekosImageData
	CreatedAt time.Time `json:"created_at"`
}

type Message struct {
	ID         int       `json:"id"`
	UserID     int       `json:"user_id"`
	CatgirlID  string    `json:"catgirl_id"`
	Content    string    `json:"content"`
	FromUser   bool      `json:"from_user"`
	CreatedAt  time.Time `json:"created_at"`
}