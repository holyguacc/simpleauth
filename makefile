postgres:
	docker run --name simpleauth -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:14-alpine
psqlexec :
	docker exec -ti simpleauth psql -U root

migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simpleauth?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simpleauth?sslmode=disable" -verbose down

sqlc:
	sqlc generate

migratecreateinit:
	migrate create -ext sql -dir db/migration -seq create_verification_table


.PHONY: postgres  migrateup migratedown sqlc psqlexec