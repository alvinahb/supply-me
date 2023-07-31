postgres:
	docker run --name postgres15 -p 5432:5432 -e POSTGRES_DB=supply-me-test -e POSTGRES_USER=root -e POSTGRES_PASSWORD=supply-me-2023 -d postgres:15-alpine

createdb:
	docker exec -it postgres15 createdb --username=root --owner=root supply-me-test

dropdb:
	docker exec -it postgres15 dropdb supply-me-test

migrateup:
	migrate -path db/migration -database "postgresql://root:supply-me-2023@localhost:5432/supply-me-test?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:supply-me-2023@localhost:5432/supply-me-test?sslmode=disable" -verbose down

sqlc:
	docker run --rm -v "C:\Users\amste\Documents\projects\supply-me:/src" -w /src kjconroy/sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

.PHONY: postgres createdb dropdb migrateup migratedown sqlc test server
