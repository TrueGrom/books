package main

import (
	"books/app/book"
	"books/app/comment"
	"books/app/common"
	"books/app/user"
	"github.com/gin-gonic/gin"
	_ "github.com/golang-migrate/migrate/source/file"

	_ "github.com/lib/pq"
)

func main() {
	db := common.Init()
	defer db.Close()

	r := gin.Default()
	r.Use(common.CORSMiddleware())

	v1 := r.Group("/api")

	userGroup := v1.Group("users/")
	user.UsersRegister(userGroup)

	bookGroup := v1.Group("books/")
	book.BooksRegister(bookGroup)

	commentGroup := v1.Group("comments/")
	comment.CommentsRegister(commentGroup)

	r.Run()
}
