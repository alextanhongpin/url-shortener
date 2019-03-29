package urlshortener

import (
	"database/sql"
	"log"
)

type Repository struct {
	stmts Statements
}

func NewRepository(db *sql.DB) *Repository {
	prepare := func(rawStmt string) *sql.Stmt {
		stmt, err := db.Prepare(rawStmt)
		if err != nil {
			log.Fatal(err)
		}
		return stmt
	}
	repo := &Repository{make(Statements)}
	for id, stmt := range rawStmts {
		repo.stmts[id] = prepare(stmt)
	}
	return repo
}

func (r *Repository) Create(it ShortURL) (int64, error) {
	res, err := r.stmts[createStmt].Exec(it.LongURL)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *Repository) WithID(shortURL string) (*ShortURL, error) {
	var res ShortURL
	err := r.stmts[withIDStmt].QueryRow(shortURL).Scan(&res.LongURL)
	return &res, err
}
