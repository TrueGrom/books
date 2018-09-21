package main

import (
	"books-backend/app/book"
	"books-backend/app/common"
	"books-backend/app/user"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"books-backend/app/comment"
)

func Migrate(db *gorm.DB) {
	user.AutoMigrate()
}

func main() {

	db := common.Init()
	//Migrate(db)
	defer db.Close()
	db.AutoMigrate(&comment.CommentModel{})

	r := gin.Default()
	r.Use(common.CORSMiddleware())

	v1 := r.Group("/api")

	userGroup := v1.Group("users/")
	user.UsersRegister(userGroup)
	user.UsersModify(userGroup)

	bookGroup := v1.Group("books/")
	book.BooksRegister(bookGroup)

	r.Run()
}
