-- +migrate Up
CREATE TABLE roles (
   id SERIAL PRIMARY KEY,
   name VARCHAR(50) NOT NULL UNIQUE,
   description TEXT,
   created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
   updated_at TIMESTAMP
);

-- +migrate Down
DROP TABLE IF EXISTS roles;
