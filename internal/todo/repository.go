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
	Create(title string) (int, error)
	GetByID(id int) (Todo, error)
	GetAll() ([]Todo, error)
	UpdateTitle(id int, title string) error
	MarkAsDone(id int) error
	Delete(id int) error
}

func NewRepository(conn *pgx.Conn) TodoRepository {
	return &Repository{
		conn: conn,
	}
}

func (r *Repository) GetAll() ([]Todo, error) {
	var todos []Todo

	rows, err := r.conn.Query(
		context.Background(),
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

func (r *Repository) GetByID(id int) (Todo, error) {
	var todo Todo

	err := r.conn.QueryRow(
		context.Background(),
		"SELECT id, title, completed FROM todos WHERE id = $1",
		id,
	).Scan(
		&todo.ID,
		&todo.Title,
		&todo.Completed,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return Todo{}, fmt.Errorf("todo not found")
	}

	return todo, err
}

func (r *Repository) Create(title string) (int, error) {
	var id int

	err := r.conn.QueryRow(
		context.Background(),
		"INSERT INTO todos (title) VALUES ($1) RETURNING id",
		title,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *Repository) Delete(id int) error {
	result, err := r.conn.Exec(
		context.Background(),
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

func (r *Repository) MarkAsDone(id int) error {
	result, err := r.conn.Exec(
		context.Background(),
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

func (r *Repository) UpdateTitle(id int, title string) error {
	result, err := r.conn.Exec(
		context.Background(),
		"UPDATE todos SET title = $1 WHERE id = $2",
		title,
		id,
	)

	if err != nil {
		return nil
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("todo not found")
	}

	return err
}
