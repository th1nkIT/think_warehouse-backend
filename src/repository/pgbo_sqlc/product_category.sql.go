// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: product_category.sql

package sqlc

import (
	"context"
	"database/sql"
	"time"
)

const deleteProductCategory = `-- name: DeleteProductCategory :exec
UPDATE product_category
SET
    deleted_at = (now() at time zone 'UTC')::TIMESTAMP,
    deleted_by = $1
WHERE
    guid = $2
  AND deleted_at IS NULL
`

type DeleteProductCategoryParams struct {
	DeletedBy sql.NullString `json:"deleted_by"`
	Guid      string         `json:"guid"`
}

func (q *Queries) DeleteProductCategory(ctx context.Context, arg DeleteProductCategoryParams) error {
	_, err := q.db.ExecContext(ctx, deleteProductCategory, arg.DeletedBy, arg.Guid)
	return err
}

const getCountListProductCategory = `-- name: GetCountListProductCategory :one
SELECT
    count(pc.id) AS total_data
FROM
    product_category pc
WHERE
    (CASE WHEN $1::bool THEN LOWER(pc.name) LIKE LOWER ($2) ELSE TRUE END)
  AND (CASE WHEN $3::bool THEN
                (pc.deleted_at IS NULL AND $4 = 'active') OR
                (pc.deleted_at IS NOT NULL AND $4 = 'inactive')
            ELSE TRUE END)
  AND pc.deleted_at IS NULL
`

type GetCountListProductCategoryParams struct {
	SetName   bool        `json:"set_name"`
	Name      string      `json:"name"`
	SetActive bool        `json:"set_active"`
	Active    interface{} `json:"active"`
}

func (q *Queries) GetCountListProductCategory(ctx context.Context, arg GetCountListProductCategoryParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, getCountListProductCategory,
		arg.SetName,
		arg.Name,
		arg.SetActive,
		arg.Active,
	)
	var total_data int64
	err := row.Scan(&total_data)
	return total_data, err
}

const getProductCategory = `-- name: GetProductCategory :one
SELECT
    pc.guid, pc.name, pc.created_at, pc.created_by,
    pc.updated_at, pc.updated_by, pc.deleted_at, pc.deleted_by,
    ub_created.name AS user_name, ub_created.guid AS user_id,
    ub_updated.name AS user_name_update, ub_updated.guid AS user_id_update
FROM
    product_category pc
        LEFT JOIN user_backoffice ub_created ON ub_created.guid = pc.created_by
        LEFT JOIN user_backoffice ub_updated ON ub_updated.guid = pc.updated_by
WHERE
    pc.guid = $1
`

type GetProductCategoryRow struct {
	Guid           string         `json:"guid"`
	Name           string         `json:"name"`
	CreatedAt      time.Time      `json:"created_at"`
	CreatedBy      string         `json:"created_by"`
	UpdatedAt      sql.NullTime   `json:"updated_at"`
	UpdatedBy      sql.NullString `json:"updated_by"`
	DeletedAt      sql.NullTime   `json:"deleted_at"`
	DeletedBy      sql.NullString `json:"deleted_by"`
	UserName       sql.NullString `json:"user_name"`
	UserID         sql.NullString `json:"user_id"`
	UserNameUpdate sql.NullString `json:"user_name_update"`
	UserIDUpdate   sql.NullString `json:"user_id_update"`
}

func (q *Queries) GetProductCategory(ctx context.Context, guid string) (GetProductCategoryRow, error) {
	row := q.db.QueryRowContext(ctx, getProductCategory, guid)
	var i GetProductCategoryRow
	err := row.Scan(
		&i.Guid,
		&i.Name,
		&i.CreatedAt,
		&i.CreatedBy,
		&i.UpdatedAt,
		&i.UpdatedBy,
		&i.DeletedAt,
		&i.DeletedBy,
		&i.UserName,
		&i.UserID,
		&i.UserNameUpdate,
		&i.UserIDUpdate,
	)
	return i, err
}

const insertProductCategory = `-- name: InsertProductCategory :one
INSERT INTO product_category
(guid, name, created_at, created_by)
VALUES
    ($1, $2, (now() at time zone 'UTC')::TIMESTAMP, $3)
RETURNING product_category.id, product_category.guid, product_category.name, product_category.created_at, product_category.created_by, product_category.updated_at, product_category.updated_by, product_category.deleted_at, product_category.deleted_by
`

type InsertProductCategoryParams struct {
	Guid      string `json:"guid"`
	Name      string `json:"name"`
	CreatedBy string `json:"created_by"`
}

func (q *Queries) InsertProductCategory(ctx context.Context, arg InsertProductCategoryParams) (ProductCategory, error) {
	row := q.db.QueryRowContext(ctx, insertProductCategory, arg.Guid, arg.Name, arg.CreatedBy)
	var i ProductCategory
	err := row.Scan(
		&i.ID,
		&i.Guid,
		&i.Name,
		&i.CreatedAt,
		&i.CreatedBy,
		&i.UpdatedAt,
		&i.UpdatedBy,
		&i.DeletedAt,
		&i.DeletedBy,
	)
	return i, err
}

