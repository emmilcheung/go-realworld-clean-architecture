package postgres

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type DB = gorm.DB

// in memory DB instance
var instance *DB

func DBInit(postgresUrl string) (*DB, error) {
	db, err := gorm.Open("postgres", postgresUrl)
	if err != nil {
		fmt.Println("db err: (Init) ", err)
	}
	db.DB().SetMaxIdleConns(10)
	//db.LogMode(true)
	instance = db
	return db, err
}

func GetDB() *DB {
	return instance
}
