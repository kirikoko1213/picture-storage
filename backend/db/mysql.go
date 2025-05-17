package db

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

var MYSQL_USERNAME = os.Getenv("MYSQL_USERNAME")
var MYSQL_PASSWORD = os.Getenv("MYSQL_PASSWORD")
var MYSQL_HOST = os.Getenv("MYSQL_HOST")
var MYSQL_PORT = os.Getenv("MYSQL_PORT")
var MYSQL_DATABASE = os.Getenv("MYSQL_DATABASE")

func init() {
	if err := InitMySQL(); err != nil {
		panic(err)
	}
}

func InitMySQL() error {
	port, _ := strconv.Atoi(MYSQL_PORT)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		MYSQL_USERNAME,
		MYSQL_PASSWORD,
		MYSQL_HOST,
		port,
		MYSQL_DATABASE,
	)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to MySQL: %v", err)
	}

	log.Default().Println("[mysql] dsn: ", dsn)

	return nil
}
