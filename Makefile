#!make
include .env

migration:
	GOOSE_DBSTRING="postgresql://$(DBUSER):$(DBPASS)@localhost:5432/$(DBDATABASE)?sslmode=disable"
	goose create -dir database/migrations $$NAME sql 

migrateup:
	GOOSE_DBSTRING="postgresql://$(DBUSER):$(DBPASS)@localhost:5432/$(DBDATABASE)?sslmode=disable" goose up -dir database/migrations

migratedown:
	GOOSE_DBSTRING="postgresql://$(DBUSER):$(DBPASS)@localhost:5432/$(DBDATABASE)?sslmode=disable" goose down -dir database/migrations

testcoverage:
	go test ./... -coverprofile
