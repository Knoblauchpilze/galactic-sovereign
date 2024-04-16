#!/bin/bash

POSTGRESQL_VERSION=14
DOCKER_HOST_MASK="172.17.0.0"

# https://stackoverflow.com/questions/31249112/allow-docker-container-to-connect-to-a-local-host-postgres-database
sudo sed -i "s/#listen_addresses = 'localhost'/listen_addresses = '*'/g" /etc/postgresql/${POSTGRESQL_VERSION}/main/postgresql.conf

# https://dba.stackexchange.com/questions/83984/connect-to-postgresql-server-fatal-no-pg-hba-conf-entry-for-host
sudo echo "host  all  all  ${DOCKER_HOST_MASK}/0  scram-sha-256" >> /etc/postgresql/${POSTGRESQL_VERSION}/main/pg_hba.conf

sudo systemctl restart postgresql

echo "Use the following set of commands:"
echo "psql"
echo "ALTER USER postgres PASSWORD 'your-password';"
echo "quit"
echo "exit"

sudo -i -u postgres
