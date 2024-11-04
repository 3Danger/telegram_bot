-- name: Set :exec
INSERT INTO users (id, user_type, first_name, last_name, surname, phone, whatsapp, telegram)
VALUES (
        @id,
        @user_type,
        @first_name,
        @last_name,
        @surname,
        @phone,
        @whatsapp,
        @telegram
)
ON CONFLICT (id) DO UPDATE SET
    user_type = EXCLUDED.user_type,
    first_name = EXCLUDED.first_name,
    last_name = EXCLUDED.last_name,
    surname = EXCLUDED.surname,
    phone = EXCLUDED.phone,
    whatsapp = EXCLUDED.whatsapp,
    telegram = EXCLUDED.telegram,
    updated_at = EXCLUDED.updated_at;

-- name: Get :one
SELECT * FROM users WHERE id = @id;

-- name: Delete :exec
DELETE FROM users WHERE id = @id;
