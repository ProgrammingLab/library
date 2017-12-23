package main

import (
	"net/http"
	"github.com/labstack/echo"
)

func CreateBook(c echo.Context) error {
	isbn := c.FormValue("isbn")
	_, err := registerBook(isbn)
	if err != nil {
		response := errorResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
		responseError(c, response)
	}

	return nil
}
