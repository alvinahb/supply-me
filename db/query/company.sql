-- name: CreateCompany :one
INSERT INTO companies (
    company_type, company_name
) VALUES (
    $1, $2
) RETURNING *;

-- name: GetCompany :one
SELECT * FROM companies WHERE id = $1 LIMIT 1;

-- name: ListCompanies :many
SELECT * FROM companies ORDER BY id LIMIT $1 OFFSET $2;

-- name: UpdateCompany :one
UPDATE companies SET company_name = $1 WHERE id = $2
RETURNING *;

-- name: DeleteCompany :exec
DELETE FROM companies where id = $1;
