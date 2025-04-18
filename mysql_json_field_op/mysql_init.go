package mysql_json_field_op

import (
	"fmt"
	"time"

	logger "github.com/achilsh/go-dev-tools-collection/base-lib/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func createDbHandle() (*gorm.DB, error) {
	LoggerMock()
	//
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&timeout=10s&readTimeout=5s&writeTimeout=5s",
		"root", "123456789",
		"172.30.170.77", 3306, "vantop_db")

	// 连接数据库
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		// Logger:                 logger.Default.LogMode(m.GetGormLogMode()),
	})
	if err != nil {
		panic("failed to connect database")
	}
	// db.Exec("PRAGMA journal_mode=WAL;")
	conn, err := db.DB()
	if err != nil {
		logger.Errorf("Get db conn error: %v", err)
		return nil, err
	}

	conn.SetMaxOpenConns(100)
	conn.SetMaxIdleConns(10)
	conn.SetConnMaxIdleTime(time.Minute * 60)

	logger.Infof("Database connection successfully established")
	return db, nil
}

var gDb *gorm.DB = nil

func InitDB() error {
	db, err := createDbHandle()
	if err != nil {
		logger.Errorf("create db handle fail")
		return err
	}
	// 初始化 全局的 db handle
	gDb = db
	return nil
}

func GetDB() *gorm.DB {
	if gDb != nil {
		return gDb
	}
	logger.Errorf("not  init db handle.")
	panic("not init db handle.")
}

func CloseDB() {
	db := GetDB()
	if db == nil {
		return
	}
	sqlDb, err := db.DB()
	if err != nil {
		logger.Errorf("close fail, as get db fail, %v", err)
		return
	}
	sqlDb.Close()
}
