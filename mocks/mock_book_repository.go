package mocks

import (
	"book-apis/domain"

	"github.com/stretchr/testify/mock"
)

type MockBookRepository struct {
	mock.Mock
}

func (m *MockBookRepository) GetAll() ([]domain.Book, error) {
	args := m.Called()
	return args.Get(0).([]domain.Book), args.Error(1)
}

func (m *MockBookRepository) GetBook(ID int) (domain.Book, error) {
	args := m.Called()
	return args.Get(0).(domain.Book), args.Error(1)
}
