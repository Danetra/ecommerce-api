-- +migrate Up
ALTER TABLE transactions
    ADD COLUMN payment_method VARCHAR(50),
ADD COLUMN reference_id VARCHAR(100),
ADD COLUMN paid_at TIMESTAMP;

-- optional index (disarankan)
CREATE INDEX idx_transactions_reference_id ON transactions(reference_id);

-- +migrate Down
DROP INDEX IF EXISTS idx_transactions_reference_id;

ALTER TABLE transactions
DROP COLUMN IF EXISTS paid_at,
DROP COLUMN IF EXISTS reference_id,
DROP COLUMN IF EXISTS payment_method;
