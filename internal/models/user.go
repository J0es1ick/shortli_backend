package models

import "time"

type User struct {
	ID            int       `db:"user_id" json:"user_id,omitempty"`
	Username      string    `db:"username" json:"username,omitempty"`
	Email         string    `db:"email" json:"email,omitempty"`
	PasswordHash  string    `db:"password" json:"password,omitempty"`
	CreatedAt     time.Time `db:"created_at" json:"created_at,omitempty"`
}