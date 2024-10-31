-- name: UpsertUser :exec
INSERT INTO users (id, has_registered, is_supplier, first_name, last_name, surname, telephone, whatsapp, telegram)
VALUES (
        @id,
        @has_registered,
        @is_supplier,
        @first_name,
        @last_name,
        @surname,
        @telephone,
        @whatsapp,
        @telegram
)
ON CONFLICT (id) DO UPDATE SET
    has_registered = EXCLUDED.has_registered,
    is_supplier    = EXCLUDED.is_supplier,
    first_name     = EXCLUDED.first_name,
    last_name      = EXCLUDED.last_name,
    surname        = EXCLUDED.surname,
    telephone      = EXCLUDED.telephone,
    whatsapp       = EXCLUDED.whatsapp,
    telegram       = EXCLUDED.telegram,
    updated_at     = EXCLUDED.updated_at;

-- name: User :one
SELECT * FROM users WHERE id = @id;


