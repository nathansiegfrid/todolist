# Load ".env" file if exists.
ifneq (,$(wildcard ./.env))
	include .env
	export
endif

# GO COMMANDS
run:
	air || go run .
update:
	go get -u -t ./...
	go mod tidy

# DB MIGRATION COMMANDS
DB_STRING=postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=$(POSTGRES_SSL_MODE)
db-up:
	goose -dir migrations postgres "$(DB_STRING)" up
db-down:
	goose -dir migrations postgres "$(DB_STRING)" down
db-redo:
	goose -dir migrations postgres "$(DB_STRING)" redo

# DOCKER COMMANDS
docker-build:
	docker build -t nathansiegfrid/todolist .
docker-push:
	docker push nathansiegfrid/todolist
docker-up:
	docker compose up -d --build
docker-down:
	docker compose down
