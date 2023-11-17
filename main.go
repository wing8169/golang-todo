package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"

	"github.com/labstack/echo/v4"
	"github.com/wing8169/golang-todo/services"
	"github.com/wing8169/golang-todo/templates"
	"github.com/wing8169/golang-todo/templates/components"
)

func main() {
	db, err := sql.Open("sqlite3", "./db/todo.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqlStmt := `
	create table if not exists todo (id text not null primary key, text text, checked bool);
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}

	todoService := &services.TodoService{
		DB: db,
	}

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		todos := todoService.GetTodos()
		component := templates.Index(todos)
		return component.Render(context.Background(), c.Response().Writer)
	})
	e.POST("/todos", func(c echo.Context) error {
		text := c.FormValue("add-todo-input")
		if text == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid text")
		}
		todoService.CreateTodo(text)
		todos := todoService.GetTodos()
		component := components.TodoCardsWithBtn(todos)
		return component.Render(context.Background(), c.Response().Writer)
	})
	e.PUT("/todos/:id", func(c echo.Context) error {
		id := c.Param("id")
		if id == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid id")
		}

		oldTodo := todoService.GetTodo(id)

		text := c.FormValue("edit-todo-input")
		if text == "" {
			text = oldTodo.Text
		}

		checkedString := c.FormValue("checked")
		var checked bool
		if checkedString == "on" {
			checked = true
		} else {
			checked = false
		}

		todo := todoService.UpdateTodo(id, text, checked)

		component := components.TodoCard(*todo)
		return component.Render(context.Background(), c.Response().Writer)
	})
	e.DELETE("/todos/:id", func(c echo.Context) error {
		id := c.Param("id")
		if id == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid id")
		}
		todoService.DeleteTodo(id)
		todos := todoService.GetTodos()
		component := components.TodoCards(todos)
		return component.Render(context.Background(), c.Response().Writer)
	})
	e.GET("/components", func(c echo.Context) error {
		t := c.QueryParam("type")
		id := c.QueryParam("id")
		switch t {
		case "add-todo":
			component := components.AddTodoInput()
			return component.Render(context.Background(), c.Response().Writer)
		case "add-todo-btn":
			component := components.AddTodoButton()
			return component.Render(context.Background(), c.Response().Writer)
		case "edit-todo-input":
			todo := todoService.GetTodo(id)
			component := components.EditTodoInput(todo)
			return component.Render(context.Background(), c.Response().Writer)
		case "edit-todo-btn":
			todo := todoService.GetTodo(id)
			component := components.TodoCard(*todo)
			return component.Render(context.Background(), c.Response().Writer)
		}
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid element")
	})
	e.Static("/css", "css")
	e.Static("/static", "static")
	e.Static("/fonts", "fonts")
	e.Logger.Fatal(e.Start(":3000"))
}
