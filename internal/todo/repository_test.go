package todo

import (
	"fmt"
	"go-service/internal/database"
	"os"
	"testing"
)

func setupRepository(t *testing.T) *Repository {
	conn, err := database.NewDB(os.Getenv("DATABASE_TEST_URL"))

	if err != nil {
		t.Fatalf("connecting to db failed %v", err)
	}

	r := NewRepository(conn)

	_, err = conn.Exec(t.Context(), "TRUNCATE todos RESTART IDENTITY")

	if err != nil {
		t.Fatalf("clean up failed %v", err)
	}

	t.Cleanup(func() {
		conn.Close(t.Context())
	})

	return r
}

func TestRepository_Create(t *testing.T) {
	r := setupRepository(t)

	id, err := r.Create("learning testing")
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}

	if id <= 0 {
		t.Fatalf("unexcpected id: %v", id)
	}
}

func TestRepository_GetById_Success(t *testing.T) {
	r := setupRepository(t)

	id, err := r.Create("learning testing")
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}

	todo, err := r.GetByID(id)
	if err != nil {
		t.Fatalf("fetch error: %v", err)
	}

	if todo.Title != "learning testing" {
		t.Fatalf("unexpected todo: %v", todo.Title)
	}
}

func TestRepository_GetById_NotFound(t *testing.T) {
	r := setupRepository(t)

	_, err := r.GetByID(999)
	if err == nil {
		t.Fatalf("should have returned error")
	}
}

func TestRepository_GetAll(t *testing.T) {
	r := setupRepository(t)

	for i := range 3 {
		_, err := r.Create(fmt.Sprint(i))
		if err != nil {
			t.Fatalf("create failed: %v", err)
		}
	}

	todos, err := r.GetAll()
	if err != nil {
		t.Fatalf("fetching todos failed: %v", err)
	}

	todosLen := len(todos)
	if todosLen != 3 {
		t.Fatalf("unappropriate amount of created todos: %v", todosLen)
	}
}

func TestRepository_UpdateTitle_Success(t *testing.T) {
	r := setupRepository(t)

	id, err := r.Create("initial title")
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}

	err = r.UpdateTitle(id, "new title")
	if err != nil {
		t.Fatalf("update title failed: %v", err)
	}

	todo, err := r.GetByID(id)
	if err != nil {
		t.Fatalf("fetch by id failed: %v", err)
	}

	if todo.Title != "new title" {
		t.Fatalf("expected title to be %s, got %s", "new title", todo.Title)
	}
}

func TestRepository_UpdateTitle_Fail(t *testing.T) {
	r := setupRepository(t)

	err := r.UpdateTitle(999, "new title")
	if err == nil {
		t.Fatalf("expected to fail")
	}
}

func TestRepository_MarkAsDone_Success(t *testing.T) {
	r := setupRepository(t)

	id, err := r.Create("learning testing")
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}

	err = r.MarkAsDone(id)
	if err != nil {
		t.Fatalf("mark as done failed: %v", err)
	}

	todo, err := r.GetByID(id)
	if err != nil {
		t.Fatalf("fetch by id failed: %v", err)
	}

	if !todo.Completed {
		t.Fatalf("expected Completed to be %t, got %t", true, todo.Completed)
	}
}

func TestRepository_MarkAsDone_Fail(t *testing.T) {
	r := setupRepository(t)

	err := r.MarkAsDone(999)
	if err == nil {
		t.Fatalf("expected to fail")
	}
}

func TestRepository_Delete_Success(t *testing.T) {
	r := setupRepository(t)

	id, err := r.Create("learning testing")
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}

	err = r.Delete(id)
	if err != nil {
		t.Fatalf("delete failed: %v", err)
	}

	_, err = r.GetByID(id)
	if err == nil {
		t.Fatalf("expected task to be deleted")
	}
}

func TestRepository_Delete_Fail(t *testing.T) {
	r := setupRepository(t)

	err := r.Delete(999)
	if err == nil {
		t.Fatalf("expected to fail")
	}
}
