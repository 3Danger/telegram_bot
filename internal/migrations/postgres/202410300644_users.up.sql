CREATE TABLE users
(
    id             BIGINT PRIMARY KEY NOT NULL,
    has_registered BOOLEAN            NOT NULL DEFAULT false,
    is_supplier    BOOLEAN            NOT NULL,
    first_name     VARCHAR(48)        NOT NULL,
    last_name      VARCHAR(48)        NOT NULL,
    surname        VARCHAR(48)        NOT NULL,
    telephone      VARCHAR(48)        NOT NULL,
    whatsapp       VARCHAR(48)        NOT NULL,
    telegram       VARCHAR(48)        NOT NULL,
    created_at     TIMESTAMPTZ        NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ        NOT NULL DEFAULT NOW()
);

COMMENT ON TABLE  users                 IS 'Все пользователи';
COMMENT ON COLUMN users.has_registered  IS 'Заполнены все полня и регистрация завершена';
COMMENT ON COLUMN users.is_supplier     IS 'Является ли поставщиком';
COMMENT ON COLUMN users.first_name      IS 'Имя';
COMMENT ON COLUMN users.last_name       IS 'Фамилия';
COMMENT ON COLUMN users.surname         IS 'Отчество';
COMMENT ON COLUMN users.Telephone       IS 'Телефон';
COMMENT ON COLUMN users.Whatsapp        IS 'Номер whatsapp';
COMMENT ON COLUMN users.Telegram        IS 'Ник в телеграмме';
