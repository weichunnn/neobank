postgres:
	docker run --name postgres12 --network bank-network -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -dp 5432:5432 postgres:12-alpine

createdb:
	docker exec postgres12 createdb --username=root --owner=root simple_bank

dropdb:
	docker exec postgres12 dropdb simple_bank

migrateup:
	migrate -path db/migration -database "postgresql://root:06PdTdEIWudH9RhiGA@neo-bank.chjh1hdu8pa4.ap-southeast-1.rds.amazonaws.com:5432/simple_bank" -verbose up

migrateup1:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down

migratedown1:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down 1

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/weichunnn/neobank/db/sqlc Store

.PHONY: postgres createdb dropdb migrateup migrateup1 migratedown migratedown1 sqlc test server mock