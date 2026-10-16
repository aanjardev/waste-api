CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE households (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    owner_name  VARCHAR(255) NOT NULL,
    address     TEXT NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TYPE waste_type AS ENUM ('organic', 'plastic', 'paper', 'electronic');

CREATE TYPE pickup_status AS ENUM ('pending', 'scheduled', 'completed', 'canceled');

CREATE TABLE waste_pickups (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    household_id  UUID NOT NULL REFERENCES households(id) ON DELETE CASCADE,
    type          waste_type NOT NULL,
    status        pickup_status NOT NULL DEFAULT 'pending',
    pickup_date   TIMESTAMPTZ,
    safety_check  BOOLEAN NOT NULL DEFAULT FALSE,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_pickups_household_id ON waste_pickups(household_id);
CREATE INDEX idx_pickups_status ON waste_pickups(status);
CREATE INDEX idx_pickups_type ON waste_pickups(type);

CREATE TYPE payment_status AS ENUM ('pending', 'paid', 'failed');

CREATE TABLE payments (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    household_id    UUID NOT NULL REFERENCES households(id) ON DELETE CASCADE,
    waste_id        UUID NOT NULL UNIQUE REFERENCES waste_pickups(id) ON DELETE CASCADE,
    amount          NUMERIC(12,2) NOT NULL,
    payment_date    TIMESTAMPTZ,
    status          payment_status NOT NULL DEFAULT 'pending',
    proof_file_url  TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_payments_household_id ON payments(household_id);
CREATE INDEX idx_payments_status ON payments(status);
CREATE INDEX idx_payments_waste_id ON payments(waste_id);