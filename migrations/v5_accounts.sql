CREATE SCHEMA IF NOT EXISTS accounts;

CREATE TABLE IF NOT EXISTS accounts.administrator (
    id          serial      not null primary key,
    id_telegram int         not null,
    name        varchar(32) not null
);

CREATE TABLE IF NOT EXISTS accounts.suppliers (
    id             serial      not null primary key,
    id_telegram    int         not null,
    id_city        int         not null references locations.city(id) ,
    created_date   timestamp   not null,
    first_name     varchar(32) not null,
    second_name    varchar(32) not null
)