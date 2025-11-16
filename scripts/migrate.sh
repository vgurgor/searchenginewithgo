#!/usr/bin/env bash
set -euo pipefail

echo "Running migrations..."
docker compose run --rm migrate
echo "Migrations completed."


