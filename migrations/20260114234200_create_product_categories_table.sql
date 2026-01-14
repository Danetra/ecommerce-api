-- +migrate Up
CREATE TABLE product_categories (
    id SERIAL PRIMARY KEY,

    name VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,

    is_active BOOLEAN DEFAULT TRUE,

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by INTEGER,

    updated_at TIMESTAMP,
    updated_by INTEGER
);

-- +migrate Down
DROP TABLE IF EXISTS product_categories;
