BEGIN;

DROP TABLE IF EXISTS sellers;
DROP TABLE IF EXISTS payouts;
DROP TABLE IF EXISTS payout_items;
DROP TABLE IF EXISTS items;
DROP TABLE IF EXISTS currencies;

DROP EXTENSION IF EXISTS "pgcrypto";

COMMIT;
