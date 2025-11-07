#!/bin/sh
set -e

# Start cron in background
service cron start

# Optional: log to stdout so you can see cron output in `docker logs`
touch /var/log/cron.log
tail -f /var/log/cron.log &

# Run the main app
exec "$@"
