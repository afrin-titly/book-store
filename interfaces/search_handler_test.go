package interfaces_test

import (
	"book-apis/application"
	"book-apis/interfaces"
	"book-apis/mocks"
	"testing"
)

func TestSearchBook(t *testing.T) {
	repo := new(mocks.MockSearchRepository)
	service := application.NewSearchService(repo)
	h := interfaces.NewSearchHandler(service)
}
