#! /bin/bash

description=$(
  echo $@ \
  | tr '[:upper:]' '[:lower:]' \
  | sed -e 's/ /_/g'
)
if [ "$description" = "" ]; then description="unnamed_database_migration"; fi

ts=$(date -u +%Y%m%d%H%M%S)
goose -dir db/migrations create ${description} sql

