CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    deleted_at TIMESTAMP NULL DEFAULT NULL
);

-- Avoid duplicate entries in the email column
ALTER TABLE users
ADD CONSTRAINT users_email_unique UNIQUE (email);

-- unique index to ensure case-insensitive uniqueness
CREATE UNIQUE INDEX IF NOT EXISTS users_email_lower_unique
ON users (LOWER(TRIM(email)))
WHERE deleted_at IS NULL;

INSERT INTO users (name, email) VALUES
('John Doe', 'john@example.com'),
('Jane Smith', 'jane@example.com');