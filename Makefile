# Variables
APP_NAME=main

DB_USER ?= Njoguu
DB_PASSWORD ?= alannjoguu
DB_NAME ?= LoyaltyPointsDB
DB_ADDRESS ?= localhost
DB_PORT ?= 5433
DB_OWNER ?= Njoguu

DB_CONN = postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable

CONTAINER_NAME ?= postgres_db
IMAGE ?= postgres:15.3-alpine3.18

MIGRATION_EXT ?= "sql"
MIGRATION_DIR ?= "db/migrations"


#================================
#== DOCKER Targets
#================================
COMPOSE := @docker compose

postgres:
	docker run --name $(CONTAINER_NAME) -p $(DB_PORT):5432 -e POSTGRES_USER=$(DB_USER) -e POSTGRES_PASSWORD=$(DB_PASSWORD) -d $(IMAGE)

createdb:
	docker exec -it $(CONTAINER_NAME) createdb --username=$(DB_USER) --owner=$(DB_OWNER) $(DB_NAME)
	
dropdb: 
	docker exec -it $(CONTAINER_NAME) dropdb -U $(DB_OWNER) $(DB_NAME) 

start-services: 
	$(COMPOSE) start db
	$(COMPOSE) start cache
	$(COMPOSE) start api

stop-services:
	$(COMPOSE) stop api
	$(COMPOSE) stop cache
	$(COMPOSE) stop db 

restart-services: stop-services start-services

dcu:
	$(COMPOSE) up -d --build

dcd:
	$(COMPOSE) down

#================================
#== GOLANG ENVIRONMENT Targets
#================================
GO := @go

install:
	${GO} get .

tidy:
	${GO} mod tidy

start:
	${GO} run main.go
	
build:
	${GO} build -o ${APP_NAME} .

migratedb:
	${GO} run migrate/migrate.go

#================================
#== DB MIGRATION Targets
#================================
create-migrations: 
	migrate create -ext $(MIGRATION_EXT) -dir $(MIGRATION_DIR) -seq $(MIGRATION_NAME)

migrateup:
	migrate -path $(MIGRATION_DIR) -database "$(DB_CONN)" -verbose up

migratedown:
	migrate -path $(MIGRATION_DIR) -database "$(DB_CONN)" -verbose down

migrateup1:
	migrate -path $(MIGRATION_DIR) -database "$(DB_CONN)" -verbose up 1

migratedown1:
	migrate -path $(MIGRATION_DIR) -database "$(DB_CONN)" -verbose down 1

.PHONY: createdb postgres dropdb migrateup migratedown create-migrations
