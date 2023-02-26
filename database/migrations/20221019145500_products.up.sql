BEGIN;

CREATE SEQUENCE products_id_seq;

CREATE TABLE products (
    id           BIGINT PRIMARY KEY NOT NULL DEFAULT pseudo_encrypt(nextval('products_id_seq')::BIGINT),
    name         TEXT NOT NULL,
    size         TEXT NOT NULL,
    quantity     INTEGER,
    in_reserve   INTEGER DEFAULT 0,
    warehouse_id INTEGER NOT NULL REFERENCES warehouse(id)
);

COMMIT;