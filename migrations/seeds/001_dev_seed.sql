-- +goose Up
-- +goose StatementBegin

-- Single dev user (password: devpassword)
INSERT INTO users (id, email, pw_hash) VALUES
    (1, 'dev@example.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy');

-- Settings
INSERT INTO settings (key, value) VALUES
    ('payday_dom', '28');

-- Tags
INSERT INTO tags (id, name) VALUES
    (1,  'groceries'),
    (2,  'rent'),
    (3,  'salary'),
    (4,  'transport'),
    (5,  'utilities'),
    (6,  'dining'),
    (7,  'entertainment'),
    (8,  'health'),
    (9,  'subscriptions'),
    (10, 'clothing');

-- Recurring rules
INSERT INTO recurring (id, user_id, amount_pence, description, frequency, interval_n, first_due_date, next_due_date, active) VALUES
    (1, 1, -85000, 'Monthly rent',        'monthly', 1, '2025-12-01', '2026-03-01', 1),
    (2, 1, 250000, 'Monthly salary',      'monthly', 1, '2025-12-28', '2026-03-28', 1),
    (3, 1,  -1799, 'Netflix',             'monthly', 1, '2025-12-15', '2026-03-15', 1),
    (4, 1,  -1199, 'Spotify',             'monthly', 1, '2025-12-10', '2026-03-10', 1);

-- Recurring tags
INSERT INTO recurring_tags (recurring_id, tag_id) VALUES
    (1, 2),  -- rent → rent
    (2, 3),  -- salary → salary
    (3, 9),  -- netflix → subscriptions
    (4, 9);  -- spotify → subscriptions

-- Transactions — recurring-generated (Dec 2025–Feb 2026)
INSERT INTO transactions (id, user_id, amount_pence, t_date, note, source_recurring) VALUES
    -- Rent
    (1,  1, -85000, '2025-12-01', 'Monthly rent',   1),
    (2,  1, -85000, '2026-01-01', 'Monthly rent',   1),
    (3,  1, -85000, '2026-02-01', 'Monthly rent',   1),
    -- Salary
    (4,  1, 250000, '2025-12-28', 'Monthly salary', 2),
    (5,  1, 250000, '2026-01-28', 'Monthly salary', 2),
    (6,  1, 250000, '2026-02-28', 'Monthly salary', 2),
    -- Netflix
    (7,  1,  -1799, '2025-12-15', 'Netflix',        3),
    (8,  1,  -1799, '2026-01-15', 'Netflix',        3),
    (9,  1,  -1799, '2026-02-15', 'Netflix',        3),
    -- Spotify
    (10, 1,  -1199, '2025-12-10', 'Spotify',        4),
    (11, 1,  -1199, '2026-01-10', 'Spotify',        4),
    (12, 1,  -1199, '2026-02-10', 'Spotify',        4);

-- Transactions — manual (Dec 2025)
INSERT INTO transactions (id, user_id, amount_pence, t_date, note, source_recurring) VALUES
    (13, 1,  -6540, '2025-12-02', 'Weekly groceries - Tesco',         NULL),
    (14, 1,  -4200, '2025-12-05', 'Bus pass top-up',                  NULL),
    (15, 1,  -3850, '2025-12-07', 'Dinner out with friends',          NULL),
    (16, 1,  -7230, '2025-12-09', 'Weekly groceries - Sainsbury''s',  NULL),
    (17, 1,  -1200, '2025-12-12', 'Pharmacy',                         NULL),
    (18, 1,  -6800, '2025-12-16', 'Weekly groceries - Tesco',         NULL),
    (19, 1,  -8500, '2025-12-18', 'Christmas gifts',                  NULL),
    (20, 1,  -4500, '2025-12-20', 'Restaurant - Christmas dinner',    NULL),
    (21, 1,  -7100, '2025-12-23', 'Weekly groceries - M&S',           NULL),
    (22, 1,  -3200, '2025-12-26', 'Boxing Day sales',                 NULL),
    (23, 1,  -6200, '2025-12-30', 'Weekly groceries - Tesco',         NULL);

