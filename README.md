# messenger-api
This is a basic message app

**REQUIREMENTS**
- You must have [Docker](https://www.docker.com/) installed and running
- You must have a .env file in project root (see below)
- You should have the [database running](https://github.com/gabbottron/messenger-db)

## Standard dev configuration (.env file)
```
APP_ENV=local-development
APP_ENDPOINT=http://localhost:3000

PORT=8080

POSTGRES_DB_USER=msgr
POSTGRES_DB_PASSWORD={{YOUR_PASSWORD_HERE}}

DB_NAME=msgr
DB_PORT=5432
DB_HOSTNAME=messenger-db
```

## To start in a container (config comes from .env)
make up

## To run locally on host (default config in makefile)
make run-local

## To run locally and supply different API port
PORT=8089 make run-local

## To run all tests locally
make test-local

## To run all tests locally and override DB port
DB_PORT=5438 make test-local
