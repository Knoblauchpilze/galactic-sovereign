#!/bin/sh

if [ "$#" -ne 1 ]; then
  echo "usage: database-restore.sh [database_dump_file]"
  exit 1
fi

DB_HOST=${DATABASE_HOST:-localhost}
DB_NAME=${DATABASE_NAME:-db_user_service}
DB_USER=${DATABASE_USER:-user_service_admin}

DB_DUMP=$1

# https://www.postgresql.org/docs/current/app-pgrestore.html
pg_restore -h ${DB_HOST} -U ${DB_USER} -cWe --if-exists -d ${DB_NAME} ${DB_DUMP}
