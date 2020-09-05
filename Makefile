.PHONY: build adminer migrate postgres

postgres:
	docker run --rm -ti --network host -e POSTGRES_PASSWORD=123 postgres

adminer:
	docker run --rm -ti --network host adminer

migrate:
	migrate -path migrations \
			-database postgres://localhost/starline?sslmode=disable up

migrate-down:
	migrate -path migrations \
		    -database postgres://localhost/starline?sslmode=disable down


build:
	go build -v ./cmd/apiserver

.PHONY: test 
test:
	go test -v -race -timeout 30s ./...	

.DEFAULT_GOAL := build	