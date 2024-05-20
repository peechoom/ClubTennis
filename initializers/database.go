package initializers

import (
	"ClubTennis/models"
	"fmt"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const TestDBName string = "ClubTennisTest"
const DBName string = "ClubTennis"

func GetDatabase() *gorm.DB {
	user := os.Getenv("DATABASE_USER")
	pass := os.Getenv("DATABASE_PASS")
	host := os.Getenv("DATABASE_HOST")
	port := os.Getenv("DATABASE_PORT")
	dbname := os.Getenv("DATABASE_DBNAME")

	createDBDsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/", user, pass, host, port)
	database, err := gorm.Open(mysql.Open(createDBDsn), &gorm.Config{})

	if err != nil {
		panic(err.Error())
	}

	_ = database.Exec("CREATE DATABASE IF NOT EXISTS " + dbname + ";")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", user, pass, host, port, dbname)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(err.Error())
	}

	err = db.AutoMigrate(models.User{}, models.Match{})
	if err != nil {
		panic(err.Error())
	}
	return db
}

func GetTestDatabase() *gorm.DB {
	// hardcoded for now, move to file at some point
	user := "root"
	pass := "1521"
	host := "localhost"
	port := "3306"
	dbname := TestDBName

	createDBDsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/", user, pass, host, port)
	database, err := gorm.Open(mysql.Open(createDBDsn), &gorm.Config{})

	if err != nil {
		panic(err.Error())
	}

	_ = database.Exec("CREATE DATABASE IF NOT EXISTS " + dbname + ";")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", user, pass, host, port, dbname)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(err.Error())
	}

	return db
}
