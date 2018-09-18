package user

import (
	"books-backend/app/book"
	"books-backend/app/common"
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
	Books        []book.BookModel `gorm:"many2many:user_books;foreignkey:ID;association_foreignkey:ID;"`
}

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

func (user *UserModel) AddBookToUser(bookId uint) error {
	db := common.GetDB()
	err := db.Model(user).Association("Books").Append(book.BookModel{ID: bookId}).Error
	//err := db.Exec(fmt.Sprintf("INSERT INTO books_user_models(user_model_id, book_model_id) VALUES (%d, %d);", user.ID, bookId)).Error
	return err
}

func (user *UserModel) DeleteBookToUser(bookId uint) error {
	db := common.GetDB()
	err := db.Model(user).Association("Books").Delete(book.BookModel{ID: bookId}).Error
	return err
}
