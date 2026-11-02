DATABASE_URL ?= postgres://postgres:postgres@127.0.0.1:5434/rent_cars_go_dev?sslmode=disable
export DATABASE_URL

.PHONY: dev db-up db-stop test check generate
dev: db-up
	go run ./cmd/api
db-up:
	docker compose up -d --wait postgres
db-stop:
	docker compose stop postgres
test:
	go test -race -coverpkg=./... -coverprofile=coverage.out ./...
check:
	@test -z "$$(gofmt -l .)" || { gofmt -l .; exit 1; }
	go vet ./...
	$(MAKE) test
generate:
	sqlc generate
