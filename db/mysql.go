package db

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"io/github/gforgame/config"
)

var (
	Db *gorm.DB
)

func init() {
	var database, err = gorm.Open(mysql.Open(config.ServerConfig.DbUrl), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	Db = database
}
