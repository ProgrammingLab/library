package main

import (
	"fmt"
	"encoding/json"

	"github.com/labstack/echo"
)

type errorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e errorResponse) Error() string {
	return fmt.Sprintf("%v:%v", e.Code, e.Message)
}

func responseError(c echo.Context, response errorResponse) {
	bytes, _ := json.Marshal(response)
	c.JSONBlob(response.Code, bytes)
}
