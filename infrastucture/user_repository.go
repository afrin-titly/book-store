package infrastucture

import (
	"book-apis/domain"
	"database/sql"
)

type UserRepositoryDB struct {
	DB *sql.DB
}

func NewUserRepositoryDB(db *sql.DB) *UserRepositoryDB {
	return &UserRepositoryDB{DB: db}
}

func (u *UserRepositoryDB) CreateUser(user *domain.User) (*domain.User, error) {
	result, err := u.DB.Exec(`INSERT INTO users (name, email, password, phone, address) VALUES (?, ?, ?, ?, ?)`, user.Name, user.Email, user.Password, user.Phone, user.Address)
	if err != nil {
		return nil, err
	}
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	newUser := &domain.User{}
	err = u.DB.QueryRow(`SELECT name, email, phone, address FROM users WHERE id = ?`, lastInsertID).
		Scan(&newUser.Name, &newUser.Email, &newUser.Phone, &newUser.Address)
	if err != nil {
		return nil, err
	}
	return newUser, nil
}
