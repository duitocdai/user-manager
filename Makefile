include .env
export

MIGRATION_DIRS = internal/db/migrations
CONN_STRING = postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)
#import database
importdb:
	docker exec -i postgres-db psql -U root -d master-golang < ./backupdb-master-golang.sql
#export database
exportdb:
	docker exec -i postgres-db pg_dump -U root -d master-golang > ./backupdb-master-golang.sql
#run server
server:
	cd cmd/api && go run . 
#Generate sqlc
sqlc:
	sqlc generate
#create a new migration
migrate-create:
	migrate create -ext sql -dir $(MIGRATION_DIRS) -seq $(NAME)
#run all pending migration
migrate-up:
	migrate -path $(MIGRATION_DIRS) -database "$(CONN_STRING)" up
migrate-down:
	migrate -path $(MIGRATION_DIRS) -database "$(CONN_STRING)" down 1	
migrate-down-n:
	migrate -path $(MIGRATION_DIRS) -database "$(CONN_STRING)" down $(N)
migrate-force:
	migrate -path $(MIGRATION_DIRS) -database "$(CONN_STRING)" force $(VERSION)
migrate-drop:
	migrate -path $(MIGRATION_DIRS) -database "$(CONN_STRING)" drop
migrate-goto:
	migrate -path $(MIGRATION_DIRS) -database "$(CONN_STRING)" goto $(VERSION)
.PHONY: importdb exportdb server migrate-create migrate-up migrate-down migrate-force migrate-drop migrate-down-n migrate-goto sqlc