#!/bin/bash

DB_PATH=$1

psql -h localhost -U postgres postgres -f ${DB_PATH}/db_user_drop.sql
