#!make
include .env

migration:
	migrate create -ext sql -dir ./database/migrations -seq $$NAME

migrateup:
	migrate -path ./database/migrations -database "postgresql://$(DBUSER):$(DBPASS)@localhost:5432/$(DBDATABASE)?sslmode=disable" up

migratedown:
	migrate -path ./database/migrations -database "postgresql://$(DBUSER):$(DBPASS)@localhost:5432/$(DBDATABASE)?sslmode=disable" down

testcoverage:
	go test ./... -coverprofile
