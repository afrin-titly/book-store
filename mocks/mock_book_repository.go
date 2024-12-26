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
	args := m.Called(ID)
	return args.Get(0).(domain.Book), args.Error(1)
}

func (m *MockBookRepository) CreateBook(book *domain.Book) (*domain.Book, error) {
	args := m.Called(book)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Book), args.Error(1)
}
