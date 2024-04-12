#!/bin/sh

if [ "$#" -ne 1 ]; then
  echo "usage: database-restore.sh [database_dump_file]"
  exit 1
fi

DB_HOST=${DATABASE_HOST:-localhost}
DB_NAME=${DATABASE_NAME:-db_user_service}
DB_USER=${DATABASE_USER:-user_service_admin}

S3_BUCKET_NAME=${S3_BUCKET:-s3://user-service-database-dumps}
IAM_ROLE_NAME=${IAM_ROLE:-user-service-database-dev}

DB_DUMP=$1

echo "Downloading from ${S3_BUCKET_NAME}/${DB_DUMP}..."
aws s3 cp ${S3_BUCKET_NAME}/${DB_DUMP} ${DB_DUMP} --profile ${IAM_ROLE_NAME}

# https://stackoverflow.com/questions/40082346/how-to-check-if-a-file-exists-in-a-shell-script
if [ ! -f "${DB_DUMP}" ]; then
  exit 1
fi

# https://www.postgresql.org/docs/current/app-pgrestore.html
echo "Restoring database ${DB_NAME}..."
pg_restore -h ${DB_HOST} -U ${DB_USER} -cWe --if-exists -d ${DB_NAME} ${DB_DUMP}
