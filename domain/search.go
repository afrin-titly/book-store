package domain

type SearchCriteria struct {
	Title  string
	Author string
	Genre  string
	Price  int
	SortBy string
	Order  string
}

type SearchRepository interface {
	ExecuteSearch(criteria SearchCriteria) ([]Book, error)
}
