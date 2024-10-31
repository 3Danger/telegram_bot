CREATE TABLE users
(
    id             BIGINT PRIMARY KEY NOT NULL,
    is_supplier    BOOLEAN            NOT NULL,
    first_name     VARCHAR(48)        NOT NULL,
    last_name      VARCHAR(48)        NOT NULL,
    surname        VARCHAR(48)        NOT NULL,
    phone          VARCHAR(48)        NOT NULL DEFAULT '',
    whatsapp       VARCHAR(48)        NOT NULL DEFAULT '',
    telegram       VARCHAR(48)        NOT NULL DEFAULT '',
    created_at     TIMESTAMPTZ        NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ        NOT NULL DEFAULT NOW()
);

COMMENT ON TABLE  users                 IS 'Все пользователи';
COMMENT ON COLUMN users.is_supplier     IS 'Является ли поставщиком';
COMMENT ON COLUMN users.first_name      IS 'Имя';
COMMENT ON COLUMN users.last_name       IS 'Фамилия';
COMMENT ON COLUMN users.surname         IS 'Отчество';
COMMENT ON COLUMN users.phone           IS 'Телефон';
COMMENT ON COLUMN users.whatsapp        IS 'Номер whatsapp';
COMMENT ON COLUMN users.telegram        IS 'Ник в телеграмме';
