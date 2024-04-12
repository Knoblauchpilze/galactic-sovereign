#!/bin/sh

# https://stackoverflow.com/questions/39296472/how-to-check-if-an-environment-variable-exists-and-get-its-value
DB_HOST=${DATABASE_HOST:-localhost}
DB_NAME=${DATABASE_NAME:-db_user_service}
DB_USER=${DATABASE_USER:-user_service_user}

S3_BUCKET_NAME=${S3_BUCKET:-s3://user-service-database-dumps}
IAM_ROLE_NAME=${IAM_ROLE:-}

# https://stackoverflow.com/questions/4018503/is-there-a-date-time-format-that-does-not-have-spaces
DATE=$(date '+%F-%T' | tr ':' '_')

DB_DUMP="${DB_NAME}_dump_${DATE}.bck"

IAM_ROLE_OPTION=""
if [ "${IAM_ROLE_NAME}" != "" ]; then
  IAM_ROLE_OPTION="--profile ${IAM_ROLE_NAME}"
  echo "Assuming role ${IAM_ROLE_NAME}"
fi

# https://www.postgresql.org/docs/current/app-pgdump.html
echo "Dumping database to ${DB_DUMP}..."
pg_dump -h ${DB_HOST} -U ${DB_USER} -F c -f ${DB_DUMP} ${DB_NAME}

echo "Uploading to ${S3_BUCKET_NAME}..."
aws s3 cp ${DB_DUMP} ${S3_BUCKET_NAME} ${IAM_ROLE_OPTION}
