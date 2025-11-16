#!/usr/bin/env bash
set -euo pipefail

echo "Running migrations..."
docker compose run --rm migrate up
echo "Migrations completed."


