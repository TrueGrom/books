package common

import (
	"books-backend/app"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"os"
)

type Database struct {
	*gorm.DB
}

var DB *gorm.DB

func getConfig() string {
	return fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		app.DB_HOST,
		app.DB_PORT,
		app.DB_USER,
		app.DB_NAME,
		app.DB_PASSWORD,
		app.DB_SSL_MODE)
}

func Init() *gorm.DB {
	db, err := gorm.Open("postgres", getConfig())
	if err != nil {
		fmt.Println("db err: ", err)
	}
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(10)
	db.LogMode(true)
	DB = db
	return DB
}

func GetDB() *gorm.DB {
	return DB
}

func TestDBInit() *gorm.DB {
	//out, err := exec.Command("/bin/sh", "../../db_init.sh testbooks").Output()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Printf("output is %s\n", out)
	//test_db, err := gorm.Open("postgres", "host=localhost port=5432 user=testbooks dbname=testbooks password=testbooks sslmode=disable")
	test_db, err := gorm.Open("sqlite3", "./../gorm_test.db")
	if err != nil {
		fmt.Println("db err: ", err)
	}
	if err != nil {
		fmt.Println("db err: ", err)
	}
	test_db.DB().SetMaxIdleConns(3)
	test_db.LogMode(true)
	DB = test_db
	return DB
}

func TestDBFree(test_db *gorm.DB) error {
	//test_db.Close()
	err := os.Remove("./../gorm_test.db")
	return err
}
