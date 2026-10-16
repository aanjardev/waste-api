DROP TABLE IF EXISTS payments;
DROP TABLE IF EXISTS waste_pickups;
DROP TABLE IF EXISTS households;

DROP TYPE IF EXISTS payment_status;
DROP TYPE IF EXISTS pickup_status;
DROP TYPE IF EXISTS waste_type;

DROP EXTENSION IF EXISTS "pgcrypto";