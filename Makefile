# Load ".env" file if exists.
ifneq (,$(wildcard ./.env))
	include .env
	export
endif

# GO COMMANDS
run:
	go run .
update:
	go get -u -t ./...
	go mod tidy

# DB MIGRATION COMMANDS
db-up:
	goose -dir migrations postgres "$(POSTGRES_URL)" up
db-down:
	goose -dir migrations postgres "$(POSTGRES_URL)" down
db-redo:
	goose -dir migrations postgres "$(POSTGRES_URL)" redo

# DOCKER COMMANDS
docker-build:
	docker build -t nathansiegfrid/todolist .
docker-push:
	docker push nathansiegfrid/todolist
docker-up:
	docker compose up -d --build
docker-down:
	docker compose down
