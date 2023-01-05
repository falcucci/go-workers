#!/bin/bash
set -e
export PGPASSWORD='localhost';
psql -v ON_ERROR_STOP=1 --username "people" --dbname "people_db" <<-EOSQL
  CREATE USER people WITH PASSWORD 'localhost';
  CREATE DATABASE people_db;
  GRANT ALL PRIVILEGES ON DATABASE people_db TO races;
  \connect people_db races
EOSQL
