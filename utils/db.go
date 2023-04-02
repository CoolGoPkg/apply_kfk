package utils

import (
	"LittleBeeMark/CoolGoPkg/apply_kfk/conf"
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// CreateDB 初始化MYSQL实例
func CreateDB(config conf.ConfigMysql) *gorm.DB {
	var (
		username = config.Username
		password = config.Password
		host     = config.Host
		port     = config.Port
		dbName   = config.DBName
		maxIdle  = config.MaxIdle
		maxOpen  = config.MaxConn
	)
	//connStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
	//	username,
	//	password,
	//	host,
	//	port,
	//	dbName,
	//)
	if config.Charset == "" {
		config.Charset = "utf8"
	}

	connStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&allowNativePasswords=true",
		username,
		password,
		host,
		port,
		dbName,
		config.Charset,
	)

	fmt.Println("Try to connect to MYSQL host: ", host, ", port: ", port, connStr)
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer（日志输出的目标，前缀和日志包含的内容——译者注）
		logger.Config{
			SlowThreshold:             10 * time.Second, // 慢 SQL 阈值
			LogLevel:                  logger.Info,      // 日志级别
			IgnoreRecordNotFoundError: false,            // 忽略ErrRecordNotFound（记录未找到）错误
			Colorful:                  false,            // 禁用彩色打印
		},
	)

	db, err := gorm.Open(mysql.Open(connStr), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		panic(fmt.Sprintf("failed to connect MYSQL %s:%d/%s: %s", host, port, dbName, err.Error()))
	}
	fmt.Println("Connected to MYSQL: ", host, ", port: ", port)

	sqlDB, err := db.DB()
	if err != nil {
		panic(fmt.Sprintf("db.DB() failed when connect to MYSQL %s:%d/%s: %s", host, port, dbName, err.Error()))
	}
	sqlDB.SetMaxIdleConns(maxIdle)
	sqlDB.SetMaxOpenConns(maxOpen)
	db.AutoMigrate()

	return db
}
