#!/bin/sh

if [ "$1" = 'migrate_and_release' ]; then
  make database-check
  make database-migration-up
  exec /app/api
elif [ "$1" = 'release' ]; then
  make database-check
  exec /app/api
fi
