DB_URL=postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable

postgres:
	docker run --name postgres12 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -dp 5432:5432 postgres:12-alpine

createdb:
	docker exec postgres12 createdb --username=root --owner=root simple_bank

dropdb:
	docker exec postgres12 dropdb simple_bank

migrateup:
	migrate -path db/migration -database "$(DB_URL)" -verbose up

migrateup1:
	migrate -path db/migration -database "$(DB_URL)" -verbose up 1

migratedown:
	migrate -path db/migration -database "$(DB_URL)" -verbose down

migratedown1:
	migrate -path db/migration -database "$(DB_URL)" -verbose down 1

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/weichunnn/neobank/db/sqlc Store

proto: 
	rm -f pb/*.go
	rm -rf docs/swagger/*.swagger.json

	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
	--go-grpc_out=pb --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=pb --grpc-gateway_opt=paths=source_relative \
	--openapiv2_out=docs/swagger --openapiv2_opt=allow_merge=true,merge_file_name=neobank \
	proto/*.proto
	statik -src=./docs/swagger -dest=./docs

evans: 
	evans --host localhost --port 9000 -r repl

redis:
	docker run --name redis -dp 6379:6379 redis:7-alpine

.PHONY: postgres createdb dropdb migrateup migrateup1 migratedown migratedown1 sqlc test server mock proto evans redis