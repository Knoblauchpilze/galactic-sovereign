#!/bin/bash

DB_PATH=$1
DB_HOST=${DATABASE_HOST:-localhost}
DB_PORT=${DATABASE_PORT:-5432}
DB_USER=${DATABASE_USER:-postgres}

if [ "${DB_PATH}" == "" ]; then
  echo "No path provided, defaulting to galactic-sovereign"
  DB_PATH="galactic-sovereign"
fi

echo "Dropping database ${DB_PATH}..."
bash drop_database.sh galactic-sovereign
echo "Creating database ${DB_PATH}..."
bash create_database.sh galactic-sovereign
echo "Migrating database ${DB_PATH}..."
make migrate
