package mysqldb

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	Db *gorm.DB
)

func init() {
	dsn := "root:123456@tcp(localhost:3306)/game_user?charset=utf8mb4&parseTime=True&loc=Local"
	var database, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	Db = database
}
