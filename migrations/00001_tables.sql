-- +goose Up

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL, --  'moderator', 'employee'
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE pvz (
    id UUID PRIMARY KEY NOT NULL,
    city VARCHAR(100) NOT NULL,
    registration_date TIMESTAMP NOT NULL
);

CREATE TABLE receptions (
    id UUID PRIMARY KEY NOT NULL,
    pvz_id UUID NOT NULL,
    started_at TIMESTAMP NOT NULL,
    closed_at  TIMESTAMP,
    status VARCHAR(50) NOT NULL, -- 'in_progress' или 'close'
    FOREIGN KEY (pvz_id) REFERENCES pvz (id) ON DELETE CASCADE
);

CREATE TABLE products (
    id UUID PRIMARY KEY NOT NULL,
    reception_id UUID NOT NULL,
    created_at TIMESTAMP NOT NULL,
    type VARCHAR(50) NOT NULL, -- 'электроника', 'одежда', 'обувь'
    FOREIGN KEY (reception_id) REFERENCES receptions (id) ON DELETE CASCADE
);



-- +goose Down
drop table users;
drop table goods;
drop table receptions;
drop table pvz;