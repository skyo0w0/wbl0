#!/bin/sh

set -e

host="$1"
shift
cmd="$@"

until PGPASSWORD=$POSTGRES_PASSWORD psql -p "$DB_PORT" -h "$host" -U "$POSTGRES_USER" -c '\q'; do
  >&2 echo "Postgres is unavailable - sleeping, {$DB_PORT}, {$host}, {$POSTGRES_USER}"
  sleep 1
done

>&2 echo "Postgres is up - executing command"
exec $cmd