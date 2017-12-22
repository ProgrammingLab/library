package main

import (
	"github.com/ProgrammingLab/library/controllers"

	"github.com/labstack/echo"
)

func main() {
	e := echo.New()

	e.POST("/books/create", controllers.CreateBook)

	e.Logger.Fatal(e.Start(":1323"))
}