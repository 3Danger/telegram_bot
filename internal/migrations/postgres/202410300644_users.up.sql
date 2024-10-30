CREATE TABLE users
(
    id          BIGINT PRIMARY KEY NOT NULL,
    is_supplier BOOLEAN            NOT NULL DEFAULT false,
    first_name  VARCHAR(48)        NOT NULL,
    last_name   VARCHAR(48)        NOT NULL,
    surname     VARCHAR(48)        NOT NULL,
    Telephone   VARCHAR(48)        NOT NULL,
    Whatsapp    VARCHAR(48)        NOT NULL,
    Telegram    VARCHAR(48)        NOT NULL
);

COMMENT ON TABLE  users             IS 'Все пользователи';
COMMENT ON COLUMN users.is_supplier IS 'Является ли поставщиком';
COMMENT ON COLUMN users.first_name  IS 'Имя';
COMMENT ON COLUMN users.last_name   IS 'Фамилия';
COMMENT ON COLUMN users.surname     IS 'Отчество';
COMMENT ON COLUMN users.Telephone   IS 'Телефон';
COMMENT ON COLUMN users.Whatsapp    IS 'Номер whatsapp';
COMMENT ON COLUMN users.Telegram    IS 'Ник в телеграмме';
