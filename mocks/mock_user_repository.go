package mocks

import (
	"book-apis/domain"

	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) CreateUser(user *domain.User) (*domain.User, error) {
	args := m.Called(user)
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) GetUser(ID int) (domain.User, error) {
	args := m.Called(ID)
	return args.Get(0).(domain.User), args.Error(1)
}

func (m *MockUserRepository) ExistsByEmail(email string) (bool, error) {
	args := m.Called(email)
	return args.Get(0).(bool), args.Error(1)
}
