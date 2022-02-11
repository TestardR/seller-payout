BEGIN;

CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE sellers (
    id            UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at    TIMESTAMPTZ DEFAULT (now()),
    updated_at    TIMESTAMPTZ,

    currency_code VARCHAR(10)
);

CREATE TABLE currencies (
    id            UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at    TIMESTAMPTZ DEFAULT (now()),
    updated_at    TIMESTAMPTZ,

    code          VARCHAR(10),
    USD_exch_rate NUMERIC
);

CREATE TABLE items (
    id            UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at    TIMESTAMPTZ DEFAULT (now()),
    updated_at    TIMESTAMPTZ,

    reference_name VARCHAR(255),
    price_amount  NUMERIC,
    paid_out      BOOLEAN DEFAULT FALSE,
    currency_code VARCHAR(10),

    seller_id     UUID NOT NULL REFERENCES sellers(id) ON DELETE CASCADE
);

CREATE TABLE payouts (
    id          UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at  TIMESTAMPTZ DEFAULT (now()),
    updated_at  TIMESTAMPTZ,
    
    price_total NUMERIC,

    currency_id UUID NOT NULL REFERENCES currencies(id),
    seller_id   UUID NOT NULL REFERENCES sellers(id)
);


CREATE TABLE payout_items (
    id         UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ DEFAULT (now()),
    updated_at TIMESTAMPTZ,

    payout_id UUID NOT NULL REFERENCES payouts(id),
    item_id  UUID NOT NULL REFERENCES items(id)
);

COMMIT;