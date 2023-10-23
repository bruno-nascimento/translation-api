-- migrate:up
ALTER TABLE languages ALTER COLUMN region DROP NOT NULL;

-- migrate:down
ALTER TABLE languages ALTER COLUMN region SET NOT NULL;
