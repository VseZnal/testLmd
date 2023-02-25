BEGIN;

CREATE SEQUENCE warehouse_id_seq;

CREATE TABLE warehouse (
    id     BIGINT PRIMARY KEY NOT NULL DEFAULT pseudo_encrypt(nextval('warehouse_id_seq')::BIGINT),
    name   TEXT NOT NULL,
    active BOOLEAN
);

COMMIT;