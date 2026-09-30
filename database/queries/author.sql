-- name: CreateAuthor :one
INSERT INTO author (id, name, bio)
VALUES ($1, $2, $3)
RETURNING id, name, bio, created_at;

-- name: ListAuthors :many
SELECT * FROM author;

-- name: CountAuthors :one
SELECT COUNT(*) FROM author;

-- name: GetAuthor :one
SELECT *
FROM author
WHERE id = $1;

-- name: GetAuthorForShare :one
-- FOR KEY SHARE locks the author row against DELETE for the rest of the transaction
-- without blocking other book inserts for the same author
SELECT *
FROM author
WHERE id = $1
FOR KEY SHARE;

-- name: UpdateAuthor :execrows
-- sqlc.narg = nullable arg, so NULL means not supply
-- COALESCE = keeps the existing column value
UPDATE author
SET name = COALESCE(sqlc.narg('name'), name),
    bio = COALESCE(sqlc.narg('bio'), bio)
WHERE id = sqlc.arg('id');

-- name: DeleteAuthor :execrows
DELETE FROM author
WHERE id = $1;

-- name: CreateBook :one
INSERT INTO book (name, author_id)
VALUES ($1, $2)
RETURNING id, name, author_id;

-- name: ListBooks :many
SELECT * FROM book
WHERE author_id = $1;

-- name: ListAuthorsAndBooks :many
SELECT sqlc.embed(author), book.id AS book_id, book.name AS book_name
FROM author
LEFT JOIN book ON author.id = book.author_id;
-- the purpose using LEFT JOIN is to not left off author with no books
