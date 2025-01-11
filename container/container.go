package container

import (
	"book-apis/application"
	"book-apis/infrastucture"
	"book-apis/interfaces"
	"database/sql"
)

type Container struct {
	DB            *sql.DB
	BookHandler   *interfaces.BookHandler
	SearchHandler *interfaces.SearchHandler
}

func NewContainer() *Container {
	connStirng := "host=localhost port=3306 user=mysql password=secret dbname=books sslmode=disable"
	db, err := sql.Open("mysql", connStirng)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	bookRepo := infrastucture.NewBookRepositoryDB(db)
	searchRepo := infrastucture.NewSearchRepositoryDB(db)

	bookService := application.NewBookService(bookRepo)
	searchService := application.NewSearchService(searchRepo)

	bookHandler := interfaces.NewBookHandler(bookService)
	searchHandler := interfaces.NewSearchHandler(searchService)

	return &Container{
		DB:            db,
		BookHandler:   bookHandler,
		SearchHandler: searchHandler,
	}
}
