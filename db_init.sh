#!/usr/bin/env bash
echo 'Preparing Postgres'
sudo su postgres -c "psql -c \"CREATE DATABASE ${1:-books}\" "
sudo su postgres -c "psql -c \"CREATE USER ${1:-books} WITH PASSWORD '${1:-books}'\" "
sudo su postgres -c "psql -c \"ALTER ROLE ${1:-books} WITH CREATEDB\" "
sudo su postgres -c "psql -c \"GRANT ALL PRIVILEGES ON DATABASE ${1:-books} to ${1:-books}\" "
