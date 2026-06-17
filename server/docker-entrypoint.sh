#!/bin/sh
set -e

echo "Running database migrations..."
while true; do
  if ./server up; then
    break
  fi
  echo "Waiting for database to be ready... retrying in 2s"
  sleep 2
done

echo "Starting application server..."
exec ./server
