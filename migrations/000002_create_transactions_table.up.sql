CREATE TABLE transactions (
  id serial PRIMARY KEY,
  amount integer NOT NULL,
  type char(1) NOT NULL,
  description varchar(10) NOT NULL,
  customer_id integer NOT NULL REFERENCES customers ON DELETE CASCADE
);

CREATE INDEX ON transactions (customer_id);
