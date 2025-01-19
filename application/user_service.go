package application

import (
	"book-apis/domain"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	service domain.UserRepository
}

func NewUserService(service domain.UserRepository) *UserService {
	return &UserService{service: service}
}

func (s *UserService) CreateUser(user *domain.User) (*domain.User, error) {
	if isExists, err := s.service.ExistsByEmail(user.Email); err != nil {
		return nil, err
	} else if isExists {
		return nil, errors.New("email already exixts")
	}
	var err error
	user.Password, err = hashPassword(user.Password)
	if err != nil {
		return nil, errors.New("failed to hashed password")
	}
	return s.service.CreateUser(user)
}

func (s *UserService) GetUser(ID int) (domain.User, error) {
	return s.service.GetUser(ID)
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}
