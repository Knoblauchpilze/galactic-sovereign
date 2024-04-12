#!/bin/bash

DB_HOST=${DATABASE_HOST:-localhost}
DB_PORT=${DATABASE_PORT:-5432}
DB_USER=${DATABASE_USER:-postgres}

# Password should be provided as environment variables.
if [[ -z "${ADMIN_PASSWORD}" ]]; then
  echo "DB admin password is not provided as environment variable, aborting"
  exit 1
fi
if [[ -z "${MANAGER_PASSWORD}" ]]; then
  echo "DB manager password is not provided as environment variable, aborting"
  exit 1
fi
if [[ -z "${USER_PASSWORD}" ]]; then
  echo "DB user password is not provided as environment variable, aborting"
  exit 1
fi

# https://stackoverflow.com/questions/8208181/create-postgres-database-using-batch-file-with-template-encoding-owner-and
psql \
  -v admin_password=${ADMIN_PASSWORD} \
  -v manager_password=${MANAGER_PASSWORD} \
  -v user_password=${USER_PASSWORD} \
  -h ${DB_HOST} \
  -p ${DB_PORT} \
  -U ${DB_USER} \
  -f db_user_create.sql
