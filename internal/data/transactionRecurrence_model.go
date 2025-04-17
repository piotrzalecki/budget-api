package data

import (
	"context"
	"time"
)

type TransactionRecurrence struct {
	Id          int	`json:"id"`
	Name        string	`json:"name"`
	Description string	`json:"description"`
	AddTime     string	`json:"add_time"`
	CreatedAt   time.Time	`json:"created_at"`
	UpdatedAt   time.Time	`json:"updated_at"`
}

func (t *TransactionRecurrence) CreateTransactionsRecurrences(trec TransactionRecurrence) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := "INSERT INTO transactions_recurrences (name, description, add_time, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)  RETURNING id"
	var id int
	err := db.QueryRowContext(ctx, stmt,
		trec.Name,
		trec.Description,
		trec.AddTime,
		time.Now(),
		time.Now(),
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (t *TransactionRecurrence) AllTransactionsRecurrences() ([]TransactionRecurrence, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "SELECT id, name, description, add_time, created_at, updated_at FROM transactions_recurrences ORDER BY name ASC"

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var trs []TransactionRecurrence

	for rows.Next() {
		var tr TransactionRecurrence
		err := rows.Scan(
			&tr.Id,
			&tr.Name,
			&tr.Description,
			&tr.AddTime,
			&tr.UpdatedAt,
			&tr.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		trs = append(trs, tr)
	}

	if err = rows.Err(); err != nil {
		return trs, err
	}

	return trs, nil
}

func (t *TransactionRecurrence) GetTransactionsRecurrencesById(id int) (TransactionRecurrence, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "SELECT id, name, description, add_time, created_at, updated_at FROM transactions_recurrences WHERE id=$1"

	var tr TransactionRecurrence

	row := db.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&tr.Id,
		&tr.Name,
		&tr.Description,
		&tr.AddTime,
		&tr.CreatedAt,
		&tr.UpdatedAt,
	)

	if err != nil {
		return tr, err
	}

	return tr, nil

}
