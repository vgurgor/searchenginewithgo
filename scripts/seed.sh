#!/usr/bin/env bash
set -euo pipefail

echo "Copying seed.sql into db container..."
docker compose cp scripts/seed.sql db:/seed.sql
echo "Running seed..."
docker compose exec -T db psql -U "${POSTGRES_USER:-postgres}" -d "${POSTGRES_DB:-searchdb}" -f /seed.sql
echo "Seed completed."


