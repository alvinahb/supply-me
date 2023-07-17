-- name: CreateCompany :one
INSERT INTO companies (
    company_type, company_name, owner
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: GetCompany :one
SELECT * FROM companies WHERE id = $1 LIMIT 1;

-- name: ListCompanies :many
SELECT * FROM companies ORDER BY id LIMIT $1 OFFSET $2;

-- name: UpdateCompany :exec
UPDATE companies SET company_name = $1 WHERE id = $2;

-- name: DeleteCompany :exec
DELETE FROM companies where id = $1;
