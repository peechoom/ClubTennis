package initializers

import (
	"ClubTennis/models"
	"fmt"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const TestDBName string = "ClubTennisTest"
const DBName string = "ClubTennis"

func GetDatabase() *gorm.DB {
	cloud := os.Getenv("CLOUD_INSTANCE")
	if len(cloud) != 0 && cloud != "false" {
		return connectUnixSocket()
	}

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

	err = db.AutoMigrate(&models.User{}, &models.Match{}, &models.Image{}, &models.Announcement{}, &models.Snippet{})
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

// conectUnixSocket initializes a Unix socket connection pool for
// a Cloud SQL instance of MySQL.
func connectUnixSocket() *gorm.DB {
	mustGetenv := func(k string) string {
		v := os.Getenv(k)
		if v == "" {
			log.Fatalf("Fatal Error in connect_unix.go: %s environment variable not set.", k)
		}
		return v
	}
	// Note: Saving credentials in environment variables is convenient, but not
	// secure - consider a more secure solution such as
	// Cloud Secret Manager (https://cloud.google.com/secret-manager) to help
	// keep secrets safe.
	var (
		dbUser         = mustGetenv("DATABASE_USER")                 // e.g. 'my-db-user'
		dbPwd          = mustGetenv("DATABASE_PASS")                 // e.g. 'my-db-password'
		dbName         = mustGetenv("DATABASE_DBNAME")               // e.g. 'my-database'
		unixSocketPath = mustGetenv("DATABASE_INSTANCE_UNIX_SOCKET") // e.g. '/cloudsql/project:region:instance'
	)

	createDBDsn := fmt.Sprintf("%s:%s@unix(%s)/",
		dbUser, dbPwd, unixSocketPath)

	// dbPool is the pool of database connections.
	database, err := gorm.Open(mysql.Open(createDBDsn), &gorm.Config{})
	if err != nil {
		panic(err.Error())
	}

	_ = database.Exec("CREATE DATABASE IF NOT EXISTS " + dbName + ";")
	dsn := fmt.Sprintf("%s:%s@unix(%s)/%s?charset=utf8&parseTime=true", dbUser, dbPwd, unixSocketPath, dbName)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(err.Error())
	}

	err = db.AutoMigrate(&models.User{}, &models.Match{}, &models.Image{}, &models.Announcement{}, &models.Snippet{})
	if err != nil {
		panic(err.Error())
	}

	return db
}
