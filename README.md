
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

## Start an EC2 instance to host the service

To start the instance on which the service will be running you need an account which can start EC2 instances on AWS. To be fully usable by the service, the CI and the back-up process the instance should fill a few criteria.

### Network interface

The instance should open the following ports:
* allow incoming SSH connections (TCP 22)
* allow outbound internet traffic (TCP 80 and 443 for http and https respectively)
* allow incoming internet traffic (TCP 80 and 443 for http and https respectively)

Note that we don't open the `postgres` server port: this is because it is generally advised (see [this SO article](https://stackoverflow.com/questions/76541185/how-to-connect-to-postgresql-server-running-on-ec2)) to not publicly expose the database to the outside world but instead to use another secured connection mechanism (e.g. SSH) to open a tunnel to the instance and then make the DB connection transit through this tunnel. This is described in more details in the [following section](#in-case-of-a-remote-environment) on how to connect to the database.

### Instance role

So as to be able to access the S3 bucket to push the back-ups of the database, the instance needs to be attached a role allowed to do so.

### Software on the instance

The service is running in a docker container and requires a `postgres` server to be available on the machine. The back-up process uses the the AWS cli to perform the S3 related operations.

### Security keys

In order to allow the CI to update and deploy the new versions of the service, the instance should be configured so that it can be accessed. This means adding the CI keys to the `authorized_keys` file.

### Convenience template

In the AWS console there's already a convenience template meeting all of these requirements. It uses the `User data` mechanism to perform some configuration on the machine when it's first being booted. Is is accessible under the `ec2-postgres-docker-aws-cli` template:

![ec2 launch template](resources/ec2-launch-template.png)

### A word on User data

The User data mechanism allows, as per the [AWS documentation](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/user-data.html?icmpid=docs_ec2_console), to:
```
Specify user data to provide commands or a command script to run when you launch your instance. Input is base64 encoded when you launch your instance unless you select the User data has already been base64 encoded check box.
```

This repository defines a script that can be copied directly in the User data section under [ec2-setup-instance.sh](scripts/ec2-setup-instance.sh). It will perform the installation of the necessary additional packages and configure the keys of the CI so that it can access the instance.

