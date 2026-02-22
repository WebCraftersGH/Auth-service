include .env

export $(shell sed 's/=.*//' .env)

DB_URL=postgres://$(DB_USER):$(DB_PASS)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable


MIGRATIONS_DIR=./migrations

create-migration:
	@if [ -z "$(name)" ]; then echo "Usage: make create-migration name=add_users"; exit 1; fi
	migrate create -ext sql -dir $(MIGRATIONS_DIR) -seq $(name)

migrate-up:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" up

migrate-test-up:
	migrate -path $(MIGRATIONS_DIR) -database "$(TEST_DB_URL)" up

migrate-down-1:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" down 1

migrate-down-all:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" down


migrate-version:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" version

migrate-up-n:

	@if [ -z "$(n)" ]; then echo "Usage: make migrate-up-n n=1"; exit 1; fi
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" up $(n)


migrate-force:
	@if [ -z "$(v)" ]; then echo "Usage: make migrate-force v=3"; exit 1; fi

	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" force $(v)
