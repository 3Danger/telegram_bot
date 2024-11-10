CREATE TYPE USER_TYPE AS ENUM ('', 'supplier', 'customer');

CREATE TABLE users
(
    id         BIGINT PRIMARY KEY NOT NULL,
    user_type  USER_TYPE          NOT NULL DEFAULT '',
    first_name VARCHAR(48)        NOT NULL,
    last_name  VARCHAR(48)        NOT NULL,
    phone      VARCHAR(48)        NOT NULL,
    tg_nick    VARCHAR(48)        NOT NULL,
    additional VARCHAR(256)       NOT NULL,
    created_at TIMESTAMPTZ        NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ        NOT NULL DEFAULT NOW()
);

COMMENT ON TABLE  users            IS 'Все пользователи';
COMMENT ON COLUMN users.user_type  IS 'Тип пользователя';
COMMENT ON COLUMN users.first_name IS 'Имя';
COMMENT ON COLUMN users.last_name  IS 'Фамилия';
COMMENT ON COLUMN users.phone      IS 'Телефон';
COMMENT ON COLUMN users.tg_nick    IS 'Ник в телеграмме';
COMMENT ON COLUMN users.additional IS 'Доп. контакты, инфо';
