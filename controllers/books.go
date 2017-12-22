package controllers

import (
	"net/http"
	"io/ioutil"

	"github.com/labstack/echo"
)



func CreateBook(c echo.Context) error {
	isbn := c.FormValue("isbn")

	content, err := getBook(isbn)
	if err != nil {
		c.Logger().Infof(err.Error())
		c.NoContent(http.StatusInternalServerError)
	}
	c.JSONBlob(http.StatusOK, content)

	return nil
}

func getBook(isbn string) ([]byte, error) {
	url := "https://www.googleapis.com/books/v1/volumes?q=isbn:" + isbn
	request, _ := http.NewRequest("GET", url, nil)
	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	content, _ := ioutil.ReadAll(response.Body)
	return content, nil
}