package common

import (
	"books-backend/app"
	"fmt"
	"github.com/jinzhu/gorm"
)

type Database struct {
	*gorm.DB
}

var DB *gorm.DB

func Init() *gorm.DB {
	configDB := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		app.DB_HOST,
		app.DB_PORT,
		app.DB_USER,
		app.DB_NAME,
		app.DB_PASSWORD,
		app.DB_SSL_MODE)
	db, err := gorm.Open("postgres", configDB)
	if err != nil {
		fmt.Println("db err: ", err)
	}
	db.DB().SetMaxIdleConns(10)
	DB = db
	return DB
}

func GetDB() *gorm.DB {
	return DB
}
