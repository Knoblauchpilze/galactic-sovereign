#!/bin/sh

DB_HOST=${DATABASE_HOST:-localhost}
DB_PORT=${DATABASE_PORT:-5432}
DB_USER=${DATABASE_USER:-postgres}

# https://stackoverflow.com/questions/8208181/create-postgres-database-using-batch-file-with-template-encoding-owner-and
psql -h ${DB_HOST} -p ${DB_PORT} -U ${DB_USER} -f db_user_create.sql
