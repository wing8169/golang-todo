package services

import (
	"database/sql"
	"log"

	"github.com/google/uuid"
	"github.com/wing8169/golang-todo/dto"
)

type TodoService struct {
	DB *sql.DB
}

func (ts *TodoService) GetTodos() []*dto.TodoCardDto {
	todos := []*dto.TodoCardDto{}
	rows, err := ts.DB.Query("select * from todo order by checked")
	if err != nil {
		log.Fatal(err)
		return todos
	}
	defer rows.Close()

	for rows.Next() {
		var id string
		var text string
		var checked bool
		err = rows.Scan(&id, &text, &checked)
		if err != nil {
			log.Fatal(err)
		}
		todos = append(todos, &dto.TodoCardDto{
			ID:      id,
			Text:    text,
			Checked: checked,
		})
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
		return todos
	}
	return todos
}

func (ts *TodoService) GetTodo(id string) *dto.TodoCardDto {
	stmt, err := ts.DB.Prepare("select text, checked from todo where id = ?")
	if err != nil {
		log.Fatal(err)
		return nil
	}
	defer stmt.Close()
	var text string
	var checked bool
	err = stmt.QueryRow(id).Scan(&text, &checked)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	return &dto.TodoCardDto{
		ID:      id,
		Text:    text,
		Checked: checked,
	}
}

func (ts *TodoService) CreateTodo(text string) *dto.TodoCardDto {
	// add to todo
	todo := &dto.TodoCardDto{
		ID:      uuid.New().String(),
		Text:    text,
		Checked: false,
	}
	sqlStmt := `
	insert into todo(id, text, checked) values(?, ?, ?)
	`
	_, err := ts.DB.Exec(sqlStmt, todo.ID, todo.Text, todo.Checked)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return nil
	}
	return todo
}

func (ts *TodoService) UpdateTodo(id string, text string, checked bool) *dto.TodoCardDto {
	todo := &dto.TodoCardDto{
		ID:      id,
		Text:    text,
		Checked: checked,
	}
	sqlStmt := `
	update todo set text=?, checked=? where id = ?
	`
	_, err := ts.DB.Exec(sqlStmt, todo.Text, todo.Checked, todo.ID)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return nil
	}
	return todo
}

func (ts *TodoService) DeleteTodo(id string) error {
	sqlStmt := `
	delete from todo where id = ?
	`
	_, err := ts.DB.Exec(sqlStmt, id)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return err
	}
	return nil
}
