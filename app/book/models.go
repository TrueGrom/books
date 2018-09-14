package book

import (
	"books-backend/app/common"
	"fmt"
)

type BookModel struct {
	ID           uint     `json:"-" gorm:"column:id;primary_key"`
	UrlId        uint64   `json:"-" gorm:"column:url_id"`
	Title        string   `json:"title" gorm:"column:title;not null"`
	Publisher    string   `json:"publisher" gorm:"column:publisher"`
	Annotation   string   `json:"annotation" gorm:"column:annotation"`
	Isbn         string   `json:"isbn" gorm:"column:isbn"`
	Authors      []string `json:"authors" gorm:"column:authors"`
	Illustrators []string `json:"illustrators" gorm:"column:illustrators"`
	Translators  []string `json:"translators" gorm:"column:translators"`
	Editors      []string `json:"editors" gorm:"column:editors"`
	Year         rune     `json:"year" gorm:"column:year"`
	Pages        rune     `json:"pages" gorm:"column:pages"`
	ImageFile    string   `json:"image_file" gorm:"column:image_file"`
}

func (BookModel) TableName() string {
	return "books"
}

func FindBooksByTitle(title string, limit rune) ([]BookModel, error) {
	db := common.GetDB()
	var books []BookModel
	str := fmt.Sprintf("SELECT * FROM  full_text_search('%s', %d);", title, limit)
	rows, err := db.Raw(str).Rows()
	if err != nil {
		return books, err
	}

	book := BookModel{}
	for rows.Next() {
		db.ScanRows(rows, &book)
		books = append(books, book)
	}
	return books, err
}
