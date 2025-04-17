package data

import (
	"context"
	"time"
)

type Transaction struct {
	Id                    int	`json:"id"`
	Name                  string	`json:"name"`
	Description           string	`json:"description"`
	Budget                Budget	`json:"budget"`
	Quote                 float32	`json:"quote"`
	TransactionRecurrence TransactionRecurrence `json:"transaction_recurrence"`
	Active                bool	`json:"active"`
	Starts                time.Time	`json:"starts"`
	Ends                  time.Time	`json:"ends"`
	Tag                   Tag	`json:"tag"`
	CreatedAt             time.Time	`json:"created_at"`
	UpdatedAt             time.Time	`json:"updated_at"`
}

func (t *Transaction) CreateTransaction(trn Transaction) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := "INSERT INTO transactions (name, description, budget, quote, transaction_recurrence, active, starts, ends, tag, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)  RETURNING id"
	var id int
	err := db.QueryRowContext(ctx, stmt,
		&trn.Name,
		&trn.Description,
		&trn.Budget.Id,
		&trn.Quote,
		&trn.TransactionRecurrence.Id,
		&trn.Active,
		&trn.Starts,
		&trn.Ends,
		&trn.Tag.Id,
		time.Now(),
		time.Now(),
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (t *Transaction) AllTransactions() ([]Transaction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `SELECT tr.id, tr.name, tr.description, tr.quote, tr.active, tr.starts, tr.ends, tr.created_at, tr.updated_at,
	tr.budget, bd.name, tr.transaction_recurrence, trs.name, tr.tag, tgs.name
	FROM transactions tr 
	LEFT JOIN budgets bd ON (tr.budget = bd.id)
	LEFT JOIN transactions_recurrences trs ON (tr.transaction_recurrence = trs.id)
	LEFT JOIN tags tgs ON (tr.tag = tgs.id)
	ORDER BY tr.starts ASC, tr.active DESC`

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var trns []Transaction

	for rows.Next() { // Zmien to bo to nie zadziala
		var tr Transaction
		err := rows.Scan(
			&tr.Id,
			&tr.Name,
			&tr.Description,
			&tr.Quote,
			&tr.Active,
			&tr.Starts,
			&tr.Ends,
			&tr.CreatedAt,
			&tr.UpdatedAt,
			&tr.Budget.Id,
			&tr.Budget.Name,
			&tr.TransactionRecurrence.Id,
			&tr.TransactionRecurrence.Name,
			&tr.Tag.Id,
			&tr.Tag.Name,
		)
		if err != nil {
			return nil, err
		}

		trns = append(trns, tr)
	}

	if err = rows.Err(); err != nil {
		return trns, err
	}

	return trns, nil
}

func (t *Transaction) GetTransactionById(id int) (Transaction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `SELECT tr.id, tr.name, tr.description, tr.quote, tr.active, tr.starts, tr.ends, tr.created_at, tr.updated_at,
	tr.budget, bd.name, bd.description,  tr.transaction_recurrence,  trs.name, trs.description, tr.tag, tgs.name, tgs.description
	FROM transactions tr
	LEFT JOIN budgets bd ON (tr.budget = bd.id)
	LEFT JOIN transactions_recurrences trs ON (tr.transaction_recurrence = trs.id)
	LEFT JOIN tags tgs ON (tr.tag = tgs.id)
	WHERE tr.id=$1`

	var tr Transaction

	row := db.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&tr.Id,
		&tr.Name,
		&tr.Description,
		&tr.Quote,
		&tr.Active,
		&tr.Starts,
		&tr.Ends,
		&tr.CreatedAt,
		&tr.UpdatedAt,
		&tr.Budget.Id,
		&tr.Budget.Name,
		&tr.Budget.Description,
		&tr.TransactionRecurrence.Id,
		&tr.TransactionRecurrence.Name,
		&tr.TransactionRecurrence.Description,
		&tr.Tag.Id,
		&tr.Tag.Name,
		&tr.Tag.Description,
	)

	if err != nil {
		return tr, err
	}

	return tr, nil

}

func (t *Transaction) AllTransactionsByBudget(id int) ([]Transaction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `SELECT tr.id, tr.name, tr.description, tr.quote, tr.active, tr.starts, tr.ends, tr.created_at, tr.updated_at,
	tr.budget, bd.name, tr.transaction_recurrence, trs.name, tr.tag, tgs.name
	FROM transactions tr 
	LEFT JOIN budgets bd ON (tr.budget = bd.id)
	LEFT JOIN transactions_recurrences trs ON (tr.transaction_recurrence = trs.id)
	LEFT JOIN tags tgs ON (tr.tag = tgs.id)
	WHERE bd.id=$1
	ORDER BY tr.starts ASC`

	rows, err := db.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var trns []Transaction

	for rows.Next() {
		var tr Transaction
		err := rows.Scan(
			&tr.Id,
			&tr.Name,
			&tr.Description,
			&tr.Quote,
			&tr.Active,
			&tr.Starts,
			&tr.Ends,
			&tr.CreatedAt,
			&tr.UpdatedAt,
			&tr.Budget.Id,
			&tr.Budget.Name,
			&tr.TransactionRecurrence.Id,
			&tr.TransactionRecurrence.Name,
			&tr.Tag.Id,
			&tr.Tag.Name,
		)
		if err != nil {
			return nil, err
		}

		trns = append(trns, tr)
	}

	if err = rows.Err(); err != nil {
		return trns, err
	}

	return trns, nil

}

func (t *Transaction) DeleteTransaction(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "DELETE FROM transactions WHERE id = $1"

	_, err := db.ExecContext(ctx, query, id)

	if err != nil {
		return err
	}

	return nil

}

func (t *Transaction) UpdateTransaction(trn Transaction) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := "UPDATE transactions SET name=$1, description=$2, budget=$3, quote=$4, transaction_recurrence=$5, active=$6, starts=$7, ends=$8, tag=$9, updated_at=$10 WHERE id=$11"

	_, err := db.ExecContext(ctx, stmt,
		trn.Name,
		trn.Description,
		trn.Budget.Id,
		trn.Quote,
		trn.TransactionRecurrence.Id,
		trn.Active,
		trn.Starts,
		trn.Ends,
		trn.Tag.Id,
		time.Now(),
		trn.Id,
	)

	if err != nil {
		return err
	}

	return nil
}

func (t *Transaction) TransactionSetStatus(id int, status bool) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := "UPDATE transactions SET active=$1, updated_at=$2 WHERE id=$3"

	_, err := db.ExecContext(ctx, stmt,
		status,
		time.Now(),
		id,
	)

	if err != nil {
		return err
	}

	return nil
}

func (t *Transaction) TransactionsSetAllActive() error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := "UPDATE transactions SET active=$1, updated_at=$2"

	_, err := db.ExecContext(ctx, stmt,
		true,
		time.Now(),
	)

	if err != nil {
		return err
	}

	return nil
}
