-- name: InsertAuthToken :one
INSERT INTO auth_token
    (name, device_id, device_type, token, token_expired, refresh_token, refresh_token_expired, is_login, user_login, created_at)
VALUES
    (@name, @device_id, @device_type, @token, @token_expired, @refresh_token, @refresh_token_expired, @is_login, @user_login, (now() at time zone 'UTC')::TIMESTAMP)
ON CONFLICT (name, device_id, device_type) DO
UPDATE
    SET token = @token,
    token_expired = @token_expired,
    refresh_token = @refresh_token,
    refresh_token_expired = @refresh_token_expired,
    is_login = @is_login,
    user_login = @user_login,
    updated_at = (now() at time zone 'UTC')::TIMESTAMP
RETURNING auth_token.*;

-- name: RecordAuthTokenUserLogin :exec
UPDATE auth_token
SET
    is_login = true,
    user_login = @user_login,
    updated_at = (now() at time zone 'UTC')::TIMESTAMP
WHERE
    name = @name
    AND device_id = @device_id
    AND device_type = @device_type;

-- name: ClearAuthTokenUserLogin :exec
UPDATE auth_token
SET
    is_login = false,
    user_login = null,
    updated_at = (now() at time zone 'UTC')::TIMESTAMP
WHERE
    name = @name
  AND device_id = @device_id
  AND device_type = @device_type;

-- name: GetAuthToken :one
SELECT t.*
FROM auth_token t
WHERE
    t.name = @name
    AND t.device_id = @device_id
    AND t.device_type = @device_type;