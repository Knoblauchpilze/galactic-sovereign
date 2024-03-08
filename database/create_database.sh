#!/bin/sh

psql -h localhost -U user_service_admin postgres -f db_create.sql
