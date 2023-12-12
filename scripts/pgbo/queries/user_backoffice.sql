-- name: InsertUserBackoffice :one
INSERT INTO user_backoffice
    (guid, name, profile_picture_image_url, phone, email, role_id, password, salt, is_active, created_at, created_by)
    VALUES
    (@guid, @name, @profile_picture_image_url, @phone, @email, @role_id, @password, @salt, @is_active,  (now() at time zone 'UTC')::TIMESTAMP, @created_by)
RETURNING user_backoffice.*;

-- name: RecordUserBackofficeLastLogin :exec
UPDATE user_backoffice
SET
    last_login = (now() at time zone 'UTC')::TIMESTAMP
WHERE
    guid = @guid
    AND deleted_at IS NULL;

-- name: UpdateUserBackoffice :one
UPDATE user_backoffice
SET
    name = @name,
    phone = @phone,
    profile_picture_image_url = @profile_picture_image_url,
    updated_at = (now() at time zone 'UTC')::TIMESTAMP,
    updated_by = @updated_by
WHERE
    guid = @guid
    AND deleted_at IS NULL
RETURNING user_backoffice.*;

-- name: UpdateUserBackofficePassword :exec
UPDATE user_backoffice
SET
    password = @password,
    salt = @salt,
    updated_at = (now() at time zone 'UTC')::TIMESTAMP,
    updated_by = @updated_by
WHERE
    guid = @guid
    AND deleted_at IS NULL
RETURNING user_backoffice.*;

-- name: UpdateUserBackofficeIsActive :exec
UPDATE user_backoffice
SET
    is_active = @is_active,
    updated_at = (now() at time zone 'UTC')::TIMESTAMP,
    updated_by = @updated_by
WHERE
    guid = @guid
    AND deleted_at IS NULL;

-- name: DeleteUserBackoffice :exec
UPDATE user_backoffice
SET
    deleted_at = (now() at time zone 'UTC')::TIMESTAMP,
    deleted_by = @deleted_by
WHERE
    guid = @guid
    AND deleted_at IS NULL;

-- name: ListUserBackoffice :many
SELECT
    ub.*,
    ubr.name as role_name,
    ubr.access as role_access,
    ubr.is_all_access as is_all_access
FROM
    user_backoffice ub
    JOIN user_backoffice_role ubr ON ubr.id = ub.role_id
WHERE
    (CASE WHEN @set_name::bool THEN LOWER(ub.name) LIKE LOWER(@name) ELSE TRUE END)
    AND(CASE WHEN @set_phone::bool THEN LOWER(ub.phone) LIKE LOWER(@phone) ELSE TRUE END)
    AND(CASE WHEN @set_email::bool THEN LOWER(ub.email) LIKE LOWER(@email) ELSE TRUE END)
    AND (CASE WHEN @set_role_id::bool THEN ub.role_id = @role_id ELSE TRUE END)
    AND (CASE WHEN @set_is_active::bool THEN ub.is_active = @is_active ELSE TRUE END)
    AND ub.deleted_at IS NULL
ORDER BY (CASE WHEN @order_param = 'id ASC' THEN ub.guid END) ASC,
         (CASE WHEN @order_param = 'id DESC' THEN ub.guid  END) DESC,
         (CASE WHEN @order_param = 'name ASC' THEN ub.name END) ASC,
         (CASE WHEN @order_param = 'name DESC' THEN ub.name  END) DESC,
         (CASE WHEN @order_param = 'role_id ASC' THEN ub.role_id END) ASC,
         (CASE WHEN @order_param = 'role_id DESC' THEN ub.role_id  END) DESC,
         (CASE WHEN @order_param = 'created_at ASC' THEN ub.created_at END) ASC,
         (CASE WHEN @order_param = 'created_at DESC' THEN ub.created_at  END) DESC,
         ub.created_at DESC
LIMIT @limit_data
OFFSET @offset_page;

-- name: GetUserBackoffice :one
SELECT
       ub.*,
       ubr.name as role_name,
       ubr.access as role_access,
       ubr.is_all_access as is_all_access
FROM user_backoffice ub
    JOIN user_backoffice_role ubr ON ubr.id = ub.role_id
WHERE
    ub.guid = @guid
    AND ub.deleted_at IS NULL;

-- name: GetUserBackofficeByEmail :one
select
    ub.*,
    ubr.name as role_name,
    ubr.access as role_access,
    ubr.is_all_access as is_all_access
FROM user_backoffice ub
    JOIN user_backoffice_role ubr ON ubr.id = ub.role_id
WHERE
    ub.email = @email;

-- name: GetCountListUserBackoffice :one
SELECT count(ub.id) FROM user_backoffice ub
WHERE
    (CASE WHEN @set_name::bool THEN LOWER(ub.name) LIKE LOWER(@name) ELSE TRUE END)
    AND(CASE WHEN @set_phone::bool THEN LOWER(ub.phone) LIKE LOWER(@phone) ELSE TRUE END)
    AND(CASE WHEN @set_email::bool THEN LOWER(ub.email) LIKE LOWER(@email) ELSE TRUE END)
    AND (CASE WHEN @set_role_id::bool THEN ub.role_id = @role_id ELSE TRUE END)
    AND (CASE WHEN @set_is_active::bool THEN ub.is_active = @is_active ELSE TRUE END)
    AND ub.deleted_at IS NULL;