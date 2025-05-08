CREATE TYPE USER_TYPE AS ENUM ('unknown', 'supplier', 'customer');
CREATE TYPE RECORDING_MODE AS ENUM ('draft', 'completed');

CREATE TABLE users
(
    id         BIGINT             NOT NULL,
    recording_mode RECORDING_MODE NOT NULL DEFAULT 'draft',
    user_type  USER_TYPE          NOT NULL,
    first_name VARCHAR(48)        NOT NULL,
    last_name  VARCHAR(48)        NOT NULL,
    phone      VARCHAR(48)        NOT NULL,
    additional VARCHAR(256)       NOT NULL,
    is_deleted BOOL               NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ        NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ        NOT NULL DEFAULT NOW(),
    CONSTRAINT users_pk PRIMARY KEY (id, recording_mode)
);

CREATE INDEX users_tm_created_at ON users USING BRIN(created_at);

CREATE INDEX users_tm_updated_at ON users (updated_at) WHERE recording_mode = 'draft';

COMMENT ON TABLE  users            IS 'Все пользователи';
COMMENT ON COLUMN users.id         IS 'ID пользователя';
COMMENT ON COLUMN users.recording_mode IS 'Статус регистрации';
COMMENT ON COLUMN users.user_type  IS 'Тип пользователя';
COMMENT ON COLUMN users.first_name IS 'Имя';
COMMENT ON COLUMN users.last_name  IS 'Фамилия';
COMMENT ON COLUMN users.phone      IS 'Телефон';
COMMENT ON COLUMN users.additional IS 'Доп. контакты, инфо';
COMMENT ON COLUMN users.created_at IS 'Дата создания записи';
COMMENT ON COLUMN users.updated_at IS 'Дата обновления записи';
