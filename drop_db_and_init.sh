#!/usr/bin/env bash

sudo su postgres -c "psql -c \"SELECT pg_terminate_backend(pg_stat_activity.pid) FROM pg_stat_activity WHERE pg_stat_activity.datname = '${1:-books}' AND pid <> pg_backend_pid();\""
sudo su postgres -c "psql -c \"DROP DATABASE ${1:-books};\" "
bash db_init.sh
