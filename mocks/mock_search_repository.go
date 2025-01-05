package mocks

import (
	"book-apis/domain"

	"github.com/stretchr/testify/mock"
)

type MockSearchRepository struct {
	mock.Mock
}

func (m *MockSearchRepository) ExecuteSearch(criteria domain.SearchCriteria) ([]domain.Book, error) {
	args := m.Called(criteria)
	return args.Get(0).([]domain.Book), args.Error(1)
}
