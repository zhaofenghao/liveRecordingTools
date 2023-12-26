package mysql

import (
	"ffmpeg_work/config"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var DB, DBLong, DBArticle, DBUseLog, DBLog, DBProxy *gorm.DB

func Set_up() {
	var err error
	DB, err = gorm.Open(config.Configs.DBDriver, fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8",
		config.Configs.DBUser,
		config.Configs.DBPass,
		config.Configs.DBHost,
		config.Configs.DBPort,
		config.Configs.DBName))

	if err != nil {
		fmt.Printf("mysql connect error %v", err)
	}
	if DB.Error != nil {
		fmt.Printf("database error %v", DB.Error)
	}
	if config.Configs.DBDebug {
		DB = DB.Debug()
	}
	//打印日志
	DB.LogMode(true)
	DB.DB().SetMaxIdleConns(10)
	DB.DB().SetMaxOpenConns(100)
}
