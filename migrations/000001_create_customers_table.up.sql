CREATE TABLE IF NOT EXISTS customers (
  id serial PRIMARY KEY,
  bank_limit integer NOT NULL,
  bank_balance integer DEFAULT 0
);
