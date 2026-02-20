#!make
include .env

migration:
	GOOSE_DBSTRING="postgresql://$(DBUSER):$(DBPASS)@localhost:5432/$(DBDATABASE)?sslmode=$(DBSSLMODE)"
	goose create -dir database/migrations $$NAME sql 

migrateup:
	GOOSE_DBSTRING="postgresql://$(DBUSER):$(DBPASS)@localhost:5432/$(DBDATABASE)?sslmode=$(DBSSLMODE)" goose up -dir database/migrations

migratedown:
	GOOSE_DBSTRING="postgresql://$(DBUSER):$(DBPASS)@localhost:5432/$(DBDATABASE)?sslmode=$(DBSSLMODE)" goose down -dir database/migrations

testcoverage:
	go test ./... -coverprofile

stresstest:
	cat scripts/script.js | docker run --network host --rm -i grafana/k6 run - 
