run_pgbouncer:
	@go run main.go

run_odyssey:
	@go run main.go -dsn "postgres://postgres@localhost:6532?dbname=db&sslmode=disable"

PHONY: run_pgbouncer run_odyssey
