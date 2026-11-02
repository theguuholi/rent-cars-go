# Rent Cars Go

A Go implementation of the Rent Cars vehicle category catalog. Clients can use the same category request as the Java API. The Go application listens on port 8081 and owns a separate PostgreSQL 17 database on port 5434, so both implementations can run side by side.

## Run locally

Install the Go version in `.tool-versions`, Docker with Compose v2, and Make. With asdf and the golang plugin installed, run `asdf install` in this directory. Start Docker, then run:

```sh
make dev
```

Compose waits for PostgreSQL. The API applies its embedded goose migrations, seeds the initial catalog once, and serves requests. Stop the API with Ctrl+C. Use `make db-stop` to stop PostgreSQL and preserve its data. To run against another database, set `DATABASE_URL` and execute `go run ./cmd/api`. Use `HTTP_ADDR` to override the listen address.

## Try the catalog

```sh
curl -i -H 'Accept: application/json' http://localhost:8081/categories
```

`GET /categories` returns a JSON array ordered by name. A fresh database has HATCHBACK, PICKUP, SEDAN and SUV. An empty catalog returns HTTP 200 with `[]`. Each item has an uppercase `name` and a `description` (a string or JSON null). For example:

```json
{"name":"HATCHBACK","description":"Compact car with a rear hatch."}
```

- Open [Swagger UI](http://localhost:8081/swagger-ui.html) and try the endpoint in your browser. Its assets are served by the application.
- Read the [OpenAPI contract](api/openapi.json), also served at [the live specification](http://localhost:8081/v3/api-docs).
- In VS Code with REST Client installed, open [requests/categories.http](requests/categories.http) and select **Send Request**. Change the base URL to port 8080 to exercise the Java implementation.

## Development

| Path | Purpose |
| --- | --- |
| `cmd/api` | Wire the PostgreSQL pool, migrations and standard-library HTTP server |
| `internal/category` | Convert query results into the public JSON response |
| `internal/db/queries` | Handwritten SQL compiled by sqlc for pgx |
| `internal/db/store` | Generated query methods and database types |
| `migrations` | Versioned PostgreSQL schema and catalog seeds embedded with goose |
| `integration` | Tests against disposable PostgreSQL 17 containers |

After changing a query or schema, install sqlc 1.31.1 and run `make generate`. Commit the generated files alongside the SQL. Add a new numbered migration for future schema changes.

## Verify

```sh
make check
go tool cover -html=coverage.out -o coverage.html
go tool govulncheck ./...
```

Docker must be running for integration tests. The suite covers successful, empty and failed catalog reads, null descriptions, documentation routes, migration replay and database uniqueness. GitHub Actions checks dependency and sqlc drift, formatting, static analysis, tests with the race detector, and reachable vulnerabilities. Download the HTML report from the **go-coverage** workflow artifact.
