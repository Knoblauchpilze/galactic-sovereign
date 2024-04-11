#!/bin/sh

# https://stackoverflow.com/questions/39296472/how-to-check-if-an-environment-variable-exists-and-get-its-value
DB_HOST=${DATABASE_HOST:-localhost}
DB_NAME=${DATABASE_NAME:-db_user_service}
DB_USER=${DATABASE_USER:-user_service_user}

S3_BUCKET_NAME=${S3_BUCKET:-s3://user-service-database-dumps}

# https://stackoverflow.com/questions/4018503/is-there-a-date-time-format-that-does-not-have-spaces
DATE=$(date '+%F-%T' | tr ':' '_')

DB_DUMP="${DB_NAME}_dump_${DATE}.bck"

# https://www.postgresql.org/docs/current/app-pgdump.html
pg_dump -h ${DB_HOST} -U ${DB_USER} -F c -f ${DB_DUMP} ${DB_NAME}
