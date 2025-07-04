CREATE TABLE IF NOT EXISTS roles (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    level INTEGER NOT NULL
);

INSERT INTO roles (name, level) VALUES ('admin', 10);
INSERT INTO roles (name, level) VALUES ('customer', 1);

CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    password BYTEA NOT NULL,
    username VARCHAR(255) NOT NULL UNIQUE,
    role_id INTEGER NOT NULL REFERENCES roles(id),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Add nullable foreign key columns
ALTER TABLE staff ADD COLUMN user_id INTEGER REFERENCES users(id);
ALTER TABLE customer ADD COLUMN user_id INTEGER REFERENCES users(id);