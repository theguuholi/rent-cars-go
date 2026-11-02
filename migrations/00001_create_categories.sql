-- +goose Up
CREATE TABLE categories (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    name varchar(50) NOT NULL UNIQUE CHECK (name = upper(name) AND name ~ '[^[:space:]]'),
    description varchar(255) CHECK (description ~ '[^[:space:]]'),
    inserted_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP
);
