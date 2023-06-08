#!/bin/bash
set -eo pipefail
sleep 10
goose -dir db/migrations postgres "host=$PGHOST user=$PGUSER password=$PGPASSWORD dbname=$PGDATABASE sslmode=$PGSSLMODE" up