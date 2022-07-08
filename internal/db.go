package internal

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func initDb() *gorm.DB {
	source, err := gorm.Open(sqlite.Open("bin/sqlitedb.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	return source
}

var (
	Db = initDb()
)
