package domain

type SearchCriteria struct {
	Title  string
	Author string
	Genre  string
	Price  string
	SortBy string
	Order  string
	Page   int
}

type SearchRepository interface {
	ExecuteSearch(criteria SearchCriteria) ([]Book, error)
}
