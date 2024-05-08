package initializers

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const TestDBName string = "ClubTennisTest"

func GetTestDatabase() *gorm.DB {
	//TODO make tempFS for this
	//TODO hardcoded for now, move to file at some point
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
