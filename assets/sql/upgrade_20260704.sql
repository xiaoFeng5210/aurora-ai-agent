-- Aurora AI Agent
-- Card module upgrade script: add title column
-- Generated at: 2026-07-04 00:00:00 Asia/Shanghai
--
-- Production execution example:
--   psql "host=<host> port=<port> user=<user> dbname=<dbname> password=<password> sslmode=disable" -f assets/sql/upgrade_20260704.sql

BEGIN;

ALTER TABLE "card"
    ADD COLUMN IF NOT EXISTS title VARCHAR(255) NOT NULL DEFAULT '';

COMMIT;
