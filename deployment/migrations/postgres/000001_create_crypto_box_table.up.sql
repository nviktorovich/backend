CREATE TABLE IF NOT EXISTS crypto_box (
    id SERIAL PRIMARY KEY,
    title TEXT,
    short_title TEXT NOT NULL,
    cost REAL,
    timestamp TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE crypto_box RENAME COLUMN timestamp TO created;