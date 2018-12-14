package user

import (
	"books/app/book"
	"books/app/common"
	"database/sql/driver"
	"errors"
	"golang.org/x/crypto/bcrypt"
)

const cost = 7

type UserModel struct {
	ID           uint             `gorm:"primary_key"`
	Username     string           `gorm:"column:username;unique_index"`
	Email        string           `gorm:"column:email;unique_index"`
	Bio          string           `gorm:"column:bio;size:1024"`
	Image        *string          `gorm:"column:image"`
	PasswordHash string           `gorm:"column:password;not null"`
	Books        []book.BookModel `gorm:"many2many:books_users_models;foreignkey:ID;association_foreignkey:ID;"`
}

type BookStatus string

const (
	NotRead   BookStatus = "Not Read"
	IsReading BookStatus = "Reading"
)

type BooksToUsers struct {
	User_id uint `gorm:"column:user_model_id"`
	Book_id uint `gorm:"column:book_model_id"`
	Rating  int8 `gorm:"column:rating"`
	//Status BookStatus `sql:"type:ENUM('Not Read', 'Reading')"`
}

func (BooksToUsers) TableName() string {
	return "books_users_models"
}

func (u *BookStatus) Scan(value interface{}) error { *u = BookStatus(value.([]byte)); return nil }
func (u BookStatus) Value() (driver.Value, error)  { return string(u), nil }

func AutoMigrate() {
	db := common.GetDB()
	db.AutoMigrate(&UserModel{})
}

func (u *UserModel) setPassword(password string) error {
	if len(password) == 0 {
		return errors.New("password should not be empty!")
	}
	bytePassword := []byte(password)
	passwordHash, _ := bcrypt.GenerateFromPassword(bytePassword, cost)
	u.PasswordHash = string(passwordHash)
	return nil
}

func (u *UserModel) checkPassword(password string) error {
	bytePassword := []byte(password)
	byteHashedPassword := []byte(u.PasswordHash)
	return bcrypt.CompareHashAndPassword(byteHashedPassword, bytePassword)
}

func FindOneUser(condition interface{}) (UserModel, error) {
	db := common.GetDB()
	var model UserModel
	err := db.Where(condition).First(&model).Error
	return model, err
}

func SaveOne(data interface{}) error {
	db := common.GetDB()
	err := db.Save(data).Error
	return err
}

func (model *UserModel) Update(data interface{}) error {
	db := common.GetDB()
	err := db.Model(model).Update(data).Error
	return err
}

func (user *UserModel) AddBooksToUser(booksID []uint) error {
	db := common.GetDB()
	query := "INSERT INTO books_users_models(user_model_id, book_model_id) VALUES (?,?) ON CONFLICT (user_model_id, book_model_id) DO NOTHING;"
	tx := db.Begin()
	for _, bookID := range booksID {
		_ = tx.Exec(query, user.ID, bookID).Error
	}
	err := tx.Commit().Error
	return err
}

func (user *UserModel) DeleteBooksToUser(booksID []uint) error {
	db := common.GetDB()
	books := make([]book.BookModel, len(booksID))
	for i, id := range booksID {
		books[i].ID = id
	}
	err := db.Model(user).Association("Books").Delete(books).Error
	return err
}

func (user *UserModel) GetAllBooksFromUser() ([]book.BookModel, error) {
	db := common.GetDB()
	books := []book.BookModel{}
	err := db.Model(user).Association("Books").Find(&books).Error
	return books, err
}

func AddRatingToBook(user *UserModel, booksRating []BookRating) error {
	db := common.GetDB()
	relation := BooksToUsers{}
	for _, bookRat := range booksRating {
		_ = db.Where(BooksToUsers{User_id: user.ID, Book_id: bookRat.Book_id}).Assign(BooksToUsers{Rating: bookRat.Rating}).FirstOrCreate(&relation).Error
	}
	return nil
}
