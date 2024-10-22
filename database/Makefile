# Default variables
# DB_PATH ?= users
# DB_NAME ?= db_user_service
# DB_HOST ?= localhost
# DB_PORT ?= 5432
# DB_USER ?= user_service_admin
# DB_PASSWORD ?= "comes-this-from-environment"
MIGRATION_STEPS ?= 1

# DB_PATH ?= galactic-sovereign
# DB_NAME ?= db_galactic_sovereign
# DB_HOST ?= localhost
# DB_PORT ?= 5432
# DB_USER ?= galactic_sovereign_admin
# DB_PASSWORD ?= W5QanVTbskfgdZ4ZPDu8LkHbA@6YS

#DB_PATH ?= stellar-sovereign
DB_NAME ?= db_stellar_dominion
DB_HOST ?= localhost
DB_PORT ?= 5432
DB_USER ?= stellar_dominion_admin
DB_PASSWORD ?= W5QanVTbskfgdZ4ZPDu8LkHbA@6YS

URL_ENCODED_PASSWORD=$(shell urlencode '${DB_PASSWORD}')

# https://stackoverflow.com/questions/6405127/how-do-i-specify-a-password-to-psql-non-interactively
connect:
	psql postgres://${DB_USER}:${URL_ENCODED_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}

.PHONY: migrate
migrate:
	migrate -path ${PWD}/${DB_PATH}/migrations -database postgres://${DB_USER}:${URL_ENCODED_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME} up

migrateOne:
	migrate -path ${PWD}/${DB_PATH}/migrations -database postgres://${DB_USER}:${URL_ENCODED_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable up ${MIGRATION_STEPS}

demigrate:
	migrate -path ${PWD}/${DB_PATH}/migrations -database postgres://${DB_USER}:${URL_ENCODED_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable down

demigrateOne:
	migrate -path ${PWD}/${DB_PATH}/migrations -database postgres://${DB_USER}:${URL_ENCODED_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable down ${MIGRATION_STEPS}
