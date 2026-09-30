DATABASE_URL = postgres://postgres:111111@localhost:5432/author?sslmode=disable

.PHONY: postgres migrate sqlc run

postgres:
	@echo "==> 1/4 starting postgres"
	@docker start author-postgres 2>/dev/null || docker run -d --name author-postgres \
		-p 5432:5432 -e POSTGRES_PASSWORD=111111 -e POSTGRES_DB=author postgres
	@sleep 3

migrate:
	@echo "==> 2/4 running migrations"
	migrate -path ./database/migrations -database "$(DATABASE_URL)" up

sqlc:
	@echo "==> 3/4 generating sqlc code"
	sqlc generate

run: postgres migrate sqlc
	@echo "==> 4/4 starting server on :1323"
	go run ./cmd/server
