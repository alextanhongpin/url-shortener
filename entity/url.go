package entity

import "time"

// URL represents the row in the `url` table.
type URL struct {
	ID        uint64    `json:"id"`
	URL       string    `json:"url"`
	URLCRC    uint64    `json:"url_crc"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}
