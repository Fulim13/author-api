package handler

import "github.com/labstack/echo/v5"

func Routes(e *echo.Echo, authorHandler *AuthorHandler) {
	e.POST("/authors", authorHandler.createAuthor)
	e.GET("/authors", authorHandler.listAuthors)
	e.GET("/authors/len", authorHandler.countAuthors)
	e.GET("/authors/:id", authorHandler.getAuthor)
	e.PATCH("/authors/:id", authorHandler.updateAuthor)
	e.DELETE("/authors/:id", authorHandler.deleteAuthor)
	e.POST("/authors/:id/books", authorHandler.createBook)
	e.GET("/authors/:id/books", authorHandler.listBooks)
	e.GET("/authors-and-books", authorHandler.listAuthorsAndBooks)
}
