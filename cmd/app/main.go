package main

import (
	"context"
	"fmt"
	"os"

	"go-service/internal/database"
	"go-service/internal/todo"
)

func main() {
	conn, err := database.NewDB(os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}

	defer conn.Close(context.Background())

	r := todo.NewRepository(conn)

	todo, err := r.GetByID(2)
	if err != nil {
		panic(err)
	}

	fmt.Println(todo)

	todos, err := r.GetAll()

	if err != nil {
		panic(err)
	}

	for _, todo := range todos {
		fmt.Printf("%v | %v | %v\n", todo.ID, todo.Title, todo.Completed)
	}
}
