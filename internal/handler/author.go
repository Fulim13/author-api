package handler

import (
	"errors"
	"net/http"

	"github.com/Fulim13/author-api/internal/service"
	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
)

type AuthorHandler struct {
	service *service.AuthorService
}

func NewAuthorHandler(service *service.AuthorService) *AuthorHandler {
	return &AuthorHandler{
		service: service,
	}
}

type CreateAuthorRequest struct {
	Name string `json:"name"`
	Bio  string `json:"bio"`
}

type UpdateAuthorRequest struct {
	Name *string `json:"name"`
	Bio  *string `json:"bio"`
}

type CreateBookRequest struct {
	Name string `json:"name"`
}

func authorID(c *echo.Context) (uuid.UUID, error) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return uuid.Nil, echo.NewHTTPError(http.StatusBadRequest, "invalid author id")
	}
	return id, nil
}

func (h *AuthorHandler) createAuthor(c *echo.Context) error {
	var reqBody CreateAuthorRequest
	if err := c.Bind(&reqBody); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	if reqBody.Name == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "name is required")
	}

	if reqBody.Bio == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "bio is required")
	}

	author, err := h.service.CreateAuthor(c.Request().Context(), reqBody.Name, reqBody.Bio)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}
	return c.JSON(http.StatusCreated, author)
}

func (h *AuthorHandler) listAuthors(c *echo.Context) error {
	authors, err := h.service.ListAuthors(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}
	return c.JSON(http.StatusOK, authors)
}

func (h *AuthorHandler) countAuthors(c *echo.Context) error {
	count, err := h.service.CountAuthors(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}
	return c.JSON(http.StatusOK, map[string]int64{"count": count})
}

func (h *AuthorHandler) getAuthor(c *echo.Context) error {
	id, err := authorID(c)
	if err != nil {
		return err
	}

	author, err := h.service.GetAuthor(c.Request().Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrAuthorNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "author not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}
	return c.JSON(http.StatusOK, author)
}

func (h *AuthorHandler) updateAuthor(c *echo.Context) error {
	id, err := authorID(c)
	if err != nil {
		return err
	}

	var reqBody UpdateAuthorRequest
	if err := c.Bind(&reqBody); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	if reqBody.Name == nil && reqBody.Bio == nil {
		return echo.NewHTTPError(http.StatusBadRequest, "name or bio is required")
	}

	if reqBody.Name != nil && *reqBody.Name == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "name cannot be empty")
	}

	err = h.service.UpdateAuthor(c.Request().Context(), id, reqBody.Name, reqBody.Bio)
	if err != nil {
		if errors.Is(err, service.ErrAuthorNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "author not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *AuthorHandler) deleteAuthor(c *echo.Context) error {
	id, err := authorID(c)
	if err != nil {
		return err
	}

	err = h.service.DeleteAuthor(c.Request().Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrAuthorNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "author not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *AuthorHandler) createBook(c *echo.Context) error {
	authorId, err := authorID(c)
	if err != nil {
		return err
	}

	var reqBody CreateBookRequest
	if err := c.Bind(&reqBody); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	if reqBody.Name == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "name is required")
	}

	book, err := h.service.CreateBook(c.Request().Context(), authorId, reqBody.Name)

	if err != nil {
		if errors.Is(err, service.ErrAuthorNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "author not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}

	return c.JSON(http.StatusCreated, book)
}

func (h *AuthorHandler) listBooks(c *echo.Context) error {
	authorId, err := authorID(c)
	if err != nil {
		return err
	}
	books, err := h.service.ListBooks(c.Request().Context(), authorId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}
	return c.JSON(http.StatusOK, books)
}

func (h *AuthorHandler) listAuthorsAndBooks(c *echo.Context) error {
	authorsAndBooks, err := h.service.ListAuthorsAndBooks(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}
	return c.JSON(http.StatusOK, authorsAndBooks)
}
