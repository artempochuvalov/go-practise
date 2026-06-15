package todo

import (
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
