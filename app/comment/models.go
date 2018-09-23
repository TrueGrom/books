package comment

import (
	"books-backend/app/book"
	"books-backend/app/common"
	"books-backend/app/user"
	"fmt"
	"time"
)

type CommentModel struct {
	ID        uint           `gorm:"column:id"`
	CreatedAt time.Time      `gorm:"column:created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at"`
	Text      string         `gorm:"column:text"`
	UserId    uint           `gorm:"column:user_id"`
	User      user.UserModel `gorm:"foreignkey:ID;association_foreignkey:UserId"`
	BookId    uint           `gorm:"column:book_id"`
	Book      book.BookModel `gorm:"foreignkey:ID;association_foreignkey:BookId"`
}

func FindOneComment(condition interface{}) (CommentModel, error) {
	db := common.GetDB()
	var model CommentModel
	err := db.Where(condition).First(&model).Error
	return model, err
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
