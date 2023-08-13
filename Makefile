DB_URL=postgresql://root:secret@localhost:5432/ross_local?sslmode=disable

run:
	go run cmd/serverd/main.go

build:
	docker build . -t ross-api

sqlc:
	sqlc generate

migratecreate:
	migrate create -ext sql -dir pkg/db/migrations -seq $(name)

migrateup:
	migrate -database "$(DB_URL)" -path pkg/db/migration -verbose up

migrateup1:
	migrate -database "$(DB_URL)" -path pkg/db/migration -verbose up 1

migratedown:
	migrate -database "$(DB_URL)" -path pkg/db/migration -verbose down

migratedown1:
	migrate -database "$(DB_URL)" -path pkg/db/migration -verbose down 1

.PHONY: run build sqlc migrate-create migrate-up migrate-down