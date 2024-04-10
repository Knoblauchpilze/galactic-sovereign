#!/bin/sh

# https://stackoverflow.com/questions/51649484/how-to-backup-postgresql-database-automatically-on-daily-basis

# https://stackoverflow.com/questions/4018503/is-there-a-date-time-format-that-does-not-have-spaces
DATE=$(date '+%F-%T')

# https://www.postgresql.org/docs/current/app-pgdump.html
pg_dump -h localhost -U user_service_user -F c -f "db_user_service_${DATE}.bck" db_user_service
