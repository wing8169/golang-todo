package main

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/wing8169/golang-todo/templates"
)

func main() {
	e := echo.New()
	// Main menu
	component := templates.Hello("John")
	e.GET("/", func(c echo.Context) error {
		return component.Render(context.Background(), c.Response().Writer)
	})
	e.Static("/static", "static")
	e.Logger.Fatal(e.Start(":3000"))
}
