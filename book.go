package main

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"strings"
	"strconv"
)

type Book struct {
	ID          int    `gorm:"type:int; primary_key;"`
	Title       string `gorm:"type:varchar(128)"`
	ISBN        string `gorm:"column:isbn; type:varchar(128)"`
	State       string `gorm:"type:varchar(128)"`
	Author      string `gorm:"type:varchar(100)"`
	Publisher   string `gorm:"type:varchar(100)"`
	ReleaseDate string `gorm:"type:varchar(100)"`
	ISBN10      string `gorm:"column:isbn10; type:varchar(100)"`
	Price       string `gorm:"type:varchar(100)"`
	Page        string `gorm:"type:varchar(100)"`
	Belong      string `gorm:"type:varchar(4)"`
}

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

type booksError struct {
	Message string
}

func (Book) TableName() string {
	return "t_books"
}

func (e booksError) Error() string {
	return e.Message
}

func registerBook(isbn string) (*Book, error) {
	if exists := findBook(isbn); exists != nil {
		return exists, nil
	}

	content, err := getBook(isbn)
	if err != nil {
		message := "something wrong(Google Books API??): " + err.Error()
		return nil, booksError{Message: message}
	}

	response := &googleBookResponse{}
	json.Unmarshal(content, response)

	if len(response.Items) == 0 || len(response.Items[0].VolumeInfo.Title) == 0 {
		return nil, booksError{Message: "book not found"}
	}

	book := toBook(&response.Items[0].VolumeInfo)
	db.Create(book)
	return book, nil
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

func findBook(isbn13 string) *Book {
	out := &Book{}
	db.Find(out, &Book{ISBN: isbn13})
	if out.ID == 0 {
		return nil
	}
	return out
}
