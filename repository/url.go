package repository

import "github.com/alextanhongpin/url-shortener/entity"

// URL repository represents the repository interface to communicate with the
// storage layer.
type URL interface {
	// Get a row by the given id.
	Get(id uint64) (*entity.URL, error)

	// Insert a new row.
	Insert(longURL string) (uint64, error)

	// Perform cleanup.
	Close() error
}
