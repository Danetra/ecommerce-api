-- +migrate Up
CREATE TABLE transactions (
  id SERIAL PRIMARY KEY,

  product_id INTEGER NOT NULL REFERENCES products(id),
  category_id INTEGER NOT NULL REFERENCES product_categories(id),

  buyer_id INTEGER NOT NULL REFERENCES users(id),
  seller_id INTEGER NOT NULL REFERENCES users(id),

  qty INTEGER NOT NULL CHECK (qty > 0),
  price NUMERIC(12,2) NOT NULL,
  total NUMERIC(12,2) NOT NULL,

  status VARCHAR(30) DEFAULT 'PENDING',

  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP
);

-- +migrate Down
DROP TABLE IF EXISTS transactions;
