package db

import (
	"LittleBeeMark/CoolGoPkg/apply_kfk/conf"
	"LittleBeeMark/CoolGoPkg/apply_kfk/model"
	"LittleBeeMark/CoolGoPkg/apply_kfk/utils"
	"fmt"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	DB = utils.CreateDB(conf.Config.LocalMysql)
	err := DB.AutoMigrate(
		&model.Stock{},
		&model.Fiance{},
	)
	if err != nil {
		panic(fmt.Sprintf("DB.AutoMigrate err : %v", err))
	}
}
