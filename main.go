package main

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/wing8169/golang-todo/dto"
	"github.com/wing8169/golang-todo/templates"
	"github.com/wing8169/golang-todo/templates/components"
)

func filterByID(todos []*dto.TodoCardDto, id string) (out []*dto.TodoCardDto) {
	for _, todo := range todos {
		if todo.ID == id {
			continue
		}
		out = append(out, todo)
	}
	return out
}

func main() {
	todos := []*dto.TodoCardDto{
		{
			ID:      uuid.New().String(),
			Text:    "First item",
			Checked: false,
		}, {
			ID:      uuid.New().String(),
			Text:    "Second item",
			Checked: false,
		},
	}

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		component := templates.Index(todos)
		return component.Render(context.Background(), c.Response().Writer)
	})
	e.POST("/todos", func(c echo.Context) error {
		text := c.FormValue("add-todo-input")
		if text == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid text")
		}
		// add to todo
		todos = append(todos,
			&dto.TodoCardDto{
				ID:      uuid.New().String(),
				Text:    text,
				Checked: false,
			},
		)
		component := components.TodoCardsWithBtn(todos)
		return component.Render(context.Background(), c.Response().Writer)
	})
	e.DELETE("/todos/:id", func(c echo.Context) error {
		id := c.Param("id")
		if id == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid id")
		}
		todos = filterByID(todos, id)
		component := components.TodoCards(todos)
		return component.Render(context.Background(), c.Response().Writer)
	})
	e.GET("/components", func(c echo.Context) error {
		t := c.QueryParam("type")
		switch t {
		case "add-todo":
			component := components.AddTodoInput()
			return component.Render(context.Background(), c.Response().Writer)
		case "add-todo-btn":
			component := components.AddTodoButton()
			return component.Render(context.Background(), c.Response().Writer)
		}
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid element")
	})
	e.Static("/css", "css")
	e.Static("/static", "static")
	e.Static("/fonts", "fonts")
	e.Logger.Fatal(e.Start(":3000"))
}
