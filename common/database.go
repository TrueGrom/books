package common

import (
	"fmt"
	"github.com/jinzhu/gorm"
)

type Database struct {
	*gorm.DB
}

var DB *gorm.DB

func Init() *gorm.DB {
	db, err := gorm.Open("postgres", "host=localhost port=5432 user=books dbname=books password=books")
	if err != nil {
		fmt.Println("db err: ", err)
	}
	//db.DB().SetMaxIdleConns(10)
	DB = db
	return DB
}

func GetDB() *gorm.DB {
	return DB
}
