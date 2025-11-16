#!/usr/bin/env bash
set -euo pipefail

echo "Running migrations..."
docker compose run --rm --entrypoint /bin/sh migrate -c 'migrate -path /migrations -database "$DATABASE_URL" up'
echo "Migrations completed."


