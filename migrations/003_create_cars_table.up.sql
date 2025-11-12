CREATE TABLE cars (
    id SERIAL PRIMARY KEY,
    make VARCHAR(100) NOT NULL,
    model VARCHAR(100) NOT NULL,
    year INTEGER,
    owner_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    previous_owners_count INTEGER NOT NULL DEFAULT 0,
    currency VARCHAR(3) NOT NULL,
    options TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT chk_currency CHECK (currency IN ('RUB', 'USD', 'EUR')),
    CONSTRAINT chk_options_not_empty CHECK (options <> ''),
    CONSTRAINT chk_previous_owners_count_positive CHECK (previous_owners_count >= 0)
);

-- пока не надо
-- CREATE INDEX idx_cars_owner_id ON cars(owner_id);