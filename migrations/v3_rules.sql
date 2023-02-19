CREATE SCHEMA IF NOT EXISTS rules;

CREATE TABLE IF NOT EXISTS rules.commissions (
    id  serial not null primary key
);

CREATE TABLE IF NOT EXISTS rules.commissions_step (
    id              serial  not null primary key,
    id_commission   int     not null references rules.commissions(id),
    before          int     not null,
    commission      int     not null
);