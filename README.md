
# user-service

The definition of a service to manage users and authentication.

[![codecov](https://codecov.io/gh/Knoblauchpilze/user-service/branch/master/badge.svg?token=WNLIZF0FBL)](https://codecov.io/gh/Knoblauchpilze/user-service)

# Installation

## Prerequisite

This projects uses:

- [golang migrate](https://github.com/golang-migrate/migrate/blob/master/cmd/migrate/README.md): following the instructions there should be enough.
- [postgresql](https://www.postgresql.org/) which can be taken from the packages with `sudo apt-get install postgresql-14` for example.

## Clone the repository

- Clone the repo: `git clone git@github.com:Knoblauchpilze/user-service.git`.
- Install Go from [here](https://go.dev/doc/install). **NOTE**: this project expects Go 1.22 to be available on the system.
- Go to the project's directory `cd ~/path/to/the/repo`.
- Compile and install: `make`.
- Execute any application with `make run app_name`.

## Setup the database

Some convenience scripts are provided in the [database](database) folder. Assuming `postgresql` is already available on the system and that you know the password for the postgres user, you can run:

```bash
cd database
./create_user.sh
./create_database.sh
make migrate
```

You should then be able to connect to the database with `make connect` and inspect the content.

# Cheat sheet

## Create new user
```bash
curl -H "Content-Type: application/json" -X POST http://localhost:60001/v1/users/users -d '{"email":"some-user@mail.com","password":"1234"}' | jq
```

## Query existing user
```bash
curl -X GET http://localhost:60001/v1/users/users/08ce96a3-3430-48a8-a3b2-b1c987a207cb | jq
```

## Query non existing user
```bash
curl -X GET http://localhost:60001/v1/users/users/08ce96a3-3430-48a8-a3b2-b1c987a207ca | jq
```

## List users
```bash
curl -X GET http://localhost:60001/v1/users/users | jq
```

## Patch existing user
```bash
curl -H "Content-Type: application/json" -X PATCH http://localhost:60001/v1/users/users/f16c01cd-a6cc-4645-b4f8-3fd2e6837af8 -d '{"email":"some-other-user@mail.com","password":"1235"}'| jq
```

## Delete user
```bash
curl -X DELETE http://localhost:60001/v1/users/users/08ce96a3-3430-48a8-a3b2-b1c987a207cc | jq
```