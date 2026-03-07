# Budget API

> Auto-generated from swagger spec. Run `make api-docs` to regenerate.

## Base URL

`http://localhost:8080/api/v1`

## Authentication

All `/api/v1/*` endpoints require: `X-API-Key: <your-api-key>` header.

## Domain Conventions

- **Amounts**: string of integer pence. `"1050"` = £10.50. Negative = expense, positive = income.
- **Dates**: ISO 8601, `YYYY-MM-DD`
- **IDs**: integer

## Endpoints

### Transactions

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| `GET` | `/transactions` | X-API-Key | Get transactions |
| `POST` | `/transactions` | X-API-Key | Create a new transaction |
| `GET` | `/transactions/by-recurring/{recurring_id}` | X-API-Key | Get transactions by recurring ID |
| `GET` | `/transactions/by-tag/{tag_id}` | X-API-Key | Get transactions by tag |
| `POST` | `/transactions/purge` | X-API-Key | Purge soft deleted transactions |
| `GET` | `/transactions/{id}` | X-API-Key | Get transaction by ID |
| `PATCH` | `/transactions/{id}` | X-API-Key | Update a transaction |

**`GET /transactions`** query parameters:

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `from` | string | no | Start date (YYYY-MM-DD format) |
| `to` | string | no | End date (YYYY-MM-DD format) |

### Tags

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| `GET` | `/tags` | X-API-Key | Get all tags |
| `POST` | `/tags` | X-API-Key | Create a new tag |

### Recurring

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| `GET` | `/recurring` | X-API-Key | Get all recurring transactions |
| `POST` | `/recurring` | X-API-Key | Create a new recurring transaction |
| `GET` | `/recurring/active` | X-API-Key | Get active recurring transactions |
| `GET` | `/recurring/by-tag/{tag_id}` | X-API-Key | Get recurring transactions by tag |
| `GET` | `/recurring/due` | X-API-Key | Get recurring transactions due on a date |
| `GET` | `/recurring/{id}` | X-API-Key | Get recurring transaction by ID |
| `PATCH` | `/recurring/{id}` | X-API-Key | Update a recurring transaction |
| `DELETE` | `/recurring/{id}` | X-API-Key | Delete a recurring transaction |
| `PATCH` | `/recurring/{id}/toggle` | X-API-Key | Toggle recurring transaction active status |

**`GET /recurring/due`** query parameters:

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `date` | string | no | Date to check (YYYY-MM-DD format, defaults to today) |

### Reports

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| `GET` | `/reports/monthly` | X-API-Key | Get monthly report |
| `GET` | `/reports/monthly/totals` | X-API-Key | Get monthly totals |

**`GET /reports/monthly`** query parameters:

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `ym` | string | no | Year-month in YYYY-MM format (defaults to current month) |

**`GET /reports/monthly/totals`** query parameters:

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `ym` | string | no | Year-month in YYYY-MM format (defaults to current month) |

### Admin

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| `POST` | `/admin/run-scheduler` | X-API-Key | Run the scheduler |

### Health

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| `GET` | `/health` |  | Health check |

## Request Schemas

### CreateRecurringRequest

| Field | Type | Required | Notes |
|-------|------|----------|-------|
| `amount` | string | yes |  |
| `description` | string | yes | len 1–255 |
| `end_date` | string | no |  |
| `first_due_date` | string | yes |  |
| `frequency` | string | yes | one of: daily, weekly, monthly, yearly |
| `interval_n` | integer | yes | range 1–365 |
| `tag_ids` | array[integer] | no |  |

### CreateTagRequest

| Field | Type | Required | Notes |
|-------|------|----------|-------|
| `name` | string | yes | len 1–100 |

### CreateTransactionRequest

| Field | Type | Required | Notes |
|-------|------|----------|-------|
| `amount` | string | yes |  |
| `note` | string | no |  |
| `t_date` | string | yes |  |
| `tag_ids` | array[integer] | no |  |

### PurgeTransactionsRequest

| Field | Type | Required | Notes |
|-------|------|----------|-------|
| `cutoff_date` | string | yes |  |

### UpdateRecurringRequest

| Field | Type | Required | Notes |
|-------|------|----------|-------|
| `active` | boolean | no |  |
| `amount` | string | no |  |
| `description` | string | no | len 1–255 |
| `end_date` | string | no |  |
| `first_due_date` | string | no |  |
| `frequency` | string | no | one of: daily, weekly, monthly, yearly |
| `interval_n` | integer | no | range 1–365 |
| `tag_ids` | array[integer] | no |  |

### UpdateTransactionRequest

| Field | Type | Required | Notes |
|-------|------|----------|-------|
| `deleted` | boolean | no |  |
| `note` | string | no |  |
| `tag_ids` | array[integer] | no |  |

## Example

### Create a transaction

**Request**

```http
POST /api/v1/transactions
X-API-Key: your-api-key
Content-Type: application/json

{
  "amount": "-1050",
  "t_date": "2026-03-07",
  "note": "Coffee",
  "tag_ids": [1, 3]
}
```

> `"amount": "-1050"` = £10.50 expense. Positive values are income.

**Response** `200 OK`

```json
{
  "id": 42,
  "amount": "-1050",
  "t_date": "2026-03-07",
  "note": "Coffee",
  "tag_ids": [1, 3],
  "created_at": "2026-03-07T10:15:30Z"
}
```
