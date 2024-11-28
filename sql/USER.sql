CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    firstname VARCHAR(100) NOT NULL,
    middlename VARCHAR(100),
    lastname VARCHAR(100) NOT NULL,
    age INTEGER CHECK (age >= 0),
    gender VARCHAR(50),
    email VARCHAR(255) UNIQUE NOT NULL,
    password BYTEA NOT NULL,
    authorities TEXT[] DEFAULT ARRAY[]::TEXT[],
    status VARCHAR(50) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX users_email_idx ON users(email);
