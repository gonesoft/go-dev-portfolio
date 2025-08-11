-- Create a sample database table for users
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Insert a sample row
INSERT INTO users (name, email)
VALUES ('John Doe', 'john@example.com')
ON CONFLICT DO NOTHING;