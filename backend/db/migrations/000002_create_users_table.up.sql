CREATE TABLE IF NOT EXISTS users(
    id UUID PRIMARY KEY, 
    email TEXT NOT NULL UNIQUE,
    first_name TEXT NOT NULL, 
    last_name TEXT NOT NULL, 
    is_admin BOOLEAN DEFAULT FALSE,
    password TEXT NOT NULL, 
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);