const listProductCategory = `-- name: ListProductCategory :many
SELECT
    pc.guid, pc.name, pc.created_at, pc.created_by,
    pc.updated_at, pc.updated_by, pc.deleted_at, pc.deleted_by,
    ub_created.name AS user_name, ub_created.guid AS user_id,
    ub_updated.name AS user_name_update, ub_updated.guid AS user_id_update
FROM
    product_category pc
        LEFT JOIN user_backoffice ub_created ON ub_created.guid = pc.created_by
        LEFT JOIN user_backoffice ub_updated ON ub_updated.guid = pc.updated_by
WHERE
    (CASE WHEN $1::bool THEN LOWER(pc.name) LIKE LOWER ($2) ELSE TRUE END)
  AND (CASE WHEN $3::bool THEN
                (pc.deleted_at IS NULL AND $4 = 'active') OR
                (pc.deleted_at IS NOT NULL AND $4 = 'inactive')
            ELSE TRUE END)
ORDER BY (CASE WHEN $5 = 'id ASC' THEN pc.guid END) ASC,
         (CASE WHEN $5 = 'id DESC' THEN pc.guid END) DESC,
         (CASE WHEN $5 = 'name ASC' THEN pc.name END) ASC,
         (CASE WHEN $5 = 'name DESC' THEN pc.name END) DESC,
         (CASE WHEN $5 = 'created_at ASC' THEN pc.created_at END) ASC,
         (CASE WHEN $5 = 'created_at DESC' THEN pc.created_at END) DESC,
         pc.created_at DESC
LIMIT $7
    OFFSET $6
`

type ListProductCategoryParams struct {
	SetName    bool        `json:"set_name"`
	Name       string      `json:"name"`
	SetActive  bool        `json:"set_active"`
	Active     interface{} `json:"active"`
	OrderParam interface{} `json:"order_param"`
	OffsetPage int32       `json:"offset_page"`
	LimitData  int32       `json:"limit_data"`
}

type ListProductCategoryRow struct {
	Guid           string         `json:"guid"`
	Name           string         `json:"name"`
	CreatedAt      time.Time      `json:"created_at"`
	CreatedBy      string         `json:"created_by"`
	UpdatedAt      sql.NullTime   `json:"updated_at"`
	UpdatedBy      sql.NullString `json:"updated_by"`
	DeletedAt      sql.NullTime   `json:"deleted_at"`
	DeletedBy      sql.NullString `json:"deleted_by"`
	UserName       sql.NullString `json:"user_name"`
	UserID         sql.NullString `json:"user_id"`
	UserNameUpdate sql.NullString `json:"user_name_update"`
	UserIDUpdate   sql.NullString `json:"user_id_update"`
}

func (q *Queries) ListProductCategory(ctx context.Context, arg ListProductCategoryParams) ([]ListProductCategoryRow, error) {
	rows, err := q.db.QueryContext(ctx, listProductCategory,
		arg.SetName,
		arg.Name,
		arg.SetActive,
		arg.Active,
		arg.OrderParam,
		arg.OffsetPage,
		arg.LimitData,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListProductCategoryRow
	for rows.Next() {
		var i ListProductCategoryRow
		if err := rows.Scan(
			&i.Guid,
			&i.Name,
			&i.CreatedAt,
			&i.CreatedBy,
			&i.UpdatedAt,
			&i.UpdatedBy,
			&i.DeletedAt,
			&i.DeletedBy,
			&i.UserName,
			&i.UserID,
			&i.UserNameUpdate,
			&i.UserIDUpdate,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const reactiveProductCategory = `-- name: ReactiveProductCategory :exec
UPDATE product_category
SET
    deleted_at = NULL,
    deleted_by = NULL,
    updated_at = (now() at time zone 'UTC')::TIMESTAMP,
    updated_by = $1
WHERE
    guid = $2
  AND deleted_at IS NOT NULL
`

type ReactiveProductCategoryParams struct {
	UpdatedBy sql.NullString `json:"updated_by"`
	Guid      string         `json:"guid"`
}

func (q *Queries) ReactiveProductCategory(ctx context.Context, arg ReactiveProductCategoryParams) error {
	_, err := q.db.ExecContext(ctx, reactiveProductCategory, arg.UpdatedBy, arg.Guid)
	return err
}

const updateProductCategory = `-- name: UpdateProductCategory :one
UPDATE product_category
SET
    name = $1,
    updated_at = (now() at time zone 'UTC')::TIMESTAMP,
    updated_by = $2
WHERE
    guid = $3
  AND deleted_at IS NULL
RETURNING product_category.id, product_category.guid, product_category.name, product_category.created_at, product_category.created_by, product_category.updated_at, product_category.updated_by, product_category.deleted_at, product_category.deleted_by
`

type UpdateProductCategoryParams struct {
	Name      string         `json:"name"`
	UpdatedBy sql.NullString `json:"updated_by"`
	Guid      string         `json:"guid"`
}

func (q *Queries) UpdateProductCategory(ctx context.Context, arg UpdateProductCategoryParams) (ProductCategory, error) {
	row := q.db.QueryRowContext(ctx, updateProductCategory, arg.Name, arg.UpdatedBy, arg.Guid)
	var i ProductCategory
	err := row.Scan(
		&i.ID,
		&i.Guid,
		&i.Name,
		&i.CreatedAt,
		&i.CreatedBy,
		&i.UpdatedAt,
		&i.UpdatedBy,
		&i.DeletedAt,
		&i.DeletedBy,
	)
	return i, err
}
