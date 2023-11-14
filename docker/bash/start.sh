#!/bin/sh

set -e
source /app/app.env

echo "Start the app"
exec "$@"