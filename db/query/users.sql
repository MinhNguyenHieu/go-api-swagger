-- Users Queries
-- name: CreateUser :one
INSERT INTO users (
    username,
    hashed_password,
    email,
    role_id
) VALUES (
    $1, $2, $3, $4
) RETURNING *;

-- name: GetUserByUsername :one
SELECT * FROM users
WHERE username = $1 LIMIT 1;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: UpdateUser :one
UPDATE users
SET
    username = $2,
    hashed_password = $3,
    email = $4,
    email_verified_at = $5,
    role_id = $6,
    remember_token_uuid = $7,
    updated_at = NOW(),
    deleted_at = $8
WHERE id = $1
RETURNING *;

-- name: UpdateUserRole :one
UPDATE users
SET
    role_id = $2,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: SoftDeleteUser :one
UPDATE users
SET
    deleted_at = NOW(),
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: RestoreUser :one
UPDATE users
SET
    deleted_at = NULL,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: VerifyUserEmail :one
UPDATE users
SET
    email_verified_at = NOW(),
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY id
LIMIT $2 OFFSET $1;

-- name: CountUsers :one
SELECT COUNT(*) FROM users;


-- Roles Queries
-- name: CreateRole :one
INSERT INTO roles (
    name,
    description
) VALUES (
    $1, $2
) RETURNING *;

-- name: GetRoleByID :one
SELECT * FROM roles
WHERE id = $1 LIMIT 1;

-- name: GetRoleByName :one
SELECT * FROM roles
WHERE name = $1 LIMIT 1;

-- name: UpdateRole :one
UPDATE roles
SET
    name = $2,
    description = $3,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteRole :exec
DELETE FROM roles
WHERE id = $1;

-- name: ListRoles :many
SELECT * FROM roles
ORDER BY id
LIMIT $2 OFFSET $1;

-- name: CountRoles :one
SELECT COUNT(*) FROM roles;


-- Password Reset Tokens Queries
-- name: CreatePasswordResetToken :one
INSERT INTO password_reset_tokens (
    email,
    token,
    created_at
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: GetPasswordResetToken :one
SELECT * FROM password_reset_tokens
WHERE email = $1 LIMIT 1;

-- name: DeletePasswordResetToken :exec
DELETE FROM password_reset_tokens
WHERE email = $1;


-- Sessions Queries
-- name: CreateSession :one
INSERT INTO sessions (
    id,
    user_id,
    ip_address,
    user_agent,
    payload,
    last_activity
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: GetSessionByID :one
SELECT * FROM sessions
WHERE id = $1 LIMIT 1;

-- name: UpdateSession :one
UPDATE sessions
SET
    user_id = $2,
    ip_address = $3,
    user_agent = $4,
    payload = $5,
    last_activity = $6
WHERE id = $1
RETURNING *;

-- name: DeleteSession :exec
DELETE FROM sessions
WHERE id = $1;

-- name: DeleteExpiredSessions :exec
DELETE FROM sessions
WHERE last_activity < $1;

