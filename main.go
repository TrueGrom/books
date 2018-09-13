package main

import (
	"books-backend/app/common"
	"books-backend/app/user"
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
	r.Use(common.CORSMiddleware())

	v1 := r.Group("/api")

	userGroup := v1.Group("users/")
	user.UsersRegister(userGroup)
	user.UsersModify(userGroup)

	r.Run()
}
