-- +migrate Up
CREATE TABLE users (
   id SERIAL PRIMARY KEY,
   role_id INTEGER NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
   username VARCHAR(100) NOT NULL UNIQUE,
   password TEXT NOT NULL,
   name VARCHAR(50) NOT NULL,
   email VARCHAR(50),
   is_active BOOLEAN DEFAULT TRUE,

   created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
   created_by INTEGER REFERENCES users(id),

   updated_at TIMESTAMP,
   updated_by INTEGER REFERENCES users(id)
);

-- +migrate Down
DROP TABLE IF EXISTS users;
