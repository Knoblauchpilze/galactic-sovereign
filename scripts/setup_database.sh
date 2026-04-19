#!/bin/bash

echo "Setting up galactic-sovereign database..."
echo "You may need to enter multiple times the postgres password"

export ADMIN_PASSWORD="admin_password"
export MANAGER_PASSWORD="manager_password"
export USER_PASSWORD="user_password"

echo "Creating users..."
/bin/bash database/create_user.sh database/galactic-sovereign

echo "Creating database..."
/bin/bash database/create_database.sh database/galactic-sovereign

echo "Seeding values..."
export DB_PATH="database/galactic-sovereign"
make -f database/Makefile migrate

echo "All done!"
