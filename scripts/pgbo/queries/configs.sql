-- name: SetConfig :one
INSERT INTO configs(
	key, description, value, created_at)
	VALUES (@key, @description, @value, (now() at time zone 'UTC')::TIMESTAMP)
    ON CONFLICT ON CONSTRAINT key_config
    DO UPDATE SET description = @description, value = @value, created_at = (now() at time zone 'UTC')::TIMESTAMP
RETURNING configs.*;

-- name: UpdateConfig :one
UPDATE configs
	SET description = @description,
        value = @value,
        updated_at = (now() at time zone 'UTC')::TIMESTAMP, 
        updated_by = @updated_by
	WHERE id = @id
RETURNING configs.*;

-- name: GetConfig :one
SELECT *
    FROM configs
WHERE key = @key;

-- name: ListConfig :many
SELECT *
    FROM configs
ORDER BY id ASC;