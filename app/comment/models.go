package comment

import (
	"time"
	"books-backend/app/user"
	"books-backend/app/book"
)

type CommentModel struct {
	ID        uint      `gorm:"column:id"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
	Text      string    `gorm:"column:text"`
	UserId    uint
	User      user.UserModel `gorm:"foreignkey:ID;association_foreignkey:UserId"`
	BookId    uint
	Book      book.BookModel `gorm:"foreignkey:ID;association_foreignkey:BookId"`
}
