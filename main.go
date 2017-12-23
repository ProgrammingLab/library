package main

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"io/ioutil"
	"encoding/json"
)

type setting struct {
	Dialect  string `json:"dialect"`
	User     string `json:"user"`
	Password string `json:"password"`
	DBName   string `json:"dbName"`
}

var db *gorm.DB

func main() {
	err := connectDB()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	e := echo.New()

	e.POST("/books/create", CreateBook)
	e.Start(":1323")
}

func connectDB() error {
	setting, err := loadSettings()
	if err != nil {
		return err
	}

	var dbErr error
	args := fmt.Sprintf("%v:%v@/%v?charset=utf8&parseTime=True&loc=Local", setting.User, setting.Password, setting.DBName)
	db, dbErr = gorm.Open("mysql", args)

	db.CreateTable(Book{})

	return dbErr
}

func loadSettings() (*setting, error) {
	bytes, err := ioutil.ReadFile("settings.json")
	if err != nil {
		return nil, err
	}

	var setting setting
	if err := json.Unmarshal(bytes, &setting); err != nil {
		return nil, err
	}

	return &setting, nil
}
