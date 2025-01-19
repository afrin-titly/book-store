package infrastucture

import "database/sql"

type UserRepositoryDB struct {
	DB *sql.DB
}

func NewUserRepositoryDB(db *sql.DB) *UserRepositoryDB {
	return &UserRepositoryDB{DB: db}
}
