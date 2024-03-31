#!/bin/sh

psql -h localhost -U postgres postgres -f db_create.sql
