package todo

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type Repository struct {
	conn *pgx.Conn
}

type TodoRepository interface {
	Create(ctx context.Context, title string) (int, error)
	GetByID(ctx context.Context, id int) (Todo, error)
	GetAll(ctx context.Context) ([]Todo, error)
	UpdateTitle(ctx context.Context, id int, title string) error
	MarkAsDone(ctx context.Context, id int) error
	Delete(ctx context.Context, id int) error
}

func NewRepository(conn *pgx.Conn) TodoRepository {
	return &Repository{
		conn: conn,
	}
}

func (r *Repository) GetAll(ctx context.Context) ([]Todo, error) {
	var todos []Todo

	rows, err := r.conn.Query(
		ctx,
		"SELECT id, title, completed FROM todos",
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var todo Todo

		err := rows.Scan(&todo.ID, &todo.Title, &todo.Completed)

		if err != nil {
			return nil, err
		}

		todos = append(todos, todo)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return todos, nil
}

func (r *Repository) GetByID(ctx context.Context, id int) (Todo, error) {
	var todo Todo

	err := r.conn.QueryRow(
		ctx,
		"SELECT id, title, completed FROM todos WHERE id = $1",
		id,
	).Scan(
		&todo.ID,
		&todo.Title,
		&todo.Completed,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return Todo{}, TodoNotFound
	}

	return todo, err
}

func (r *Repository) Create(ctx context.Context, title string) (int, error) {
	var id int

	err := r.conn.QueryRow(
		ctx,
		"INSERT INTO todos (title) VALUES ($1) RETURNING id",
		title,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *Repository) Delete(ctx context.Context, id int) error {
	result, err := r.conn.Exec(
		ctx,
		"DELETE FROM todos WHERE id = $1",
		id,
	)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("todo not found")
	}

	return err
}

func (r *Repository) MarkAsDone(ctx context.Context, id int) error {
	result, err := r.conn.Exec(
		ctx,
		"UPDATE todos SET completed = true WHERE id = $1",
		id,
	)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("todo not found")
	}

	return err
}

func (r *Repository) UpdateTitle(ctx context.Context, id int, title string) error {
	result, err := r.conn.Exec(
		ctx,
		"UPDATE todos SET title = $1 WHERE id = $2",
		title,
		id,
	)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("todo not found")
	}

	return err
}
