package data

import (
	"context"
	"time"
)

type Budget struct {
	Id          int	`json:"id"`
	Name        string	`json:"name"`
	Description string	`json:"description"`
	CreatedAt   time.Time	`json:"created_at"`
	UpdatedAt   time.Time	`json:"updated_at"`
}

func (b *Budget) CreateBudget(bud Budget) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := "INSERT INTO budgets (name, description, created_at, updated_at) VALUES ($1, $2, $3, $4)  RETURNING id"
	var id int
	err := db.QueryRowContext(ctx, stmt,
		bud.Name,
		bud.Description,
		time.Now(),
		time.Now(),
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (b *Budget) AllBudgets() ([]Budget, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "SELECT id, name, description, created_at, updated_at FROM budgets ORDER BY name ASC"

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var budgets []Budget

	for rows.Next() {
		var b Budget
		err := rows.Scan(
			&b.Id,
			&b.Name,
			&b.Description,
			&b.UpdatedAt,
			&b.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		budgets = append(budgets, b)
	}

	if err = rows.Err(); err != nil {
		return budgets, err
	}

	return budgets, nil
}

func (b *Budget) GetBudgetById(id int) (Budget, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "SELECT id, name, description, created_at, updated_at FROM budgets WHERE id=$1"

	var bud Budget

	row := db.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&bud.Id,
		&bud.Name,
		&bud.Description,
		&bud.CreatedAt,
		&bud.UpdatedAt,
	)

	if err != nil {
		return bud, err
	}

	return bud, nil

}

func (b *Budget) DeleteBudget(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "DELETE FROM budgets WHERE id = $1"

	_, err := db.ExecContext(ctx, query, id)

	if err != nil {
		return err
	}

	return nil

}

func (b *Budget) UpdateBudget(bud Budget) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := "UPDATE budgets SET name=$1, description=$2, updated_at=$3 WHERE id=$4"

	_, err := db.ExecContext(ctx, stmt,
		bud.Name,
		bud.Description,
		time.Now(),
		bud.Id,
	)

	if err != nil {
		return err
	}

	return nil
}
