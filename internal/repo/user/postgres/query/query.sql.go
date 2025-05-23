// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: query.sql

package query

import (
	"context"
)

const approveChanges = `-- name: ApproveChanges :exec
WITH changes AS (
    DELETE FROM users
    WHERE users.id = $1
      AND recording_mode = 'draft'
    RETURNING id, recording_mode, user_type, first_name, last_name, phone, additional, is_deleted, created_at, updated_at
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
    updated_at = NOW()
`

func (q *Queries) ApproveChanges(ctx context.Context, id int64) error {
	_, err := q.db.Exec(ctx, approveChanges, id)
	return err
}

const delete = `-- name: Delete :exec
DELETE FROM users WHERE id = $1 AND recording_mode = $2
`

type DeleteParams struct {
	ID            int64         `json:"id"`
	RecordingMode RecordingMode `json:"recording_mode"`
}

func (q *Queries) Delete(ctx context.Context, arg *DeleteParams) error {
	_, err := q.db.Exec(ctx, delete, arg.ID, arg.RecordingMode)
	return err
}

const get = `-- name: Get :one
SELECT id, recording_mode, user_type, first_name, last_name, phone, additional, is_deleted, created_at, updated_at FROM users WHERE id = $1 AND recording_mode = $2
`

type GetParams struct {
	ID            int64         `json:"id"`
	RecordingMode RecordingMode `json:"recording_mode"`
}

func (q *Queries) Get(ctx context.Context, arg *GetParams) (*User, error) {
	row := q.db.QueryRow(ctx, get, arg.ID, arg.RecordingMode)
	var i User
	err := row.Scan(
		&i.ID,
		&i.RecordingMode,
		&i.UserType,
		&i.FirstName,
		&i.LastName,
		&i.Phone,
		&i.Additional,
		&i.IsDeleted,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return &i, err
}

const getDraft = `-- name: GetDraft :one
SELECT id, recording_mode, user_type, first_name, last_name, phone, additional, is_deleted, created_at, updated_at FROM users WHERE id = $1 AND recording_mode = 'draft'
`

func (q *Queries) GetDraft(ctx context.Context, id int64) (*User, error) {
	row := q.db.QueryRow(ctx, getDraft, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.RecordingMode,
		&i.UserType,
		&i.FirstName,
		&i.LastName,
		&i.Phone,
		&i.Additional,
		&i.IsDeleted,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return &i, err
}

const setAdditional = `-- name: SetAdditional :exec
UPDATE users SET additional = $1 WHERE id = $2 AND recording_mode = 'draft'
`

type SetAdditionalParams struct {
	Additional string `json:"additional"`
	ID         int64  `json:"id"`
}

func (q *Queries) SetAdditional(ctx context.Context, arg *SetAdditionalParams) error {
	_, err := q.db.Exec(ctx, setAdditional, arg.Additional, arg.ID)
	return err
}

const setFirstName = `-- name: SetFirstName :exec
UPDATE users SET first_name = $1 WHERE id = $2 AND recording_mode = 'draft'
`

type SetFirstNameParams struct {
	FirstName string `json:"first_name"`
	ID        int64  `json:"id"`
}

func (q *Queries) SetFirstName(ctx context.Context, arg *SetFirstNameParams) error {
	_, err := q.db.Exec(ctx, setFirstName, arg.FirstName, arg.ID)
	return err
}

const setLastName = `-- name: SetLastName :exec
UPDATE users SET last_name = $1 WHERE id = $2 AND recording_mode = 'draft'
`

type SetLastNameParams struct {
	LastName string `json:"last_name"`
	ID       int64  `json:"id"`
}

func (q *Queries) SetLastName(ctx context.Context, arg *SetLastNameParams) error {
	_, err := q.db.Exec(ctx, setLastName, arg.LastName, arg.ID)
	return err
}

const setPhone = `-- name: SetPhone :exec
UPDATE users SET phone = $1 WHERE id = $2 AND recording_mode = 'draft'
`

type SetPhoneParams struct {
	Phone string `json:"phone"`
	ID    int64  `json:"id"`
}

func (q *Queries) SetPhone(ctx context.Context, arg *SetPhoneParams) error {
	_, err := q.db.Exec(ctx, setPhone, arg.Phone, arg.ID)
	return err
}

const setUserType = `-- name: SetUserType :exec
UPDATE users SET user_type = $1 WHERE id = $2 AND recording_mode = 'draft'
`

type SetUserTypeParams struct {
	UserType UserType `json:"user_type"`
	ID       int64    `json:"id"`
}

func (q *Queries) SetUserType(ctx context.Context, arg *SetUserTypeParams) error {
	_, err := q.db.Exec(ctx, setUserType, arg.UserType, arg.ID)
	return err
}

const upsert = `-- name: Upsert :exec
INSERT INTO users (id, recording_mode, user_type, first_name, last_name, phone, additional)
VALUES ($1,
        $2,
        $3,
        $4,
        $5,
        $6,
        $7)
ON CONFLICT (id, recording_mode) DO UPDATE SET
    user_type  = EXCLUDED.user_type,
    first_name = EXCLUDED.first_name,
    last_name  = EXCLUDED.last_name,
    phone      = EXCLUDED.phone,
    additional = EXCLUDED.additional,
    updated_at = NOW()
`

type UpsertParams struct {
	ID            int64         `json:"id"`
	RecordingMode RecordingMode `json:"recording_mode"`
	UserType      UserType      `json:"user_type"`
	FirstName     string        `json:"first_name"`
	LastName      string        `json:"last_name"`
	Phone         string        `json:"phone"`
	Additional    string        `json:"additional"`
}

func (q *Queries) Upsert(ctx context.Context, arg *UpsertParams) error {
	_, err := q.db.Exec(ctx, upsert,
		arg.ID,
		arg.RecordingMode,
		arg.UserType,
		arg.FirstName,
		arg.LastName,
		arg.Phone,
		arg.Additional,
	)
	return err
}