-- Transactions — manual (Jan 2026)
INSERT INTO transactions (id, user_id, amount_pence, t_date, note, source_recurring) VALUES
    (24, 1,  -6900, '2026-01-03', 'Weekly groceries - Tesco',         NULL),
    (25, 1,  -2800, '2026-01-06', 'Bus pass top-up',                  NULL),
    (26, 1,  -4200, '2026-01-08', 'Dinner out',                       NULL),
    (27, 1,  -7500, '2026-01-11', 'Weekly groceries - Waitrose',      NULL),
    (28, 1,  -9800, '2026-01-13', 'New shoes',                        NULL),
    (29, 1,  -6100, '2026-01-17', 'Weekly groceries - Tesco',         NULL),
    (30, 1,  -3500, '2026-01-19', 'Gym membership',                   NULL),
    (31, 1,  -7200, '2026-01-24', 'Weekly groceries - Sainsbury''s',  NULL),
    (32, 1,  -2900, '2026-01-26', 'Coffee & lunch',                   NULL),
    (33, 1,  -6700, '2026-01-31', 'Weekly groceries - Tesco',         NULL);

-- Transactions — manual (Feb 2026)
INSERT INTO transactions (id, user_id, amount_pence, t_date, note, source_recurring) VALUES
    (34, 1,  -7100, '2026-02-03', 'Weekly groceries - Tesco',         NULL),
    (35, 1,  -2800, '2026-02-05', 'Bus pass top-up',                  NULL),
    (36, 1,  -5500, '2026-02-07', 'Valentine''s dinner',              NULL),
    (37, 1,  -6800, '2026-02-11', 'Weekly groceries - Sainsbury''s',  NULL),
    (38, 1,  -1500, '2026-02-13', 'Pharmacy',                         NULL),
    (39, 1,  -7300, '2026-02-18', 'Weekly groceries - Tesco',         NULL),
    (40, 1,  -3200, '2026-02-20', 'Cinema tickets',                   NULL),
    (41, 1,  -6900, '2026-02-25', 'Weekly groceries - M&S',           NULL),
    (42, 1,  -4100, '2026-02-27', 'Dinner out',                       NULL);

-- Transaction tags
INSERT INTO transaction_tags (transaction_id, tag_id) VALUES
    -- rent
    (1, 2), (2, 2), (3, 2),
    -- salary
    (4, 3), (5, 3), (6, 3),
    -- subscriptions
    (7, 9), (8, 9), (9, 9),
    (10, 9), (11, 9), (12, 9),
    -- groceries
    (13, 1), (16, 1), (18, 1), (21, 1), (23, 1),
    (24, 1), (27, 1), (29, 1), (31, 1), (33, 1),
    (34, 1), (37, 1), (39, 1), (41, 1),
    -- transport
    (14, 4), (25, 4), (35, 4),
    -- dining
    (15, 6), (20, 6), (26, 6), (32, 6), (36, 6), (42, 6),
    -- health
    (17, 8), (30, 8), (38, 8),
    -- clothing
    (19, 10), (28, 10),
    -- entertainment
    (22, 7), (40, 7);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DELETE FROM transaction_tags  WHERE transaction_id BETWEEN 1 AND 42;
DELETE FROM recurring_tags    WHERE recurring_id   BETWEEN 1 AND 4;
DELETE FROM transactions      WHERE id             BETWEEN 1 AND 42;
DELETE FROM recurring         WHERE id             BETWEEN 1 AND 4;
DELETE FROM tags              WHERE id             BETWEEN 1 AND 10;
DELETE FROM settings          WHERE key IN ('payday_dom');
DELETE FROM users             WHERE email = 'dev@example.com';

-- +goose StatementEnd
