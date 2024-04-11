#!/bin/sh

DB_HOST="localhost"
DB_NAME="db_user_service"
DB_USER="user_service_user"

# https://stackoverflow.com/questions/4018503/is-there-a-date-time-format-that-does-not-have-spaces
DATE=$(date '+%F-%T' | tr ':' '_')

# https://www.postgresql.org/docs/current/app-pgdump.html
pg_dump -h ${DB_HOST} -U ${DB_USER} -F c -f "${DB_NAME}_dump_${DATE}.bck" ${DB_NAME}
