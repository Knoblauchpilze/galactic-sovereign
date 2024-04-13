#!/bin/bash

if [ "$#" -ne 1 ]; then
  echo "usage: database-restore.sh [database_dump_file]"
  exit 1
fi

DB_HOST=${DATABASE_HOST:-localhost}
DB_PORT=${DATABASE_PORT:-5432}
DB_NAME=${DATABASE_NAME:-db_user_service}
DB_USER=${DATABASE_USER:-user_service_admin}
DB_PASSWORD=${DATABASE_PASSWORD:-}

S3_BUCKET_NAME=${S3_BUCKET:-s3://user-service-database-dumps}
IAM_ROLE_NAME=${IAM_ROLE:-}

DB_DUMP=$1

if [ "${DB_PASSWORD}" == "" ]; then
  echo "DB password is not configured, please define environment variable DATABASE_PASSWORD, aborting"
  exit 1
fi
URL_ENCODED_PASSWORD=$(urlencode ${DB_PASSWORD})

IAM_ROLE_OPTION=""
if [ "${IAM_ROLE_NAME}" != "" ]; then
  IAM_ROLE_OPTION="--profile ${IAM_ROLE_NAME}"
  echo "Assuming role ${IAM_ROLE_NAME}"
fi

echo "Downloading from ${S3_BUCKET_NAME}/${DB_DUMP}..."
aws s3 cp ${S3_BUCKET_NAME}/${DB_DUMP} ${DB_DUMP} ${IAM_ROLE_OPTION}

# https://stackoverflow.com/questions/40082346/how-to-check-if-a-file-exists-in-a-shell-script
if [ ! -f "${DB_DUMP}" ]; then
  echo "Failed to download ${DB_DUMP}, aborting"
  exit 1
fi

# https://www.postgresql.org/docs/current/app-pgrestore.html
echo "Restoring database ${DB_NAME}..."
pg_restore -ce --if-exists -d postgres://${DB_USER}:${URL_ENCODED_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME} ${DB_DUMP}
