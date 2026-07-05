DB_URL ?= mysql://root:secret@tcp(127.0.0.1:3306)/coffee_pos
MIGRATIONS_PATH=migrations

.PHONY: migrate-up migrate-down migrate-down-all migrate-version migrate-create run tidy

migrate-up:
	migrate -database "$(DB_URL)" -path $(MIGRATIONS_PATH) up

migrate-down:
	migrate -database "$(DB_URL)" -path $(MIGRATIONS_PATH) down 1

migrate-down-all:
	migrate -database "$(DB_URL)" -path $(MIGRATIONS_PATH) down

migrate-version:
	migrate -database "$(DB_URL)" -path $(MIGRATIONS_PATH) version

migrate-create:
	migrate create -ext sql -dir $(MIGRATIONS_PATH) -seq $(name)

run:
	go run cmd/api/main.go

tidy:
	go mod tidy
