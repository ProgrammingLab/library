package main

type Book struct {
	ID          int    `gorm:"primary_key"`
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

func (Book) TableName() string {
	return "t_books"
}
