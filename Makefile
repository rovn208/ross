DB_URL=postgresql://root:secret@localhost:5432/ross_local?sslmode=disable

run:
	DATABASE_URL=$(DB_URL) go run cmd/serverd/main.go

build:
	docker build . -t ross-api

db:
	docker-compose up -d db

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

swagger:
	swag init -d pkg/api -o pkg/docs -g ../../cmd/serverd/main.go

test:
	go test ./... -cover

.PHONY: run build sqlc migratecreate migrateup migrateup1 migratedown migratedown1 swagger test
