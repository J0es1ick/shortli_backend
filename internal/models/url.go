package models

import "time"

type URL struct {
	ID           int       `db:"url_id" json:"url_id,omitempty"`
	OriginalURL  string    `db:"original_url" json:"original_url,omitempty"`
	ShortCode    string    `db:"short_code" json:"short_code,omitempty"`
	UserId 		 int 	   `db:"user_id" json:"user_id,omitempty"`
	ClickCount   int       `db:"click_count" json:"click_count,omitempty"`
	QRClickCount int       `db:"qr_click_count" json:"qr_click_count,omitempty"`
	CreatedAt    time.Time `db:"created_at" json:"created_at,omitempty"`
}