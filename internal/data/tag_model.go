package data

import (
	"context"
	"time"
)

type Tag struct {
	Id          int
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (t *Tag) CreateTag(tag Tag) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := "INSERT INTO tags (name, description, created_at, updated_at) VALUES ($1, $2, $3, $4)  RETURNING id"
	var id int
	err := db.QueryRowContext(ctx, stmt,
		tag.Name,
		tag.Description,
		time.Now(),
		time.Now(),
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (t *Tag) AllTags() ([]Tag, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "SELECT id, name, description, created_at, updated_at FROM tags ORDER BY name ASC"

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []Tag

	for rows.Next() {
		var t Tag
		err := rows.Scan(
			&t.Id,
			&t.Name,
			&t.Description,
			&t.UpdatedAt,
			&t.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		tags = append(tags, t)
	}

	if err = rows.Err(); err != nil {
		return tags, err
	}

	return tags, nil
}

func (t *Tag) GetTagById(id int) (Tag, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "SELECT id, name, description, created_at, updated_at FROM tags WHERE id=$1"

	var tag Tag

	row := db.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&tag.Id,
		&tag.Name,
		&tag.Description,
		&tag.CreatedAt,
		&tag.UpdatedAt,
	)

	if err != nil {
		return tag, err
	}

	return tag, nil

}

func (t *Tag) DeleteTag(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "DELETE FROM tags WHERE id = $1"

	_, err := db.ExecContext(ctx, query, id)

	if err != nil {
		return err
	}

	return nil

}

func (t *Tag) UpdateTag(tag Tag) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := "UPDATE tags SET name=$1, description=$2, updated_at=$3 WHERE id=$4"

	_, err := db.ExecContext(ctx, stmt,
		tag.Name,
		tag.Description,
		time.Now(),
		tag.Id,
	)

	if err != nil {
		return err
	}

	return nil
}
