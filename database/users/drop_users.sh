#!/bin/bash

psql -h localhost -U postgres postgres -f db_user_drop.sql
