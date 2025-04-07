#!/bin/bash
set -e

# Function to execute a SQL file
execute_sql_file() {
  local file="$1"
  echo "Executing SQL file: $file"
  psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" -f "$file"
}

# First run all migration files in order
echo "Running database migrations..."
for migration in /docker-entrypoint-initdb.d/migrations/*.sql; do
  if [ -f "$migration" ]; then
    execute_sql_file "$migration"
  fi
done

# Then run all seed files
echo "Running seed data scripts..."
for seed in /docker-entrypoint-initdb.d/seed/*.sql; do
  if [ -f "$seed" ]; then
    execute_sql_file "$seed"
  fi
done

echo "Database initialization complete!"

