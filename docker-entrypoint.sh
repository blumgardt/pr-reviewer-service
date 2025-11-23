#!/bin/sh
set -e

echo "Waiting for postgres at ${DB_HOST}:${DB_PORT}..."

until nc -z "$DB_HOST" "$DB_PORT"; do
  sleep 0.5
done

echo "Postgres is up. Running migrations..."

psql "postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable" \
  -f /app/internal/migrations/0001_init.up.sql

echo "Migrations done. Starting app..."

exec ./pr-reviewer-service