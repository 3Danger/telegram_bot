-- name: CreateUser :exec
INSERT INTO users (id, is_supplier, first_name, last_name, surname)
VALUES (
        @id,
        @is_supplier,
        @first_name,
        @last_name,
        @surname
);

-- name: UpdateUserContactTelegram :execrows
UPDATE users SET telegram = @telegram WHERE id = @id;

-- name: UpdateUserContactWhatsapp :execrows
UPDATE users SET whatsapp = @whatsapp WHERE id = @id;

-- name: UpdateUserContactPhone :execrows
UPDATE users SET phone = @phone WHERE id = @id;

-- name: User :one
SELECT * FROM users WHERE id = @id;


