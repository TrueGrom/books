package common

import (
	"books-backend/app"
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
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
	test_db_name := "test_books"
	connectionString := fmt.Sprintf(
		"port=%s user=%s password=%s dbname=%s sslmode=disable",
		app.DB_PORT,
		app.DB_USER,
		app.DB_PASSWORD,
		"postgres")
	fmt.Println(connectionString)
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		panic(err)
	}
	_, err = db.Exec("CREATE DATABASE test_books")
	db.Close()
	connectionString = fmt.Sprintf(
		"port=%s user=%s password=%s dbname=%s sslmode=disable",
		app.DB_PORT,
		app.DB_USER,
		app.DB_PASSWORD,
		test_db_name)
	db, err = sql.Open("postgres", connectionString)
	defer db.Close()
	if err != nil {
		panic(err)
	}
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		panic(err)
	}
	//dir, err := os.Getwd()
	//if err != nil {
	//	panic(err)
	//}
	m, err := migrate.NewWithDatabaseInstance(
		`file://../../migrations`,
		"postgres", driver)
	if err != nil {
		panic(err)
	}
	m.Up()
	test_db, err := gorm.Open("postgres", connectionString)
	DB = test_db
	return DB
}

func TestDBFree() error {
	//test_db.Close()
	connectionString := fmt.Sprintf(
		"port=%s user=%s password=%s dbname=%s sslmode=disable",
		app.DB_PORT,
		app.DB_USER,
		app.DB_PASSWORD,
		"postgres")
	db, err := sql.Open("postgres", connectionString)
	fmt.Println(err)
	db.Exec("SELECT pg_terminate_backend(pg_stat_activity.pid) " +
		"FROM pg_stat_activity " +
		"WHERE pg_stat_activity.datname = 'test_books' AND pid <> pg_backend_pid();")
	_, err = db.Exec("DROP DATABASE test_books")
	return err
}
