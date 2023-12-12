-- name: InsertUserBackofficeRole :one
INSERT INTO user_backoffice_role
        (name, access, is_all_access, created_at, created_by)
    VALUES
        (@name, @access, @is_all_access, (now() at time zone 'UTC')::TIMESTAMP, @created_by)
RETURNING user_backoffice_role.*;

-- name: UpdateUserBackofficeRole :one
UPDATE user_backoffice_role
SET
    name = @name,
    access = @access,
    is_all_access = @is_all_access,
    updated_at = (now() at time zone 'UTC')::TIMESTAMP,
    updated_by = @updated_by
WHERE
    id = @id
    AND deleted_at IS NULL
RETURNING user_backoffice_role.*;

-- name: DeleteUserBackofficeRole :exec
UPDATE user_backoffice_role
SET
    deleted_at = (now() at time zone 'UTC')::TIMESTAMP,
    deleted_by = @updated_by
WHERE
    id = @id
  AND deleted_at IS NULL;

-- name: ListUserBackofficeRole :many
SELECT ubr.* FROM user_backoffice_role ubr
WHERE
    (CASE WHEN @set_name::bool THEN LOWER(ubr.name) LIKE LOWER(@name) ELSE TRUE END)
    AND ubr.deleted_at IS NULL
ORDER BY (CASE WHEN @order_param = 'id ASC' THEN ubr.id END) ASC,
         (CASE WHEN @order_param = 'id DESC' THEN ubr.id  END) DESC,
         (CASE WHEN @order_param = 'name ASC' THEN ubr.name END) ASC,
         (CASE WHEN @order_param = 'name DESC' THEN ubr.name  END) DESC,
         (CASE WHEN @order_param = 'created_at ASC' THEN ubr.created_at END) ASC,
         (CASE WHEN @order_param = 'created_at DESC' THEN ubr.created_at  END) DESC,
         ubr.created_at DESC
LIMIT @limit_data
OFFSET @offset_page;

-- name: GetUserBackofficeRole :one
SELECT
       ubr.*
FROM
     user_backoffice_role ubr
WHERE ubr.id = @id;

-- name: GetCountListUserBackofficeRole :one
SELECT count(ubr.id) FROM user_backoffice_role ubr
WHERE
    (CASE WHEN @set_name::bool THEN LOWER(ubr.name) LIKE LOWER(@name) ELSE TRUE END)
    AND ubr.deleted_at IS NULL;