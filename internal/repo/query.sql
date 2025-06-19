-- name: CreateUser :one
INSERT INTO users (email, pw_hash)
VALUES (?, ?)
RETURNING *;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = ?;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = ?;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY created_at DESC;

-- name: UpdateUser :one
UPDATE users
SET email = ?, pw_hash = ?
WHERE id = ?
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = ?;

-- name: CreateSetting :one
INSERT INTO settings (key, value)
VALUES (?, ?)
ON CONFLICT(key) DO UPDATE SET value = excluded.value
RETURNING *;

-- name: GetSetting :one
SELECT * FROM settings
WHERE key = ?;

-- name: ListSettings :many
SELECT * FROM settings
ORDER BY key;

-- name: UpdateSetting :one
UPDATE settings
SET value = ?
WHERE key = ?
RETURNING *;

-- name: DeleteSetting :exec
DELETE FROM settings
WHERE key = ?;

-- name: CreateTag :one
INSERT INTO tags (name)
VALUES (?)
RETURNING *;

-- name: GetTagByID :one
SELECT * FROM tags
WHERE id = ?;

-- name: GetTagByName :one
SELECT * FROM tags
WHERE name = ?;

-- name: ListTags :many
SELECT * FROM tags
ORDER BY name;

-- name: UpdateTag :one
UPDATE tags
SET name = ?
WHERE id = ?
RETURNING *;

-- name: DeleteTag :exec
DELETE FROM tags
WHERE id = ?;

-- name: CreateTransaction :one
INSERT INTO transactions (user_id, amount_pence, t_date, note, source_recurring)
VALUES (?, ?, ?, ?, ?)
RETURNING *;

-- name: GetTransactionByID :one
SELECT * FROM transactions
WHERE id = ? AND deleted_at IS NULL;

-- name: ListTransactions :many
SELECT * FROM transactions
WHERE user_id = ? AND deleted_at IS NULL
  AND (t_date >= ? OR ? IS NULL)
  AND (t_date <= ? OR ? IS NULL)
ORDER BY t_date DESC, created_at DESC;

-- name: ListTransactionsByDateRange :many
SELECT * FROM transactions
WHERE user_id = ? AND deleted_at IS NULL
  AND t_date BETWEEN ? AND ?
ORDER BY t_date DESC, created_at DESC;

-- name: UpdateTransaction :one
UPDATE transactions
SET amount_pence = ?, t_date = ?, note = ?
WHERE id = ? AND deleted_at IS NULL
RETURNING *;

-- name: SoftDeleteTransaction :exec
UPDATE transactions
SET deleted_at = CURRENT_TIMESTAMP
WHERE id = ? AND deleted_at IS NULL;

-- name: HardDeleteTransaction :exec
DELETE FROM transactions
WHERE id = ?;

-- name: GetTransactionsByRecurringID :many
SELECT * FROM transactions
WHERE source_recurring = ? AND deleted_at IS NULL
ORDER BY t_date DESC;

-- name: CreateTransactionTag :exec
INSERT INTO transaction_tags (transaction_id, tag_id)
VALUES (?, ?)
ON CONFLICT(transaction_id, tag_id) DO NOTHING;

-- name: GetTransactionTags :many
SELECT t.* FROM tags t
JOIN transaction_tags tt ON t.id = tt.tag_id
WHERE tt.transaction_id = ?
ORDER BY t.name;

-- name: DeleteTransactionTag :exec
DELETE FROM transaction_tags
WHERE transaction_id = ? AND tag_id = ?;

-- name: DeleteAllTransactionTags :exec
DELETE FROM transaction_tags
WHERE transaction_id = ?;

-- name: CreateRecurring :one
INSERT INTO recurring (user_id, amount_pence, description, frequency, interval_n, first_due_date, next_due_date, end_date, active)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: GetRecurringByID :one
SELECT * FROM recurring
WHERE id = ?;

-- name: ListRecurring :many
SELECT * FROM recurring
WHERE user_id = ?
ORDER BY next_due_date ASC;

-- name: ListActiveRecurring :many
SELECT * FROM recurring
WHERE user_id = ? AND active = 1
ORDER BY next_due_date ASC;

-- name: GetRecurringDueOnDate :many
SELECT * FROM recurring
WHERE active = 1 AND next_due_date <= ?
ORDER BY next_due_date ASC;

-- name: UpdateRecurring :one
UPDATE recurring
SET amount_pence = ?, description = ?, frequency = ?, interval_n = ?, 
    first_due_date = ?, next_due_date = ?, end_date = ?, active = ?
WHERE id = ?
RETURNING *;

-- name: UpdateRecurringNextDue :exec
UPDATE recurring
SET next_due_date = ?
WHERE id = ?;

-- name: ToggleRecurringActive :exec
UPDATE recurring
SET active = CASE WHEN active = 1 THEN 0 ELSE 1 END
WHERE id = ?;

-- name: DeleteRecurring :exec
DELETE FROM recurring
WHERE id = ?;

-- name: CreateRecurringTag :exec
INSERT INTO recurring_tags (recurring_id, tag_id)
VALUES (?, ?)
ON CONFLICT(recurring_id, tag_id) DO NOTHING;

-- name: GetRecurringTags :many
SELECT t.* FROM tags t
JOIN recurring_tags rt ON t.id = rt.tag_id
WHERE rt.recurring_id = ?
ORDER BY t.name;

-- name: DeleteRecurringTag :exec
DELETE FROM recurring_tags
WHERE recurring_id = ? AND tag_id = ?;

-- name: DeleteAllRecurringTags :exec
DELETE FROM recurring_tags
WHERE recurring_id = ?;

-- name: PurgeSoftDeletedTransactions :exec
DELETE FROM transactions
WHERE deleted_at IS NOT NULL AND deleted_at < ?;

-- name: GetMonthlyReport :many
SELECT 
    t.name as tag_name,
    SUM(CASE WHEN tx.amount_pence > 0 THEN tx.amount_pence ELSE 0 END) as total_in_pence,
    SUM(CASE WHEN tx.amount_pence < 0 THEN ABS(tx.amount_pence) ELSE 0 END) as total_out_pence,
    COUNT(*) as transaction_count
FROM transactions tx
LEFT JOIN transaction_tags tt ON tx.id = tt.transaction_id
LEFT JOIN tags t ON tt.tag_id = t.id
WHERE tx.user_id = ? 
  AND tx.deleted_at IS NULL
  AND strftime('%Y-%m', tx.t_date) = ?
GROUP BY t.id, t.name
ORDER BY total_out_pence DESC;

-- name: GetMonthlyTotals :one
SELECT 
    SUM(CASE WHEN amount_pence > 0 THEN amount_pence ELSE 0 END) as total_in_pence,
    SUM(CASE WHEN amount_pence < 0 THEN ABS(amount_pence) ELSE 0 END) as total_out_pence,
    COUNT(*) as transaction_count
FROM transactions
WHERE user_id = ? 
  AND deleted_at IS NULL
  AND strftime('%Y-%m', t_date) = ?;

-- name: GetTransactionsByTag :many
SELECT tx.* FROM transactions tx
JOIN transaction_tags tt ON tx.id = tt.transaction_id
WHERE tt.tag_id = ? AND tx.deleted_at IS NULL
ORDER BY tx.t_date DESC;

-- name: GetRecurringByTag :many
SELECT r.* FROM recurring r
JOIN recurring_tags rt ON r.id = rt.recurring_id
WHERE rt.tag_id = ?
ORDER BY r.next_due_date ASC;
