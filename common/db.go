package common

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func Init() *gorm.DB {
	db, err := gorm.Open("mysql", "root@/transcode?parseTime=true")

	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&Video{})
	return db
}
