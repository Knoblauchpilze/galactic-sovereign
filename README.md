
# user-service

The definition of a service to manage users and authentication. This service comes with a working CI allowing to deploy it on an EC2 instance. The configuration also includes an automatic periodic backup of the database to a S3 storage to prevent data losses.

[![codecov](https://codecov.io/gh/Knoblauchpilze/user-service/branch/master/badge.svg?token=WNLIZF0FBL)](https://codecov.io/gh/Knoblauchpilze/user-service)

# Installation

## Prerequisite

This projects uses:

- [golang migrate](https://github.com/golang-migrate/migrate/blob/master/cmd/migrate/README.md): following the instructions there should be enough.
- [postgresql](https://www.postgresql.org/) which can be taken from the packages with `sudo apt-get install postgresql-14` for example.
- [docker](https://docs.docker.com/engine/install/ubuntu/#install-using-the-repository) which can be installed following the instructions of the previous link.
- [aws cli](https://docs.aws.amazon.com/cli/latest/userguide/getting-started-install.html#getting-started-install-instructions) which can be installed using the instructions in the previous link.

## Clone the repository

- Clone the repo: `git clone git@github.com:Knoblauchpilze/user-service.git`.
- Install Go from [here](https://go.dev/doc/install). **NOTE**: this project expects Go 1.22 to be available on the system.
- Go to the project's directory: `cd ~/path/to/the/repo`.
- Go to the server's directory: `cd cmd/users`.
- Compile and install the server: `make run`.

## Preparing for setting up the project

In the rest of the installation instructions, the project assumes that:
* all necessary software is installed on the local machine.
* 3 passwords have been generated and are available (for the db users).
* a RSA key pair is available along with an identity file wich authorization to log on the EC2 instance.

### Generate a RSA key pair

In order to generate the RSA key pair and the `pem` file you can first follow this [github link](https://docs.github.com/en/authentication/connecting-to-github-with-ssh/generating-a-new-ssh-key-and-adding-it-to-the-ssh-agent#generating-a-new-ssh-key) to generate a new ssh key and then this [server fault link](https://serverfault.com/questions/706336/how-to-get-a-pem-file-from-ssh-key-pair) to create a `pem` file from it.

### Allow the key pair to connect to EC2 instance

In order to allow a user to SSH onto an instance, their keys need to be added to the `~/.ssh/authorized_keys` file of the EC2 instance. In the AWS console, this can be achieved when starting the instance through the _User data_ as shown in the screenshot below:

![ec2 User data](resources/ec2-user-data.png)

The idea is to include (among other things) the public key of the RSA key pair to automatically be installed when starting the instance.

## Setup the database

In order to properly execute them, make sure that you know the `postgres` password (see [this SO link](https://stackoverflow.com/questions/27107557/what-is-the-default-password-for-postgres) to alter it if needed).

### In case of a remote environment

Some convenience scripts are provided to create a tunnel to the EC2 instance: the idea is to connect through SSH to the instance and then forward the local `postgres` port locally so that it's 'as if' a local database server was up and running.

To do so, you can run the [ec2-establish-db-tunnel.sh](scripts/ec2-establish-db-tunnel.sh) script:
```bash
./ec2-establish-db-tunnel.sh ec2-ip-address path/to/identify/file
```

The ip address can be fetched from the AWS console and the identity file is the `pem` file mentioned in [Generate a RSA key pair](generate-a-rsa-key-pair).

This command will by default create a tunnel allowing to access the remote `postgres` server on the EC2 instance on the local port `5000`.

### Creating the database

In both remote and local case, the configuration of the database happens in 3 steps:
* creation of the users to manage the database
* creation of the database
* creating the schema of the database and the relations

With the `postgres` password at hand, you can navigate to the [database](database) directory and use the following scripts defined there.

Each script requires a few environment variables to be defined in order to properly work. A common one is `DATABASE_PORT` which should point to the port to reach the `postgres` server. Should be `5432` (default) locally and `5000` in case of remote.

Additionally, [create_user.sh](database/create_user.sh) requires `ADMIN_PASSWORD`, `MANAGER_PASSWORD` and `USER_PASSWORD` which correspond to the three passwords you should have generated beforehand.

Finally the [Makefile](database/Makefile) to perform the migration requires to define the variables `DB_PORT` (same value as `DATABASE_PORT`) and `DB_PASSWORD` which should correspond to `ADMIN_PASSWORD`.

As a recap, here are the required steps for the remote case:
```bash
export DATABASE_PORT=5000
export ADMIN_PASSWORD='admin-password'
export MANAGER_PASSWORD='manager-password'
export USER_PASSWORD='user-password'

cd database
./create_user.sh
./create_database.sh

DB_PORT=5000 DB_PASSWORD='admin-password' make migrate
```

If everything goes well for the migration, you should obtain something like this:
![DB migration success](resources/db-migration-success.png)

### Connect and inspect the database

If you want to connect to the database to inspect its content, you can use the `connect` target defined in the [Makefile](database/Makefile). Just like for the migration, it expects the `DB_PORT` and `DB_PASSWORD` environment variables to be available. An example command would be:

```bash
DB_PORT=5000 DB_PASSWORD='admin-password' make connect
```

### Log out from the remote environment

In case you set up a tunnel using [for remote access to the database](#in-case-of-a-remote-environment), don't forget to close the tunnel once you're done using the corresponding script:
```bash
cd scripts
./ec2-close-db-tunnel.sh
```

# Cheat sheet

## Create new user
```bash
curl -H "Content-Type: application/json" -X POST http://localhost:60001/v1/users/users -d '{"email":"some-user@mail.com","password":"1234"}' | jq
```

## Query existing user
```bash
curl -X GET http://localhost:60001/v1/users/users/0463ed3d-bfc9-4c10-b6ee-c223bbca0fab | jq
```

## Query non existing user
```bash
curl -X GET http://localhost:60001/v1/users/users/0463ed3d-bfc9-4c10-b6ee-c223bbca0fac | jq
```

## List users
```bash
curl -X GET http://localhost:60001/v1/users/users | jq
```

## Patch existing user
```bash
curl -H "Content-Type: application/json" -X PATCH http://localhost:60001/v1/users/users/0463ed3d-bfc9-4c10-b6ee-c223bbca0fab -d '{"email":"some-other-user@mail.com","password":"1235"}'| jq
```

## Delete user
```bash
curl -X DELETE http://localhost:60001/v1/users/users/0463ed3d-bfc9-4c10-b6ee-c223bbca0fab | jq
```

## Run the docker container
```bash
docker run --network=host user-service
```
