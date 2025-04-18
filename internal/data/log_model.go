package data

import (
	"context"
	"time"
)

type Log struct {
	Id        int	`json:"id"`
	Log       string	`json:"log"`	
	CreatedAt time.Time	`json:"created_at"`
	UpdatedAt time.Time	`json:"updated_at"`
}

func (l *Log) AllLogs() ([]Log, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "SELECT id, log, created_at, updated_at FROM logs ORDER BY id DESC LIMIT 50"

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []Log

	for rows.Next() {
		var l Log
		err := rows.Scan(
			&l.Id,
			&l.Log,
			&l.CreatedAt,
			&l.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		logs = append(logs, l)
	}

	if err = rows.Err(); err != nil {
		return logs, err
	}

	return logs, nil
}

func (l *Log) AddLog(log string) (error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := "INSERT INTO logs (log, created_at, updated_at) VALUES ($1, $2, $3)  RETURNING id"


	_, err := db.ExecContext(ctx, stmt,
		log,
		time.Now(),
		time.Now(),
	)

	if err != nil {
		return err
	}

	return nil

}
