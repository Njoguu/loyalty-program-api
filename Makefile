# Variables
APP_NAME=loyalty-program-api

#================================
#== DOCKER Targets
#================================
COMPOSE := @docker-compose

dcb:
	${COMPOSE} build

dcu:
	${COMPOSE} up -d build

dcd:
	${COMPOSE} down

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
