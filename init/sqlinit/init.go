package sqlinit

import (
	"TM-Spike/model"
	"database/sql"
	"fmt"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"gorm.io/plugin/opentelemetry/logging/logrus"
	"gorm.io/plugin/opentelemetry/tracing"
)

// test
// 測試
func TestInit(username, password, host, port, dbname string) *gorm.DB {

	fmt.Println("222")
	DBconn := InitMySQL(username, password, host, port, dbname)
	fmt.Println("333")
	return DBconn
}

// 建立資料庫
func creatDataBase(username, password, host, port, dbname string) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/",
		username,
		password,
		host,
		port)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("無法打開 mysql", err.Error())
	}
	_, err = db.Exec("CREATE DATABASE " + dbname)
	if err != nil {
		return
	}

}

func InitMySQL(username, password, host, port, dbname string) *gorm.DB {

	creatDataBase(username, password, host, port, dbname)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		username, password, host, port, dbname)

	logger := logger.New(
		logrus.NewWriter(),
		logger.Config{
			SlowThreshold: time.Microsecond,
			LogLevel:      logger.Warn,
			Colorful:      false,
		},
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: logger})
	if err != nil {
		fmt.Println("asdaaaaa")
		log.Println("asd", dsn)
		log.Fatal("連接數據庫失敗111, err: ", err.Error())
	}

	if err = db.Use(tracing.NewPlugin()); err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&model.Product{}, &model.Order{})

	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(100)
	sqlDB.SetMaxOpenConns(5000)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db
}
