#!/bin/sh
set -e

DB_HOST="${DB_HOST:-db}"
DB_PORT="${DB_PORT:-5432}"
DB_USER="${DB_USER:-postgres}"
DB_NAME="${DB_NAME:-waste_db}"

echo "waiting for postgres at ${DB_HOST}:${DB_PORT}..."

MAX_RETRIES=30
RETRY=0

until pg_isready -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" > /dev/null 2>&1; do
  RETRY=$((RETRY + 1))
  if [ "$RETRY" -ge "$MAX_RETRIES" ]; then
    echo "error: postgres not ready after ${MAX_RETRIES} attempts"
    exit 1
  fi
  sleep 1
done

echo "postgres is ready"

TABLE_EXISTS=$(psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -tAc "
  SELECT EXISTS (
    SELECT 1 FROM information_schema.tables
    WHERE table_schema = 'public' AND table_name = 'households'
  );")

if [ "$TABLE_EXISTS" = "t" ]; then
  echo "migrations already applied, skipping"
else
  echo "applying migrations..."
  psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -f /app/migrations/000001_init.up.sql
  echo "applying seed..."
  psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -f /app/migrations/000002_seed.up.sql
  echo "migrations applied"
fi

exec "$@"