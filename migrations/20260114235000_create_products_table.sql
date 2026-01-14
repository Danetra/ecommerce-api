-- +migrate Up
CREATE TABLE products (
  id SERIAL PRIMARY KEY,

  category_id INTEGER NOT NULL REFERENCES product_categories(id) ON DELETE RESTRICT,

  name VARCHAR(150) NOT NULL,
  description TEXT,

  price NUMERIC(12,2) NOT NULL,
  stock INTEGER DEFAULT 0,

  image TEXT,

  is_active BOOLEAN DEFAULT TRUE,

  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  created_by INTEGER,

  updated_at TIMESTAMP,
  updated_by INTEGER
);

-- +migrate Down
DROP TABLE IF EXISTS products;
