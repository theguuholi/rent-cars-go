-- name: ListCategories :many
SELECT name, description FROM categories ORDER BY name ASC;
