CREATE TABLE cars (
    id SERIAL PRIMARY KEY,
    make VARCHAR(100) NOT NULL,
    model VARCHAR(100) NOT NULL,
    year INTEGER,
    owner_id INTEGER NOT NULL,
    previous_owners_count INTEGER NOT NULL DEFAULT 0,
    currency VARCHAR(3) NOT NULL,
    price INTEGER NOT NULL CHECK (price > 0),
    options TEXT[],
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- пока не надо
-- CREATE INDEX idx_cars_owner_id ON cars(owner_id);