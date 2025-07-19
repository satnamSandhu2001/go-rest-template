# Go REST API Template

A minimal Go HTTP Rest API server starter using Chi router with PostgreSQL(pgx driver), sqlc for SQL queries and custom DB migration management.

## Instructions

- **Set environment variables:**

  copy `.env.example` to `.env` and update values

  ```bash
  GO_ENV="production"
  DEBUG=false
  PORT=8080
  MIGRATE_ON_START=true
  DB_URL="postgres://postgres:postgres@localhost:5432/myDb?sslmode=disable"
  API_URL="http://localhost:8080"
  CLIENT_URL="http://localhost:3000"
  COOKIE_DOMAIN="" # leave empty for localhost else ".domain.com"
  COOKIE_AGE_HOURS=48 # 2 days
  JWT_TOKEN="your secret token"
  ```

- **Create new migration files:**

  ```bash
  # Using make:
  make create-migration name=create_users_table

  # Manual:
  go run cmd/scripts/create_migration.go create_users_table
  ```

- **Run the server:**

  ```bash
  # Using make:
  make run

  # Manual:
  go run .
  ```

- **Apply migrations:**

  ```bash
  # Using make:
  make migrate

  # Manual:
  go run . --migrate
  ```

- **Rollback last (n) migrations:**

  ```bash
  # Using make:
  make rollback

  # Manual:
  go run . --rollback=1
  ```

- **Show current migration version:**

  ```bash
  # Using make:
  make version

  # Manual:
  go run . --version
  ```

- **Auto-run migrations on server start (optional):**

  ```bash
  # Using make:
  make auto_migrate

  # Manual:
  MIGRATE_ON_START=true go run .
  ```

- **Build binary:**

  ```bash
  # Using make:
  make build

  # Manual:
  mkdir -p dist
  go build -o dist/app.exe main.go
  ```
