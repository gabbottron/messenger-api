# messenger-api
This is a basic message app

After checkout
- Run dep init / dep ensure
- Create a .env file and fill it out

.env file should look like this:
```
APP_ENV=local-development
APP_ENDPOINT=http://localhost:3000

PORT=8083

POSTGRES_DB_USER=msgr
POSTGRES_DB_PASSWORD={{YOUR_PASSWORD_HERE}}

DB_NAME=msgr
DB_PORT=5439
DB_HOSTNAME=127.0.0.1
```

go run main.go

To run the tests with a clean database it can be helpful to use a run_tests script:
```
#!/bin/bash

psql -h 127.0.0.1 -p 5439 -U msgr -f /...path-to-db-repo.../messenger-db/sql/drop_all_tables.sql
psql -h 127.0.0.1 -p 5439 -U msgr -f /...path-to-db-repo.../messenger-db/sql/init_db.sql
psql -h 127.0.0.1 -p 5439 -U msgr -f /...path-to-db-repo.../messenger-db/sql/fixtures.sql
```