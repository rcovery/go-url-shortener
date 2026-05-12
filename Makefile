#!make
include .env

migration:
	GOOSE_DBSTRING="postgresql://$(DBUSER):$(DBPASS)@localhost:5432/$(DBDATABASE)?sslmode=$(DBSSLMODE)"
	goose create -dir internal/infra/postgres/migrations $$NAME sql 

migrateup:
	GOOSE_DBSTRING="postgresql://$(DBUSER):$(DBPASS)@localhost:5432/$(DBDATABASE)?sslmode=$(DBSSLMODE)" goose up -dir internal/infra/postgres/migrations

migratedown:
	GOOSE_DBSTRING="postgresql://$(DBUSER):$(DBPASS)@localhost:5432/$(DBDATABASE)?sslmode=$(DBSSLMODE)" goose down -dir internal/infra/postgres/migrations

testcoverage:
	go test ./... -cover profile.cov

stresstest:
	cat scripts/script.js | docker run --network host --rm -i grafana/k6 run - 
