DB_URL=postgres://root:secret@localhost:5432/coin?sslmode=disable

migrateup:
	migrate -path db/migrations -database "$(DB_URL)" -verbose up

migratedown:
	migrate -path db/migrations -database "$(DB_URL)" -verbose down

new_migration:
	migrate create -ext sql -dir db/migrations -seq $(table)

sqlc:
	sqlc -f db/sqlc.yaml generate
	
start:
	docker-compose -p coin -f docker/docker-compose.yml up -d

down:
	docker-compose -p coin -f docker/docker-compose.yml down

.PHONY: migrateup migratedown new_migration