postgres:
	docker run --name popo -p 5435:5432 -e POSTGRES_PASSWORD=pricolush1 -d postgres:16-alpine

createdb:
	docker exec -it popo createdb --username=postgres simple_bank

dropdb:
	docker exec -it popo dropdb --username=postgres simple_bank
migrateup:
	migrate -path db/migration -database "postgresql://postgres:pricolush1@localhost:5435/simple_bank?sslmode=disable" up

migrateup1:
	migrate -path db/migration -database "postgresql://postgres:pricolush1@localhost:5435/simple_bank?sslmode=disable" up 1

migratedown:
	migrate -path db/migration -database "postgresql://postgres:pricolush1@localhost:5435/simple_bank?sslmode=disable" down

migratedown1:
	migrate -path db/migration -database "postgresql://postgres:pricolush1@localhost:5435/simple_bank?sslmode=disable" down 1
sqlc:
	sqlc generate
mock:
	mockgen -package mockdb -destination db/mock/store.go Backend/db/sqlc Store

.PHONY:postgres createdb dropdb migrateup migratedown migrateup1 migratedown1 sqlc mock