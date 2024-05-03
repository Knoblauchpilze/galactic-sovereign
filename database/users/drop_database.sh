#!/bin/bash

psql -h localhost -U postgres postgres -f db_drop.sql
