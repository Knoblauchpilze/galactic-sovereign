#!/bin/bash

# https://stackoverflow.com/questions/39296472/how-to-check-if-an-environment-variable-exists-and-get-its-value
DB_HOST=${DATABASE_HOST:-localhost}
DB_PORT=${DATABASE_PORT:-5432}
DB_NAME=${DATABASE_NAME:-db_user_service}
DB_USER=${DATABASE_USER:-user_service_user}
DB_PASSWORD=${DATABASE_PASSWORD:-}

S3_BUCKET_NAME=${S3_BUCKET:-s3://user-service-database-backups}
IAM_ROLE_NAME=${IAM_ROLE:-}

# https://stackoverflow.com/questions/4018503/is-there-a-date-time-format-that-does-not-have-spaces
DATE=$(date '+%F-%T' | tr ':' '_')

DB_DUMP="/tmp/${DB_NAME}_dump_${DATE}.bck"

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

# https://www.postgresql.org/docs/current/app-pgdump.html
echo "Dumping database to ${DB_DUMP}..."
pg_dump postgres://${DB_USER}:${URL_ENCODED_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME} -Fc -f ${DB_DUMP}

if [ ! -f "${DB_DUMP}" ]; then
  echo "Failed to create ${DB_DUMP}, aborting"
  exit 1
fi

echo "Uploading to ${S3_BUCKET_NAME}..."
aws s3 cp ${DB_DUMP} ${S3_BUCKET_NAME} ${IAM_ROLE_OPTION}
