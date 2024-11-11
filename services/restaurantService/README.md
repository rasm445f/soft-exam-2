# go-startkode

Scaffolding for a go api project.

<!--toc:start-->

- [Run:](#run)
- [Prerequisites:](#prerequisites)
- [Setup:](#setup)
- [Usage:](#usage)
<!--toc:end-->

## Run:

To simply run the application only docker and docker compose is needed. Firstly rename the `.example.env` file to `.env` then run:

```shell
docker-compose --profile app up --build
```

Then the swagger UI can be found on: `localhost:8080/api/docs/index.html`

## Prerequisites:

- [Golang](https://go.dev/doc/install)
- [Make](https://www.gnu.org/software/make/)
- [docker](https://www.docker.com/) & [docker compose](https://docs.docker.com/compose/)

## Setup:

- Run the `install-tools.sh` scripts if you don't already have the tools. (Ensure your `PATH` is set up correctly).

```shell
./install-tools.sh
```

- Rename or copy the `.example.env` file to `.env` and populate the fields as needed.

- Run `docker compose up` to start the PostgreSQL Docker container or `docker compose up -d` to run it in the background.

- Run `make migrate-up` to setup the database with the tables etc. specified in the `db/migration/` folder.

## Usage:

The Makefile defines different task to run. When in doubt you can just run `make` to find out which tasks are available or read the `Makefile` for more information

```
$ make
Available targets:
  run                   # Run with air
  build                 # Build the project
  docs                  # Generate docs with swag
  migrate-new           # Create new sql migration
  migrate-up            # Migrate the db
  migrate-down          # Roll back the db
  test                  # Run tests
  test-verbose          # Run test with verbose flag
  test-cover            # Show test coverage
  test-cover-html       # Show test coverage in the browser
```
