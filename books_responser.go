package main

import (
	"strings"
	"net/http"
	"io/ioutil"
	"encoding/json"

	"github.com/labstack/echo"
)

type googleBook struct {
	Title               string               `json:"title"`
	Authors             []string             `json:"authors"`
	PublishedDate       string               `json:"publishedDate"`
	IndustryIdentifiers []industryIdentifier `json:"industryIdentifiers"`
	PageCount           int                  `json:"pageCount"`
}

type googleBookResponse struct {
	Items []googleBook `json:"items"`
}

type industryIdentifier struct {
	Type       string `json:"type"`
	Identifier string `json:"identifier"`
}

func CreateBook(c echo.Context) error {
	isbn := c.FormValue("isbn")

	content, err := getBook(isbn)
	if err != nil {
		c.Logger().Infof(err.Error())
		response := errorResponse{
			Code:    http.StatusInternalServerError,
			Message: "something wrong(Google Books API??)",
		}
		responseError(c, response)
		return response
	}

	response := &googleBookResponse{}
	json.Unmarshal(content, response)
	if len(response.Items) == 0 {
		response := errorResponse{
			Code:    http.StatusNotFound,
			Message: "book not found",
		}
		responseError(c, response)
		return nil
	}

	db.NewRecord(toBook(&response.Items[0]))

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

func toBook(book *googleBook) Book {
	var isbn, isbn10 string
	for _, identifier := range book.IndustryIdentifiers {
		switch identifier.Type {
		case "ISBN_10":
			isbn10 = identifier.Identifier
		case "ISBN_13":
			isbn = identifier.Identifier
		}
	}

	authors := strings.Join(book.Authors, "/")
	return Book{
		Title:       book.Title,
		ISBN:        isbn,
		State:       "部室",
		Author:      authors,
		ReleaseDate: book.PublishedDate,
		ISBN10:      isbn10,
		Page:        string(book.PageCount),
	}
}
