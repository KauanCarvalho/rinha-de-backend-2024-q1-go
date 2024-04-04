#!/bin/sh

if [ "$1" = 'migrate_and_release' ]; then
  make db/check
  make db/drop
  make db/create
  make db/migrations/up
  exec /app/api
elif [ "$1" = 'release' ]; then
  make db/check
  exec /app/api
fi
