
-- name: Upsert :exec
INSERT INTO users (id, user_type, first_name, last_name, phone, tg_nick, additional)
VALUES (@id,
        @user_type,
        @first_name,
        @last_name,
        @phone,
        @tg_nick,
        @additional)
ON CONFLICT (id) DO UPDATE SET user_type  = EXCLUDED.user_type,
                               first_name = EXCLUDED.first_name,
                               last_name  = EXCLUDED.last_name,
                               phone      = EXCLUDED.phone,
                               tg_nick    = EXCLUDED.tg_nick,
                               additional = EXCLUDED.additional,
                               updated_at = NOW();

-- name: Get :one
SELECT * FROM users WHERE id = @id;

-- name: Delete :exec
DELETE FROM users WHERE id = @id;
