package main

import (
	"strings"
	"strconv"
	"net/http"
	"io/ioutil"
	"encoding/json"

	"github.com/labstack/echo"
)

type googleBook struct {
	VolumeInfo volumeInfo `json:"volumeInfo"`
}

type volumeInfo struct {
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
	if exists(isbn) {
		return nil
	}

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

	if len(response.Items) == 0 || len(response.Items[0].VolumeInfo.Title) == 0 {
		response := errorResponse{
			Code:    http.StatusNotFound,
			Message: "book not found",
		}
		responseError(c, response)
		return nil
	}

	db.Create(toBook(&response.Items[0].VolumeInfo))

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

func toBook(volume *volumeInfo) *Book {
	var isbn, isbn10 string
	for _, identifier := range volume.IndustryIdentifiers {
		switch identifier.Type {
		case "ISBN_10":
			isbn10 = identifier.Identifier
		case "ISBN_13":
			isbn = identifier.Identifier
		}
	}

	authors := strings.Join(volume.Authors, "/")
	return &Book{
		Title:       volume.Title,
		ISBN:        isbn,
		State:       "部室",
		Author:      authors,
		ReleaseDate: volume.PublishedDate,
		ISBN10:      isbn10,
		Page:        strconv.Itoa(volume.PageCount),
	}
}

func exists(isbn13 string) bool {
	out := &Book{}
	db.Find(out, &Book{ISBN: isbn13})
	return out.ID != 0
}
