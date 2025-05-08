
-- name: Upsert :exec
INSERT INTO users (id, recording_mode, user_type, first_name, last_name, phone, additional)
VALUES (@id,
        @recording_mode,
        @user_type,
        @first_name,
        @last_name,
        @phone,
        @additional)
ON CONFLICT (id, recording_mode) DO UPDATE SET
    user_type  = EXCLUDED.user_type,
    first_name = EXCLUDED.first_name,
    last_name  = EXCLUDED.last_name,
    phone      = EXCLUDED.phone,
    additional = EXCLUDED.additional,
    updated_at = NOW();

-- name: SetFirstName :exec
UPDATE users SET first_name = @first_name WHERE id = @id AND recording_mode = 'draft';

-- name: SetLastName :exec
UPDATE users SET last_name = @last_name WHERE id = @id AND recording_mode = 'draft';

-- name: SetUserType :exec
UPDATE users SET user_type = @user_type WHERE id = @id AND recording_mode = 'draft';

-- name: SetPhone :exec
UPDATE users SET phone = @phone WHERE id = @id AND recording_mode = 'draft';

-- name: SetAdditional :exec
UPDATE users SET additional = @additional WHERE id = @id AND recording_mode = 'draft';

-- name: ApproveChanges :exec
WITH changes AS (
    DELETE FROM users
    WHERE users.id = @id
      AND recording_mode = 'draft'
    RETURNING *
)
INSERT INTO users (id, recording_mode, user_type, first_name, last_name, phone, additional)
SELECT id, 'completed', user_type, first_name, last_name, phone, additional
FROM changes
ON CONFLICT (id, recording_mode) DO UPDATE SET
    user_type  = EXCLUDED.user_type,
    first_name = EXCLUDED.first_name,
    last_name  = EXCLUDED.last_name,
    phone      = EXCLUDED.phone,
    additional = EXCLUDED.additional,
    updated_at = NOW();

-- name: Get :one
SELECT * FROM users WHERE id = @id AND recording_mode = @recording_mode;

-- name: GetDraft :one
SELECT * FROM users WHERE id = @id AND recording_mode = 'draft';

-- name: Delete :exec
DELETE FROM users WHERE id = @id AND recording_mode = @recording_mode;
