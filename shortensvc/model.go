package shortensvc

import "time"

// ShortURL represents the row in the `url` table.
type ShortURL struct {
	ID uint64 `json:"id,omitempty"`
	// Owner      string    `json:"owner,omitempty"`
	LongURL    string    `json:"long_url,omitempty"`
	LongURLCRC uint64    `json:"long_url_crc,omitempty"`
	ExpiredAt  time.Time `json:"expired_at,omitempty"`
	CreatedAt  time.Time `json:"created_at,omitempty"`
	UpdatedAt  time.Time `json:"updated_at,omitempty"`
	DeletedAt  time.Time `json:"deleted_at,omitempty"`
}
