## Article

https://fulim.tech/blog/rest-api-with-echo-sqlc-migrate-postgres/

## How to Run Application

```sh
make run
```

## API Structure

```sh
POST /authors
GET /authors
GET /authors/len
GET /authors/:id
PATCH /authors/:id
DELETE /authors/:id
POST /authors/:id/books
GET /authors/:id/books
GET /authors-and-books
```

## Curl Command

```sh
# Create Author
curl -s localhost:1323/authors \
	-H 'Content-Type: application/json' \
	-d '{"name": "fulim", "bio": "An Author of this tutorial"}' | jq

# List Authors
curl -s localhost:1323/authors | jq

# Count Authors
curl -s localhost:1323/authors/len | jq

# Get Author
curl -s localhost:1323/authors/author-id | jq

# Update Author
curl -s -i -X PATCH localhost:1323/authors/author-id \
	-H 'Content-Type: application/json' \
	-d '{"name": "Fu Lim Wong", "bio": "New Author"}'

# Delete Author
curl -s -i -X DELETE localhost:1323/authors/author-id

# Create Book
curl -s localhost:1323/authors/author-id/books \
	-H 'Content-Type: application/json' \
	-d '{"name": "Go Web Tutorial"}' | jq

# List Book
curl -s localhost:1323/authors/author-id/books | jq

# List Authors and Books
curl -s localhost:1323/authors-and-books | jq

# Create/Update Author / Create Book - invalid request body
curl -s -i localhost:1323/authors \
	-H 'Content-Type: application/json' \
	-d '{"na}'
curl -s -i -X PATCH localhost:1323/authors/98843538-6cd4-4811-b16e-cb58dbb30218 \
	-H 'Content-Type: application/json' \
	-d '{"na}'
curl -s -i localhost:1323/authors/98843538-6cd4-4811-b16e-cb58dbb30218/books \
	-H 'Content-Type: application/json' \
	-d '{"na}'

# Get/Update/Delete Author / Create/List Book - invalid author id
curl -s -i localhost:1323/authors/author-id
curl -s -i -X PATCH localhost:1323/authors/author-id \
	-H 'Content-Type: application/json' \
	-d '{"name": "Fu Lim Wong", "bio": "New Author"}'
curl -s -i -X DELETE localhost:1323/authors/author-id
curl -s -i localhost:1323/authors/author-id/books \
	-H 'Content-Type: application/json' \
	-d '{"name": "Go Web Tutorial"}'
curl -s -i localhost:1323/authors/author-id/books

# Get/Update/Delete Author / Create Book - author not found
curl -s -i localhost:1323/authors/98843538-6cd4-4811-b16e-cb58dbb30218
curl -s -i localhost:1323/authors/98843538-6cd4-4811-b16e-cb58dbb30218/books \
	-H 'Content-Type: application/json' \
	-d '{"name": "Go Web Tutorial"}'
curl -s -i localhost:1323/authors/98843538-6cd4-4811-b16e-cb58dbb30218/books \
	-H 'Content-Type: application/json' \
	-d '{"name": "Go Web Tutorial"}'
curl -s -i localhost:1323/authors/98843538-6cd4-4811-b16e-cb58dbb30218/books \
	-H 'Content-Type: application/json' \
	-d '{"name": "Go Web Tutorial"}'
```

## Database Structure

```sql
CREATE TABLE author (
	id UUID PRIMARY KEY,
	name VARCHAR NOT NULL,
	bio TEXT NOT NULL DEFAULT '',
	created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE book (
	id SERIAL PRIMARY KEY,
	name VARCHAR NOT NULL,
	author_id UUID NOT NULL REFERENCES author(id)
);
```

## Tool Installation

```sh
brew install golang-migrate
brew install sqlc
```

## Module Package Installation

```sh
go get github.com/labstack/echo/v5
go get github.com/jackc/pgx/v5
go get github.com/joho/godotenv
go get github.com/google/uuid
```

## Postgres in Docker Setup

```sh
docker run -d -p 5432:5432 -e POSTGRES_PASSWORD=111111 -e POSTGRES_DB=author postgres
```

## Create Migration file

```sh
migrate create -ext sql -seq -digits 5 -dir database/migrations initial
```

## Migrate

```sh
migrate -path database/migrations -database postgresql://postgres:111111@localhost:5432/author?sslmode=disable up
```
