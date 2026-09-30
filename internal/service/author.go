package service

import (
	"context"
	"errors"

	"github.com/Fulim13/author-api/internal/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

var (
	ErrAuthorNotFound = errors.New("author not found")
)

type DB interface {
	repository.DBTX
	Begin(ctx context.Context) (pgx.Tx, error)
}

type AuthorService struct {
	db   DB
	repo *repository.Queries
}

func NewAuthorService(db DB) *AuthorService {
	return &AuthorService{
		db:   db,
		repo: repository.New(db),
	}
}

func (s *AuthorService) withTx(ctx context.Context, fn func(q *repository.Queries) error) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if err := fn(s.repo.WithTx(tx)); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (s *AuthorService) CreateAuthor(ctx context.Context, name string, bio string) (repository.Author, error) {
	return s.repo.CreateAuthor(ctx, repository.CreateAuthorParams{
		ID:   uuid.New(),
		Name: name,
		Bio:  bio,
	})
}

func (s *AuthorService) ListAuthors(ctx context.Context) ([]repository.Author, error) {
	authors, err := s.repo.ListAuthors(ctx)
	if err != nil {
		return nil, err
	}
	if authors == nil {
		authors = []repository.Author{}
	}
	return authors, nil
}

func (s *AuthorService) CountAuthors(ctx context.Context) (int64, error) {
	return s.repo.CountAuthors(ctx)
}

func (s *AuthorService) GetAuthor(ctx context.Context, id uuid.UUID) (repository.Author, error) {
	author, err := s.repo.GetAuthor(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		return repository.Author{}, ErrAuthorNotFound
	}
	return author, err
}

func (s *AuthorService) UpdateAuthor(ctx context.Context, id uuid.UUID, name, bio *string) error {
	rows, err := s.repo.UpdateAuthor(ctx, repository.UpdateAuthorParams{
		ID:   id,
		Name: name,
		Bio:  bio,
	})
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrAuthorNotFound
	}
	return nil
}

func (s *AuthorService) DeleteAuthor(ctx context.Context, id uuid.UUID) error {
	rows, err := s.repo.DeleteAuthor(ctx, id)
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrAuthorNotFound
	}
	return nil
}

func (s *AuthorService) CreateBook(ctx context.Context, authorId uuid.UUID, name string) (repository.Book, error) {
	var book repository.Book

	err := s.withTx(ctx, func(q *repository.Queries) error {
		if _, err := q.GetAuthorForShare(ctx, authorId); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return ErrAuthorNotFound
			}
			return err
		}

		var err error
		book, err = q.CreateBook(ctx, repository.CreateBookParams{
			Name:     name,
			AuthorID: authorId,
		})
		return err
	})

	if err != nil {
		return repository.Book{}, err
	}
	return book, nil
}

func (s *AuthorService) ListBooks(ctx context.Context, authorId uuid.UUID) ([]repository.Book, error) {
	books, err := s.repo.ListBooks(ctx, authorId)
	if err != nil {
		return nil, err
	}
	if books == nil {
		books = []repository.Book{}
	}
	return books, nil
}

func (s *AuthorService) ListAuthorsAndBooks(ctx context.Context) ([]repository.ListAuthorsAndBooksRow, error) {
	booksAndAuthors, err := s.repo.ListAuthorsAndBooks(ctx)
	if err != nil {
		return nil, err
	}
	if booksAndAuthors == nil {
		booksAndAuthors = []repository.ListAuthorsAndBooksRow{}
	}
	return booksAndAuthors, nil
}
