#!/bin/sh

set -e # exit immediately if non zero status returned

echo "run db migration"
/app/migrate -path /app/migration -database "$DB_SOURCE" -verbose up

echo "start app"
exec "$@"
