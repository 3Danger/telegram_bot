CREATE SCHEMA IF NOT EXISTS products;

CREATE TABLE IF NOT EXISTS products.brands(
	id              serial not null primary key,
	name            varchar(48) not null,
	id_commission   int not null references rules.commissions(id)
);

CREATE TABLE IF NOT EXISTS products.products (
	id             serial           not null primary key,
	upload_date    timestamp        not null default now(),

	id_supplier    int              not null references accounts.suppliers(id),
	id_city        int              not null references locations.city(id),
	id_brand       int              not null references products.brands(id),

	is_box_weight  bool             not null default false,
	weight         numeric(14,2)    not null default 0,
	count          int              not null default 0,
	price          numeric(14,2)    not null,
	description    text             not null
);
