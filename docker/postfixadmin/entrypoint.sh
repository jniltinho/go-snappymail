#!/bin/sh
set -e

DB_HOST="${DB_HOST:-mariadb}"
DB_PORT="${DB_PORT:-3306}"
DB_PASS="${DB_PASS:-postfixPassword}"
SESSION_SECRET="${SESSION_SECRET:-change-me-lab-session-secret-32chars-min}"

sed -e "s|@@DB_PASS@@|${DB_PASS}|g" \
    -e "s|@@SESSION_SECRET@@|${SESSION_SECRET}|g" \
    /app/config.toml.in > /app/config.toml

if [ -n "$DB_HOST" ] && [ -n "$DB_PORT" ]; then
    echo "Waiting for database at $DB_HOST:$DB_PORT..."
    while ! nc -z "$DB_HOST" "$DB_PORT"; do
      sleep 1
    done
    echo "Database is up!"
fi

echo "Running database migrations..."
./postfixadmin migrate

if [ -n "$ADMIN_EMAIL" ] && [ -n "$ADMIN_PASSWORD" ]; then
    echo "Checking/Creating initial admin: $ADMIN_EMAIL"
    ./postfixadmin admin --add-superadmin "$ADMIN_EMAIL:$ADMIN_PASSWORD" || true
fi

exec "$@"
