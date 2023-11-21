DB_URL=postgres://root:secret@localhost:5432/coin?sslmode=disable

migrateup:
	migrate -path db/migrations -database "$(DB_URL)" -verbose up

migratedown:
	migrate -path db/migrations -database "$(DB_URL)" -verbose down

new_migration:
	migrate create -ext sql -dir db/migrations -seq $(table)

sqlc:
	sqlc -f db/sqlc.yaml generate

mock:
	mockgen --package dbmock --destination db/mock/store.go db Store

build:
	docker-compose -p coin -f docker/docker-compose.yml build app

start:
	docker-compose -p coin -f docker/docker-compose.yml up -d

down:
	docker-compose -p coin -f docker/docker-compose.yml down

test:
	go test -v -cover api
	go test -v -cover db
	go test -v -cover token
	go test -v -cover utils

proto:
	rm -f pb/*.go
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
    --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
    proto/*.proto

.PHONY: migrateup migratedown new_migration sqlc mock build start down test proto