-- name: InsertUserHandheld :one
INSERT INTO user_handheld
    (guid, name, profile_picture_image_url, phone, email, gender, address, salt, password, is_active, fcm_token, created_at)
VALUES
    (@guid, @name, @profile_picture_image_url, @phone, @email, @gender, @address, @salt, @password, true, @fcm_token, (now() at time zone 'UTC')::TIMESTAMP)
RETURNING user_handheld.*;

-- name: RecordUserHandheldLastLogin :exec
UPDATE user_handheld
SET
    last_login = (now() at time zone 'UTC')::TIMESTAMP
WHERE
    guid = @guid
    AND deleted_at IS NULL;

-- name: UpdateUserHandheld :one
UPDATE user_handheld
SET
    name = @name,
    profile_picture_image_url = @profile_picture_image_url,
    phone = @phone,
    gender = @gender,
    address = @address,
    updated_at = (now() at time zone 'UTC')::TIMESTAMP
WHERE
    guid = @guid
    AND deleted_at IS NULL
RETURNING user_handheld.*;

-- name: UpdateUserHandheldFcmToken :one
UPDATE user_handheld
SET
    fcm_token = @fcm_token,
    updated_at = (now() at time zone 'UTC')::TIMESTAMP
WHERE
    guid = @guid
    AND deleted_at IS NULL
RETURNING user_handheld.*;

-- name: UpdateUserHandheldPassword :exec
UPDATE user_handheld
SET
    salt = @salt,
    password = @password,
    updated_at = (now() at time zone 'UTC')::TIMESTAMP
WHERE
    guid = @guid
    AND deleted_at IS NULL
RETURNING user_handheld.*;

-- name: UpdateUserHandheldIsActive :exec
UPDATE user_handheld
SET
    is_active = @is_active,
    updated_at = (now() at time zone 'UTC')::TIMESTAMP
WHERE
    guid = @guid
    AND deleted_at IS NULL;

-- name: DeleteUserHandheld :exec
UPDATE user_handheld
SET
    deleted_at = (now() at time zone 'UTC')::TIMESTAMP
WHERE
    guid = @guid
    AND deleted_at IS NULL;

-- name: ListUserHandheld :many
SELECT
    uh.*
FROM
    user_handheld uh
WHERE
    (CASE WHEN @set_name::bool THEN LOWER(uh.name) LIKE LOWER(@name) ELSE TRUE END)
    AND (CASE WHEN @set_phone::bool THEN LOWER(uh.phone) LIKE LOWER(@phone) ELSE TRUE END)
    AND (CASE WHEN @set_email::bool THEN LOWER(uh.email) LIKE LOWER(@email) ELSE TRUE END)
    AND (CASE WHEN @set_gender::bool THEN LOWER(uh.gender) LIKE LOWER(@gender) ELSE TRUE END)
    AND (CASE WHEN @set_address::bool THEN LOWER(uh.address) LIKE LOWER(@address) ELSE TRUE END)
    AND (CASE WHEN @set_is_active::bool THEN uh.is_active = @is_active ELSE TRUE END)
    AND uh.deleted_at IS NULL
ORDER BY (CASE WHEN @order_param = 'id ASC' THEN uh.guid END) ASC,
         (CASE WHEN @order_param = 'id DESC' THEN uh.guid  END) DESC,
         (CASE WHEN @order_param = 'name ASC' THEN uh.name END) ASC,
         (CASE WHEN @order_param = 'name DESC' THEN uh.name  END) DESC,
         (CASE WHEN @order_param = 'phone ASC' THEN uh.phone END) ASC,
         (CASE WHEN @order_param = 'phone DESC' THEN uh.phone  END) DESC,
         (CASE WHEN @order_param = 'email ASC' THEN uh.email END) ASC,
         (CASE WHEN @order_param = 'email DESC' THEN uh.email  END) DESC,
         (CASE WHEN @order_param = 'gender ASC' THEN uh.gender END) ASC,
         (CASE WHEN @order_param = 'gender DESC' THEN uh.gender  END) DESC,
         (CASE WHEN @order_param = 'address ASC' THEN uh.address END) ASC,
         (CASE WHEN @order_param = 'address DESC' THEN uh.address  END) DESC,
         (CASE WHEN @order_param = 'is_active ASC' THEN uh.is_active END) ASC,
         (CASE WHEN @order_param = 'is_active DESC' THEN uh.is_active  END) DESC,
         (CASE WHEN @order_param = 'created_at ASC' THEN uh.created_at END) ASC,
         (CASE WHEN @order_param = 'created_at DESC' THEN uh.created_at  END) DESC,
         uh.created_at DESC
LIMIT @limit_data
OFFSET @offset_page;


-- name: GetUserHandheld :one
SELECT
    uh.*
FROM user_handheld uh
WHERE
    uh.guid = @guid;

-- name: GetUserHandheldByEmail :one
SELECT
    uh.*
FROM user_handheld uh
WHERE
        uh.email = @email;

-- name: GetCountUserHandheld :one
SELECT
    count(uh.id)
FROM
    user_handheld uh
WHERE
    (CASE WHEN @set_name::bool THEN LOWER(uh.name) LIKE LOWER(@name) ELSE TRUE END)
  AND (CASE WHEN @set_phone::bool THEN LOWER(uh.phone) LIKE LOWER(@phone) ELSE TRUE END)
  AND (CASE WHEN @set_email::bool THEN LOWER(uh.email) LIKE LOWER(@email) ELSE TRUE END)
  AND (CASE WHEN @set_gender::bool THEN LOWER(uh.gender) LIKE LOWER(@gender) ELSE TRUE END)
  AND (CASE WHEN @set_address::bool THEN LOWER(uh.address) LIKE LOWER(@address) ELSE TRUE END)
  AND (CASE WHEN @set_is_active::bool THEN uh.is_active = @is_active ELSE TRUE END)
  AND uh.deleted_at IS NULL;