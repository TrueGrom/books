package book

import (
	"books/app/common"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type BookModel struct {
	ID          uint   `json:"-" gorm:"column:id;primary_key"`
	UrlId       uint64 `json:"-" gorm:"column:url_id"`
	Title       string `json:"title" gorm:"column:title;not null"`
	Publisher   string `json:"publisher" gorm:"column:publisher"`
	Annotation  string `json:"annotation" gorm:"column:annotation"`
	Isbn        string `json:"isbn" gorm:"column:isbn"`
	Cover       string `json:"cover" gorm:"column:cover"`
	Size        string `json:"size" gorm:"column:size"`
	Rating      string `json:"rating" gorm:"column:rating"`
	Image       string `json:"image" gorm:"column:image"`
	Year        rune   `json:"year" gorm:"column:year"`
	Pages       rune   `json:"pages" gorm:"column:pages"`
	Weight      rune   `json:"weight" gorm:"column:weight"`
	Circulation rune   `json:"circulation" gorm:"column:circulation"`
}

func (BookModel) TableName() string {
	return "books"
}

func FindBooksByTitle(title string, limit rune) ([]BookModel, error) {
	db := common.GetDB()
	rows, err := db.Raw("SELECT * FROM  full_text_search(?, ?);", title, limit).Rows()
	books := make([]BookModel, 0, limit)
	if err != nil {
		return books, err
	}

	book := BookModel{}
	for rows.Next() {
		db.ScanRows(rows, &book)
		books = append(books, book)
	}
	return books, nil
}

func IsExist(condition interface{}) bool {
	db := common.GetDB()
	var model BookModel
	err := db.Where(condition).First(&model).Error
	if err != nil {
		return false
	}
	return true
}
