# =========================================================
# SHELL CONFIG
# =========================================================

# Dùng PowerShell thay vì cmd.exe.
# Điều này phù hợp hơn khi bạn dùng GNU Make trên Windows.
SHELL := powershell.exe

# Không load PowerShell profile.
# Giúp command chạy nhanh và ổn định hơn.
.SHELLFLAGS := -NoProfile -Command


# =========================================================
# ENV
# =========================================================

# Load các biến môi trường từ file .env
include .env

# Export các biến trong Makefile thành environment variables
# để các command bên dưới có thể sử dụng.
export


# =========================================================
# DATABASE
# =========================================================

# Thư mục chứa migration.
MIGRATION_DIRS = internal/db/migrations
ENV_FILE =.env
# Connection string dùng cho migrate.
CONN_STRING = postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)

PROD_COMPOSE=docker-compose.prod.yml
NOAPP_COMPOSE=docker-compose.noapp.yml
DEV_COMPOSE=docker-compose.dev.yml

# =========================================================
# IMPORT DATABASE
# =========================================================

importdb:
	docker exec -i postgres-db psql -U root -d master-golang < ./backupdb-master-golang.sql


# =========================================================
# EXPORT DATABASE
# =========================================================

exportdb:
	docker exec -i postgres-db pg_dump -U root -d master-golang > ./backupdb-master-golang.sql

# =========================================================
# BUILD A BINARY FILE
# =========================================================
build:
	go build -o .\bin\myapp.exe .\cmd\api
# =========================================================
# RUN BINARY
# =========================================================
run-binary:
	.\bin\myapp.exe


# =========================================================
# RUN SERVER
# =========================================================

server:
	go run ./cmd/api


# =========================================================
# GENERATE SQLC
# =========================================================

sqlc:
	sqlc generate


# =========================================================
# CREATE MIGRATION
# =========================================================

migrate-create:
	migrate create -ext sql -dir $(MIGRATION_DIRS) -seq $(NAME)


# =========================================================
# MIGRATION UP
# =========================================================

migrate-up:
	migrate -path $(MIGRATION_DIRS) -database "$(CONN_STRING)" up


# =========================================================
# MIGRATION DOWN
# =========================================================

migrate-down:
	migrate -path $(MIGRATION_DIRS) -database "$(CONN_STRING)" down 1


# =========================================================
# MIGRATION DOWN N
# =========================================================

migrate-down-n:
	migrate -path $(MIGRATION_DIRS) -database "$(CONN_STRING)" down $(N)


# =========================================================
# MIGRATION FORCE
# =========================================================

migrate-force:
	migrate -path $(MIGRATION_DIRS) -database "$(CONN_STRING)" force $(VERSION)


# =========================================================
# MIGRATION DROP
# =========================================================

migrate-drop:
	migrate -path $(MIGRATION_DIRS) -database "$(CONN_STRING)" drop


# =========================================================
# MIGRATION GOTO
# =========================================================

migrate-goto:
	migrate -path $(MIGRATION_DIRS) -database "$(CONN_STRING)" goto $(VERSION)

# =========================================================
# DOCKER NOAPP_COMPOSE
# =========================================================
noapp:
	docker compose -f $(NOAPP_COMPOSE) down	
	docker compose -f $(NOAPP_COMPOSE) --env-file $(ENV_FILE) up -d --build

stop-noapp: 
	docker compose -f $(NOAPP_COMPOSE) down		

# =========================================================
# DOCKER DEV
# =========================================================
dev:
	docker compose -f $(DEV_COMPOSE) down
	docker compose -f $(DEV_COMPOSE) --env-file $(ENV_FILE) up --build

stop-dev:
	docker compose -f $(DEV_COMPOSE) down	
# =========================================================
# DOCKER PROD_COMPOSE
# =========================================================
prod:
	docker compose -f $(PROD_COMPOSE) down	
	docker compose -f $(PROD_COMPOSE) --env-file $(ENV_FILE) up -d --build

stop-prod: 
	docker compose -f $(PROD_COMPOSE) down	

logs-prod:
	docker compose -f $(PROD_COMPOSE) logs -f --tail=100

bash:
	docker exec -it golang-api /bin/sh
# =========================================================
# PHONY
# =========================================================

.PHONY: importdb exportdb server migrate-create migrate-up migrate-down migrate-force migrate-drop migrate-down-n migrate-goto sqlc build run-binary noapp stop-noapp dev prod stop-prod logs-prod bash