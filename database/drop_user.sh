#!/bin/bash

DB_PATH="galactic-sovereign"
if [[ $# -ge 1 ]]; then
  DB_PATH="$1"
fi

DB_HOST=${DATABASE_HOST:-localhost}
DB_PORT=${DATABASE_PORT:-5432}
DB_USER=${DATABASE_USER:-postgres}

psql -h ${DB_HOST} -p ${DB_PORT} -U ${DB_USER} -f ${DB_PATH}/db_user_drop.sql
