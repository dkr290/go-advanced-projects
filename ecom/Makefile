run-local:
	go run cmd/main.go

start-db:
	docker start mysql

stop-db:
	docker stop mysql

test:
	@go test -v ./...


migration:
	@migrate create -ext sql -dir cmd/migrate/migrations $(filter-out $@,$(MAKECMDGOALS))

migrate-up:
	@go run cmd/migrate/main.go -up

migrate-down:
	@go run cmd/migrate/main.go -down
