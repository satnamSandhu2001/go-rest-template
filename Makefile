.PHONY: run migrate rollback version auto_migrate build init_air watch

BINARY_NAME=app
BUILD_DIR=dist

run:
	@go run .

migrate:
	@go run . --migrate

rollback:
	@go run . --rollback=1

version:
	@go run . --version

auto_migrate:
	@MIGRATE_ON_START=true go run .

build:
	@mkdir -p $(BUILD_DIR)
	@OS=$(shell go env GOOS); \
	ARCH=$(shell go env GOARCH); \
	EXT=$$( [ "$$OS" = "windows" ] && echo ".exe" ); \
	OUT_NAME="$(BUILD_DIR)/$(BINARY_NAME)-$$OS-$$ARCH$$EXT"; \
	echo "Building for $$OS/$$ARCH..."; \
	GOOS=$$OS GOARCH=$$ARCH go build -o $$OUT_NAME main.go; \
	echo "Built binary at $$OUT_NAME"

# Usage: make create-migration name=create_
create-migration:
	@if [ -z "$(name)" ]; then \
		echo "	Missing migration name!"; \
		echo "	Usage: make create-migration name=create_users_table"; \
		exit 1; \
	fi; \
	go run cmd/scripts/create_migration.go $(name)
