.PHONY: run
run:
	go run cmd/main.go

.PHONY: test_db
test_db:
	docker run -d --rm -e POSTGRES_PASSWORD=qwerty123 -e POSTGRES_USER=postgres -e POSTGRES_DB=postgres --name=pg -p 5432:5432 postgres

.PHONY: create_migrations
create_migrations:
	migrate create -ext sql -dir migrations create_users_table

.PHONY: migrate_up
migrate_up:
	migrate -path ./migrations -database "postgres://postgres:qwerty123@localhost:5432/postgres?sslmode=disable" up

.PHONY: migrate_down
migrate_down:
	migrate -path ./migrations -database "postgres://postgres:qwerty123@localhost:5432/postgres?sslmode=disable" down

.PHONY: swag
swag:
	swag init --parseDependency --parseInternal -g cmd/main.go

.PHONY: mockGen
mockGen:
	mockgen -source="./handler.go" -destination="./mocks/mock_service.go" -package=mocks
