# galactic-sovereign

This repository contains the backend service powering the [galactic-sovereign.gasteropo.de](https://galactic-sovereign.gasteropo.de) website. The frontend for this project is defined in the [galactic-sovereign-frontend](https://github.com/Knoblauchpilze/galactic-sovereign-frontend) repository.

On this website you can find an online multiplayer web-browser game titled **Galactic Sovereign**. This is a strategy game where the players can mine resources on their respective planets and improve the infrastructure to do so by upgrading their mines. It is largely inspired by the basics of [OGame](https://en.wikipedia.org/wiki/OGame), a famous strategy game.

Below is a screenshot of the welcome page:

![Welcome page of the Galactic Sovereign website](resources/galactic-sovereign-welcome-page.png)

Here is a view of the lobby:

![Lobby of the Galactic Sovereign website](resources/galactic-sovereign-lobby.png)

And a screenshot of the actual game:

![In-game view for the Galactic Sovereign website](resources/galactic-sovereign-game-view.png)

# Overview

This project uses the following technologies:

- [postgres](https://www.postgresql.org/) for the databases.
- [go](https://go.dev/) as the server backend language.
- [docker](https://www.docker.com/) as the containerization tool to deploy services.
- [dockerhub](https://hub.docker.com/) to host the images of services and make them available.

# Badges

[![codecov](https://codecov.io/gh/Knoblauchpilze/galactic-sovereign/branch/master/badge.svg?token=WNLIZF0FBL)](https://codecov.io/gh/Knoblauchpilze/galactic-sovereign)

[![Build services](https://github.com/Knoblauchpilze/galactic-sovereign/actions/workflows/build-and-push.yml/badge.svg)](https://github.com/Knoblauchpilze/galactic-sovereign/actions/workflows/build-and-push.yml)

[![Database migration tests](https://github.com/Knoblauchpilze/galactic-sovereign/actions/workflows/database-migration-tests.yml/badge.svg)](https://github.com/Knoblauchpilze/galactic-sovereign/actions/workflows/database-migration-tests.yml)

# Installation

## Prerequisites

The tools described below are directly used by the project. It is mandatory to install them in order to build the project locally.

See the following links:

- [golang](https://go.dev/doc/install): this project was developed using go `1.23.2`.
- [golang migrate](https://github.com/golang-migrate/migrate/blob/master/cmd/migrate/README.md): following the instructions there should be enough.
- [postgresql](https://www.postgresql.org/) which can be taken from the packages with `sudo apt-get install postgresql-14` for example.
- [docker](https://docs.docker.com/engine/install/ubuntu/#install-using-the-repository) which can be installed following the instructions of the previous link.

We also assume that this repository is cloned locally and available to use. To achieve this, just use the following command:

```bash
git clone git@github.com:Knoblauchpilze/galactic-sovereign.git`
```

## Secrets in the CI

The CI workflows define several secrets that are expected to be created for the repository when cloned/forked/used. Each secret should be self-explanatory based on its name. Most of them require to setup an account on one or the other service mentioned in this README.

## Process to install

The following sections describe how to setup the local postgres server to be able to host the databases needed by the project. This includes:

- altering the postgres password if needed
- setting up the database
- connecting to the database

In addition, the following sections expect that you have access to several (secured) passwords to be used for the various users needed by the databases. The current architecture requires **3** passwords per database. It is recommended to have separate passwords and not reuse them.

## The postgres password

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

## Creating the database

⚠️ In this section we will focus on setting up the galactic-sovereign database but this approach can be extended to the other databases needed by the project.

In both remote and local case, the configuration of the database happens in 3 steps:

- creation of the users to manage the database
- creation of the database
- creating the schema of the database and the relations

With the `postgres` password at hand, you can navigate to the [database](database) directory and use the following scripts defined there.

Each script requires a few environment variables to be defined in order to properly work. The very first one defines the path to reach the files related to the database to deploy. It is expected as an argument to the script. A common one is `DATABASE_PORT` which should point to the port to reach the `postgres` server. Should be `5432` (default) locally and `5000` in case of remote.

Additionally, [create_user.sh](database/create_user.sh) requires `ADMIN_PASSWORD`, `MANAGER_PASSWORD` and `USER_PASSWORD` which correspond to the three passwords you should have generated beforehand.

Finally the [Makefile](database/Makefile) to perform the migration requires to define the variables `DB_PORT` (same value as `DATABASE_PORT`) and `DB_PASSWORD` which should correspond to `ADMIN_PASSWORD`.

As a recap, here are the required steps for the remote case:

```bash
export DATABASE_PORT=5000
export ADMIN_PASSWORD='admin-password'
export MANAGER_PASSWORD='manager-password'
export USER_PASSWORD='user-password'

cd database
./create_user.sh galactic-sovereign
./create_database.sh galactic-sovereign

DB_PATH=galactic-sovereign DB_PORT=5000 DB_PASSWORD='admin-password' make migrate
```

If everything goes well for the migration, you should obtain something like this:
![DB migration success](resources/db-migration-success.png)

### Connect and inspect the database

If you want to connect to the database to inspect its content, you can use the `connect` target defined in the [Makefile](database/Makefile). Just like for the migration, it expects the `DB_PORT` and `DB_PASSWORD` environment variables to be available. An example command would be:

```bash
DB_PORT=5000 DB_PASSWORD='admin-password' make connect
```

## Generate mocks for testing

This project uses generated mocks for testing purposes. The mocks live alongside the package needing them and can be generated by running the following command from the root of the repository:

```bash
go generate ./...
```

Or by using the convenience make target with:


```bash
make generate-mocks
```

It is mandatory to re-generate the mocks in case of a modification of the interfaces supporting the mocks. This is enforced by a check in the CI (see [check-generated-mocks.yml](.github/workflows/check-generated-mocks.yml)).

## Linting

This project uses [golangci-lint](https://github.com/golangci/golangci-lint) as a collection of linters to perform static code analysis and detect issues. It can be installed locally through the procedure described [here](https://golangci-lint.run/docs/welcome/install/local/#binaries)

```bash
# binary will be $(go env GOPATH)/bin/golangci-lint
curl -sSfL https://golangci-lint.run/install.sh | sh -s -- -b $(go env GOPATH)/bin v2.12.2

golangci-lint --version
```

This will install the executable in `/home/$USER/go/bin` which might or might not be in your path. To add it you can add `/home/knoblauch/$USER/go/bin` to the path.

Once this is done, you can run the linters on the code by using the dedicated make target:
```bash
make lint
```

This target should download the executable automatically and not require any installation step. A target `fix-lint` is also available to fix the fixable lint errors.

## Local development

In case you're using Visual Studio Code as an IDE you can copy the following template in `.vscode/launch.json`:

```json
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Launch server",
      "type": "go",
      "request": "launch",
      "mode": "debug",
      "program": "${workspaceFolder}/cmd/galactic-sovereign/main.go",
      "args": ["galactic-sovereign-dev"]
    }
  ]
}
```

This will allow you to debug the galactic-sovereign server directly in your IDE.

The documentation from VS code indicates that you can just hover over existing attributes to get a description of their purpose. You can also use `Ctrl + Space` to ask Intellisense to give you information about other arguments.

## Generate API specification

You can generate the Swagger specification from the annotated handlers with:

```bash
make generate-api-spec
```

This creates `api/swagger.yaml`.

## Using the data generation scripts

Scripts are provided to make easy to test common scenarios of the game. They live under [scripts/game](scripts/game). Those scripts allow to create players and building actions in a semi-automated way. To ensure that the scripts are working properly, it is recommended to first run once the make target `setup` to create the `sandbox` folder (or create it manually).

The [create-player.sh](scripts/game/create-player.sh) script allows to create a player in the Oberon universe with a generated name. When the creation is successful, the data is saved in a shared folder (`sandbox`) in a file named `player.json`. This folder is added to the ignore list for git and allows to persist local data so that subsequent scripts can use it.

The [create-building-action.sh](scripts/game/create-building-action.sh) script allows to create a building action for a planet. The action is hardcoded for the metal mine and the player/planet can either be provided as inputs to the script or derived from the data available in the `sandbox`. The script persists a `building_action.json` file to the sandbox for subsequent use.

The [create-planet.sh](scripts/game/create-planet.sh) script allows to create a planet for a player. The player is expected to either be provided as an input to the script or is derived from the data availablei n the `sandbox`. The script persists a `planet.json` file to the sandbox for subsequent use.

A typical use case is visible below:

```bash
the-pc:/galactic-sovereign$ ./scripts/game/create-player.sh
No player name provided, using toto-2026-07-04-10:55:00
Created player 0dcbf740-96d0-4b73-8350-0ba90b569202!
Homeworld: 7ccda1c0-3f48-477d-908f-dd95b7594c07
the-pc:/galactic-sovereign$ ./scripts/game/create-building-action.sh
Using player and planet from file 0dcbf740-96d0-4b73-8350-0ba90b569202 (planet: 7ccda1c0-3f48-477d-908f-dd95b7594c07)
Created building action 0dcbf740-96d0-4b73-8350-0ba90b569202!
Completion time: 2026-07-04T08:56:48.640723Z
the-pc:/galactic-sovereign$ ./scripts/game/create-planet.sh
Using player from file 0dcbf740-96d0-4b73-8350-0ba90b569202!
Created planet 10651ef1-fd91-4114-941f-167cd67f2a24!
```
