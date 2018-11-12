package shortener

import (
	"database/sql"
	"errors"
	"log"

	"github.com/alextanhongpin/url-shortener/entity"
)

// RepositoryImpl return an implementation of the URL repository.
type RepositoryImpl struct {
	db    *sql.DB
	stmts map[string]*sql.Stmt
}

// NewRepository returns a new URL repository.
func NewRepository(db *sql.DB) *RepositoryImpl {
	repo := RepositoryImpl{
		db:    db,
		stmts: make(map[string]*sql.Stmt),
	}
	err := repo.initStatements()
	if err != nil {
		log.Fatal(err)
	}
	return &repo
}

func (r *RepositoryImpl) initStatements() error {
	stmtInsert, err := r.db.Prepare("INSERT INTO url (url, url_crc) VALUES (?, CRC32(?))")
	if err != nil {
		return err
	}
	r.stmts["insert"] = stmtInsert

	stmtGet, err := r.db.Prepare("SELECT (url) FROM url WHERE id = ?")
	if err != nil {
		return err
	}
	r.stmts["get"] = stmtGet
	return nil
}

// Close performs cleanup on the repository.
func (r *RepositoryImpl) Close() error {
	for key, stmt := range r.stmts {
		stmt.Close()
		delete(r.stmts, key)
	}
	return nil
}

// Get returns a row from the URL table by the given row id.
func (r *RepositoryImpl) Get(id uint64) (*entity.URL, error) {
	var result entity.URL
	if err := r.stmts["get"].QueryRow(id).Scan(&result.URL); err != nil {
		return nil, err
	}
	return &result, nil
}

// Insert creates a new row in the URL table and returns the created row id.
func (r *RepositoryImpl) Insert(longURL string) (uint64, error) {
	result, err := r.stmts["insert"].Exec(longURL, longURL)
	if err != nil {
		return 0, err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	if rowsAffected != 1 {
		return 0, errors.New("insert failed")
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	// Id is returned as int64, convert it to uint64.
	return uint64(id), nil
}
