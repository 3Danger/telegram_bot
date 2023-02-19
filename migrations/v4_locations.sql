CREATE SCHEMA IF NOT EXISTS locations;

CREATE TABLE IF NOT EXISTS locations.countries(
    id      serial      not null primary key,
    name    varchar(64) not null
);

CREATE TABLE IF NOT EXISTS locations.regions (
    id          serial      not null primary key,
    id_country  int         not null references locations.countries(id),
    name        varchar(64) not null
);

CREATE TABLE IF NOT EXISTS locations.city (
    id          serial not null primary key,
    id_region   int not null references locations.regions(id),
    name        varchar(48) not null
)