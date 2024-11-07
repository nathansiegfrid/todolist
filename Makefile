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

# CODE GEN COMMANDS
sqlc:
	sqlc generate

# DEV TOOLS INSTALL COMMANDS
install-goose:
	go install github.com/pressly/goose/v3/cmd/goose@latest
install-sqlc:
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

# DB MIGRATION COMMANDS
db-up:
	goose -dir migrations postgres "$(POSTGRES_URL)" up
db-down:
	goose -dir migrations postgres "$(POSTGRES_URL)" down
db-redo:
	goose -dir migrations postgres "$(POSTGRES_URL)" redo

# DOCKER COMMANDS
docker-build:
	docker build -t siegfrid/todolist .
docker-push:
	docker push siegfrid/todolist
