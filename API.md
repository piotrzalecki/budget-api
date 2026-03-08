# Budget API

> Auto-generated from swagger spec. Run `make api-docs` to regenerate.

## Base URL

`http://localhost:8080/api/v1`

## Authentication

| Scope | Method | Header |
|-------|--------|--------|
| `/api/v1/*` | Session token | `Authorization: Bearer <token>` |
| `/admin/*` | Static API key | `X-API-Key: <your-api-key>` |
| `POST /api/v1/auth/login` | None (public) | — |

Obtain a token via `POST /api/v1/auth/login`. Tokens expire after 30 days. Service accounts use a permanent token seeded from `SERVICE_USER_TOKEN` env var.

## Domain Conventions

- **Amounts**: string of integer pence. `"1050"` = £10.50. Negative = expense, positive = income.
- **Dates**: ISO 8601, `YYYY-MM-DD`
- **IDs**: integer

## Endpoints

### Transactions

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| `GET` | `/transactions` | Bearer | Get transactions |
| `POST` | `/transactions` | Bearer | Create a new transaction |
| `GET` | `/transactions/by-recurring/{recurring_id}` | Bearer | Get transactions by recurring ID |
| `GET` | `/transactions/by-tag/{tag_id}` | Bearer | Get transactions by tag |
| `POST` | `/transactions/purge` | Bearer | Purge soft deleted transactions |
| `GET` | `/transactions/{id}` | Bearer | Get transaction by ID |
| `PATCH` | `/transactions/{id}` | Bearer | Update a transaction |

**`GET /transactions`** query parameters:

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `from` | string | no | Start date (YYYY-MM-DD format) |
| `to` | string | no | End date (YYYY-MM-DD format) |

### Tags

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| `GET` | `/tags` | Bearer | Get all tags |
| `POST` | `/tags` | Bearer | Create a new tag |
| `PATCH` | `/tags/{id}` | Bearer | Update a tag |
| `DELETE` | `/tags/{id}` | Bearer | Delete a tag |

### Recurring

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| `GET` | `/recurring` | Bearer | Get all recurring transactions |
| `POST` | `/recurring` | Bearer | Create a new recurring transaction |
| `GET` | `/recurring/active` | Bearer | Get active recurring transactions |
| `GET` | `/recurring/by-tag/{tag_id}` | Bearer | Get recurring transactions by tag |
| `GET` | `/recurring/due` | Bearer | Get recurring transactions due on a date |
| `GET` | `/recurring/{id}` | Bearer | Get recurring transaction by ID |
| `PATCH` | `/recurring/{id}` | Bearer | Update a recurring transaction |
| `DELETE` | `/recurring/{id}` | Bearer | Delete a recurring transaction |
| `PATCH` | `/recurring/{id}/toggle` | Bearer | Toggle recurring transaction active status |

**`GET /recurring/due`** query parameters:

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `date` | string | no | Date to check (YYYY-MM-DD format, defaults to today) |

### Reports

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| `GET` | `/reports/monthly` | Bearer | Get monthly report |
| `GET` | `/reports/monthly/totals` | Bearer | Get monthly totals |

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

### Auth

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| `POST` | `/auth/login` | None | Login |
| `POST` | `/auth/logout` | Bearer | Logout |

### Health

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| `GET` | `/health` |  | Health check |

### Users

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| `GET` | `/users` | Bearer | List users |
| `POST` | `/users` | Bearer | Create user |
| `GET` | `/users/{id}` | Bearer | Get user |
| `PATCH` | `/users/{id}` | Bearer | Update user |
| `DELETE` | `/users/{id}` | Bearer | Delete user |

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

### CreateUserRequest

| Field | Type | Required | Notes |
|-------|------|----------|-------|
| `email` | string | yes |  |
| `is_service` | boolean | no |  |
| `password` | string | yes | min len 8 |

### ErrorResponse

| Field | Type | Required | Notes |
|-------|------|----------|-------|
| `data` | object | no |  |
| `error` | string | no |  |

### LoginRequest

| Field | Type | Required | Notes |
|-------|------|----------|-------|
| `email` | string | yes |  |
| `password` | string | yes | min len 1 |

### LoginResponse

| Field | Type | Required | Notes |
|-------|------|----------|-------|
| `email` | string | no |  |
| `expires_at` | string | no |  |
| `token` | string | no |  |
| `user_id` | integer | no |  |

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

### UpdateTagRequest

| Field | Type | Required | Notes |
|-------|------|----------|-------|
| `name` | string | yes | len 1–100 |

### UpdateTransactionRequest

| Field | Type | Required | Notes |
|-------|------|----------|-------|
| `deleted` | boolean | no |  |
| `note` | string | no |  |
| `tag_ids` | array[integer] | no |  |

### UpdateUserRequest

| Field | Type | Required | Notes |
|-------|------|----------|-------|
| `email` | string | no |  |
| `password` | string | no | min len 8 |

### UserResponse

| Field | Type | Required | Notes |
|-------|------|----------|-------|
| `created_at` | string | no |  |
| `email` | string | no |  |
| `id` | integer | no |  |
| `is_service` | boolean | no |  |

## Example

### Login and make a request

**1. Login**

```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "yourpassword"
}
```

**Response** `200 OK`

```json
{
  "data": {
    "token": "a3f8c2...",
    "expires_at": "2026-04-06T10:00:00Z",
    "user_id": 1,
    "email": "user@example.com"
  },
  "error": null
}
```

**2. Use the token**

```http
GET /api/v1/transactions
Authorization: Bearer a3f8c2...
```

---

### Create a transaction

**Request**

```http
POST /api/v1/transactions
Authorization: Bearer a3f8c2...
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
