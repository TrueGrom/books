package comment

import (
	"books/app/book"
	"books/app/common"
	"books/app/user"
	"fmt"
	"time"
)

type CommentModel struct {
	ID        uint           `json:"id" gorm:"column:id"`
	CreatedAt time.Time      `json:"created_at" gorm:"column:created_at"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"column:updated_at"`
	Text      string         `json:"text" gorm:"column:text"`
	UserId    uint           `json:"user_id" gorm:"column:user_id"`
	User      user.UserModel `json:"-" gorm:"foreignkey:ID;association_foreignkey:UserId"`
	BookId    uint           `json:"book_id" gorm:"column:book_id"`
	Book      book.BookModel `json:"-" gorm:"foreignkey:ID;association_foreignkey:BookId"`
}

func FindOneComment(condition interface{}) (CommentModel, error) {
	db := common.GetDB()
	var model CommentModel
	err := db.Where(condition).First(&model).Error
	return model, err
}

func FindManyComments(condition interface{}) ([]CommentModel, error) {
	db := common.GetDB()
	comments := []CommentModel{}
	err := db.Where(condition).Find(&comments).Error
	return comments, err
}

func SaveOne(data interface{}) error {
	db := common.GetDB()
	err := db.Save(data).Error
	return err
}

func DeleteOneComment(comment CommentModel, date time.Time) (int64, error) {
	db := common.GetDB()
	str := fmt.Sprintf("created_at > '%s'", date.Format("2006-01-02 15:04:05"))
	q := db.Delete(&comment, str)
	return q.RowsAffected, q.Error
}
