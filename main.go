package main

import (
	"books-backend/common"
	"books-backend/user"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func Migrate(db *gorm.DB) {
	user.AutoMigrate()
}

func main() {

	db := common.Init()
	Migrate(db)
	defer db.Close()

	r := gin.Default()

	v1 := r.Group("/api")

	user.UsersRegister(v1.Group("user/"))

	r.Run()
}
