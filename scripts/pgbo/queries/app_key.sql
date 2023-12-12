-- name: GetAppKeyByName :one
SELECT
    ak.id,
    ak.name,
    ak.key
FROM app_key ak
WHERE
        ak.name = @name;