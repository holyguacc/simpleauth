package db

import "database/sql"

type Store interface {
	Querier
}

func NewStore(db *sql.DB) Store {
	return &SqlStore{
		db:      db,
		Queries: New(db),
	}
}

type SqlStore struct {
	*Queries
	db *sql.DB
}