The script will **not** setup the database (see the [creating the database](#creating-the-database) section for this). It will also not make sure that the service's docker container can properly access the local database server. This can be enabled by going to the section [just after](#allow-docker-containers-to-access-the-local-postgres-server) this one.

### Allow docker containers to access the local postgres server

When the service is up and running, it will attempt to access a remote `postgres` server from within the docker container. During local development, it is most common to use `localhost` as the host of the database to connect to a local `postgres` server.

In the production case, the docker container is for now also connecting to a local `postgres` server on the same EC2 instance: the only difference is that it does so from within a docker container.

The way docker exposes the `localhost` to application running within it is a bit different to how it works outside of it. This SO post describes [how to connect to localhost from inside a docker container](https://stackoverflow.com/questions/24319662/from-inside-of-a-docker-container-how-do-i-connect-to-the-localhost-of-the-mach) and goes into a bit of details on how it works. The TL;DR is that in the default network mode (bridge), the docker engine will associate an IP with a pattern like `172.17.X.Y` to reach localhost. This is already reflected in the [default config](cmd/users/internal/config.go) file.

### Allow incoming connections to the postgres server from docker host 

The previous step is just half of the work though. Now that the application within the docker container is successfully able to contact `localhost`, we need to instruct the `postgres` server to accept connections from this IP. By default, only connection from the local host are checked which does not work with the IP ranges assigned by the docker engine.

Once again we found a nice explanatory [article on SO](https://stackoverflow.com/questions/31249112/allow-docker-container-to-connect-to-a-local-host-postgres-database) to solve this problem. The solution involves two steps.

First, we need to instruct the `postgres` server to accept connections not only from local host but on a broader range of IPs. Reading what's indicated in the article, it can be done by editing a specific file:
```bash
POSTGRESQL_VERSION=14
vim /etc/postgresql/${POSTGRESQL_VERSION}/main/postgresql.conf
```

Once the editor is started, you can look if the following line already exists: if not create it.
```
listen_addresses = '*'
```
This will instruct the `postgres` server to listen on all addresses instead of just on the local host (the default).

Once this is done, we still need to instruct the `postgres` server about which authentication method it should expect for the new possible hosts. This can be achieved as described in [this SO article](https://dba.stackexchange.com/questions/83984/connect-to-postgresql-server-fatal-no-pg-hba-conf-entry-for-host).

First, open the following file with your favourite editor:
```bash
POSTGRESQL_VERSION=14
vim /etc/postgresql/${POSTGRESQL_VERSION}/main/pg_hba.conf
```

You can now add an entry for the IP range assigned by the docker host so as to allow connections from docker containers. We can try to keep things a bit more secure by not allowing all IPs but just the ones assigned by docker:
```
host  all  all  172.17.0.0/0  scram-sha-256
```

These modifications require to restart the `postgres` system to take effect. This can be done by running the following:
```bash
sudo systemctl restart postgresql
```

With all of this, the service running in the docker container should be able to connect to the database.

## Setup the database

In order to properly execute them, make sure that you know the `postgres` password (see [this SO link](https://stackoverflow.com/questions/27107557/what-is-the-default-password-for-postgres) to alter it if needed).

The general principle is to first start a shell for the `postgres` user with:
```bash
sudo -i -u postgres
```

After this you can start a `postgres` shell with:
```bash
psql
```

Once in the shell you can alter the password of the `postgres` user with:
```sql
ALTER USER postgres PASSWORD 'your-password';
```

After exiting the shells, you should know the password for the `postgres` password and can proceed with the rest of the setup process.

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

## Database back-up

### General principle

When the user-service is up and running it will perform some interactions with the database such as creating users or updating them. This means that in the event of a failure of this instance, we would loose some data.

To defend against such cases this service relies on an automatic back-up task which periodically dumps the content of the database and saves it to a dedicated S3 bucket. In the event of a failure of the EC2 instance hosting the database we can restore to the latest available back-up.

The back-up process is composed of two main components:
* a cron job running regularly to dump the database
* a S3 bucket to which the dumps are uploaded

The bucket is named `user-service-database-dumps` and has a policy to only keep the back-ups for a couple of days: after that they anyway lose their relevance.

### The back-up script

A convenience script is provided to perform the back-up: [database-backup.sh](scripts/database-backup.sh). This script will call `pg_dump` with a user allowed to view the schema and data of the database and then proceed to upload the generated file to the S3 bucket. By default, the script is configured to try to connect to a local database with the correct properties.

**It is necessary to copy this script on the EC2 instance hosting the database.**

The script expects the database password to be provided as environment variable under `DATABASE_PASSWORD`. Alternatively and in order to make it simpler to use the script with a cron job (see corresponding [section](#setting-up-the-database-backup)), in the remote case it can be beneficial to modify the copied version of the script to directly include the password in it. So change the line as follows (do not change the initial `:-` sequence):
```bash
DB_PASSWORD=${DATABASE_PASSWORD:-the-actual-password}
```

A significant difference between the local and remote environment is the behavior of the `IAM_ROLE` environment variable. When starting the EC2 instance hosting the service, we recommend assigning it a role that allows to access the S3 bucket by default (see the [starting an EC2 instance](#start-an-ec2-instance-to-host-the-service) section). In a local environment this is most likely not the case.

Assuming the user has an AWS profile and after the admin allowed the user's role to assume the right role, it is recommended to alter the `~/.aws/config` file to register the new profile, like so:

![DB user role](resources/db-user-role.png)

Once this is done, you can export an environment variable indicating to the script which role it should use to perform the upload as follows:
```bash
export IAM_ROLE=your-role-name
```

The script will then try to assume this role before uploading the backup to the bucket. By default the back-ups are named after the date at which they were taken:

![S3 back-ups example](resources/db-s3-bucket-back-ups.png)

### Setting up the database backup

In order to automatically configure the back-up, we recommend using [cron](https://en.wikipedia.org/wiki/Cron). There are numerous articles explaining on how to make it work. We used the following article: [How do I set up a cron job](https://askubuntu.com/questions/2368/how-do-i-set-up-a-cron-job) on `askubuntu`.

The idea here is to run and pick your favourite editor:
```bash
crontab -e
```

Once this is done, insert the following line:
```bash
*/20 * * * * /path/to/database-backup.sh
```

Most likely locally this will be the path to where the user cloned the repository and then `scripts/database-backup.sh` while in an EC2 instance it will be `/home/ubuntu/scripts/database-backup.sh` for example.

This will make sure that the back-up is performed automatically every 20 minutes.

### Restoring the database from a saved file

In case of a failure of the EC2 instance or in general for any situation that would require to restore the database to an earlier state, we can use the [database-restore.sh](scripts/database-restore.sh) script.

It works in a very similar way to the back-up one and expects a single argument representing the name of the database back-up to use to restore. The command line looks like so:
```bash
./database-restore.sh db_user_service_dump_2024-04-25-22_20_01.bck
```

Just like for the database backup case the script expects the database password to be available in the `DATABASE_PASSWORD` environment variable and the name of the IAM profile to use to access the data in the S3 bucket to be under `IAM_ROLE`. Note that the password for the database should be the one attached to the admin role.

The script will attempt to assume the role to access to the S3 bucket, download the dump locally and restore the database to its previous state.

A couple of useful considerations:
* in case the download fails, no restoration will be attempted
* in case the restore fails, the script stops at the first error

In case something fails during the restore operation, the best course of action is probably to drop the database entirely and start fresh from the [create database](#creating-the-database) section before attempting the restore again.

## Deploying the service to an EC2 instance

In order to deploy the service to run on an EC2 instance, we dockerize it by building a docker image for it. This process is automated in the CI and runs on each push to the master branch after the tests and the build of the image completed successfully.

The image is pushed to a dedicated dockerhub repository with a tag corresponding to the short SHA of the commit which generated the build:

![Dockerhub repository](resources/dockerhub-repo.png)

Images are publicly available and can be downloaded with a simple `docker pull` command such as:
```bash
docker up totocorpsoftwareinc/user-service:YOUR-TAG
```

To build the CI pipeline, we took a lot of inspiration from the content of [this article](https://medium.com/ryanjang-devnotes/ci-cd-hands-on-github-actions-docker-hub-aws-ec2-ba09f80297e1) by Minho Jang. It explains the steps to connect to the instance and update the running docker image.

The deployment process offers a small amount of flexibility. As it is it will try to reach an EC2 instance with a known IP address (if the IP changes we need to update the CI) and start the container to listen on port `80`.

## Monitoring the service

In order to monitor the health of the service while deployed, it is possible to use the script [service-monitoring.sh](scripts/service-monitoring.sh). This script is meant to be ran as a cron job in a similar way to the database [backup script](#setting-up-the-database-backup).

The script will attempt to curl the endpoint providing the healtcheck for the service and check the return code. For a healthy service this shoud be 200. Anything else will be interepreted as a unhealthy service and the script will attempt to stop the running docker container (if any) and start it again.

**Note:** the script is expected to run as root or alternatively to run with a user having the permissions to run `docker` without `sudo`.

To install the script, you first need to open the crontab:
```bash
crontab -e
```

Once this is done, insert the following line:
```bash
*/5 * * * * /path/to/service-monitoring.sh
```

Most likely locally this will be the path to where the user cloned the repository and then `scripts/service-monitoring.sh` while in an EC2 instance it will be `/home/ubuntu/scripts/service-monitoring.sh` for example.

As it is presented here the cron job will trigger every 5 minutes and take actions appropriately.

# How does the service work?

## General design

This service exposes CRUD operations to manage users. A user is defined as a simple collection of an email and a password. Upon creation, a user is associated with an API key, which is required to perform interactions with the server.

## API keys

It is only possible to access the API key attached to a user on the response returned by the `POST` operation. All subsequent operations will erase the keys from the response. This is to make sure that only the client who initially created the user knows the keys. A typical response will look like so:

![Create response](resources/service-api-keys.png)

And a typical GET call will return:

![Get response](resources/service-get-call.png)

Even though no endpoint currently allows to do this, it is possible to de/activate API keys and each user may have more than one. Only the active ones are considered when returning them (in the CREATE operation).

## Web framework

We chose to use [echo](https://echo.labstack.com/) as a web framework for this project. This is used to instantiate the web server running the application and for interpreting the parameters provided by the input requests.

## Throttling

Because this server is supposed to be deployed on a internet facing instance, we want to make sure that we don't risk being subjected to too many requests from undesired attackers. In an attempt to mitigate this we implemented a throttling mechanism which by default allows 10 requests per second from any sources before returning `429` on new requests.

# Future work

As of now, **all** operations require to be using an API key to succeed. This is of course not optimal because how are you supposed to get an API key in the first place?

Currently we also listen on `http://...` and don't provide anything in regards to `https`. This should be changed for enhanced security.

# Cheat sheet

## Create new user
```bash
curl -X POST -H "Content-Type: application/json" -H 'X-Api-Key: 2da3e9ec-7299-473a-be0f-d722d870f51a' http://localhost:60001/v1/users -d '{"email":"user-1@mail.com","password":"password-for-user-1"}' | jq
```

## Query existing user
```bash
curl -X GET -H "Content-Type: application/json" -H 'X-Api-Key: 2da3e9ec-7299-473a-be0f-d722d870f51a' http://localhost:60001/v1/users/4f26321f-d0ea-46a3-83dd-6aa1c6053aaf | jq
```

## Query non existing user
```bash
curl -X GET -H "Content-Type: application/json" -H 'X-Api-Key: 2da3e9ec-7299-473a-be0f-d722d870f51a' http://localhost:60001/v1/users/4f26321f-d0ea-46a3-83dd-6aa1c6053aae | jq
```

## Query without API key
```bash
curl -X GET -H "Content-Type: application/json" http://localhost:60001/v1/users/4f26321f-d0ea-46a3-83dd-6aa1c6053aae | jq
```

## List users
```bash
curl -X GET -H "Content-Type: application/json" -H 'X-Api-Key: 2da3e9ec-7299-473a-be0f-d722d870f51a' http://localhost:60001/v1/users | jq
```

## Patch existing user
```bash
curl -X PATCH -H "Content-Type: application/json" -H 'X-Api-Key: 2da3e9ec-7299-473a-be0f-d722d870f51a' http://localhost:60001/v1/users/0463ed3d-bfc9-4c10-b6ee-c223bbca0fab -d '{"email":"test-user@real-provider.com","password":"strong-password"}'| jq
```

## Delete user
```bash
curl -X DELETE -H "Content-Type: application/json" -H 'X-Api-Key: 2da3e9ec-7299-473a-be0f-d722d870f51a' http://localhost:60001/v1/users/0463ed3d-bfc9-4c10-b6ee-c223bbca0fab | jq
```

## Build the docker container
```bash
GIT_COMMIT_HASH=$(git rev-parse --short HEAD) ENV_DATABASE_PASSWORD='password' make user-service-build
```

## Run the docker container
```bash
GIT_COMMIT_HASH=$(git rev-parse --short HEAD) ENV_DATABASE_PASSWORD='password' make user-service-run
```
