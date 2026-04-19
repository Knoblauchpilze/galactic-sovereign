#!/bin/bash

if [[ $# -lt 1 ]]; then
  echo "Usage: drops the database defined in the provided folder"
  echo "Examples:"
  echo "./drop_database.sh /path/to/database/folder"
  exit 1
fi

DB_PATH=$1
DB_HOST=${DATABASE_HOST:-localhost}
DB_PORT=${DATABASE_PORT:-5432}
DB_USER=${DATABASE_USER:-postgres}

psql -h ${DB_HOST} -p ${DB_PORT} -U ${DB_USER} -f ${DB_PATH}/db_drop.sql
